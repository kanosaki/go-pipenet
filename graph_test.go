package pipenet

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/kanosaki/go-pipenet/component"
	"github.com/kanosaki/go-pipenet/core"
	"fmt"
	"strings"
)

func SimplePacket(data interface{}) *core.Packet {
	pkt := core.NewPacket()
	pkt.Set("data", data)
	return pkt
}

func TestJson(t *testing.T) {
	graphDef :=
	`{
		"inlets": ["in0", "in1"],
		"outlets": ["out"],
		"joints": {
			"j1": {
				"type": "merge",
				"inlets": ["in0", "in1"],
				"outlets": ["out"],
				"param": {}
			}
		},
		"pipes": [
			[":in0", "j1:in0"],
			[":in1", "j1:in1"],
			["j1:out", ":out"]
		]
	}`
	assert := assert.New(t)
	mGraph, err := FromJson(strings.NewReader(graphDef))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	sink := core.NewBufferTerminator()
	mGraph.Sink("out", sink)
	mGraph.Concrete()
	mGraph.Push("in0", SimplePacket("foo"))
	mGraph.Push("in1", SimplePacket("bar"))
	assert.Equal([]*core.Packet{
		SimplePacket("foo"),
		SimplePacket("bar"),
	}, sink.ToArray())
}

func TestMultiHop(t *testing.T) {
	assert := assert.New(t)
	mGraph := Create()
	merge1, err := mGraph.AddJointByComponent("", &component.MergeParam{})
	if err != nil {
		t.FailNow()
	}
	merge2, err := mGraph.AddJointByComponent("", &component.MergeParam{})
	if err != nil {
		t.FailNow()
	}
	merge3, err := mGraph.AddJointByComponent("", &component.MergeParam{})
	if err != nil {
		t.FailNow()
	}
	merge1.DefineInlet("m_in1", "m_in2")
	merge2.DefineInlet("m_in1", "m_in2")
	merge3.DefineInlet("m_in1", "m_in2")
	mGraph.AddBridge(core.GRAPH, "in0", merge1.Key, "m_in1")
	mGraph.AddBridge(core.GRAPH, "in1", merge1.Key, "m_in2")
	mGraph.AddBridge(core.GRAPH, "in2", merge2.Key, "m_in1")
	mGraph.AddBridge(core.GRAPH, "in3", merge2.Key, "m_in2")
	mGraph.AddBridge(merge1.Key, "out", merge3.Key, "m_in1")
	mGraph.AddBridge(merge2.Key, "out", merge3.Key, "m_in2")
	mGraph.AddBridge(merge3.Key, core.PORT_DEFAULT_OUT, core.GRAPH, "out")
	sink := core.NewBufferTerminator()
	mGraph.Sink("out", sink)
	if mGraph.Concrete() != nil {
		t.FailNow()
	}
	mGraph.Push("in0", SimplePacket("foo"))
	mGraph.Push("in1", SimplePacket("bar"))
	mGraph.Push("in2", SimplePacket("baz"))
	mGraph.Push("in3", SimplePacket("hoge"))
	assert.Equal([]*core.Packet{
		SimplePacket("foo"),
		SimplePacket("bar"),
		SimplePacket("baz"),
		SimplePacket("hoge"),
	}, sink.ToArray())
}

const DOUBLE_STEP_MERGE =
`{
	"inlets": ["in0", "in1", "in2", "in3"],
	"outlets": ["out"],
	"joints": {
		"j1": {
			"type": "merge",
			"inlets": ["in0", "in1"],
			"outlets": ["out"],
			"param": {}
		},
		"j2": {
			"type": "merge",
			"inlets": ["in0", "in1"],
			"outlets": ["out"]
		},
		"j3": {
			"type": "merge",
			"inlets": ["in0", "in1"],
			"outlets": ["out"],
			"param": {}
		}
	},
	"pipes": [
		[":in0", "j1:in0"],
		[":in1", "j1:in1"],
		[":in2", "j2:in0"],
		[":in3", "j2:in1"],
		["j1:out", "j3:in0"],
		["j2:out", "j3:in1"],
		["j3:out", ":out"]
	]
}`

const DOUBLE_STEP_MERGE_V2 =
`{
	"inlets": {
		"in0": ["j1", "in0"],
		"in1": ["j1", "in1"]
	}
	"outlets": ["out"],
	"joints": {
		"j1": {
			"type": "merge",
			"inlets": ["in0", "in1"],
			"outlets": {
				"out": ["j3", "in0"]
			},
			"param": {}
		},
		"j2": {
			"type": "merge",
			"inlets": ["in0", "in1"],
			"outlets": {
				"out": ["j3", "in1"]
			},
			"param": {}
		},
		"j3": {
			"type": "merge",
			"inlets": ["in0", "in1"],
			"outlets": {
				"out": ["", "out"]
			},
			"param": {}
		}
	},
	"pipes": [
		[["", "in0"], ["j1", "in0"]],
		[["", "in1"], ["j1", "in1"]]
	]
}`

func TestMultiHopFromJson(t *testing.T) {
	graphDef := DOUBLE_STEP_MERGE
	assert := assert.New(t)
	mGraph, err := FromJson(strings.NewReader(graphDef))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	sink := core.NewBufferTerminator()
	mGraph.Sink("out", sink)
	if mGraph.Concrete() != nil {
		t.FailNow()
	}
	mGraph.Push("in0", SimplePacket("foo"))
	mGraph.Push("in1", SimplePacket("bar"))
	mGraph.Push("in2", SimplePacket("baz"))
	mGraph.Push("in3", SimplePacket("hoge"))
	assert.Equal([]*core.Packet{
		SimplePacket("foo"),
		SimplePacket("bar"),
		SimplePacket("baz"),
		SimplePacket("hoge"),
	}, sink.ToArray())
}

func BenchmarkMultiHop(b *testing.B) {
	graphDef := DOUBLE_STEP_MERGE
	mGraph, err := FromJson(strings.NewReader(graphDef))
	if err != nil {
		b.FailNow()
	}
	sink := core.NewBufferTerminator()
	mGraph.Sink("out", sink)
	if mGraph.Concrete() != nil {
		b.FailNow()
	}
	packets := make([]*core.Packet, 0, 400)
	for num := 0; num < 400; num++ {
		packets = append(packets, SimplePacket(fmt.Sprintf("packet_%d", num)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			mGraph.Push("in0", packets[4 * j])
			mGraph.Push("in1", packets[4 * j + 1])
			mGraph.Push("in2", packets[4 * j + 2])
			mGraph.Push("in3", packets[4 * j + 3])
		}
		sink.Clear()
	}
	b.StopTimer()
}
