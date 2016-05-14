package core

type Universe struct {
	Components map[ComponentKey]Component
}

func NewUniverse(components []Component) *Universe {
	compMap := make(map[ComponentKey]Component, len(components))
	for _, comp := range components {
		compMap[comp.Name()] = comp
	}
	return &Universe{
		Components: compMap,
	}
}