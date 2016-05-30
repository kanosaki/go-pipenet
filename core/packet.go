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

func NewPacket_Single(data interface{}) *Packet {
	pkt := NewPacket()
	pkt.Set("data", data)
	return pkt
}

func NewPacket_OK() *Packet {
	pkt := NewPacket()
	pkt.Set("status", "ok")
	return pkt
}

func NewPacket_Error(message string) *Packet {
	pkt := NewPacket()
	pkt.Set("status", "error")
	pkt.Set("message", message)
	return pkt
}
