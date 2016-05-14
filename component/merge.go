package component

import (
	"github.com/kanosaki/go-pipenet/core"
	"fmt"
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

func (self *Merge) ConfigureJoint(mJoint *core.MetaJoint, param interface{}) (*core.MetaJoint, error) {
	mJoint.Forward(func(param core.FactoryParams) (core.PacketHandler, error) {
		out0, ok := param.Outlets().Lookup(core.PORT_DEFAULT_OUT)
		if !ok {
			return nil, fmt.Errorf("Unable find to output port %s", core.PORT_DEFAULT_OUT)
		}
		return func(from core.PortKey, data *core.Packet) {
			out0.Send(data)
		}, nil
	})
	return mJoint, nil
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
