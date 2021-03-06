// Code generated by MockGen. DO NOT EDIT.
// Source: main.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	discordgo "github.com/bwmarrin/discordgo"
	gomock "github.com/golang/mock/gomock"
	xivdata "github.com/kyoukaya/catte/internal/xivdata"
)

// MockDiscordSession is a mock of DiscordSession interface.
type MockDiscordSession struct {
	ctrl     *gomock.Controller
	recorder *MockDiscordSessionMockRecorder
}

// MockDiscordSessionMockRecorder is the mock recorder for MockDiscordSession.
type MockDiscordSessionMockRecorder struct {
	mock *MockDiscordSession
}

// NewMockDiscordSession creates a new mock instance.
func NewMockDiscordSession(ctrl *gomock.Controller) *MockDiscordSession {
	mock := &MockDiscordSession{ctrl: ctrl}
	mock.recorder = &MockDiscordSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDiscordSession) EXPECT() *MockDiscordSessionMockRecorder {
	return m.recorder
}

// AddHandler mocks base method.
func (m *MockDiscordSession) AddHandler(handler interface{}) func() {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddHandler", handler)
	ret0, _ := ret[0].(func())
	return ret0
}

// AddHandler indicates an expected call of AddHandler.
func (mr *MockDiscordSessionMockRecorder) AddHandler(handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHandler", reflect.TypeOf((*MockDiscordSession)(nil).AddHandler), handler)
}

// ApplicationCommandCreate mocks base method.
func (m *MockDiscordSession) ApplicationCommandCreate(appID, guildID string, cmd *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplicationCommandCreate", appID, guildID, cmd)
	ret0, _ := ret[0].(*discordgo.ApplicationCommand)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplicationCommandCreate indicates an expected call of ApplicationCommandCreate.
func (mr *MockDiscordSessionMockRecorder) ApplicationCommandCreate(appID, guildID, cmd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplicationCommandCreate", reflect.TypeOf((*MockDiscordSession)(nil).ApplicationCommandCreate), appID, guildID, cmd)
}

// FollowupMessageCreate mocks base method.
func (m *MockDiscordSession) FollowupMessageCreate(appID string, interaction *discordgo.Interaction, wait bool, data *discordgo.WebhookParams) (*discordgo.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FollowupMessageCreate", appID, interaction, wait, data)
	ret0, _ := ret[0].(*discordgo.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FollowupMessageCreate indicates an expected call of FollowupMessageCreate.
func (mr *MockDiscordSessionMockRecorder) FollowupMessageCreate(appID, interaction, wait, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FollowupMessageCreate", reflect.TypeOf((*MockDiscordSession)(nil).FollowupMessageCreate), appID, interaction, wait, data)
}

// InteractionRespond mocks base method.
func (m *MockDiscordSession) InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InteractionRespond", interaction, resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// InteractionRespond indicates an expected call of InteractionRespond.
func (mr *MockDiscordSessionMockRecorder) InteractionRespond(interaction, resp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InteractionRespond", reflect.TypeOf((*MockDiscordSession)(nil).InteractionRespond), interaction, resp)
}

// MockbuffdebuffInterface is a mock of buffdebuffInterface interface.
type MockbuffdebuffInterface struct {
	ctrl     *gomock.Controller
	recorder *MockbuffdebuffInterfaceMockRecorder
}

// MockbuffdebuffInterfaceMockRecorder is the mock recorder for MockbuffdebuffInterface.
type MockbuffdebuffInterfaceMockRecorder struct {
	mock *MockbuffdebuffInterface
}

// NewMockbuffdebuffInterface creates a new mock instance.
func NewMockbuffdebuffInterface(ctrl *gomock.Controller) *MockbuffdebuffInterface {
	mock := &MockbuffdebuffInterface{ctrl: ctrl}
	mock.recorder = &MockbuffdebuffInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockbuffdebuffInterface) EXPECT() *MockbuffdebuffInterfaceMockRecorder {
	return m.recorder
}

// EndTs mocks base method.
func (m *MockbuffdebuffInterface) EndTs() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EndTs")
	ret0, _ := ret[0].(int64)
	return ret0
}

// EndTs indicates an expected call of EndTs.
func (mr *MockbuffdebuffInterfaceMockRecorder) EndTs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EndTs", reflect.TypeOf((*MockbuffdebuffInterface)(nil).EndTs))
}

// GetID mocks base method.
func (m *MockbuffdebuffInterface) GetID() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetID")
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetID indicates an expected call of GetID.
func (mr *MockbuffdebuffInterfaceMockRecorder) GetID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetID", reflect.TypeOf((*MockbuffdebuffInterface)(nil).GetID))
}

// StartTs mocks base method.
func (m *MockbuffdebuffInterface) StartTs() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartTs")
	ret0, _ := ret[0].(int64)
	return ret0
}

// StartTs indicates an expected call of StartTs.
func (mr *MockbuffdebuffInterfaceMockRecorder) StartTs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTs", reflect.TypeOf((*MockbuffdebuffInterface)(nil).StartTs))
}

// String mocks base method.
func (m *MockbuffdebuffInterface) String(ds *xivdata.DataSource) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String", ds)
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockbuffdebuffInterfaceMockRecorder) String(ds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockbuffdebuffInterface)(nil).String), ds)
}
