package core

import (
	"encoding/json"
	"github.com/ugorji/go/codec"
)

type ComponentKey string

type ComponentParam interface {
	Name() ComponentKey
}

type Component interface {
	DecodeParam(decoder *codec.Decoder, data json.RawMessage) (ComponentParam, error)
	Name() ComponentKey
	CreateController(metaJoint *MetaJoint, param interface{}) (JointController, error)
	Save(joint *MetaJoint)
	Restore()
}

type JointController interface {
	Push(port PortKey, data *Packet)
	Pull(port PortKey, param *Packet) *Packet
	Concrete(self *MetaJoint, graph *MetaGraph) error
}
