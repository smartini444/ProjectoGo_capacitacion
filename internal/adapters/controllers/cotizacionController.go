package controllers

import (
	"log"
	"net/http"
	"primerProjecto/internal/entities/criptomonedas"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Save a quotation
// @Description Save a new quotation
// @Tags cryptocurrencies
// @Accept json
// @Produce json
// @Param cotizacion body criptomonedas.Cotizacion true "Cotizacion to be saved"
// @Success 200 {object} gin.H "Successful response with a message"
// @Failure 400 {object} gin.H "Bad Request"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /cotizaciones [post]
func (c *CryptoController) RegistrarCotizacion(ctx *gin.Context) {
	var cotizacion criptomonedas.Cotizacion
	err := ctx.ShouldBindJSON(&cotizacion)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos de moneda inválidos"})
		return
	}
	err = c.serv.SaveCotizacion(cotizacion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar la cotizacion"})
		log.Printf("Error al registrar la cotizacion: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "cotizacion registrada correctamente"})
}

func (c *CryptoController) FindAllByFilter(ctx *gin.Context) {
	var filter criptomonedas.CriptoMonedaFilter

	if nombre := ctx.Query("nombre"); nombre != "" {
		filter.Nombre = &nombre
	}
	if minCotizacion := ctx.Query("min_cotizacion"); minCotizacion != "" {
		min, err := strconv.ParseFloat(minCotizacion, 64)
		if err == nil {
			filter.MinCotizacion = &min
		}
	}
	if maxCotizacion := ctx.Query("max_cotizacion"); maxCotizacion != "" {
		max, err := strconv.ParseFloat(maxCotizacion, 64)
		if err == nil {
			filter.MaxCotizacion = &max
		}
	}
	if startDate := ctx.Query("start_date"); startDate != "" {
		start, err := time.Parse(time.RFC3339, startDate)
		if err == nil {
			filter.StartDate = &start
		}
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		end, err := time.Parse(time.RFC3339, endDate)
		if err == nil {
			filter.EndDate = &end
		}
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10 // Default page size
	}
	filter.PageSize = pageSize

	pageNumber, err := strconv.Atoi(ctx.Query("page_number"))
	if err != nil || pageNumber <= 0 {
		pageNumber = 1 // Default page number
	}
	filter.PageNumber = pageNumber

	monedas, summary, err := c.serv.FindAllByFilter(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las criptomonedas"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"summary": summary,
		"data":    monedas,
	})
}

// @Summary Get latest quote by cryptocurrency name
// @Description Get the latest quote for a given cryptocurrency name
// @Tags cryptocurrencies
// @Accept json
// @Produce json
// @Param nombre path string true "Cryptocurrency name"
// @Success 200 {object} criptomonedas.Cotizacion "Successful response with the latest quote"
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /cryptocurrencies/lastcotization/{nombre} [get]
func (c CryptoController) FindUltimaCotizacion(ctx *gin.Context) {
	nombre := ctx.Param("nombre")

	Cotizacion, err := c.serv.FindUltimaCotizacion(nombre)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la criptomoneda"})
		log.Printf("Error al obtener la criptomoneda: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, Cotizacion)
}

func (c CryptoController) SaveCotizacionExterna(ctx *gin.Context) {
	monedaNombre := ctx.Query("nombre")
	api := ctx.Query("api")

	err := c.serv.GuardarCotizacionExterna(monedaNombre, api)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar la cotizacion"})
		log.Printf("Error al registrar la cotizacion: %s", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "cotizacion guardada correctamente"})
}

// @Summary Find all cryptocurrencies by filter
// @Description Find all cryptocurrencies by filter for a specific user
// @Tags cryptocurrencies
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param nombre query string false "Nombre"
// @Param min_cotizacion query number false "Minimum Cotizacion"
// @Param max_cotizacion query number false "Maximum Cotizacion"
// @Param start_date query string false "Start Date in RFC3339 format"
// @Param end_date query string false "End Date in RFC3339 format"
// @Param page_size query int true "Page Size"
// @Param page_number query int true "Page Number"
// @Success 200 {object} map[string]interface{} "Successful response with summary and data"
// @Router /usuarios/{id}/cotizaciones [get]
func (c *CryptoController) FindAllByFilterUsuario(ctx *gin.Context) {
	// Obtener el ID del usuario desde el parámetro de la URL
	usuarioId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})

		return
	}

	// Crear y poblar el filtro
	var filter criptomonedas.CriptoMonedaFilter

	if nombre := ctx.Query("nombre"); nombre != "" {
		filter.Nombre = &nombre
	}
	if minCotizacion := ctx.Query("min_cotizacion"); minCotizacion != "" {
		min, err := strconv.ParseFloat(minCotizacion, 64)
		if err == nil {
			filter.MinCotizacion = &min
		}
	}
	if maxCotizacion := ctx.Query("max_cotizacion"); maxCotizacion != "" {
		max, err := strconv.ParseFloat(maxCotizacion, 64)
		if err == nil {
			filter.MaxCotizacion = &max
		}
	}
	if startDate := ctx.Query("start_date"); startDate != "" {
		start, err := time.Parse(time.RFC3339, startDate)
		if err == nil {
			filter.StartDate = &start
		}
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		end, err := time.Parse(time.RFC3339, endDate)
		if err == nil {
			filter.EndDate = &end
		}
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10 // Default page size
	}
	filter.PageSize = pageSize

	pageNumber, err := strconv.Atoi(ctx.Query("page_number"))
	if err != nil || pageNumber <= 0 {
		pageNumber = 1 // Default page number
	}
	filter.PageNumber = pageNumber

	// Llamar al servicio con el filtro y el ID del usuario
	cotizaciones, summary, err := c.serv.FindAllByFilterForUser(filter, usuarioId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las criptomonedas"})
		log.Println("Error al obtener las criptomonedas:", err)
		return
	}

	// Enviar la respuesta
	ctx.JSON(http.StatusOK, gin.H{
		"summary": summary,
		"data":    cotizaciones,
	})
}


