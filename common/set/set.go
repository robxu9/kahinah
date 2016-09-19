package set

import (
	"sync"
)

type Set struct {
	m    map[string]byte
	lock *sync.RWMutex
}

func NewSet() *Set {
	return &Set{
		m:    make(map[string]byte),
		lock: &sync.RWMutex{},
	}
}

func (s *Set) Add(element string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.m[element]; ok {
		return false
	}
	s.m[element] = 0
	return true
}

func (s *Set) AddSlice(elements []string) {
	for _, v := range elements {
		s.Add(v)
	}
}

func (s *Set) Contains(element string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.m[element]
	return ok
}

// iterate over keys
func (s *Set) Iterator() <-chan string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	ch := make(chan string, s.Size())
	defer close(ch)
	for k, _ := range s.m {
		ch <- k
	}
	return ch
}

func (s *Set) Remove(element string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.m, element)
}

func (s *Set) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.m)
}
