package core

type Universe struct {
	Components map[ComponentKey]Component
	Storage    Storage
}

func NewUniverse(components []Component, storage Storage) *Universe {
	compMap := make(map[ComponentKey]Component, len(components))
	for _, comp := range components {
		compMap[comp.Name()] = comp
	}
	return &Universe{
		Components: compMap,
		Storage: storage,
	}
}

func (self *Universe) Save(key string, graph *MetaGraph) error {
	return self.Storage.Save(key, graph, self)
}

func (self *Universe) Load(key string) (*MetaGraph, error) {
	return self.Storage.Load(key, self)
}

