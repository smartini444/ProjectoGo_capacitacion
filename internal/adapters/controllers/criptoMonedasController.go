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
