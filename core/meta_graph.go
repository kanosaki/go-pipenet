package core

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

type PortKey string

const (
	PORT_ANY PortKey = "_"
	PORT_DEFAULT_OUT PortKey = "out"
	PORT_DEFAULT_IN PortKey = "in"
)

type Endpoint struct {
	Joint JointKey
	Port  PortKey
}

func (self Endpoint) Repr() string {
	return fmt.Sprintf("[%s:%s]", self.Joint, self.Port)
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
	Inputs   map[PortKey]Pipe
	IdGen    func() JointKey
}

func NewMetaGraph(univ *Universe) *MetaGraph {
	return &MetaGraph{
		Joints: make(map[JointKey]*MetaJoint),
		Flavor: FlavorBetterLatency,
		sinks: make(map[PortKey]Pipe),
		Inputs: make(map[PortKey]Pipe),
		Universe: univ,
	}
}

func (self *MetaGraph) nextID() JointKey {
	if self.IdGen == nil {
		self.IdGen = numericIDGenerator()
	}
	prevID := JOINT_KEY_AUTO
	for {
		candidate := self.IdGen()
		if candidate == prevID {
			// prevent infinite loop
			panic(fmt.Sprintf("Invalid ID gen! %s == %s", candidate, prevID))
		}
		if _, ok := self.Joints[candidate]; !ok {
			return candidate
		}
		prevID = candidate
	}
	panic("Never here")
}

func (self *MetaGraph) NewJoint(component ComponentKey, key JointKey) *MetaJoint {
	joint := NewMetaJoint(self, component, key)
	return joint
}

func (self *MetaGraph) AddJoint(joint *MetaJoint) (*MetaJoint, error) {
	if len(joint.Key) == 0 {
		joint.Key = self.nextID()
	} else if _, ok := self.Joints[joint.Key]; ok {
		return nil, fmt.Errorf("Duplicate JointKey! %s", joint.Key)
	}
	if len(joint.outlets) == 0 {
		joint.outlets[PORT_DEFAULT_OUT] = nil
	}
	if len(joint.inlets) == 0 {
		joint.inlets[PORT_DEFAULT_IN] = nil
	}
	self.Joints[joint.Key] = joint
	return joint, nil
}

func (self *MetaGraph) AddJointByComponent(key JointKey, param ComponentParam) (*MetaJoint, error) {
	component := param.Name()
	if comp, ok := self.Universe.Components[component]; !ok {
		return nil, fmt.Errorf("Undefined component %s", component)
	} else {
		joint := self.NewJoint(component, key)
		joint, err := comp.ConfigureJoint(joint, param)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create joint!")
		} else {
			return self.AddJoint(joint)
		}
	}
}

func (self *MetaGraph) AddPipeBridge(from Endpoint, to Endpoint) error {
	self.Pipes = append(self.Pipes, &JointBridge{
		Source: from,
		Destination: to,
		Mode: FlavorToMode(self.Flavor),
	})
	if targetJoint, ok := self.Joints[from.Joint]; !ok {
		return fmt.Errorf("Undefined joint! %s", from.Joint)
	} else {
		targetJoint.SetOutlet(from.Port, NewDelegatePipe(self, to))
		return nil
	}
}

func (self *MetaGraph) AddBridge(fromJoint JointKey, fromPort PortKey, toJoint JointKey, toPort PortKey) {
	if fromJoint == GRAPH {
		self.AddInputBridge(fromPort, toJoint, toPort)
	} else if toJoint == GRAPH {
		self.AddOutputBridge(fromJoint, fromPort, toPort)
	} else {
		self.AddPipeBridge(Endpoint{fromJoint, fromPort}, Endpoint{toJoint, toPort})
	}
}

func (self *MetaGraph) AddInputBridge(graphPort PortKey, joint JointKey, jointPort PortKey) error {
	if _, ok := self.Joints[joint]; !ok {
		return fmt.Errorf("Undefined Joint! %s", joint)
	} else {
		self.Inputs[graphPort] = NewDelegatePipe(self, Endpoint{joint, jointPort})
		return nil
	}
}

func (self *MetaGraph) AddOutputBridge(joint JointKey, jointPort PortKey, graphPort PortKey) error {
	if targetJoint, ok := self.Joints[joint]; !ok {
		return fmt.Errorf("Undefined Joint! %s", joint)
	} else {
		targetJoint.SetOutlet(jointPort, NewDelegatePipe(self, Endpoint{GRAPH, graphPort}))
		return nil
	}
}

func (self *MetaGraph) SinkHandler(port PortKey, handler func(*Packet)) {
	self.Sink(port, NewFuncTerminator(handler))
}

func (self *MetaGraph) Sink(port PortKey, handler Pipe) {
	self.sinks[port] = handler
}

// dynamic routing
func (self *MetaGraph) Dispatch(ep Endpoint, data *Packet) {
	if ep.Joint == GRAPH {
		self.dispatchOutlet(ep.Port, data)
	} else {
		if dst, ok := self.Joints[ep.Joint]; !ok {
			self.TellError(nil, &DispatchFailed{
				Destination: ep.Joint,
				Data: data,
			})
		} else {
			dst.Push(ep.Port, data)
		}
	}
}

func (self *MetaGraph) dispatchOutlet(port PortKey, data *Packet) {
	if out, ok := self.sinks[port]; ok {
		out.Send(data)
	} else {
		self.TellError(nil, fmt.Errorf("Undefined outlet! %s", port))
	}
}

func (self *MetaGraph) TellError(on Node, err error) {
	if on == nil {
		fmt.Fprintf(os.Stderr, "ERR(GRAPH): %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "ERR(JOINT: %s): %v\n", on, err)
	}
}

func (self *MetaGraph) Push(port PortKey, data *Packet) {
	if initPipe, ok := self.Inputs[port]; ok {
		if initPipe != nil {
			initPipe.Send(data)
		} else {
			self.TellError(nil, fmt.Errorf("Destination unreachable %s", port))
		}
	} else {
		self.TellError(nil, fmt.Errorf("Undefined port %s", port))
	}
}

func (self *MetaGraph) Concrete() error {
	for _, j := range self.Joints {
		err := j.Concrete(self)
		if err != nil {
			return err
		}
	}
	return nil
}
