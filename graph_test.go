package pipenet

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/kanosaki/go-pipenet/component"
	"github.com/kanosaki/go-pipenet/core"
	"fmt"
	"strings"
	"github.com/kanosaki/go-pipenet/storage"
)

var univ = core.NewUniverse(component.Builtins, storage.NewNullStorage())

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
					"type": "merge"
				}
			},
			"pipes": [
				[":in0", "j1:in0"],
				[":in1", "j1:in1"],
				["j1:out", ":out"]
			]
		}`
	assert := assert.New(t)
	mGraph, err := storage.FromJson(strings.NewReader(graphDef), univ)
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
	//merge1.DefineInlet("m_in1", "m_in2")
	//merge2.DefineInlet("m_in1", "m_in2")
	//merge3.DefineInlet("m_in1", "m_in2")
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

func TestMultiHopFromJson(t *testing.T) {
	graphDef := DOUBLE_STEP_MERGE
	assert := assert.New(t)
	mGraph, err := storage.FromJson(strings.NewReader(graphDef), univ)
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

func TestMultiHopFromJsonDrainJust(t *testing.T) {
	graphDef := DOUBLE_STEP_MERGE
	assert := assert.New(t)
	mGraph, err := storage.FromJson(strings.NewReader(graphDef), univ)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	mGraph.Source("in0", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("foo1"),
			SimplePacket("foo2")}))
	mGraph.Source("in1", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("bar1"),
			SimplePacket("bar2")}))
	mGraph.Source("in2", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("hoge1"),
			SimplePacket("hoge2")}))
	mGraph.Source("in3", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("fuga1"),
			SimplePacket("fuga2")}))
	if mGraph.Concrete() != nil {
		t.FailNow()
	}
	res := mGraph.Pull(core.PortKey("out"), &core.DrainRequest{8})
	assert.NotNil(res, "Empty response")
	assert.Equal([]*core.Packet{
		SimplePacket("foo1"),
		SimplePacket("foo2"),
		SimplePacket("bar1"),
		SimplePacket("bar2"),
		SimplePacket("hoge1"),
		SimplePacket("hoge2"),
		SimplePacket("fuga1"),
		SimplePacket("fuga2"),
	}, res.Items)
}

func TestMultiHopFromJsonDrainEnoughSource(t *testing.T) {
	graphDef := DOUBLE_STEP_MERGE
	assert := assert.New(t)
	mGraph, err := storage.FromJson(strings.NewReader(graphDef), univ)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	mGraph.Source("in0", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("foo1"),
			SimplePacket("foo2"),
			SimplePacket("foo3"),
			SimplePacket("foo4")}))
	mGraph.Source("in1", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("bar1"),
			SimplePacket("bar2")}))
	mGraph.Source("in2", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("hoge1"),
			SimplePacket("hoge2"),
			SimplePacket("hoge3")}))
	mGraph.Source("in3", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("fuga1"),
			SimplePacket("fuga2"),
			SimplePacket("fuga3")}))
	if mGraph.Concrete() != nil {
		t.FailNow()
	}
	res := mGraph.Pull(core.PortKey("out"), &core.DrainRequest{8})
	assert.NotNil(res, "Empty response")
	assert.Equal([]*core.Packet{
		SimplePacket("foo1"),
		SimplePacket("foo2"),
		SimplePacket("foo3"),
		SimplePacket("foo4"),
		SimplePacket("bar1"),
		SimplePacket("bar2"),
		SimplePacket("hoge1"),
		SimplePacket("hoge2"),
	}, res.Items)
	res = mGraph.Pull(core.PortKey("out"), &core.DrainRequest{8})
	assert.NotNil(res, "Empty response")
	assert.Equal([]*core.Packet{
		SimplePacket("hoge3"),
		SimplePacket("fuga1"),
		SimplePacket("fuga2"),
		SimplePacket("fuga3"),
	}, res.Items)
}

func TestMultiHopFromJsonDrainShortSource(t *testing.T) {
	graphDef := DOUBLE_STEP_MERGE
	assert := assert.New(t)
	mGraph, err := storage.FromJson(strings.NewReader(graphDef), univ)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	mGraph.Source("in0", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("foo2")}))
	mGraph.Source("in1", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("bar2")}))
	mGraph.Source("in2", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("hoge2")}))
	mGraph.Source("in3", core.NewBufferSource(
		[]*core.Packet{
			SimplePacket("fuga2")}))
	if mGraph.Concrete() != nil {
		t.FailNow()
	}
	res := mGraph.Pull(core.PortKey("out"), &core.DrainRequest{8})
	assert.NotNil(res, "Empty response")
	assert.Equal([]*core.Packet{
		SimplePacket("foo2"),
		SimplePacket("bar2"),
		SimplePacket("hoge2"),
		SimplePacket("fuga2"),
	}, res.Items)
}

func BenchmarkMultiHop(b *testing.B) {
	graphDef := DOUBLE_STEP_MERGE
	mGraph, err := storage.FromJson(strings.NewReader(graphDef), univ)
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
