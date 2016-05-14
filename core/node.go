package core

type NodeBase struct {
	inlets  map[PortKey]Pipe
	outlets map[PortKey]Pipe
}

func newNodeBase() NodeBase {
	return NodeBase{
		inlets: map[PortKey]Pipe{
			PORT_DEFAULT_IN: nil,
		},
		outlets: map[PortKey]Pipe{
			PORT_DEFAULT_OUT: nil,
		},
	}
}

func (self *NodeBase) String() string {
	panic("Must override")
}

func (self *NodeBase) DefineInlet(keys ...PortKey) {
	for _, key := range keys {
		self.inlets[key] = nil
	}
}

func (self *NodeBase) DefineOutlet(keys ...PortKey) {
	for _, key := range keys {
		self.outlets[key] = nil
	}
}

func (self *NodeBase) SetInlet(port PortKey, pipe Pipe) error {
	if _, ok := self.inlets[port]; ok {
		self.inlets[port] = pipe
		return nil
	} else {
		return &UndefinedPort{self.String(), port}
	}
}

func (self *NodeBase) SetOutlet(port PortKey, pipe Pipe) error {
	if _, ok := self.outlets[port]; ok {
		self.outlets[port] = pipe
		return nil
	} else {
		return &UndefinedPort{self.String(), port}
	}
}

func (self *NodeBase) Inlets() PipeSpace {
	return self.inlets
}

func (self *NodeBase) Inlet(key PortKey) (Pipe, bool) {
	pipe, ok := self.inlets[key]
	return pipe, ok
}

func (self *NodeBase) Outlets() PipeSpace {
	return self.outlets
}

func (self *NodeBase) Outlet(key PortKey) (Pipe, bool) {
	pipe, ok := self.outlets[key]
	return pipe, ok
}
