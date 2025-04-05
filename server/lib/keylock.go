package lib

import (
	"sync"

	"github.com/artifact-space/ArtiSpace/log"
)

type KeyLock struct {
	mu sync.Mutex
	locks map[string]*sync.Mutex
}

func NewKeyLock() *KeyLock {
	return &KeyLock{
		locks: make(map[string]*sync.Mutex),
	}
}

func (kl *KeyLock) Lock(key string) {
	kl.mu.Lock()
	if _, exists := kl.locks[key]; !exists {
		kl.locks[key] = &sync.Mutex{}
	}

	lock := kl.locks[key]

	kl.mu.Unlock()

	lock.Lock()
}

func (kl *KeyLock) Unlock(key string) {
	kl.mu.Lock()

	lock, exists := kl.locks[key]
	if !exists {
		log.Logger().Warn().Msgf("unlock was called for non-existent key: %s", key)
		return
	}
	kl.mu.Unlock()

	lock.Unlock()
}