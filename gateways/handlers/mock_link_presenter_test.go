// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package handlers

import (
	http "net/http"

	entities "github.com/thewizardplusplus/go-link-shortener-backend/entities"

	mock "github.com/stretchr/testify/mock"
)

// MockLinkPresenter is an autogenerated mock type for the LinkPresenter type
type MockLinkPresenter struct {
	mock.Mock
}

// PresentLink provides a mock function with given fields: writer, request, link
func (_m *MockLinkPresenter) PresentLink(writer http.ResponseWriter, request *http.Request, link entities.Link) {
	_m.Called(writer, request, link)
}