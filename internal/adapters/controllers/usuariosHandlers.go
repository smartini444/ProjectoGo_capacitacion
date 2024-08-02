package controllers

import (
	"net/http"
	"primerProjecto/internal/entities/criptomonedas"
	"primerProjecto/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UsuarioHandler struct {
	serv *services.UsuarioService
}

func NewUsuarioHandler(service *services.UsuarioService) *UsuarioHandler {
	return &UsuarioHandler{serv: service}
}

// CreateUsuario crea un nuevo usuario junto con sus criptomonedas favoritas.
// @Summary Create a new user
// @Description Create a new user along with their favorite cryptocurrencies
// @Tags users
// @Accept json
// @Produce json
// @Param user body criptomonedas.UsuarioRequest true "User and Favorite Cryptocurrencies"
// @Success 200 {object} map[string]string "message": "Usuario creado exitosamente"
// @Failure 400 {object} map[string]string "error": "Bad Request"
// @Failure 500 {object} map[string]string "error": "Internal Server Error"
// @Router /usuarios [post]
func (h *UsuarioHandler) CreateUsuario(c *gin.Context) {
	var request criptomonedas.UsuarioRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !esMayorDeEdad(request.Usuario.Fecha_Nacimiento) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El usuario debe ser mayor de edad"})
		return
	}

	err := h.serv.CreateUsuario(request.Usuario, request.MonedasFavoritas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario creado exitosamente"})
}

func esMayorDeEdad(fechaNacimiento time.Time) bool {
	hoy := time.Now()
	edad := hoy.Year() - fechaNacimiento.Year()
	if hoy.YearDay() < fechaNacimiento.YearDay() {
		edad--
	}
	return edad >= 18
}

// @Summary Update a user by ID
// @Description Update the details of an existing user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body criptomonedas.Usuario true "User details to update"
// @Success 200 {object} map[string]string "message": "Usuario actualizado exitosamente"
// @Failure 400 {object} map[string]string "error": "ID inválido"
// @Failure 400 {object} map[string]string "error": "Bad Request"
// @Failure 500 {object} map[string]string "error": "Internal Server Error"
// @Router /usuarios/{id} [put]
func (h *UsuarioHandler) UpdateUsuarioByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var usuario criptomonedas.Usuario
	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.serv.UpdateUsuarioById(id, usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado exitosamente"})
}

func (h *UsuarioHandler) UpsertUsuario(c *gin.Context) {
	var request criptomonedas.UsuarioRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Usuario.Id != 0 {
		if err := h.serv.UpdateUsuarioById(request.Usuario.Id, request.Usuario); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := h.serv.CreateUsuario(request.Usuario, request.MonedasFavoritas); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Usuario upserteado exitosamente"})
		return
	}

	for _, moneda := range request.MonedasFavoritas {
		if err := h.serv.GuardarMonedaFavorita(moneda, request.Usuario.Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario upserteado exitosamente"})
}

// @Summary Find user by ID
// @Description Get the details of a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} criptomonedas.Usuario
// @Failure 400 {object} map[string]string "error": "ID inválido"
// @Failure 404 {object} map[string]string "message": "Usuario no encontrado"
// @Failure 500 {object} map[string]string "error": "Internal Server Error"
// @Router /usuarios/{id} [get]
func (h *UsuarioHandler) FindUsuarioByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	usuario, err := h.serv.FindUsuarioByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if usuario == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, usuario)
}

// @Summary Find favorite cryptocurrencies by user ID
// @Description Get the list of favorite cryptocurrencies for a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} criptomonedas.CriptoMoneda
// @Failure 400 {object} map[string]string "error": "ID inválido"
// @Failure 500 {object} map[string]string "error": "Internal Server Error"
// @Router /usuarios/{id}/monedas [get]
func (h *UsuarioHandler) FindMonedasByUsuarioID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	monedas, err := h.serv.FindMonedasByUsuarioID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, monedas)
}

// @Summary Partially update a user by ID
// @Description Partially update the details of an existing user by their ID. This can include updating their favorite cryptocurrencies.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param updates body map[string]interface{} true "User details to update, including favorite cryptocurrencies"
// @Success 200 {object} map[string]string "message": "Usuario actualizado exitosamente"
// @Failure 400 {object} map[string]string "error": "Bad Request"
// @Failure 500 {object} map[string]string "error": "Internal Server Error"
// @Router /usuarios/{id} [patch]
func (h *UsuarioHandler) PatchUsuarioByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extraer las monedas de interés si están presentes en la solicitud
	var monedas []string
	eliminarMonedas := false

	if m, exists := updates["monedas"]; exists {
		if m == nil {
			// Eliminar todas las monedas de interés
			eliminarMonedas = true
		} else {
			monedasInterface, ok := m.([]interface{})
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido para monedas"})
				return
			}
			for _, mi := range monedasInterface {
				monedaCodigo, ok := mi.(string) // JSON numbers are float64
				if !ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": "ID de moneda inválido"})
					return
				}
				monedas = append(monedas, string(monedaCodigo))
			}
		}
		delete(updates, "monedas")
	}

	err = h.serv.PatchUsuarioByID(id, updates, monedas, eliminarMonedas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado exitosamente"})
}

// @Summary Add favorite cryptocurrency to user
// @Description Add a favorite cryptocurrency to a user's list by their ID and the cryptocurrency's name
// @Tags users
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Param nombre query string true "Cryptocurrency Name"
// @Success 200 {object} map[string]string "message": "Moneda favorita guardada exitosamente"
// @Failure 400 {object} map[string]string "error": "ID inválido"
// @Failure 500 {object} map[string]string "error": "Internal Server Error"
// @Router /usuarios/monedaFavorita [post]
func (c *UsuarioHandler) GuardarMonedaFavorita(ctx *gin.Context) {
	monedaNombre := ctx.Query("nombre")
	usuarioId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	c.serv.GuardarMonedaFavorita(monedaNombre, usuarioId)

}

// @Summary Register manual cryptocurrency quote
// @Description Register a manual quote for a cryptocurrency for a specific user by their ID
// @Tags quotes
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Param cotizacion body criptomonedas.Cotizacion true "Cryptocurrency Quote"
// @Success 200 {object} map[string]string "message": "Cotización registrada exitosamente"
// @Failure 400 {object} map[string]string "error": "ID inválido" or "Datos de moneda inválidos"
// @Failure 500 {object} map[string]string "error": "Internal Server Error"
// @Router /cotization/manual [post]
func (c *UsuarioHandler) RegistrarCotizacionManual(ctx *gin.Context) {

	var cotizacion criptomonedas.Cotizacion
	err := ctx.ShouldBindJSON(&cotizacion)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos de moneda inválidos"})
		return
	}
	usuarioId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	c.serv.GuardarCotizacionManual(usuarioId, cotizacion)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cotización manual registrada exitosamente",
	})

}

// @Summary Update manual cryptocurrency quote
// @Description Update a manual quote for a cryptocurrency for a specific user by their ID
// @Tags quotes
// @Accept json
// @Produce json
// @Param usuarioId path int true "User ID"
// @Param cotizacionId path int true "Quote ID"
// @Param cotizacion body criptomonedas.Cotizacion true "Cryptocurrency Quote"
// @Success 200 {object} map[string]string "message": "Cotización actualizada exitosamente"
// @Failure 400 {object} map[string]string "error": "ID inválido" or "Datos de cotización inválidos"
// @Failure 500 {object} map[string]string "error": "Internal Server Error"
// @Router /cotizacion/manual/{usuarioId}/{cotizacionId} [put]
func (c *UsuarioHandler) ActualizarCotizacionManual(ctx *gin.Context) {
	var cotizacion criptomonedas.Cotizacion
	err := ctx.ShouldBindJSON(&cotizacion)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos de cotización inválidos"})
		return
	}

	usuarioId, err := strconv.Atoi(ctx.Param("usuarioId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	cotizacionId, err := strconv.Atoi(ctx.Param("cotizacionId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de cotización inválido"})
		return
	}

	cotizacion.Id = cotizacionId // Asegúrate de asignar el ID de la cotización
	_, err = c.serv.ActualizarCotizacionManual(usuarioId, cotizacion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cotización"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cotización actualizada exitosamente"})
}

// @Summary Delete manual cryptocurrency quote
// @Description Delete a manual quote for a cryptocurrency for a specific user by their ID
// @Tags quotes
// @Accept json
// @Produce json
// @Param id path int true "Quote ID"
// @Success 200 {object} map[string]string "message": "Cotización eliminada exitosamente"
// @Failure 400 {object} map[string]string "error": "ID inválido"
// @Failure 500 {object} map[string]string "error": "Internal Server Error"
// @Router /cotizacion/manual/{id} [delete]
func (c *UsuarioHandler) BorrarCotizacionManual(ctx *gin.Context) {

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	c.serv.BorrarCotizacionManual(id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cotización manual borrada exitosamente",
	})
}
