/*
 *
 * strings.go
 * utils
 *
 * Created by lintao on 2020/8/8 4:20 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package utils

import (
	"strings"
)

func IndexNoCase(s, sep string) int {
	return strings.Index(strings.ToLower(s), strings.ToLower(sep))
}

func SplitNoCase(s, sep string) []string {
	idx := IndexNoCase(s, sep)
	if idx < 0 {
		return []string{s}
	}
	return strings.Split(s, s[idx:idx+len(sep)])
}

func SplitNNoCase(s, sep string, n int) []string {
	idx := IndexNoCase(s, sep)
	if idx < 0 {
		return []string{s}
	}
	return strings.SplitN(s, s[idx:idx+len(sep)], n)
}
