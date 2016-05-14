package core

import (
	"fmt"
)

type JointKey string

const (
	JOINT_KEY_AUTO JointKey = ""
	GRAPH JointKey = ""
)

type PortDirection int

const (
	DIRECTION_FORWARD PortDirection = iota
	DIRECTION_BACKWARD
)

type Node interface {
	Push(port PortKey, data *Packet)
}

type PacketHandler func(from PortKey, data *Packet)

type FactoryParams interface {
	Inlets() PipeSpace
	Outlets() PipeSpace
}

type PacketHandlerFactory func(param FactoryParams) (PacketHandler, error)

var (
//DEFAULT_BACKWARD_HANDLER = func(joint *Joint, gr *Graph) (PacketHandler, error) {
//	return func(from PortKey, data *Packet) {
//		for _, input := range joint.Inputs() {
//			input.Send(data)
//		}
//	}, nil
//}
)

type MetaJoint struct {
	Component       ComponentKey
	Key             JointKey
	Graph           *MetaGraph
	forwardHandler  PacketHandlerFactory
	backwardHandler PacketHandlerFactory
	fwd, back       PacketHandler
	NodeBase
}

func NewMetaJoint(graph *MetaGraph, component ComponentKey, key JointKey) *MetaJoint {
	return &MetaJoint{
		Graph: graph,
		Key: key,
		Component: component,
		NodeBase: newNodeBase(),
	}
}

func (self *MetaJoint) Forward(fn PacketHandlerFactory) {
	self.forwardHandler = fn
}

func (self *MetaJoint) Backward(fn PacketHandlerFactory) {
	self.backwardHandler = fn
}

func (self *MetaJoint) String() string {
	return fmt.Sprintf("<%s(%s)>", self.Component, self.Key)
}

func (self *MetaJoint) Push(port PortKey, data *Packet) {
	if _, ok := self.Inlet(port); ok {
		self.fwd(port, data)
	} else {
		fmt.Printf("UNREACHABLE: undefined port %s\n", port)
	}
}

func (self *MetaJoint) Concrete(graph *MetaGraph) error {
	if self.forwardHandler != nil {
		fwd, err := self.forwardHandler(self)
		if err != nil {
			return err
		} else {
			self.fwd = fwd
		}
	}
	if self.backwardHandler != nil {
		back, err := self.backwardHandler(self)
		if err != nil {
			return err
		} else {
			self.back = back
		}
	}
	return nil
}
