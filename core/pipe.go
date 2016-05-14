package core

import (
	"fmt"
	"container/list"
)

type PipeMode int

const (
	PIPE_DIRECT PipeMode = iota
	PIPE_CHANNEL
	PIPE_ROUTINE
)

type Pipe interface {
	Send(data *Packet)
}

type PipeSpace map[PortKey]Pipe

func (self PipeSpace) Lookup(key PortKey) (Pipe, bool) {
	p, b := self[key]
	return p, b
}

func (self PipeSpace) Register(key PortKey, pipe Pipe) error {
	if _, ok := self[key]; ok {
		return fmt.Errorf("Key %s already exists", key)
	} else {
		self[key] = pipe
		return nil
	}
}

func (self PipeSpace) Define(key PortKey) error {
	if _, ok := self[key]; ok {
		return fmt.Errorf("Key %s already exists", key)
	} else {
		self[key] = &DirectPipe{

		}
		return nil
	}
}

func FlavorToMode(flavor PerformanceFlavor) PipeMode {
	throughput := flavor & FlavorBetterLatency != 0
	latency := flavor & FlavorBetterLatency != 0
	footprint := flavor & FlavorBetterFootprint != 0
	switch {
	// more intelligent flow control?
	case latency:
		return PIPE_DIRECT
	case footprint:
		return PIPE_CHANNEL
	case throughput:
		return PIPE_ROUTINE
	default:
		return PIPE_CHANNEL
	}
}

type JointBridge struct {
	Source      Endpoint
	Destination Endpoint
	Mode        PipeMode
}

func (self *JointBridge) Repr() string {
	return fmt.Sprintf("[%s:%s-%s:%s]", self.Source.Joint, self.Source.Port, self.Destination.Joint, self.Destination.Port)
}

func NewMetaPipe(flavor PerformanceFlavor) *JointBridge {
	return &JointBridge{
		Mode: FlavorToMode(flavor),
	}
}

type DelegatePipe struct {
	delegate *MetaGraph
	target   Endpoint
}

func NewDelegatePipe(graph *MetaGraph, target Endpoint) *DelegatePipe {
	return &DelegatePipe{
		delegate: graph,
		target: target,
	}
}

func (self *DelegatePipe) Send(data *Packet) {
	self.delegate.Dispatch(self.target, data)
}

// call destination's method directly
type DirectPipe struct {
	dstJoint Node
	dstPort  PortKey
}

func NewDirectPipe(target Node, port PortKey) DirectPipe {
	return DirectPipe{
		dstJoint: target,
		dstPort: port,
	}
}

func (self DirectPipe) Send(data *Packet) {
	self.dstJoint.Push(self.dstPort, data)
}

func (self DirectPipe) Close() {
}

func (self DirectPipe) Repr() string {
	return fmt.Sprintf("--> %v(%v)", self.dstJoint, self.dstPort)
}
//
// Utility pipes
//
type FuncTerminator struct {
	handler func(*Packet)
}

func (self *FuncTerminator) Send(data *Packet) {
	self.handler(data)
}

func NewFuncTerminator(handler func(*Packet)) *FuncTerminator {
	return &FuncTerminator{
		handler: handler,
	}
}

type BufferTerminator struct {
	buf *list.List
}

func NewBufferTerminator() *BufferTerminator {
	return &BufferTerminator{
		buf: list.New(),
	}
}

func (self *BufferTerminator) Send(data *Packet) {
	self.buf.PushBack(data)
}

func (self *BufferTerminator) ToArray() []*Packet {
	ret := make([]*Packet, 0, self.buf.Len())
	for e := self.buf.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*Packet))
	}
	return ret
}

func (self *BufferTerminator) Clear() {
	self.buf = list.New()
}

func (self *BufferTerminator) Len() int {
	return self.buf.Len()
}

