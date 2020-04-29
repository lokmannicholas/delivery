// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DistanceFounder is an autogenerated mock type for the DistanceFounder type
type DistanceFounder struct {
	mock.Mock
}

// CountDistance provides a mock function with given fields: start, end
func (_m *DistanceFounder) CountDistance(start []string, end []string) int {
	ret := _m.Called(start, end)

	var r0 int
	if rf, ok := ret.Get(0).(func([]string, []string) int); ok {
		r0 = rf(start, end)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}
