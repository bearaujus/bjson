package bjson

import (
	"fmt"
)

type tracer struct {
	origin    []string
	remaining []string
	passed    []string
}

func newTracer(targets []string) *tracer {
	return &tracer{
		origin:    targets,
		remaining: targets,
		passed:    nil,
	}
}

func (t *tracer) next() bool {
	if t.isTail() {
		return false
	}

	t.passed = append(t.passed, t.remaining[0])
	t.remaining = t.remaining[1:]
	return true
}

func (t *tracer) currTarget() string {
	if len(t.passed) == 0 {
		return ""
	}

	return t.passed[len(t.passed)-1]
}

func (t *tracer) isTail() bool {
	return len(t.remaining) == 0
}

func (t *tracer) passedPath() string {
	return parseTracerPath(t.passed)
}

func (t *tracer) originPath() string {
	return parseTracerPath(t.origin)
}

func parseTracerPath(v []string) string {
	ret := `'JSON`
	for _, v := range v {
		ret += fmt.Sprintf(`[%v]`, v)
	}
	ret += `'`

	return ret
}
