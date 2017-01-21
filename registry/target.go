package registry

import "strings"

type Target struct {
	Address  string
	Failures int
}

type OrderedTargets []Target

func (s OrderedTargets) Len() int {
	return len(s)
}

func (s OrderedTargets) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s OrderedTargets) Less(i, j int) bool {
	less := strings.Compare(s[i].Address, s[j].Address)
	if less < 0 {
		return true
	} else {
		return false
	}
}
