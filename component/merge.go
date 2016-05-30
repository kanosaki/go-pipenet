package component

import (
	"github.com/kanosaki/go-pipenet/core"
	"encoding/json"
	"github.com/ugorji/go/codec"
)

const (
	KEY_MERGE core.ComponentKey = "merge"
)

type Merge struct {
}

type MergeParam struct {
}

func (self *MergeParam) Name() core.ComponentKey {
	return KEY_MERGE
}

func (self *Merge) Name() core.ComponentKey {
	return KEY_MERGE
}

func (self *Merge) CreateController(metaJoint *core.MetaJoint, param interface{}) (core.JointController, error) {
	return NewDelegateController(
		func(port core.PortKey, data *core.Packet) {
			if out, ok := metaJoint.Outlet(core.PORT_DEFAULT_OUT); ok {
				out.Send(data)
			} else {
				panic("Unreachable")
			}
		}, func(port core.PortKey, param *core.Packet) *core.Packet {
			panic("NIE")
		}), nil
}

func (self *Merge) Save(joint *core.MetaJoint) {

}
func (self *Merge) Restore() {

}

func (self *Merge) DecodeParam(decoder *codec.Decoder, data json.RawMessage) (core.ComponentParam, error) {
	ret := &MergeParam{}
	err := decoder.Decode(ret)
	return ret, err
}
