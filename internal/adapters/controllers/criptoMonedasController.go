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

// @Summary Get all cryptocurrencies
// @Description Retrieve a list of all cryptocurrencies
// @Tags cryptocurrencies
// @Accept  json
// @Produce  json
// @Success 200 {array} criptomonedas.CriptoMoneda
// @Router /cryptocurrencies [get]
func (c *CryptoController) FindAll(ctx *gin.Context) {
	criptomonedas, err := c.serv.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las criptomonedas"})
		log.Printf("Error al obtener las criptomonedas: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, criptomonedas)
}

// @Summary Get all cryptocurrencies cotizadas entre dos fechas dadas
// @Tags cryptocurrencies
// @Accept  json
// @Produce  json
// @Success 200 {array} criptomonedas.CriptoMoneda
// @Router /cryptocurrencies/fechas [get]
//deprecado?
/*
func (c *CryptoController) FindAllByDate(ctx *gin.Context) {

	fechaInicioString := ctx.Param("fechaInicio")
	fechaFinString := ctx.Param("fechaFin")

	const layout = "2006-01-02T15:04:05Z07:00"

	// Parsear las fechas
	fechaInicio, err := time.Parse(layout, fechaInicioString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fecha de inicio inválida"})
		return
	}

	fechaFin, err := time.Parse(layout, fechaFinString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fecha de fin inválida"})
		return
	}
	criptomonedas, err := c.serv.FindCotizacionesPorFecha(fechaInicio, fechaFin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las criptomonedas"})
		log.Printf("Error al obtener las criptomonedas: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, criptomonedas)
}*/

// @Summary Save una criptomoneda
// @Description Save una nueva criptomoneda
// @Tags cryptocurrencies
// @Accept  json
// @Produce  json
// @Success 200 {array} criptomonedas.CriptoMoneda
// @Router /cryptocurrencies [post]
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
// @Accept  json
// @Produce  json
// @Param id query int true "Cryptocurrency ID"
// @Success 200 {object} criptomonedas.CriptoMoneda
// @Router /cryptocurrencies/cryptocurrency/{id} [get]
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
// @Success 200 {object} gin.H "Successful response with a message"
// @Failure 400 {object} gin.H "Bad Request"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /cryptocurrency/{id} [put]
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

// @Summary Save cryptocurrency with quote
// @Description Save a cryptocurrency along with its quote using the provided name and API
// @Tags cryptocurrencies
// @Accept json
// @Produce json
// @Param nombre query string true "Cryptocurrency name"
// @Param api query string true "API for obtaining the quote"
// @Success 200 {object} gin.H "Successful response with a message"
// @Failure 400 {object} gin.H "Bad Request"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /cotization [post]
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
