package scoll

import (
	"reflect"

	"github.com/ahmetb/go-linq/v3"
)

type collector struct {
	supplier    func() reflect.Value
	accumulator func(a reflect.Value, v interface{}) reflect.Value
}

func ToList(typ interface{}) collector {
	t := reflect.TypeOf(typ)
	if t.Kind() != reflect.Slice {
		panic("typ should be slice")
	}

	supplier := func() reflect.Value { return reflect.Indirect(reflect.New(t)) }
	accumulator := func(a reflect.Value, v interface{}) reflect.Value {
		return reflect.Append(a, reflect.ValueOf(v))
	}
	return collector{
		supplier:    supplier,
		accumulator: accumulator,
	}
}

func GroupingBy(typ interface{}, classifier func(interface{}) interface{}, downstream collector) collector {
	t := reflect.TypeOf(typ)
	if t.Kind() != reflect.Map {
		panic("typ should be map")
	}
	supplier := func() reflect.Value { return reflect.Indirect(reflect.MakeMap(t)) }
	accumulator := func(a reflect.Value, v interface{}) reflect.Value {
		key := classifier(v)
		keyV := reflect.ValueOf(key)
		container := a.MapIndex(keyV)
		if !container.IsValid() {
			container = downstream.supplier()
		}
		a.SetMapIndex(keyV, downstream.accumulator(container, v))
		return a
	}
	return collector{
		supplier:    supplier,
		accumulator: accumulator,
	}
}

func (s Stream) collect(c collector) reflect.Value {
	return linq.Query(s).AggregateWithSeedT(c.supplier(), c.accumulator).(reflect.Value)
}

func (s Stream) Collect(c collector) interface{} {
	return s.collect(c).Interface()
}
