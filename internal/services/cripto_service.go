package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	cotizadores "primerProjecto/internal/adapters/cotizadores"
	repositories "primerProjecto/internal/adapters/repositories"
	criptomonedas "primerProjecto/internal/entities/criptomonedas"
	"strconv"
	"sync"
	"time"
)

type CryptoService struct {
	repo         repositories.CryptoRepository
	getCotizador func(name string) (cotizadores.Cotizador, error) // Función para obtener el cotizador
	tasks        map[string]TaskStatusEntry
	mu           sync.Mutex
}

type TaskStatus struct {
	Status string // "pending", "completed", "failed"
	Data   []byte
}

type TaskStatusEntry struct {
	StatusChan chan TaskStatus
	Status     TaskStatus
}

// NewCryptoService crea una nueva instancia del servicio de criptomonedas con un cotizador
func NewCryptoService(repo repositories.CryptoRepository, getCotizador func(name string) (cotizadores.Cotizador, error)) *CryptoService {
	return &CryptoService{
		repo:         repo,
		getCotizador: getCotizador,
		tasks:        make(map[string]TaskStatusEntry),
	}
}

type CryptoServiceInterface interface {
	GetCotizacion(api, moneda, fiat string) (criptomonedas.Cotizacion, error)
	FindMonedaByID(id int) (*criptomonedas.CriptoMoneda, error)
	SaveMoneda(cripto criptomonedas.CriptoMoneda)
	UpdateMoneda(id int, cripto criptomonedas.CriptoMoneda)
	FindCriptoByNombre(nombre string) (*criptomonedas.CriptoMoneda, error)
	SaveMonedaConCotizacion(nombre, api string) error
	GenerateCSV() ([]byte, error)
	GenerateCSVAsync(taskID string) chan TaskStatus
	GetTaskStatus(taskID string) (TaskStatus, bool)
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

func (s *CryptoService) GenerateCSV() ([]byte, error) {
	monedas, err := s.repo.FindAllMonedas()
	if err != nil {
		return nil, err
	}
	log.Println("Monedas obtenidas:", monedas)

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	defer writer.Flush()

	headers := []string{"ID", "Nombre", "Codigo", "Cotización"}
	if err := writer.Write(headers); err != nil {
		log.Println("Error al escribir encabezados CSV:", err)
		return nil, fmt.Errorf("error al escribir encabezados CSV: %w", err)
	}

	log.Println("Encabezados CSV escritos:", headers)

	for _, moneda := range monedas {
		UltimaCotizacion, err := s.repo.FindUltimaCotizacion(moneda.Nombre)
		if err != nil {
			log.Println("Error al obtener última cotización para", moneda.Nombre, ":", err)
			UltimaCotizacion.Cotizacion = 0
		}

		record := []string{
			strconv.Itoa(moneda.Id),
			moneda.Nombre,
			moneda.Codigo,
			strconv.FormatFloat(UltimaCotizacion.Cotizacion, 'f', 2, 64),
		}

		if err := writer.Write(record); err != nil {
			log.Println("Error al escribir datos CSV:", err)
			return nil, fmt.Errorf("error al escribir datos CSV: %w", err)
		}

		log.Println("Registro CSV escrito:", record)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Println("Error al flushear el writer:", err)
		return nil, err
	}

	log.Println("Datos CSV escritos correctamente")

	return buffer.Bytes(), nil
}

func (s *CryptoService) generateCSVAsync(taskID string) {
	statusChan := make(chan TaskStatus, 1)
	s.mu.Lock()
	s.tasks[taskID] = TaskStatusEntry{StatusChan: statusChan}
	s.mu.Unlock()

	defer close(statusChan)
	status := TaskStatus{Status: "In Progress"}

	csvData, err := s.GenerateCSV()
	if err != nil {
		status.Status = "Failed"
		status.Data = nil
		statusChan <- status
		s.mu.Lock()
		s.tasks[taskID] = TaskStatusEntry{Status: status}
		s.mu.Unlock()
		return
	}

	status.Status = "Completed"
	status.Data = csvData
	statusChan <- status
	s.mu.Lock()
	s.tasks[taskID] = TaskStatusEntry{Status: status}
	s.mu.Unlock()
}

func (s *CryptoService) GetTaskStatus(taskID string) (TaskStatus, bool) {
	s.mu.Lock()
	entry, exists := s.tasks[taskID]
	s.mu.Unlock()

	if !exists {
		return TaskStatus{}, false
	}

	if entry.StatusChan != nil {
		status := <-entry.StatusChan
		return status, true
	}

	return entry.Status, true
}

func (s *CryptoService) GetCSVFile(taskID string) ([]byte, error) {
	s.mu.Lock()
	entry, exists := s.tasks[taskID]
	s.mu.Unlock()

	if !exists {
		return nil, fmt.Errorf("tarea no encontrada")
	}

	if entry.Status.Status != "Completed" {
		return nil, fmt.Errorf("archivo no encontrado o tarea no completada")
	}

	return entry.Status.Data, nil
}

// Función para generar un ID único para la tarea
func generateUniqueID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (s *CryptoService) StartCSVTask() string {
	taskID := generateUniqueID()
	go s.generateCSVAsync(taskID)
	return taskID
}
