package services

import (
	"fmt"
	cotizadores "primerProjecto/internal/adapters/cotizadores"
	repositories "primerProjecto/internal/adapters/repositories"
	criptomonedas "primerProjecto/internal/entities/criptomonedas"
)

type CryptoService struct {
	repo         repositories.CryptoRepository
	getCotizador func(name string) (cotizadores.Cotizador, error) // Función para obtener el cotizador
}

// NewCryptoService crea una nueva instancia del servicio de criptomonedas con un cotizador
func NewCryptoService(repo repositories.CryptoRepository, getCotizador func(name string) (cotizadores.Cotizador, error)) *CryptoService {
	return &CryptoService{repo: repo, getCotizador: getCotizador}
}


type CryptoServiceInterface interface {
	GetCotizacion(api, moneda, fiat string) (criptomonedas.Cotizacion, error)
	FindMonedaByID(id int) (*criptomonedas.CriptoMoneda, error)
	SaveMoneda(cripto criptomonedas.CriptoMoneda)
	UpdateMoneda(id int, cripto criptomonedas.CriptoMoneda)
	FindCriptoByNombre(nombre string) (*criptomonedas.CriptoMoneda, error)
	SaveMonedaConCotizacion(nombre, api string) error
}




// Método para encontrar una criptomoneda por ID
func (s *CryptoService) FindMonedaByID(id int) (*criptomonedas.CriptoMoneda, error) {
	return s.repo.FindByMonedaID(id)
}

// guardar moneda normal
func (s *CryptoService) SaveMoneda(cripto criptomonedas.CriptoMoneda) error {
	return s.repo.SaveMoneda(cripto)
}

// Método para actualizar una criptomoneda por ID
func (s *CryptoService) UpdateMoneda(id int, cripto criptomonedas.CriptoMoneda) error {
	return s.repo.UpdateMoneda(id, cripto)
}

func (s *CryptoService) FindCriptoByNombre(nombre string) (*criptomonedas.CriptoMoneda, error) {
	return s.repo.FindCryptoByName(nombre)
}

// guardar moneda y buscar cotizacion en la api especificada
func (s *CryptoService) SaveMonedaConCotizacion(nombre, api string) error {

	// Buscar la criptomoneda por nombre
	cripto, err := s.repo.FindCryptoByName(nombre)
	if err != nil {
		return err
	}

	if cripto != nil {
		return fmt.Errorf("la criptomoneda %s ya está registrada en la base de datos", nombre)
	}
	// Guardar la nueva criptomoneda
	cripto = &criptomonedas.CriptoMoneda{Nombre: nombre}
	if err := s.repo.SaveMoneda(*cripto); err != nil {
		return fmt.Errorf("la criptomoneda %s no se pudo guardar", nombre)
	}

	// Obtener la cotización utilizando el handler apropiado
	cotizacion, Error := s.GetCotizacion(api, nombre, "USD")
	if Error != nil {

		return fmt.Errorf("no se pudo guardar la cotizacion externa para moneda %s", nombre)
	}
	cotizacion.CriptoMoneda_ID = cripto.Id
	s.repo.SaveCotizacion(cotizacion)
	return nil
}
