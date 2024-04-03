package modmake_docker

import (
	"fmt"
	"strings"
)

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

type strmap = map[string]*string

func anyBlankPanic(data strmap) {
	for k, v := range data {
		if v == nil {
			panicf("%s: nil value", k)
		}
		*v = strings.TrimSpace(*v)
		if len(*v) == 0 {
			panicf("%s: blank string", k)
		}
	}
}
