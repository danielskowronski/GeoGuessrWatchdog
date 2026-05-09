package main

type StringSet struct {
	values []string
	seen   map[string]struct{}
}

func NewStringSet() *StringSet {
	return &StringSet{
		seen: make(map[string]struct{}),
	}
}

func (s *StringSet) Add(value string) {
	if _, exists := s.seen[value]; exists {
		return
	}
	s.seen[value] = struct{}{}
	s.values = append(s.values, value)
}
