package storage

import (
	"github.com/kanosaki/go-pipenet/core"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

type NullStorage struct {

}

func NewNullStorage() *NullStorage {
	return &NullStorage{}
}

func (self *NullStorage) Save(key string, graph *core.MetaGraph, univ *core.Universe) error {
	log.Warn("Writing to NullStorage")
	return nil
}

func (self *NullStorage) Load(key string, univ *core.Universe) (*core.MetaGraph, error) {
	return nil, fmt.Errorf("Read from NullStorage")
}
