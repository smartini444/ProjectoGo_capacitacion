package criptomonedas

import "time"

// CriptoMoneda representa una criptomoneda.
// @Description Estructura que define una criptomoneda.
type CriptoMoneda struct {
	// ID es el identificador único de la criptomoneda.
	// @example 1
	Id int `json:"id"`

	// Nombre es el nombre de la criptomoneda.
	// @example Bitcoin
	Nombre string `json:"nombre"`

	// Codigo es el código de la criptomoneda.
	// @example BTC
	Codigo string `json:"codigo"`
}

// Cotizacion representa una cotización de criptomoneda.
// @Description Estructura que define una cotización de criptomoneda.
type Cotizacion struct {
	// ID es el identificador único de la cotización.
	// @example 123
	Id int `json:"id"`

	// CriptoMoneda_ID es el identificador de la criptomoneda asociada.
	// @example 1
	CriptoMoneda_ID int `json:"cripto_id"`

	// Cotizacion es el valor de la criptomoneda en un momento específico.
	// @example 50000.00
	Cotizacion float64 `json:"cotizacion"`

	// Fecha es la fecha y hora en que se registró la cotización.
	// @example 2024-07-29T12:00:00Z
	Fecha time.Time `json:"fecha"`

	// Manual indica si la cotización fue ingresada manualmente.
	// @example true
	Manual bool `json:"manual"`

	// UsuarioId es el identificador del usuario que ingresó la cotización.
	// @example 42
	UsuarioId *int `json:"usuario_id,omitempty"`
}

// CotizacionCompleta representa una cotización completa de criptomoneda.
// @Description Estructura que define una cotización completa de criptomoneda.
type CotizacionCompletas struct {
	// ID es el identificador único de la cotización.
	// @example 123
	Id int `json:"id"`

	// CriptoMoneda_ID es el identificador de la criptomoneda asociada.
	// @example 1
	CriptoMoneda_ID int `json:"cripto_id"`

	// Cotizacion es el valor de la criptomoneda en un momento específico.
	// @example 50000.00
	Cotizacion float64 `json:"cotizacion"`

	// Fecha es la fecha y hora en que se registró la cotización.
	// @example 2024-07-29T12:00:00Z
	Fecha time.Time `json:"fecha"`

	// Manual indica si la cotización fue ingresada manualmente.
	// @example true
	Manual bool `json:"manual"`

	// UsuarioId es el identificador del usuario que ingresó la cotización.
	// @example 42
	UsuarioId *int `json:"usuario_id,omitempty"`
}

// TipoDocumento representa un tipo de documento.
type TipoDocumento string

const (
	// DNI representa un documento nacional de identidad.
	DNI TipoDocumento = "DNI"

	// Pasaporte representa un pasaporte.
	Pasaporte TipoDocumento = "pasaporte"

	// Cedula representa una cédula de identidad.
	Cedula TipoDocumento = "cedula"
)

// Usuario representa un usuario del sistema.
// @Description Estructura que define a un usuario del sistema.
type Usuario struct {
	// ID es el identificador único del usuario.
	// @example 1
	Id int `json:"id"`

	// Nombre es el nombre del usuario.
	// @example Juan
	Nombre string `json:"nombre"`

	// Apellidos son los apellidos del usuario.
	// @example Perez
	Apellidos string `json:"apellido"`

	// Fecha_Nacimiento es la fecha de nacimiento del usuario.
	// @example 1990-01-01T00:00:00Z
	Fecha_Nacimiento time.Time `json:"fecha_Nacimiento"`

	// CodigoUsuario es el código de usuario.
	// @example JP1990
	CodigoUsuario string `json:"codigoUsuario"`

	// Email es el correo electrónico del usuario.
	// @example juan.perez@example.com
	Email string `json:"email"`

	// TipoDocumento es el tipo de documento del usuario.
	// @example DNI
	TipoDocumento TipoDocumento `json:"tipoDocumento"`

	// Fecha_registro es la fecha en que se registró el usuario.
	// @example 2024-07-29T12:00:00Z
	Fecha_registro time.Time `json:"fecha_registro"`

	// Esta_activo indica si el usuario está activo.
	// @example true
	Esta_activo bool `json:"esta_activo"`
}

// CriptoMonedaFilter representa los filtros para buscar criptomonedas.
// @Description Estructura que define los filtros para buscar criptomonedas.
type CriptoMonedaFilter struct {
	// Nombre es el nombre de la criptomoneda.
	// @example Bitcoin
	Nombre *string

	// MinCotizacion es el valor mínimo de la cotización.
	// @example 30000.00
	MinCotizacion *float64

	// MaxCotizacion es el valor máximo de la cotización.
	// @example 60000.00
	MaxCotizacion *float64

	// StartDate es la fecha de inicio del periodo de búsqueda.
	// @example 2024-01-01T00:00:00Z
	StartDate *time.Time

	// EndDate es la fecha de fin del periodo de búsqueda.
	// @example 2024-12-31T23:59:59Z
	EndDate *time.Time

	// PageSize es el tamaño de la página de resultados.
	// @example 10
	PageSize int

	// PageNumber es el número de la página de resultados.
	// @example 1
	PageNumber int
}

// Summary representa un resumen de los resultados de búsqueda.
// @Description Estructura que define un resumen de los resultados de búsqueda.
type Summary struct {
	// TotalResults es el número total de resultados.
	// @example 100
	TotalResults int `json:"totalResults"`

	// PageNumber es el número de la página de resultados.
	// @example 1
	PageNumber int `json:"pageNumber"`

	// PageSize es el tamaño de la página de resultados.
	// @example 10
	PageSize int `json:"pageSize"`

	// CotizacionesValores es una lista de los valores de cotización.
	// @example [30000.00, 35000.00, 40000.00]
	CotizacionesValores []float64 `json:"cotizacionesValores"`

	// CotizacionesFechas es una lista de las fechas de cotización.
	// @example ["2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z", "2024-01-03T00:00:00Z"]
	CotizacionesFechas []string `json:"cotizacionesFechas"`

	// CriptoNombres es una lista de los nombres de las criptomonedas.
	// @example ["Bitcoin", "Ethereum", "Ripple"]
	CriptoNombres []string `json:"criptoNombres"`
}

// @description Estructura de la solicitud para crear un nuevo usuario con sus criptomonedas favoritas.
// @example { "usuario": { "id": 1, "nombre": "John", "apellido": "Doe", "fecha_Nacimiento": "2000-01-01T00:00:00Z", "codigoUsuario": "jdoe", "email": "john.doe@example.com", "tipoDocumento": "DNI", "fecha_registro": "2024-07-29T00:00:00Z", "esta_activo": true }, "monedasFavoritas": [1, 2, 3] }
type UsuarioRequest struct {
	// Usuario contiene la información del usuario.
	// @description Información del usuario.
	Usuario Usuario `json:"usuario"`

	// MonedasFavoritas contiene una lista de IDs de las criptomonedas favoritas del usuario.
	// @description Lista de IDs de las criptomonedas favoritas del usuario.
	MonedasFavoritas []string `json:"monedasFavoritas"`
}
