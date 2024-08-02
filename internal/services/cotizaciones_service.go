package services

import (
	"fmt"
	criptomonedas "primerProjecto/internal/entities/criptomonedas"
)

// Método para guardar una nueva criptomoneda
func (s *CryptoService) SaveCotizacion(cripto criptomonedas.Cotizacion) error {
	return s.repo.SaveCotizacion(cripto)
}

// Método para actualizar una criptomoneda por ID
func (s *CryptoService) UpdateCotizacion(id int, cripto criptomonedas.Cotizacion) error {
	return s.repo.UpdateCotizacion(id, cripto)
}

func (s *CryptoService) FindUltimaCotizacion(nombre string) (*criptomonedas.Cotizacion, error) {
	return s.repo.FindUltimaCotizacion(nombre)
}

// Método para encontrar todas las cotizaciones
func (s *CryptoService) FindAll() ([]*criptomonedas.Cotizacion, error) {
	return s.repo.FindAllCotizaciones()
}

func (s *CryptoService) FindAllByFilter(filter criptomonedas.CriptoMonedaFilter) ([]criptomonedas.Cotizacion, criptomonedas.Summary, error) {
	return s.repo.FindAllByFilter(filter)
}

func (s *CryptoService) FindAllByFilterForUser(filter criptomonedas.CriptoMonedaFilter, userId int) ([]criptomonedas.Cotizacion, criptomonedas.Summary, error) {
	return s.repo.FindAllByFilterForUser(filter, userId)
}

// guardar cotizacion externa
func (s *CryptoService) GuardarCotizacionExterna(nombreMoneda, api string) error {
	cotizacion, err := s.GetCotizacion(api, nombreMoneda, "USD")
	if err != nil {
		return fmt.Errorf("no se pudo guardar la cotizacion externa para moneda %s", nombreMoneda)
	}

	cripto, err := s.repo.FindCryptoByName(nombreMoneda)
	if err != nil {
		return fmt.Errorf("error al buscar la criptomoneda %s en la base de datos", nombreMoneda)
	}
	if cripto == nil {
		return fmt.Errorf("la criptomoneda %s no está registrada en la base de datos", nombreMoneda)
	}

	cotizacion.CriptoMoneda_ID = cripto.Id
	s.repo.SaveCotizacion(cotizacion)
	return nil
}

func (s *CryptoService) GetCotizacion(api, moneda, fiat string) (criptomonedas.Cotizacion, error) {
	cotizador, err := s.getCotizador(api)
	if err != nil {
		return criptomonedas.Cotizacion{}, fmt.Errorf("el Cotizador %s no es soportado", api)
	}
	monedaEnbase, err := s.repo.FindCryptoByName(moneda)
	if err != nil {
		return criptomonedas.Cotizacion{}, fmt.Errorf("la criptomoneda %s no está registrada en la base de datos", moneda)
	}
	if monedaEnbase == nil {
		return criptomonedas.Cotizacion{}, fmt.Errorf("la criptomoneda %s no está registrada en la base de datos", moneda)
	}
	return cotizador.GetCotizacionExterna(monedaEnbase.Nombre, monedaEnbase.Codigo, fiat)
}
