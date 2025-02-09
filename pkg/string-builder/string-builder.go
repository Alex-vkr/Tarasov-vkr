package string_builder

import (
	"strings"
)

func BuildConditionalStringForWord(str string) string {
	strs := strings.Split(str, "\\ ")

	b := strings.Builder{}
	var buffSize int
	for _, s := range strs {
		buffSize += len(s) + len("~* ")
	}
	b.Grow(buffSize)

	for _, v := range strs {
		b.WriteString("+")
		b.WriteString(v)
		b.WriteString("~* ")
	}

	return b.String()
}

func BuildConditionalString(str string) string {
	strs := strings.Split(str, "\\ ")

	b := strings.Builder{}
	var buffSize int
	for _, s := range strs {
		buffSize += len(s) + len("~* ")
	}
	b.Grow(buffSize)

	for _, v := range strs {
		//b.WriteString("+")
		b.WriteString(v)
		b.WriteString("~* ")
	}

	return b.String()
}

func BuildStrings(strs []string) string {
	b := strings.Builder{}
	var buffSize int
	for _, s := range strs {
		buffSize += len(s) + len(" ")
	}
	b.Grow(buffSize)

	for _, v := range strs {
		b.WriteString(v)
	}

	return b.String()
}
