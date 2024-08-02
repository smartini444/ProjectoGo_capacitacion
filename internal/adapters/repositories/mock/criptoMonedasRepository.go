// Code generated by MockGen. DO NOT EDIT.
// Source: ./criptoMonedasRepository.go
//
// Generated by this command:
//
//	mockgen -source=./criptoMonedasRepository.go -destination=./mock/criptoMonedasRepository.go -package mock
//

// Package mock is a generated GoMock package.
package mock

import (
	criptomonedas "primerProjecto/internal/entities/criptomonedas"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCryptoRepository is a mock of CryptoRepository interface.
type MockCryptoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCryptoRepositoryMockRecorder
}

// MockCryptoRepositoryMockRecorder is the mock recorder for MockCryptoRepository.
type MockCryptoRepositoryMockRecorder struct {
	mock *MockCryptoRepository
}

// NewMockCryptoRepository creates a new mock instance.
func NewMockCryptoRepository(ctrl *gomock.Controller) *MockCryptoRepository {
	mock := &MockCryptoRepository{ctrl: ctrl}
	mock.recorder = &MockCryptoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCryptoRepository) EXPECT() *MockCryptoRepositoryMockRecorder {
	return m.recorder
}

// ActualizarCotizacionManual mocks base method.
func (m *MockCryptoRepository) ActualizarCotizacionManual(usuarioId int, cotizacion criptomonedas.Cotizacion) (criptomonedas.Cotizacion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActualizarCotizacionManual", usuarioId, cotizacion)
	ret0, _ := ret[0].(criptomonedas.Cotizacion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ActualizarCotizacionManual indicates an expected call of ActualizarCotizacionManual.
func (mr *MockCryptoRepositoryMockRecorder) ActualizarCotizacionManual(usuarioId, cotizacion any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActualizarCotizacionManual", reflect.TypeOf((*MockCryptoRepository)(nil).ActualizarCotizacionManual), usuarioId, cotizacion)
}

// BorrarCotizacionById mocks base method.
func (m *MockCryptoRepository) BorrarCotizacionById(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BorrarCotizacionById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// BorrarCotizacionById indicates an expected call of BorrarCotizacionById.
func (mr *MockCryptoRepositoryMockRecorder) BorrarCotizacionById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BorrarCotizacionById", reflect.TypeOf((*MockCryptoRepository)(nil).BorrarCotizacionById), id)
}

// BorrarCotizacionManual mocks base method.
func (m *MockCryptoRepository) BorrarCotizacionManual(cotizacion criptomonedas.Cotizacion) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BorrarCotizacionManual", cotizacion)
	ret0, _ := ret[0].(error)
	return ret0
}

// BorrarCotizacionManual indicates an expected call of BorrarCotizacionManual.
func (mr *MockCryptoRepositoryMockRecorder) BorrarCotizacionManual(cotizacion any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BorrarCotizacionManual", reflect.TypeOf((*MockCryptoRepository)(nil).BorrarCotizacionManual), cotizacion)
}

// FindAllByFilter mocks base method.
func (m *MockCryptoRepository) FindAllByFilter(filter criptomonedas.CriptoMonedaFilter) ([]criptomonedas.Cotizacion, criptomonedas.Summary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByFilter", filter)
	ret0, _ := ret[0].([]criptomonedas.Cotizacion)
	ret1, _ := ret[1].(criptomonedas.Summary)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindAllByFilter indicates an expected call of FindAllByFilter.
func (mr *MockCryptoRepositoryMockRecorder) FindAllByFilter(filter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByFilter", reflect.TypeOf((*MockCryptoRepository)(nil).FindAllByFilter), filter)
}

// FindAllByFilterForUser mocks base method.
func (m *MockCryptoRepository) FindAllByFilterForUser(filter criptomonedas.CriptoMonedaFilter, usuarioId int) ([]criptomonedas.Cotizacion, criptomonedas.Summary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByFilterForUser", filter, usuarioId)
	ret0, _ := ret[0].([]criptomonedas.Cotizacion)
	ret1, _ := ret[1].(criptomonedas.Summary)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// FindAllByFilterForUser indicates an expected call of FindAllByFilterForUser.
func (mr *MockCryptoRepositoryMockRecorder) FindAllByFilterForUser(filter, usuarioId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByFilterForUser", reflect.TypeOf((*MockCryptoRepository)(nil).FindAllByFilterForUser), filter, usuarioId)
}

// FindAllCotizaciones mocks base method.
func (m *MockCryptoRepository) FindAllCotizaciones() ([]*criptomonedas.Cotizacion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllCotizaciones")
	ret0, _ := ret[0].([]*criptomonedas.Cotizacion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllCotizaciones indicates an expected call of FindAllCotizaciones.
func (mr *MockCryptoRepositoryMockRecorder) FindAllCotizaciones() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllCotizaciones", reflect.TypeOf((*MockCryptoRepository)(nil).FindAllCotizaciones))
}

// FindAllMonedas mocks base method.
func (m *MockCryptoRepository) FindAllMonedas() ([]*criptomonedas.CriptoMoneda, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllMonedas")
	ret0, _ := ret[0].([]*criptomonedas.CriptoMoneda)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllMonedas indicates an expected call of FindAllMonedas.
func (mr *MockCryptoRepositoryMockRecorder) FindAllMonedas() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllMonedas", reflect.TypeOf((*MockCryptoRepository)(nil).FindAllMonedas))
}

// FindByCotizacionID mocks base method.
func (m *MockCryptoRepository) FindByCotizacionID(id int) (*criptomonedas.Cotizacion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCotizacionID", id)
	ret0, _ := ret[0].(*criptomonedas.Cotizacion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCotizacionID indicates an expected call of FindByCotizacionID.
func (mr *MockCryptoRepositoryMockRecorder) FindByCotizacionID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCotizacionID", reflect.TypeOf((*MockCryptoRepository)(nil).FindByCotizacionID), id)
}

// FindByMonedaID mocks base method.
func (m *MockCryptoRepository) FindByMonedaID(id int) (*criptomonedas.CriptoMoneda, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByMonedaID", id)
	ret0, _ := ret[0].(*criptomonedas.CriptoMoneda)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByMonedaID indicates an expected call of FindByMonedaID.
func (mr *MockCryptoRepositoryMockRecorder) FindByMonedaID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByMonedaID", reflect.TypeOf((*MockCryptoRepository)(nil).FindByMonedaID), id)
}

// FindCryptoByCode mocks base method.
func (m *MockCryptoRepository) FindCryptoByCode(codigo string) (*criptomonedas.CriptoMoneda, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCryptoByCode", codigo)
	ret0, _ := ret[0].(*criptomonedas.CriptoMoneda)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCryptoByCode indicates an expected call of FindCryptoByCode.
func (mr *MockCryptoRepositoryMockRecorder) FindCryptoByCode(codigo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCryptoByCode", reflect.TypeOf((*MockCryptoRepository)(nil).FindCryptoByCode), codigo)
}

// FindCryptoByName mocks base method.
func (m *MockCryptoRepository) FindCryptoByName(name string) (*criptomonedas.CriptoMoneda, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCryptoByName", name)
	ret0, _ := ret[0].(*criptomonedas.CriptoMoneda)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCryptoByName indicates an expected call of FindCryptoByName.
func (mr *MockCryptoRepositoryMockRecorder) FindCryptoByName(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCryptoByName", reflect.TypeOf((*MockCryptoRepository)(nil).FindCryptoByName), name)
}

// FindUltimaCotizacion mocks base method.
func (m *MockCryptoRepository) FindUltimaCotizacion(nombre string) (*criptomonedas.Cotizacion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUltimaCotizacion", nombre)
	ret0, _ := ret[0].(*criptomonedas.Cotizacion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUltimaCotizacion indicates an expected call of FindUltimaCotizacion.
func (mr *MockCryptoRepositoryMockRecorder) FindUltimaCotizacion(nombre any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUltimaCotizacion", reflect.TypeOf((*MockCryptoRepository)(nil).FindUltimaCotizacion), nombre)
}

// GuardarCotizacionManual mocks base method.
func (m *MockCryptoRepository) GuardarCotizacionManual(usuarioId int, cotizacion criptomonedas.Cotizacion) (criptomonedas.Cotizacion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GuardarCotizacionManual", usuarioId, cotizacion)
	ret0, _ := ret[0].(criptomonedas.Cotizacion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GuardarCotizacionManual indicates an expected call of GuardarCotizacionManual.
func (mr *MockCryptoRepositoryMockRecorder) GuardarCotizacionManual(usuarioId, cotizacion any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GuardarCotizacionManual", reflect.TypeOf((*MockCryptoRepository)(nil).GuardarCotizacionManual), usuarioId, cotizacion)
}

// SaveCotizacion mocks base method.
func (m *MockCryptoRepository) SaveCotizacion(cripto criptomonedas.Cotizacion) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCotizacion", cripto)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCotizacion indicates an expected call of SaveCotizacion.
func (mr *MockCryptoRepositoryMockRecorder) SaveCotizacion(cripto any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCotizacion", reflect.TypeOf((*MockCryptoRepository)(nil).SaveCotizacion), cripto)
}

// SaveMoneda mocks base method.
func (m *MockCryptoRepository) SaveMoneda(cripto criptomonedas.CriptoMoneda) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMoneda", cripto)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveMoneda indicates an expected call of SaveMoneda.
func (mr *MockCryptoRepositoryMockRecorder) SaveMoneda(cripto any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMoneda", reflect.TypeOf((*MockCryptoRepository)(nil).SaveMoneda), cripto)
}

// UpdateCotizacion mocks base method.
func (m *MockCryptoRepository) UpdateCotizacion(id int, cotizacion criptomonedas.Cotizacion) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCotizacion", id, cotizacion)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCotizacion indicates an expected call of UpdateCotizacion.
func (mr *MockCryptoRepositoryMockRecorder) UpdateCotizacion(id, cotizacion any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCotizacion", reflect.TypeOf((*MockCryptoRepository)(nil).UpdateCotizacion), id, cotizacion)
}

// UpdateMoneda mocks base method.
func (m *MockCryptoRepository) UpdateMoneda(id int, moneda criptomonedas.CriptoMoneda) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMoneda", id, moneda)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMoneda indicates an expected call of UpdateMoneda.
func (mr *MockCryptoRepositoryMockRecorder) UpdateMoneda(id, moneda any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMoneda", reflect.TypeOf((*MockCryptoRepository)(nil).UpdateMoneda), id, moneda)
}
