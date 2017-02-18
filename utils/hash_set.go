package utils

import (
	"bytes"
	"fmt"
)

type Set interface {
	Add(e interface{}) bool
	Remove(e interface{})
	Clear()
	Contains(e interface{}) bool
	Len() int
	Same(other Set) bool
	Elements() []interface{}
	String() string
}

// 判断集合 one 是否是集合 other 的超集
func IsSuperset(one Set, other Set) bool {
	if one == nil || other == nil {
		return false
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0 || oneLen == otherLen {
		return false
	}
	if oneLen > 0 && otherLen == 0 {
		return true
	}
	for _, v := range other.Elements() {
		if !one.Contains(v) {
			return false
		}
	}
	return true
}

// 生成集合 one 和集合 other 的并集
func Union(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	unionedSet := NewSimpleSet()
	for _, v := range one.Elements() {
		unionedSet.Add(v)
	}
	if other.Len() == 0 {
		return unionedSet
	}
	for _, v := range other.Elements() {
		unionedSet.Add(v)
	}
	return unionedSet
}

// 生成集合 one 和集合 other 的交集
func Intersect(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	intersectedSet := NewSimpleSet()
	if other.Len() == 0 {
		return intersectedSet
	}
	if one.Len() < other.Len() {
		for _, v := range one.Elements() {
			if other.Contains(v) {
				intersectedSet.Add(v)
			}
		}
	} else {
		for _, v := range other.Elements() {
			if one.Contains(v) {
				intersectedSet.Add(v)
			}
		}
	}
	return intersectedSet
}

// 生成集合 one 对集合 other 的差集
func Difference(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	differencedSet := NewSimpleSet()
	if other.Len() == 0 {
		for _, v := range one.Elements() {
			differencedSet.Add(v)
		}
		return differencedSet
	}
	for _, v := range one.Elements() {
		if !other.Contains(v) {
			differencedSet.Add(v)
		}
	}
	return differencedSet
}

// 生成集合 one 和集合 other 的对称差集
func SymmetricDifference(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	diffA := Difference(one, other)
	if other.Len() == 0 {
		return diffA
	}
	diffB := Difference(other, one)
	return Union(diffA, diffB)
}

func NewSimpleSet() Set {
	return NewHashSet()
}

func IsSet(value interface{}) bool {
	if _, ok := value.(Set); ok {
		return true
	}
	return false
}

type HashSet struct {
	m map[interface{}]bool
}

func NewHashSet() *HashSet {
	return &HashSet{m: make(map[interface{}]bool)}
}

func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

func (set *HashSet) Remove(e interface{}) {
	delete(set.m, e)
}

func (set *HashSet) Clear() {
	set.m = make(map[interface{}]bool)
}

func (set *HashSet) Contains(e interface{}) bool {
	return set.m[e]
}

func (set *HashSet) Len() int {
	return len(set.m)
}

func (set *HashSet) Same(other Set) bool {
	if other == nil {
		return false
	}
	if set.Len() != other.Len() {
		return false
	}
	for key := range set.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (set *HashSet) Elements() []interface{} {
	initialLen := len(set.m)
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	for key := range set.m {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else {
			snapshot = append(snapshot, key)
		}
		actualLen++
	}
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}

func (set *HashSet) String() string {
	var buf bytes.Buffer
	buf.WriteString("HashSet{")
	first := true
	for key := range set.m {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	buf.WriteString("}")
	return buf.String()
}

