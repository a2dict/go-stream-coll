package scoll

import (
	"github.com/ahmetb/go-linq/v3"
)

type Stream linq.Query

func From(source interface{}) Stream {
	return Stream(linq.From(source))
}

func (s Stream) Map(fn func(interface{}) interface{}) Stream {
	return Stream(linq.Query(s).Select(fn))
}

func (s Stream) FlatMap(fn func(interface{}) Stream) Stream {
	selector := func(v interface{}) linq.Query {
		return linq.Query(fn(v))
	}
	return Stream(linq.Query(s).SelectMany(selector))
}

func (s Stream) Filter(predicate func(interface{}) bool) Stream {
	return Stream(linq.Query(s).Where(predicate))
}

func (s Stream) SortedBy(fn func(interface{}) interface{}) Stream {
	return Stream(linq.Query(s).OrderBy(fn).Query)
}

func (s Stream) ForEach(action func(interface{})) {
	linq.Query(s).ForEach(action)
}
