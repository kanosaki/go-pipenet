package core

type Storage interface {
	Save(key string, graph *MetaGraph, univ *Universe) error
	Load(key string, univ *Universe) (*MetaGraph, error)
}

