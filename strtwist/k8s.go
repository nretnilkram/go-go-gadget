// Package strtwist implements additional functions UTF-8 encoded strings.
package strtwist

import (
	"fmt"
	"strings"
)

func k8s_transform(s string) string {
	if len(s) == 0 {
		return s
	}
	var first string = string(s[0])
	var length int = len(s)
	var last string = string(s[length-1])
	if len(s) > 2 {
		return first + fmt.Sprint(length-2) + last
	}
	return s
}

// Kubernetes returns input string in the Kubernetes k8s shortened version.
func K8s(s string) string {
	list := strings.Split(s, " ")
	for i := 0; i < len(list); i++ {
		list[i] = k8s_transform(list[i])
	}
	return strings.Join(list, " ")
}
