package tests

import (
	"errors"
	"fmt"
	"primerProjecto/internal/adapters/cotizadores"
	mockCotizador "primerProjecto/internal/adapters/cotizadores/mock"
	mockRepo "primerProjecto/internal/adapters/repositories/mock"
	"primerProjecto/internal/entities/criptomonedas"
	"primerProjecto/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSaveCotizacion_succes(t *testing.T) {
	cotizacion := criptomonedas.Cotizacion{
		CriptoMoneda_ID: 11,
		Cotizacion:      100,
		Fecha:           time.Now(),
	}
	ctrl := gomock.NewController(t)

	repoCritpo := mockRepo.NewMockCryptoRepository(ctrl)

	repoCritpo.EXPECT().SaveCotizacion(cotizacion).Return(nil)

	cs := services.NewCryptoService(repoCritpo, cotizadores.GetCotizador)

	err := cs.SaveCotizacion(cotizacion)

	assert.Nil(t, err)
}

func TestSaveCotizacion_fail(t *testing.T) {
	cotizacion := criptomonedas.Cotizacion{
		CriptoMoneda_ID: 11,
		Cotizacion:      100,
		Fecha:           time.Now(),
	}

	ctrl := gomock.NewController(t)

	repoCritpo := mockRepo.NewMockCryptoRepository(ctrl)

	repoCritpo.EXPECT().SaveCotizacion(cotizacion).Return(errors.New("error"))

	cs := services.NewCryptoService(repoCritpo, cotizadores.GetCotizador)

	err := cs.SaveCotizacion(cotizacion)

	assert.NotNil(t, err)
}

func TestGuardarCotizacionExterna_Succes(t *testing.T) {

	ctrl := gomock.NewController(t)
	repoCripto := mockRepo.NewMockCryptoRepository(ctrl)
	cotizador := mockCotizador.NewMockCotizador(ctrl)
	repoCripto.EXPECT().FindCryptoByName("Bitcoin").Return(&criptomonedas.CriptoMoneda{Nombre: "A", Codigo: "B"}, nil).Times(1)
	cotizador.EXPECT().GetCotizacionExterna("A", "B", "USD").Return(criptomonedas.Cotizacion{}, nil).Times(1)
	repoCripto.EXPECT().FindCryptoByName("Bitcoin").Return(&criptomonedas.CriptoMoneda{Nombre: "A", Codigo: "B"}, nil).Times(1)
	repoCripto.EXPECT().SaveCotizacion(gomock.Any()).Return(nil)
	getCotizador := func(name string) (cotizadores.Cotizador, error) {
		if name == "criptoya" {
			return cotizador, nil
		}
		return nil, fmt.Errorf("cotizador %s no soportado", name)
	}
	cs := services.NewCryptoService(repoCripto, getCotizador)

	err := cs.GuardarCotizacionExterna("Bitcoin", "criptoya")
	assert.Nil(t, err)
}

func TestGuardarCotizacionExterna_Fail(t *testing.T) {

	testCases := []struct {
		expectedError error
		service       *services.CryptoService
		name          string
		cotizacion    criptomonedas.Cotizacion
		userId        string
		api           string
	}{
		{
			expectedError: fmt.Errorf("no se pudo guardar la cotizacion externa para moneda Bitcoin"),
			service: func() *services.CryptoService {
				ctrl := gomock.NewController(t)
				repoCripto := mockRepo.NewMockCryptoRepository(ctrl)
				cs := services.NewCryptoService(repoCripto, cotizadores.GetCotizador)
				return cs
			}(),
			name:       "error de nombre cotizador",
			cotizacion: criptomonedas.Cotizacion{},
			userId:     "userId",
			api:        "criptoLLa"},

		{
			expectedError: fmt.Errorf("error al buscar la criptomoneda Bitcoin en la base de datos"),
			service: func() *services.CryptoService {
				ctrl := gomock.NewController(t)
				repoCripto := mockRepo.NewMockCryptoRepository(ctrl)
				cotizador := mockCotizador.NewMockCotizador(ctrl)
				repoCripto.EXPECT().FindCryptoByName("Bitcoin").Return(&criptomonedas.CriptoMoneda{Nombre: "A", Codigo: "B"}, nil).Times(1)
				cotizador.EXPECT().GetCotizacionExterna("A", "B", "USD").Return(criptomonedas.Cotizacion{}, nil).Times(1)
				repoCripto.EXPECT().FindCryptoByName("Bitcoin").Return(nil, errors.New("error al buscar la criptomoneda")).Times(1)
				getCotizador := func(name string) (cotizadores.Cotizador, error) {
					if name == "criptoya" {
						return cotizador, nil
					}
					return nil, fmt.Errorf("cotizador %s no soportado", name)
				}
				cs := services.NewCryptoService(repoCripto, getCotizador)
				return cs
			}(),
			name:       "error nombre invalido de criptomoneda",
			cotizacion: criptomonedas.Cotizacion{},
			userId:     "userId",
			api:        "criptoya"},

		{
			expectedError: fmt.Errorf("la criptomoneda Bitcoin no está registrada en la base de datos"),
			service: func() *services.CryptoService {
				ctrl := gomock.NewController(t)
				repoCripto := mockRepo.NewMockCryptoRepository(ctrl)
				cotizador := mockCotizador.NewMockCotizador(ctrl)
				repoCripto.EXPECT().FindCryptoByName("Bitcoin").Return(&criptomonedas.CriptoMoneda{Nombre: "A", Codigo: "B"}, nil).Times(1)
				cotizador.EXPECT().GetCotizacionExterna("A", "B", "USD").Return(criptomonedas.Cotizacion{}, nil).Times(1)
				repoCripto.EXPECT().FindCryptoByName("Bitcoin").Return(nil, nil).Times(1)
				getCotizador := func(name string) (cotizadores.Cotizador, error) {
					if name == "criptoya" {
						return cotizador, nil
					}
					return nil, fmt.Errorf("cotizador %s no soportado", name)
				}
				cs := services.NewCryptoService(repoCripto, getCotizador)
				return cs
			}(),
			name:       "error moneda retornada nill cuando se busca en base",
			cotizacion: criptomonedas.Cotizacion{},
			userId:     "userId",
			api:        "criptoya"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.service.GuardarCotizacionExterna("Bitcoin", tc.api)
			assertions := assert.New(t)
			assertions.True(err.Error() == tc.expectedError.Error())
		})
	}
}

func TestGetCotizacion_Succes(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoCripto := mockRepo.NewMockCryptoRepository(ctrl)
	cotizador := mockCotizador.NewMockCotizador(ctrl)
	repoCripto.EXPECT().FindCryptoByName("Bitcoin").Return(&criptomonedas.CriptoMoneda{Nombre: "A", Codigo: "B"}, nil).Times(1)
	cotizador.EXPECT().GetCotizacionExterna("A", "B", "USD").Return(criptomonedas.Cotizacion{}, nil).Times(1)
	getCotizador := func(name string) (cotizadores.Cotizador, error) {
		if name == "criptoya" {
			return cotizador, nil
		}
		return nil, fmt.Errorf("cotizador %s no soportado", name)
	}
	cs := services.NewCryptoService(repoCripto, getCotizador)
	cripto, err := cs.GetCotizacion("criptoya", "Bitcoin", "USD")
	assert.Nil(t, err)
	assert.NotNil(t, cripto)
}

func TestGetCotizacion_Fail(t *testing.T) {

	testCases := []struct {
		expectedError error
		service       *services.CryptoService
		name          string
		cotizacion    criptomonedas.Cotizacion
		api           string
	}{
		{
			expectedError: fmt.Errorf("la criptomoneda Bitcoin no está registrada en la base de datos"),
			service: func() *services.CryptoService {
				ctrl := gomock.NewController(t)
				repoCripto := mockRepo.NewMockCryptoRepository(ctrl)
				cotizador := mockCotizador.NewMockCotizador(ctrl)
				repoCripto.EXPECT().FindCryptoByName("Bitcoin").Return(&criptomonedas.CriptoMoneda{}, errors.New("error al buscar la criptomoneda")).Times(1)
				getCotizador := func(name string) (cotizadores.Cotizador, error) {
					if name == "criptoya" {
						return cotizador, nil
					}
					return nil, fmt.Errorf("cotizador %s no soportado", name)
				}
				cs := services.NewCryptoService(repoCripto, getCotizador)
				return cs
			}(),
			name:       "error de FindCryptoByName, err != nil ",
			cotizacion: criptomonedas.Cotizacion{},
			api:        "criptoya"},

			{
				expectedError: fmt.Errorf("la criptomoneda Bitcoin no está registrada en la base de datos"),
				service: func() *services.CryptoService {
					ctrl := gomock.NewController(t)
					repoCripto := mockRepo.NewMockCryptoRepository(ctrl)
					cotizador := mockCotizador.NewMockCotizador(ctrl)
					repoCripto.EXPECT().FindCryptoByName("Bitcoin").Return(nil, nil).Times(1)
					getCotizador := func(name string) (cotizadores.Cotizador, error) {
						if name == "criptoya" {
							return cotizador, nil
						}
						return nil, fmt.Errorf("cotizador %s no soportado", name)
					}
					cs := services.NewCryptoService(repoCripto, getCotizador)
					return cs
				}(),
				name:       "error de FindCryptoByName, la moneda retornada es null ",
				cotizacion: criptomonedas.Cotizacion{},
				api:        "criptoya"},

				{
					expectedError: fmt.Errorf("el Cotizador criptoLLa no es soportado"),
					service: func() *services.CryptoService {
						ctrl := gomock.NewController(t)
						repoCripto := mockRepo.NewMockCryptoRepository(ctrl)
						cotizador := mockCotizador.NewMockCotizador(ctrl)
						getCotizador := func(name string) (cotizadores.Cotizador, error) {
							if name == "criptoya" {
								return cotizador, nil
							}
							return nil, fmt.Errorf("cotizador %s no soportado", name)
						}
						cs := services.NewCryptoService(repoCripto, getCotizador)
						return cs
					}(),
					name:       "error de FindCotizador cotizador no soportado ",
					cotizacion: criptomonedas.Cotizacion{},
					api:        "criptoLLa"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cotizacion,err := tc.service.GetCotizacion(tc.api,"Bitcoin","USD")
			assertions := assert.New(t)
			assertions.True(err.Error() == tc.expectedError.Error())
			assertions.True(cotizacion == tc.cotizacion)
		})
	}

}

/*
testCases := []struct {
		expectedError error
		service       *services.CryptoService
		name          string
		cotizacion    criptomonedas.Cotizacion
		userId        string
	}{
		expectedError
		name "error al guardar cotizacion"

	}

	cotizacion := criptomonedas.Cotizacion{
		CriptoMoneda_ID: 11,
		Cotizacion:      100,
		Fecha:           time.Now(),
	}*/
