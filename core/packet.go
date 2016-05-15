package core

type Packet struct {
	value map[string]interface{}
}

func NewPacket() *Packet {
	return &Packet{
		value: make(map[string]interface{}),
	}
}

func (self *Packet) Set(key string, value interface{}) {
	self.value[key] = value
}

func (self *Packet) Get(key string) (interface{}, bool) {
	v, ok := self.value[key]
	return v, ok
}
