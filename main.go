package pipenet

import (
	"github.com/kanosaki/go-pipenet/core"
	"github.com/kanosaki/go-pipenet/component"
	"github.com/kanosaki/go-pipenet/storage"
)

func Port(jointKey core.JointKey, portKey core.PortKey) core.Endpoint {
	return core.Endpoint{
		Joint: jointKey,
		Port: portKey,
	}
}

func Create(components ...core.Component) *core.MetaGraph {
	comps := make([]core.Component, 0, len(component.Builtins) + len(components))
	comps = append(comps, components...)
	comps = append(comps, component.Builtins...)
	univ := core.NewUniverse(comps, &storage.NullStorage{})
	ret := core.NewMetaGraph(univ)
	return ret
}

