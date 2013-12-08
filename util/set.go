package util

type Set map[string]byte

func NewSet() Set {
	return Set(make(map[string]byte))
}

func (s Set) Add(element string) bool {
	if _, ok := s[element]; ok {
		return false
	}
	s[element] = 0
	return true
}

func (s Set) AddSlice(elements []string) {
	for _, v := range elements {
		s.Add(v)
	}
}

func (s Set) Contains(element string) bool {
	_, ok := s[element]
	return ok
}

// iterate over keys
func (s Set) KeySet() map[string]byte {
	return s
}

func (s Set) Remove(element string) {
	delete(s, element)
}

func (s Set) Size() int {
	return len(s)
}
