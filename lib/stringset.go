package lib

func New(strings ...string) StringSet {
	ss := make(StringSet)
	for _, str := range strings {
		ss[str] = struct{}{}
	}

	return ss
}

type StringSet map[string]struct{}
