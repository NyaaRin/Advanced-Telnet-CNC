package database

import (
	"fmt"
	"strconv"
	"strings"
)

func (user *UserProfile) MethodGloballyBlacklisted(method int) bool {
	for _, mtd := range GlobalRestricted {
		if method == mtd {
			if user.Admin {
				return false
			}

			return true
		}
	}

	return false
}

func (user *UserProfile) MethodGloballyBlacklistedForAdmins(method int) bool {
	for _, mtd := range GlobalRestricted {
		if method == mtd {
			return true
		}
	}

	return false
}

func (user *UserProfile) methodAllowed(method int) bool {
	if user.MethodGloballyBlacklisted(method) {
		return false
	}

	if method == -1 {
		return true
	}

	for _, mtd := range user.Methods {
		if method == mtd {
			return true
		}
	}
	return false
}

func MethodsToStr(methods []int) (str string) {
	for _, method := range methods {
		str += fmt.Sprintf("%d,", method)
	}
	return strings.TrimSuffix(str, ",")
}

func MethodsToInt(methods string) (str []int) {
	mtd := strings.Split(methods, ",")
	for _, s := range mtd {
		atoi, err := strconv.Atoi(s)
		if err != nil {
			return make([]int, 0)
		}
		str = append(str, atoi)
	}
	return str
}
