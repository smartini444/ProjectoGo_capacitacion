package controllers

import (
	"log"
	"net/http"
	"primerProjecto/internal/entities/criptomonedas"
	"primerProjecto/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CryptoController struct {
	serv *services.CryptoService
}

func NewCryptoController(service *services.CryptoService) *CryptoController {
	return &CryptoController{serv: service}
}

func (c *CryptoController) FindAll(ctx *gin.Context) {
	criptomonedas, err := c.serv.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las criptomonedas"})
		log.Printf("Error al obtener las criptomonedas: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, criptomonedas)
}

func (c *CryptoController) RegistrarCriptoMoneda(ctx *gin.Context) {
	var moneda criptomonedas.CriptoMoneda
	err := ctx.ShouldBindJSON(&moneda)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos de moneda inválidos"})
		return
	}
	err = c.serv.SaveMoneda(moneda)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar la moneda"})
		log.Printf("Error al registrar la moneda: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Moneda registrada correctamente"})
}

// @Summary Get cryptocurrency by ID
// @Description Retrieve a cryptocurrency by its ID
// @Tags cryptocurrencies
// @Accept json
// @Produce json
// @Param id path int true "Cryptocurrency ID"
// @Success 200 {object} criptomonedas.CriptoMoneda
// @Failure 400 {object} map[string]string "error": "ID inválido"
// @Failure 500 {object} map[string]string "error": "Error al obtener la criptomoneda"
// @Router /cryptocurrencies/{id} [get]
func (c *CryptoController) FindMonedaByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	criptoMoneda, err := c.serv.FindMonedaByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la criptomoneda"})
		log.Printf("Error al obtener la criptomoneda: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, criptoMoneda)
}

// @Summary Update cryptocurrency by ID
// @Description Update the details of a cryptocurrency by its ID
// @Tags cryptocurrencies
// @Accept json
// @Produce json
// @Param id path int true "Cryptocurrency ID"
// @Param cryptocurrency body criptomonedas.CriptoMoneda true "Cryptocurrency Data"
// @Success 200 {object} map[string]string "message": "Successful response with a message"
// @Failure 400 {object} map[string]string "error": "Bad Request"
// @Failure 500 {object} map[string]string "error": "Internal Server Error"
// @Router /cryptocurrencies/{id} [put]
func (c *CryptoController) HandleUpdateCryptoByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var moneda criptomonedas.CriptoMoneda
	err = ctx.ShouldBindJSON(&moneda)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos de moneda inválidos"})
		return
	}
	err = c.serv.UpdateMoneda(id, moneda)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la moneda"})
		log.Printf("Error al actualizar la moneda: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Moneda actualizada correctamente"})
}

func (c *CryptoController) FindMondaByNombre(ctx *gin.Context) {
	nombre := ctx.Param("nombre")
	moneda, err := c.serv.FindCriptoByNombre(nombre)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar moneda"})
		log.Printf("Error al buscar la moneda: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, moneda)
}

func (c CryptoController) SaveMonedaConCotizacion(ctx *gin.Context) {
	monedaNombre := ctx.Query("nombre")
	api := ctx.Query("api")

	err := c.serv.SaveMonedaConCotizacion(monedaNombre, api)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar la moneda"})
		log.Printf("Error al registrar la moneda: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Moneda registrada correctamente"})
}

// DownloadCSV godoc
// @Summary      Generar CSV sincrónico
// @Description  Genera un archivo CSV con datos de criptomonedas de forma sincrónica
// @Tags         csv
// @Produce      text/csv
// @Success      200  {file}  file
// @Failure      500  {string}  string "Error al generar el archivo CSV"
// @Router       /csv/sync/generate [get]
func (c *CryptoController) DownloadCSV(ctx *gin.Context) {
	data, err := c.serv.GenerateCSV()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el archivo CSV"})
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename=monedas.csv")
	ctx.Header("Content-Type", "text/csv")
	ctx.Writer.Write(data)
}

// StartCSVTask godoc
// @Summary      Iniciar tarea asíncrona de generación de CSV
// @Description  Inicia una tarea para generar un archivo CSV con datos de criptomonedas de forma asíncrona
// @Tags         csv
// @Produce      json
// @Success      200  {string}  string "task_id"
// @Router       /csv/async/generate [post]
func (c *CryptoController) StartCSVTask(ctx *gin.Context) {
	taskID := c.serv.StartCSVTask()
	ctx.JSON(http.StatusOK, gin.H{"task_id": taskID})
}

// GetTaskStatus godoc
// @Summary      Obtener el estado de una tarea de generación de CSV
// @Description  Obtiene el estado de una tarea asíncrona de generación de CSV mediante el ID de la tarea
// @Tags         csv
// @Produce      json
// @Param        task_id  path      string  true  "ID de la tarea"
// @Success      200  {string}  string "status"
// @Failure      404  {string}  string "Tarea no encontrada"
// @Router       /csv/async/status/{task_id} [get]
func (c *CryptoController) GetTaskStatus(ctx *gin.Context) {
	taskID := ctx.Param("task_id")
	status, exists := c.serv.GetTaskStatus(taskID)
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": status.Status})
}
// DownloadCSVFile godoc
// @Summary      Descargar archivo CSV generado
// @Description  Descarga el archivo CSV generado asíncronamente mediante el ID de la tarea
// @Tags         csv
// @Produce      text/csv
// @Param        task_id  path      string  true  "ID de la tarea"
// @Success      200  {file}  file
// @Failure      404  {string}  string "Error al descargar el archivo CSV"
// @Router       /csv/async/download/{task_id} [get]
func (c *CryptoController) DownloadCSVFile(ctx *gin.Context) {
	taskID := ctx.Param("task_id")
	data, err := c.serv.GetCSVFile(taskID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename=monedas.csv")
	ctx.Header("Content-Type", "text/csv")
	ctx.Writer.Write(data)
}
