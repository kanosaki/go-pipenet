package pipenet

import (
	"io"
	"github.com/kanosaki/go-pipenet/core"
	"github.com/ugorji/go/codec"
	"github.com/pkg/errors"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	ENDPOINT_SEPARATOR = ":"
)

var (
	jsonHandle *codec.JsonHandle = &codec.JsonHandle{}
)

type GraphInfo struct {
	Joints  map[core.JointKey]*JointInfo `codec:"joints"`
	Pipes   []*PipeInfo `codec:"pipes"`
	Inlets  []core.PortKey `codec:"inlets"`
	Outlets []core.PortKey `codec:"outlets"`
}

type JointInfo struct {
	Component core.ComponentKey `codec:"type"`
	Inlets    []core.PortKey `codec:"inlets"`
	Outlets   []core.PortKey `codec:"outlets"`
	Param     json.RawMessage `codec:"param"`
}

type PipeInfo struct {
	_struct     bool `codec:",toarray"`
	Source      EndpointInfo
	Destination EndpointInfo
	Mode        string
}

type EndpointInfo string

func (self EndpointInfo) Joint() core.JointKey {
	return core.JointKey(self[:strings.Index(string(self), ENDPOINT_SEPARATOR)])
}

func (self EndpointInfo) Port() core.PortKey {
	return core.PortKey(self[strings.Index(string(self), ENDPOINT_SEPARATOR) + 1:])
}

type ComponentParamInfo struct {
}

// Construct MetaGraph from Document

func FromDocument(reader io.Reader, handle codec.Handle) (*core.MetaGraph, error) {
	dec := codec.NewDecoder(reader, handle)
	info := &GraphInfo{}
	err := dec.Decode(info)
	if err != nil {
		return nil, errors.Wrap(err, "JSON Decode failed")
	}
	mGraph := Create()
	for jKey, jInfo := range info.Joints {
		var param core.ComponentParam
		if component, ok := mGraph.Universe.Components[jInfo.Component]; !ok {
			return nil, fmt.Errorf("Undefined component %s", jInfo.Component)
		} else {
			paramDecoder := codec.NewDecoderBytes(jInfo.Param, handle)
			if len(jInfo.Param) != 0 {
				param, err = component.DecodeParam(paramDecoder, jInfo.Param)
				if err != nil {
					return nil, errors.Wrapf(err, "Failed to decode compoennt param for %s at %s", jInfo.Component, jKey)
				}
			} else {
				param = &core.EmptyComponentParam{jInfo.Component}
			}
		}
		mJoint, err := mGraph.AddJointByComponent(jKey, param)
		if err != nil {
			return nil, errors.Wrapf(err, "Error during adding joint %s", jKey)
		}
		mJoint.DefineInlet(jInfo.Inlets...)
		mJoint.DefineOutlet(jInfo.Outlets...)
	}
	for _, pInfo := range info.Pipes {
		mGraph.AddBridge(
			pInfo.Source.Joint(),
			pInfo.Source.Port(),
			pInfo.Destination.Joint(),
			pInfo.Destination.Port())
	}
	return mGraph, nil
}

func FromJson(reader io.Reader) (*core.MetaGraph, error) {
	return FromDocument(reader, jsonHandle)
}
