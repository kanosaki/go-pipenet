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
type PacketBackwardHandler func(from PortKey, param *Packet) *Packet

type FactoryParams interface {
	Inlets() PipeSpace
	Outlets() PipeSpace
}

type PacketHandlerFactory func(param FactoryParams) (PacketHandler, error)
type PacketBackwardHandlerFactory func(param FactoryParams) (PacketBackwardHandler, error)

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
	Component  ComponentKey
	Key        JointKey
	graph      *MetaGraph
	controller JointController
	NodeBase
}

func NewMetaJoint(graph *MetaGraph, component ComponentKey, key JointKey) *MetaJoint {
	return &MetaJoint{
		graph: graph,
		Key: key,
		Component: component,
		NodeBase: newNodeBase(),
	}
}


func (self *MetaJoint) String() string {
	return fmt.Sprintf("<%s(%s)>", self.Component, self.Key)
}

func (self *MetaJoint) Push(port PortKey, data *Packet) {
	self.controller.Push(port, data)
}

func (self *MetaJoint) Pull(port PortKey, param *Packet) *Packet {
	return self.controller.Pull(port, param)
}

func (self *MetaJoint) Concrete(graph *MetaGraph) error {
	return self.controller.Concrete(self, graph)
}
