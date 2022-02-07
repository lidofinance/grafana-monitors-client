// Code generated by MockGen. DO NOT EDIT.
// Source: ./grafana.go

// Package grafana is a generated GoMock package.
package grafana

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockGrafana is a mock of Grafana interface.
type MockGrafana struct {
	ctrl     *gomock.Controller
	recorder *MockGrafanaMockRecorder
}

// MockGrafanaMockRecorder is the mock recorder for MockGrafana.
type MockGrafanaMockRecorder struct {
	mock *MockGrafana
}

// NewMockGrafana creates a new mock instance.
func NewMockGrafana(ctrl *gomock.Controller) *MockGrafana {
	mock := &MockGrafana{ctrl: ctrl}
	mock.recorder = &MockGrafanaMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGrafana) EXPECT() *MockGrafanaMockRecorder {
	return m.recorder
}

// GetGrafanaPanel mocks base method.
func (m *MockGrafana) GetGrafanaPanel(panelName, dashboardID string) (*Panel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGrafanaPanel", panelName, dashboardID)
	ret0, _ := ret[0].(*Panel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGrafanaPanel indicates an expected call of GetGrafanaPanel.
func (mr *MockGrafanaMockRecorder) GetGrafanaPanel(panelName, dashboardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGrafanaPanel", reflect.TypeOf((*MockGrafana)(nil).GetGrafanaPanel), panelName, dashboardID)
}

// GetPanelPicture mocks base method.
func (m *MockGrafana) GetPanelPicture(url string) (PanelPicture, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPanelPicture", url)
	ret0, _ := ret[0].(PanelPicture)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPanelPicture indicates an expected call of GetPanelPicture.
func (mr *MockGrafanaMockRecorder) GetPanelPicture(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPanelPicture", reflect.TypeOf((*MockGrafana)(nil).GetPanelPicture), url)
}

// Panels mocks base method.
func (m *MockGrafana) Panels(ctx context.Context, dashboardUid string) ([]Panel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Panels", ctx, dashboardUid)
	ret0, _ := ret[0].([]Panel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Panels indicates an expected call of Panels.
func (mr *MockGrafanaMockRecorder) Panels(ctx, dashboardUid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Panels", reflect.TypeOf((*MockGrafana)(nil).Panels), ctx, dashboardUid)
}

// PanelsFiltered mocks base method.
func (m *MockGrafana) PanelsFiltered(ctx context.Context, dashboardUid string, filterPanelNames []string) ([]Panel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PanelsFiltered", ctx, dashboardUid, filterPanelNames)
	ret0, _ := ret[0].([]Panel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PanelsFiltered indicates an expected call of PanelsFiltered.
func (mr *MockGrafanaMockRecorder) PanelsFiltered(ctx, dashboardUid, filterPanelNames interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PanelsFiltered", reflect.TypeOf((*MockGrafana)(nil).PanelsFiltered), ctx, dashboardUid, filterPanelNames)
}
