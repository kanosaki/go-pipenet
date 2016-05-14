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
	ConfigureJoint(metaJoint *MetaJoint, param interface{}) (*MetaJoint, error)
	Save(joint *MetaJoint)
	Restore()
}
