package core

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

type PortKey string

const (
	PORT_ANY PortKey = "_"
	PORT_EMPTY PortKey = ""
	PORT_DEFAULT_OUT PortKey = "out"
	PORT_DEFAULT_IN PortKey = "in"
)

type Endpoint struct {
	Joint JointKey
	Port  PortKey
}

func (mg Endpoint) String() string {
	return fmt.Sprintf("[%s:%s]", mg.Joint, mg.Port)
}

type PerformanceFlavor int

const (
	// less latency
	FlavorBetterLatency PerformanceFlavor = 1 << iota
	// more throughput
	FlavorBetterThroughput
	// less memory consumption
	FlavorBetterFootprint
)

type MetaGraph struct {
	Universe *Universe
	Flavor   PerformanceFlavor
	Pipes    []*JointBridge
	Joints   map[JointKey]*MetaJoint
	sinks    map[PortKey]Pipe
	pools    map[PortKey]Pipe
	IdGen    func() JointKey
}

func NewMetaGraph(univ *Universe) *MetaGraph {
	return &MetaGraph{
		Joints: make(map[JointKey]*MetaJoint),
		Flavor: FlavorBetterLatency,
		sinks: make(map[PortKey]Pipe),
		pools: make(map[PortKey]Pipe),
		Universe: univ,
	}
}

func (mg *MetaGraph) nextID() JointKey {
	if mg.IdGen == nil {
		mg.IdGen = numericIDGenerator()
	}
	prevID := JOINT_ANY
	for {
		candidate := mg.IdGen()
		if candidate == prevID {
			// prevent infinite loop
			panic(fmt.Sprintf("Invalid ID gen! %s == %s", candidate, prevID))
		}
		if _, ok := mg.Joints[candidate]; !ok {
			return candidate
		}
		prevID = candidate
	}
	panic("Never here")
}

func (mg *MetaGraph) NewJoint(component ComponentKey, key JointKey) *MetaJoint {
	joint := NewMetaJoint(mg, component, key)
	return joint
}

func (mg *MetaGraph) AddJoint(joint *MetaJoint) (*MetaJoint, error) {
	if len(joint.Key) == 0 {
		joint.Key = mg.nextID()
	} else if _, ok := mg.Joints[joint.Key]; ok {
		return nil, fmt.Errorf("Duplicate JointKey! %s", joint.Key)
	}
	mg.Joints[joint.Key] = joint
	return joint, nil
}

func (mg *MetaGraph) AddJointByComponent(key JointKey, param ComponentParam) (*MetaJoint, error) {
	component := param.Name()
	if comp, ok := mg.Universe.Components[component]; !ok {
		return nil, fmt.Errorf("Undefined component %s", component)
	} else {
		joint := mg.NewJoint(component, key)
		jc, err := comp.CreateController(joint, param, mg)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create joint!")
		} else {
			joint.controller = jc
			return mg.AddJoint(joint)
		}
	}
}

func (mg *MetaGraph) AddPipeBridge(from Endpoint, to Endpoint) error {
	mg.Pipes = append(mg.Pipes, &JointBridge{
		Source: from,
		Destination: to,
		Mode: FlavorToMode(mg.Flavor),
	})
	return nil
}

func (mg *MetaGraph) AddBridge(fromJoint JointKey, fromPort PortKey, toJoint JointKey, toPort PortKey) {
	mg.AddPipeBridge(Endpoint{fromJoint, fromPort}, Endpoint{toJoint, toPort})
}

func (mg *MetaGraph) SinkHandler(port PortKey, handler func(*Packet)) {
	mg.Sink(port, NewFuncTerminator(handler))
}

func (mg *MetaGraph) Sink(port PortKey, handler Pipe) {
	mg.sinks[port] = handler
}

func (mg *MetaGraph) Source(port PortKey, handler Pipe) {
	mg.pools[port] = handler
}

// dynamic routing
// for internal use
func (mg *MetaGraph) SendToNode(ep Endpoint, data *Packet) {
	if ep.Joint == GRAPH {
		mg.dispatchOutlet(ep.Port, data)
	} else {
		if downNode, ok := mg.Joints[ep.Joint]; !ok {
			mg.TellError(nil, &DispatchFailed{
				Destination: ep.Joint,
				Data: data,
			})
		} else {
			downNode.Push(ep.Port, data)
		}
	}
}

// for internal use
func (mg *MetaGraph) DrainFromNode(ep Endpoint, data *DrainRequest) *DrainResponse {
	if ep.Joint == GRAPH {
		return mg.pullPool(ep.Port, data)
	} else {
		if upNode, ok := mg.Joints[ep.Joint]; !ok {
			mg.TellError(nil, &DispatchFailed{
				Destination: ep.Joint,
				Data: data,
			})
			return nil // TODO return error object
		} else {
			return upNode.Pull(ep.Port, data)
		}
	}
}

func (mg *MetaGraph) dispatchOutlet(port PortKey, data *Packet) {
	if out, ok := mg.sinks[port]; ok {
		out.Send(data)
	} else {
		mg.TellError(nil, fmt.Errorf("Undefined outlet! %s", port))
	}
}

func (mg *MetaGraph) pullPool(port PortKey, param *DrainRequest) *DrainResponse {
	if pool, ok := mg.pools[port]; ok {
		return pool.Drain(param)
	} else {
		mg.TellError(nil, fmt.Errorf("Missing pool! %s", port))
		return nil
	}
}

func (mg *MetaGraph) TellError(on Node, err error) {
	if on == nil {
		fmt.Fprintf(os.Stderr, "ERR(GRAPH): %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "ERR(JOINT: %s): %v\n", on, err)
	}
}

// External -- push --> Internal
func (mg *MetaGraph) Push(inlet PortKey, data *Packet) {
	initBridges := mg.SelectBridges(GRAPH, inlet, JOINT_ANY, PORT_ANY)
	if len(initBridges) > 0 {
		mg.SendToNode(initBridges[0].Destination, data)
	} else {
		mg.TellError(nil, fmt.Errorf("Destination unreachable %s", inlet))
	}
}

// This API basically for sending control messages.
// 1 Internal <-- request -- External
// 2 Internal -- response --> External
func (mg *MetaGraph) Pull(outlet PortKey, param *DrainRequest) *DrainResponse {
	initBridges := mg.SelectBridges(JOINT_ANY, PORT_ANY, GRAPH, outlet)
	if len(initBridges) > 0 {
		return mg.DrainFromNode(initBridges[0].Source, param)
	} else {
		mg.TellError(nil, fmt.Errorf("Destination unreachable %s", outlet))
		return nil
	}
}

func (mg *MetaGraph) SelectBridges(fromJoint JointKey, fromPort PortKey, toJoint JointKey, toPort PortKey) []*JointBridge {
	var ret []*JointBridge
	for _, br := range mg.Pipes {
		if (fromJoint == JOINT_ANY || fromJoint == br.Source.Joint) &&
			(fromPort == PORT_ANY || fromPort == br.Source.Port) &&
			(toJoint == JOINT_ANY || toJoint == br.Destination.Joint) &&
			(toPort == PORT_ANY || toPort == br.Destination.Port) {
			ret = append(ret, br)
		}
	}
	return ret
}

func (mg *MetaGraph) JointOutlets(jointKey JointKey) []Pipe {
	bridges := mg.SelectBridges(jointKey, PORT_ANY, JOINT_ANY, PORT_ANY)
	ret := make([]Pipe, 0, len(bridges))
	for _, br := range bridges {
		ret = append(ret, NewDelegatePipe(mg, br.Source, br.Destination))
	}
	return ret
}

func (mg *MetaGraph) JointInlets(jointKey JointKey) []Pipe {
	bridges := mg.SelectBridges(JOINT_ANY, PORT_ANY, jointKey, PORT_ANY)
	ret := make([]Pipe, 0, len(bridges))
	for _, br := range bridges {
		ret = append(ret, NewDelegatePipe(mg, br.Source, br.Destination))
	}
	return ret
}

func (mg *MetaGraph) Concrete() error {
	for _, j := range mg.Joints {
		err := j.Concrete(mg)
		if err != nil {
			return err
		}
	}
	return nil
}
