package scoll

import (
	"reflect"
	"testing"
)

func TestCollectToList(t *testing.T) {
	v := []int{1, 2, 3, 4, 5}
	got := From(v).Collect(ToList([]int{}))

	if !reflect.DeepEqual(v, got) {
		t.Errorf("%v != %v", got, v)
	}
}

type Cargo struct {
	ID     int
	Name   string
	City   string
	Status int
}

func TestCollectGroupBy(t *testing.T) {
	res := map[string][]*Cargo{
		"shenzhen": {{
			ID:     1,
			Name:   "foo",
			City:   "shenzhen",
			Status: 1,
		}, {
			ID:     2,
			Name:   "bar",
			City:   "shenzhen",
			Status: 0,
		}},
		"guangzhou": {{
			ID:     3,
			Name:   "zhang",
			City:   "guangzhou",
			Status: 1,
		}},
	}

	v := []*Cargo{{
		ID:     1,
		Name:   "foo",
		City:   "shenzhen",
		Status: 1,
	}, {
		ID:     2,
		Name:   "bar",
		City:   "shenzhen",
		Status: 0,
	}, {
		ID:     3,
		Name:   "zhang",
		City:   "guangzhou",
		Status: 1,
	}}

	getCity := func(v interface{}) interface{} {
		return v.(*Cargo).City
	}
	got := From(v).Collect(GroupingBy(map[string][]*Cargo{}, getCity, ToList([]*Cargo{}))).(map[string][]*Cargo)

	if !reflect.DeepEqual(res, got) {
		t.Errorf("%v != %v", got, res)
	}
}

func TestMultiGroupBy(t *testing.T) {
	v := []*Cargo{{
		ID:     1,
		Name:   "foo",
		City:   "shenzhen",
		Status: 1,
	}, {
		ID:     2,
		Name:   "bar",
		City:   "shenzhen",
		Status: 0,
	}, {
		ID:     3,
		Name:   "zhang",
		City:   "guangzhou",
		Status: 1,
	}}

	// group by status, city
	res := map[int]map[string][]*Cargo{
		1: {
			"shenzhen": {
				{
					ID:     1,
					Name:   "foo",
					City:   "shenzhen",
					Status: 1,
				},
			},
			"guangzhou": {
				{
					ID:     3,
					Name:   "zhang",
					City:   "guangzhou",
					Status: 1,
				},
			},
		},
		0: {
			"shenzhen": {
				{
					ID:     2,
					Name:   "bar",
					City:   "shenzhen",
					Status: 0,
				},
			},
		},
	}

	getStatus := func(v interface{}) interface{} { return v.(*Cargo).Status }
	getCity := func(v interface{}) interface{} { return v.(*Cargo).City }

	got := From(v).Collect(
		GroupingBy(map[int]map[string][]*Cargo{}, getStatus,
			GroupingBy(map[string][]*Cargo{}, getCity,
				ToList([]*Cargo{}))))

	if !reflect.DeepEqual(res, got) {
		t.Errorf("%v != %v", got, res)
	}
}
