package main

import (
	"fmt"
	"strings"
)

type StringArray []string

func (a *StringArray) Set(s string) error {
	for _, ss := range strings.Split(s, ",") {
		*a = append(*a, ss)
	}
	return nil
}

func (a *StringArray) String() string {
	return fmt.Sprint(*a)
}
