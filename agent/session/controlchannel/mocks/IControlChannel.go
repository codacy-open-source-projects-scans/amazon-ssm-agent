// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "github.com/aws/amazon-ssm-agent/agent/context"
	contracts "github.com/aws/amazon-ssm-agent/agent/session/contracts"

	log "github.com/aws/amazon-ssm-agent/agent/log"

	mock "github.com/stretchr/testify/mock"

	service "github.com/aws/amazon-ssm-agent/agent/session/service"
)

// IControlChannel is an autogenerated mock type for the IControlChannel type
type IControlChannel struct {
	mock.Mock
}

// Close provides a mock function with given fields: _a0
func (_m *IControlChannel) Close(_a0 log.T) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(log.T) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Initialize provides a mock function with given fields: _a0, mgsService, instanceId, agentMessageIncomingMessageChan
func (_m *IControlChannel) Initialize(_a0 context.T, mgsService service.Service, instanceId string, agentMessageIncomingMessageChan chan contracts.AgentMessage) {
	_m.Called(_a0, mgsService, instanceId, agentMessageIncomingMessageChan)
}

// Open provides a mock function with given fields: _a0, ableToOpenMGSConnection
func (_m *IControlChannel) Open(_a0 context.T, ableToOpenMGSConnection *uint32) error {
	ret := _m.Called(_a0, ableToOpenMGSConnection)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.T, *uint32) error); ok {
		r0 = rf(_a0, ableToOpenMGSConnection)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Reconnect provides a mock function with given fields: _a0, ableToOpenMGSConnection
func (_m *IControlChannel) Reconnect(_a0 context.T, ableToOpenMGSConnection *uint32) error {
	ret := _m.Called(_a0, ableToOpenMGSConnection)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.T, *uint32) error); ok {
		r0 = rf(_a0, ableToOpenMGSConnection)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendMessage provides a mock function with given fields: _a0, input, inputType
func (_m *IControlChannel) SendMessage(_a0 log.T, input []byte, inputType int) error {
	ret := _m.Called(_a0, input, inputType)

	var r0 error
	if rf, ok := ret.Get(0).(func(log.T, []byte, int) error); ok {
		r0 = rf(_a0, input, inputType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetWebSocket provides a mock function with given fields: _a0, mgsService, ableToOpenMGSConnection
func (_m *IControlChannel) SetWebSocket(_a0 context.T, mgsService service.Service, ableToOpenMGSConnection *uint32) error {
	ret := _m.Called(_a0, mgsService, ableToOpenMGSConnection)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.T, service.Service, *uint32) error); ok {
		r0 = rf(_a0, mgsService, ableToOpenMGSConnection)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}