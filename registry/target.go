package registry

import "strings"

type Target struct {
	Address  string
	Failures int
}

func (t *Target) setAddress(address string) {
	t.Address = address
}

func (t *Target) getAddress() string {
	return t.Address
}

func (t *Target) incrementFailures(value int) int {
	t.Failures += value
	return t.Failures
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
