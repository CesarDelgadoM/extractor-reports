package requests

import (
	"crypto/md5"
	"encoding/hex"
	"sync"

	"github.com/CesarDelgadoM/extractor-reports/pkg/logger/zap"
)

// Data structure to handle no repeated requests
type ISet interface {
	Set(key string)
	Exist(key string) bool
	Delete(key string)
	Print()
}

// Stores key requests by a md5 of the fields
type storeRequests struct {
	store map[string]bool
	m     sync.Mutex
}

func NewStoreRequests() ISet {
	return &storeRequests{
		store: make(map[string]bool),
	}
}

func (s *storeRequests) Set(key string) {
	s.m.Lock()
	s.store[s.md5(key)] = true
	s.m.Unlock()
}

func (s *storeRequests) Exist(key string) bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.store[s.md5(key)]
}

func (s *storeRequests) Delete(key string) {
	s.m.Lock()
	delete(s.store, s.md5(key))
	s.m.Unlock()
}

func (s *storeRequests) Print() {
	for k := range s.store {
		zap.Log.Info("key: ", k)
	}
}

func (s *storeRequests) md5(key string) string {
	hash := md5.New()
	hash.Write([]byte(key))

	return hex.EncodeToString(hash.Sum(nil))
}
