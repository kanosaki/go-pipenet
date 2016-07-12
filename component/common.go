package component

import "github.com/kanosaki/go-pipenet/core"

type DelegateController struct {
	push func(port core.PortKey, data *core.Packet)
	pull func(port core.PortKey, param *core.DrainRequest) *core.DrainResponse
}

func NewDelegateController(
push func(port core.PortKey, data *core.Packet),
pull func(port core.PortKey, param *core.DrainRequest) *core.DrainResponse) *DelegateController {
	return &DelegateController{
		push,
		pull,
	}
}

func (self *DelegateController) Concrete(joint *core.MetaJoint, graph *core.MetaGraph) error {
	return nil
}
func (self *DelegateController) Push(port core.PortKey, data *core.Packet) {
	self.push(port, data)
}
func (self *DelegateController) Pull(port core.PortKey, param *core.DrainRequest) *core.DrainResponse {
	return self.pull(port, param)
}

