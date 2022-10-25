package utils

import (
	"regexp"
	"strings"
)

func PString(s string) *string {
	return &s
}

func PFloat(i int) *float32 {
	f := float32(i)
	return &f
}

func ToCamelCase(s string) string {
    re, _ := regexp.Compile(`[-_]\w`)
    res := re.ReplaceAllStringFunc(s, func(m string) string {
        return strings.ToUpper(m[1:])
    })
    return res
}