// Code generated by mockery v1.0.0. DO NOT EDIT.

package handlers

import entities "github.com/thewizardplusplus/go-link-shortener/entities"
import mock "github.com/stretchr/testify/mock"

// MockLinkGetter is an autogenerated mock type for the LinkGetter type
type MockLinkGetter struct {
	mock.Mock
}

// GetLink provides a mock function with given fields: code
func (_m *MockLinkGetter) GetLink(code string) (entities.Link, error) {
	ret := _m.Called(code)

	var r0 entities.Link
	if rf, ok := ret.Get(0).(func(string) entities.Link); ok {
		r0 = rf(code)
	} else {
		r0 = ret.Get(0).(entities.Link)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
