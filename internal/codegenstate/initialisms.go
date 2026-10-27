package codegenstate

import (
	"sync"
	_ "unsafe"
)

var (
	initialismsOnce sync.Once
	initialisms     []string
)

func Initialisms() []string {
	initialismsOnce.Do(func() {
		initialisms = make([]string, 0, len(globalState.initialismsMap))

		for initialism := range globalState.initialismsMap {
			initialisms = append(initialisms, initialism)
		}
	})

	return initialisms
}
