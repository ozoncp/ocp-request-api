// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ozoncp/ocp-request-api/internal/producer (interfaces: Producer)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	messaging "github.com/ozoncp/ocp-request-api/internal/producer"
)

// MockProducer is a mock of Producer interface.
type MockProducer struct {
	ctrl     *gomock.Controller
	recorder *MockProducerMockRecorder
}

// MockProducerMockRecorder is the mock recorder for MockProducer.
type MockProducerMockRecorder struct {
	mock *MockProducer
}

// NewMockProducer creates a new mock instance.
func NewMockProducer(ctrl *gomock.Controller) *MockProducer {
	mock := &MockProducer{ctrl: ctrl}
	mock.recorder = &MockProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducer) EXPECT() *MockProducerMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockProducer) Send(arg0 context.Context, arg1 messaging.EventMsg) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Send", arg0, arg1)
}

// Send indicates an expected call of Send.
func (mr *MockProducerMockRecorder) Send(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockProducer)(nil).Send), arg0, arg1)
}