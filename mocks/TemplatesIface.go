// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	templates "github.com/x86-Yantras/code-gen/internal/adapters/templates"
)

// TemplatesIface is an autogenerated mock type for the TemplatesIface type
type TemplatesIface struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *TemplatesIface) Create(_a0 *templates.FileCreateParams) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*templates.FileCreateParams) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateDir provides a mock function with given fields: name
func (_m *TemplatesIface) CreateDir(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateMany provides a mock function with given fields: service, files
func (_m *TemplatesIface) CreateMany(service interface{}, files ...*templates.FileCreateParams) error {
	_va := make([]interface{}, len(files))
	for _i := range files {
		_va[_i] = files[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, service)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, ...*templates.FileCreateParams) error); ok {
		r0 = rf(service, files...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTemplatesIface interface {
	mock.TestingT
	Cleanup(func())
}

// NewTemplatesIface creates a new instance of TemplatesIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTemplatesIface(t mockConstructorTestingTNewTemplatesIface) *TemplatesIface {
	mock := &TemplatesIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}