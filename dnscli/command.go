package main

import (
	"regexp"
	"fmt"
	"errors"
)

type Command interface {
	Name() []string

	Run(Options, []string) error
}

func validateIp(s string) (string, error) {
	n := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	r := regexp.MustCompile(fmt.Sprintf("^%s\\.%s\\.%s\\.%s$", n, n, n, n))
	if r.Match([]byte(s)) {
		return s, nil
	}
	return "", errors.New(fmt.Sprintf("invalid ip: %v", s))
}

func validateHost(s string) (string, error) {
	r := regexp.MustCompile("^[a-zA-Z0-9\\-_.]+$")
	if r.Match([]byte(s)) {
		if _, err := validateIp(s); err != nil {
			return s, nil
		}
	}
	return "", errors.New(fmt.Sprintf("invalid host: %v", s))
}

func matchName(c Command, s string) bool {
	for _, n := range c.Name() {
		if n == s {
			return true
		}
	}
	return false
}

func allCommands() []Command {
	return []Command{Add{}, Remove{}, Get{}}
}
