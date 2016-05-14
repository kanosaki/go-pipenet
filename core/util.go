package core

import "fmt"

func numericIDGenerator() (func() JointKey) {
	counter := 0
	return func() JointKey {
		counter += 1
		return JointKey(fmt.Sprintf("_%d", counter))
	}
}


