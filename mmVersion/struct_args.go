package mmVersion

import (
	"GoWebcam/Only"
	"strings"
)


type ExecArgs []string

func (r *ExecArgs) ToString() string {
	return strings.Join(*r, " ")
}

func (r *ExecArgs) String() string {
	return strings.Join(*r, " ")
}

func (r *ExecArgs) Set(a ...string) *ExecArgs {
	for range Only.Once {
		*r = a
	}
	return r
}

func (r *ExecArgs) Add(a ...string) *ExecArgs {
	*r = append(*r, a...)
	return r
}

func (r *ExecArgs) Append(a ...string) *ExecArgs {
	*r = append(*r, a...)
	return r
}

func (r *ExecArgs) Get(index int) string {
	var ret string

	for range Only.Once {
		if index > len(*r)-1 {
			break
		}

		ret = (*r)[index]
	}

	return ret
}

func (r *ExecArgs) GetAll() []string {
	return *r
}

func (r *ExecArgs) Sprintf() string {
	return strings.Join(*r, " ")
}

func (r *ExecArgs) Range(lower int, upper int) []string {
	var ret []string

	for range Only.Once {
		if lower < 0 {
			break
		}

		if upper <= 0 {
			break
		}

		upper++

		if lower > upper {
			break
		}

		as := len(*r) - 1

		if lower > as {
			break
		}

		if upper > as {
			// @TODO - Should we pad out this array to the full 'size' if it's less?
			ret = (*r)[lower:]
			break
		}

		if lower == upper {
			ret = []string{(*r)[lower]}
			break
		}

		ret = (*r)[lower:upper]
	}

	return ret
}

func (r *ExecArgs) SprintfRange(lower int, upper int) string {
	return strings.Join(r.Range(lower, upper), " ")
}

func (r *ExecArgs) GetFromSize(begin int, size int) []string {
	return r.Range(begin, begin+size-1)
}

func (r *ExecArgs) SprintfFromSize(lower int, upper int) string {
	return strings.Join(r.GetFromSize(lower, upper), " ")
}

func (r *ExecArgs) GetFrom(lower int) []string {
	return r.Range(lower, len(*r))
}

func (r *ExecArgs) SprintfFrom(lower int) string {
	return strings.Join(r.GetFrom(lower), " ")
}
