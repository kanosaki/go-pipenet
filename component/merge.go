package component

import (
	"github.com/kanosaki/go-pipenet/core"
	"encoding/json"
	"github.com/ugorji/go/codec"
	"fmt"
)

const (
	KEY_MERGE core.ComponentKey = "merge"
)

type Merge struct {
}

type MergeParam struct {
}

func (m *MergeParam) Name() core.ComponentKey {
	return KEY_MERGE
}

func (m *Merge) Name() core.ComponentKey {
	return KEY_MERGE
}

func (m *Merge) CreateController(metaJoint *core.MetaJoint, param interface{}, graph *core.MetaGraph) (core.JointController, error) {
	return &MergeController{}, nil
}

func (m *Merge) Save(joint *core.MetaJoint) {

}
func (m *Merge) Restore() {

}

func (m *Merge) DecodeParam(decoder *codec.Decoder, data json.RawMessage) (core.ComponentParam, error) {
	ret := &MergeParam{}
	err := decoder.Decode(ret)
	return ret, err
}

type MergeController struct {
	inlets         []core.Pipe
	outlets        []core.Pipe
	currentOutput  int
	currentInlet   int
	oddsEndsBuffer []*core.Packet
}

func (mc *MergeController) Push(port core.PortKey, data *core.Packet) {
	// Round robbin
	// set next outlet
	mc.currentOutput = (mc.currentOutput + 1) % len(mc.outlets)
	mc.outlets[mc.currentOutput].Send(data)
}

func (mc *MergeController) Pull(port core.PortKey, param *core.DrainRequest) *core.DrainResponse {
	ret := mc.oddsEndsBuffer
	for len(ret) < param.Count && len(mc.inlets) > mc.currentInlet {
		resFromUpstream := mc.inlets[mc.currentInlet].Drain(param)
		if resFromUpstream != nil && len(resFromUpstream.Items) > 0 {
			if len(resFromUpstream.Items) + len(ret) > param.Count {
				cutAt := param.Count - len(ret)
				ret = append(ret, resFromUpstream.Items[:cutAt]...)
				mc.oddsEndsBuffer = resFromUpstream.Items[cutAt:]
			} else {
				ret = append(ret, resFromUpstream.Items...)
			}
		} else {
			mc.currentInlet += 1
		}
	}
	return &core.DrainResponse{
		Items: ret,
	}
}

func (mc *MergeController) Concrete(metaJoint *core.MetaJoint, graph *core.MetaGraph) error {
	// initialize inlets and outlets
	outlets := graph.JointOutlets(metaJoint.Key)
	inlets := graph.JointInlets(metaJoint.Key)
	if len(inlets) == 0 {
		return fmt.Errorf("MergeController requires one or more inlets")
	}
	if len(outlets) == 0 {
		return fmt.Errorf("MergeController requires one or more outlets")
	}
	mc.inlets = inlets
	mc.outlets = outlets
	mc.currentOutput = -1
	mc.currentInlet = 0
	return nil
}
