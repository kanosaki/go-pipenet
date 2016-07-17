package pipenet

import (
	"github.com/kanosaki/go-pipenet/core"
	"github.com/kanosaki/go-pipenet/component"
	"github.com/kanosaki/go-pipenet/storage"
	"strings"
)

func Port(port string) core.Endpoint {
	chunks := strings.Split(port, ":")
	switch len(chunks) {
	case 1:
		return core.Endpoint{
			Joint: core.GRAPH,
			Port: core.PortKey(chunks[0]),
		}
	case 2:
		return core.Endpoint{
			Joint: core.JointKey(chunks[0]),
			Port: core.PortKey(chunks[1]),
		}
	default:
		panic("Invalid port repr")
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

