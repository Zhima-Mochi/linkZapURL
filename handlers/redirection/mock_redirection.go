// Code generated by MockGen. DO NOT EDIT.
// Source: redirection.go

// Package redirection is a generated GoMock package.
package redirection

import (
	context "context"
	reflect "reflect"

	models "github.com/Zhima-Mochi/linkZapURL/models"
	gomock "github.com/golang/mock/gomock"
)

// MockRedirection is a mock of Redirection interface.
type MockRedirection struct {
	ctrl     *gomock.Controller
	recorder *MockRedirectionMockRecorder
}

// MockRedirectionMockRecorder is the mock recorder for MockRedirection.
type MockRedirectionMockRecorder struct {
	mock *MockRedirection
}

// NewMockRedirection creates a new mock instance.
func NewMockRedirection(ctrl *gomock.Controller) *MockRedirection {
	mock := &MockRedirection{ctrl: ctrl}
	mock.recorder = &MockRedirectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedirection) EXPECT() *MockRedirectionMockRecorder {
	return m.recorder
}

// Redirect mocks base method.
func (m *MockRedirection) Redirect(ctx context.Context, code string) (*models.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Redirect", ctx, code)
	ret0, _ := ret[0].(*models.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Redirect indicates an expected call of Redirect.
func (mr *MockRedirectionMockRecorder) Redirect(ctx, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Redirect", reflect.TypeOf((*MockRedirection)(nil).Redirect), ctx, code)
}
