package main

import (
	"fmt"
)

type KeyValueStore struct {
	store map[string]struct{}
}

func MakeSet() *KeyValueStore {
	return &KeyValueStore{
		store: make(map[string]struct{}),
	}
}

func (s *KeyValueStore) Exists(key string) bool {
	_, exists := s.store[key]
	return exists
}

func (s *KeyValueStore) Add(key string) {
	s.store[key] = struct{}{}
}

func (s *KeyValueStore) Remove(key string) error {
	_, exists := s.store[key]
	if !exists {
		return fmt.Errorf("Remove Error: Item doesn't exist in set")
	}
	delete(s.store, key)
	return nil
}

func (s *KeyValueStore) Size() int {
	return len(s.store)
}
