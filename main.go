package pipenet

import (
	"github.com/kanosaki/go-pipenet/core"
	"github.com/kanosaki/go-pipenet/component"
)

func Port(jointKey core.JointKey, portKey core.PortKey) core.Endpoint {
	return core.Endpoint{
		Joint: jointKey,
		Port: portKey,
	}
}

func Create() *core.MetaGraph {
	univ := core.NewUniverse(component.Builtins)
	ret := core.NewMetaGraph(univ)
	return ret
}

