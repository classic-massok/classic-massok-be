package lib

func NewStringset(strings ...string) StringSet {
	ss := make(StringSet)
	for _, str := range strings {
		ss[str] = struct{}{}
	}

	return ss
}

type StringSet map[string]struct{}

func (ss StringSet) Contains(text string) bool {
	_, ok := ss[text]
	return ok
}
