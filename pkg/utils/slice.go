package utils

import (
	"strings"
)

func SliceStrContains(s []string, find string) bool {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if strings.TrimSpace(v) == strings.TrimSpace(find) {
			return true
		}
	}
	return false
}

func SliceStrUnion(s1, s2 []string) []string {
	m := make(map[string]int)
	for _, v := range s1 {
		m[v]++
	}

	for _, v := range s2 {
		times, _ := m[v]
		if times == 0 {
			s1 = append(s1, v)
		}
	}
	return s1
}

func SliceStrInter(s1, s2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range s1 {
		m[v]++
	}

	for _, v := range s2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

func SliceStrDiff(s1, s2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := SliceStrInter(s1, s2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range s1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}
