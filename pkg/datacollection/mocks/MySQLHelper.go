// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	sql "database/sql"
)

// MySQLHelper is an autogenerated mock type for the MySQLHelper type
type MySQLHelper struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MySQLHelper) Close() {
	_m.Called()
}

// DB provides a mock function with given fields:
func (_m *MySQLHelper) DB() *sql.DB {
	ret := _m.Called()

	var r0 *sql.DB
	if rf, ok := ret.Get(0).(func() *sql.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.DB)
		}
	}

	return r0
}

// Tx provides a mock function with given fields: ctx, f
func (_m *MySQLHelper) Tx(ctx context.Context, f func(*sql.Tx) error) error {
	ret := _m.Called(ctx, f)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(*sql.Tx) error) error); ok {
		r0 = rf(ctx, f)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}