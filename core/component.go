package core

import (
	"encoding/json"
	"github.com/ugorji/go/codec"
)

type ComponentKey string

type ComponentParam interface {
	Name() ComponentKey
}

type EmptyComponentParam struct {
	ComponentName ComponentKey
}

func (self *EmptyComponentParam) Name() ComponentKey {
	return self.ComponentName
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
	Pull(port PortKey, param *DrainRequest) *DrainResponse
	Concrete(self *MetaJoint, graph *MetaGraph) error
}
