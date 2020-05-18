// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Azure/ARO-RP/pkg/util/azureclient/mgmt/msi (interfaces: UserAssignedIdentitiesClient)

// Package mock_msi is a generated GoMock package.
package mock_msi

import (
	context "context"
	reflect "reflect"

	msi "github.com/Azure/azure-sdk-for-go/services/msi/mgmt/2018-11-30/msi"
	gomock "github.com/golang/mock/gomock"
)

// MockUserAssignedIdentitiesClient is a mock of UserAssignedIdentitiesClient interface
type MockUserAssignedIdentitiesClient struct {
	ctrl     *gomock.Controller
	recorder *MockUserAssignedIdentitiesClientMockRecorder
}

// MockUserAssignedIdentitiesClientMockRecorder is the mock recorder for MockUserAssignedIdentitiesClient
type MockUserAssignedIdentitiesClientMockRecorder struct {
	mock *MockUserAssignedIdentitiesClient
}

// NewMockUserAssignedIdentitiesClient creates a new mock instance
func NewMockUserAssignedIdentitiesClient(ctrl *gomock.Controller) *MockUserAssignedIdentitiesClient {
	mock := &MockUserAssignedIdentitiesClient{ctrl: ctrl}
	mock.recorder = &MockUserAssignedIdentitiesClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserAssignedIdentitiesClient) EXPECT() *MockUserAssignedIdentitiesClientMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockUserAssignedIdentitiesClient) Get(arg0 context.Context, arg1, arg2 string) (msi.Identity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(msi.Identity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockUserAssignedIdentitiesClientMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserAssignedIdentitiesClient)(nil).Get), arg0, arg1, arg2)
}
