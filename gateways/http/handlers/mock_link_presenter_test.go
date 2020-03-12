// Code generated by mockery v1.0.0. DO NOT EDIT.

package handlers

import entities "github.com/thewizardplusplus/go-link-shortener-backend/entities"
import http "net/http"
import mock "github.com/stretchr/testify/mock"

// MockLinkPresenter is an autogenerated mock type for the LinkPresenter type
type MockLinkPresenter struct {
	mock.Mock
}

// PresentLink provides a mock function with given fields: writer, request, link
func (_m *MockLinkPresenter) PresentLink(writer http.ResponseWriter, request *http.Request, link entities.Link) {
	_m.Called(writer, request, link)
}