package core

import "fmt"

// Packet routing error
// by INTERNAL problem
type DispatchFailed struct {
	Destination JointKey
	Data        *Packet
}

func (self *DispatchFailed) Error() string {
	return fmt.Sprintf("DispatchFailed: dst: %v, data: %v", self.Destination, self.Data)
}

// Packet routing fail
// caused by graph topology or port missing
type PacketUnreachable struct {
}

func (self *PacketUnreachable) Error() string {
	return fmt.Sprintf("Unreachable: ") // TODO: decide spec of this error
}

type UndefinedPort struct {
	At   string
	Port PortKey
}

func (self *UndefinedPort) Error() string {
	return fmt.Sprintf("Port %s undefined at %s", self.Port, self.At)
}
