package repositories

//go:generate echo $GOPACKAGE/$GOFILE
//go:generate mockgen -source=./$GOFILE -destination=./mock/$GOFILE -package mock

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"primerProjecto/internal/entities/criptomonedas"
	"time"
)

func (r *MySQLCryptoRepository) SaveCotizacion(cripto criptomonedas.Cotizacion) error {
	_, err := r.db.Exec("INSERT INTO cotizaciones (cripto_id, cotizacion, fecha) VALUES (?, ?, ?)", cripto.CriptoMoneda_ID, cripto.Cotizacion, cripto.Fecha)
	if err != nil {
		log.Println("Error al guardar cotizacion:", err)
		return err
	}
	return nil
}

func (r *MySQLCryptoRepository) FindByCotizacionID(id int) (*criptomonedas.Cotizacion, error) {
	query := `
	SELECT c.id, m.nombre, c.cotizacion, c.fecha , c.manual , c.usuario_id
	FROM cotizaciones c
	JOIN monedas m ON c.cripto_id = m.id
	WHERE c.id = ?
`
	row := r.db.QueryRow(query, id)
	moneda := criptomonedas.Cotizacion{}
	var fecha string

	err := row.Scan(&moneda.Id, &moneda.CriptoMoneda_ID, &moneda.Cotizacion, &fecha, &moneda.Manual, &moneda.UsuarioId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no se encontro moneda con id %d", id)
			return nil, err
		}
	}
	// Convertir fecha a time.Time si es necesario
	moneda.Fecha, err = time.Parse("2006-01-02 15:04:05", fecha)
	if err != nil {
		log.Println("Error al convertir fecha:", err)
	}
	return &moneda, nil
}

func (r *MySQLCryptoRepository) FindAllCotizaciones() ([]*criptomonedas.Cotizacion, error) {
	query := "SELECT id, cripto_id, cotizacion, fecha FROM cotizaciones"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("Error al ejecutar la consulta:", err)
		return nil, err
	}
	defer rows.Close()

	var cotizaciones []*criptomonedas.Cotizacion

	for rows.Next() {
		var cotizacion criptomonedas.Cotizacion
		var fecha string

		err := rows.Scan(&cotizacion.Id, &cotizacion.CriptoMoneda_ID, &cotizacion.Cotizacion, &fecha)
		if err != nil {
			log.Println("Error al escanear fila:", err)
			continue
		}

		// Convertir fecha a time.Time si es necesario
		cotizacion.Fecha, err = time.Parse("2006-01-02 15:04:05", fecha)
		if err != nil {
			log.Println("Error al convertir fecha:", err)
			continue
		}

		cotizaciones = append(cotizaciones, &cotizacion)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error al iterar sobre las filas:", err)
		return nil, err
	}

	return cotizaciones, nil
}

func (r *MySQLCryptoRepository) UpdateCotizacion(id int, cotizacion criptomonedas.Cotizacion) error {
	query := "UPDATE cotizaciones SET cotizacion = ?, fecha = ? WHERE id = ?"
	_, err := r.db.Exec(query, cotizacion.Cotizacion, cotizacion.Fecha, id)
	if err != nil {
		log.Println("Error al actualizar la moneda:", err)
		return err
	}
	return nil
}

func (r *MySQLCryptoRepository) FindAllByFilter(filter criptomonedas.CriptoMonedaFilter) ([]criptomonedas.Cotizacion, criptomonedas.Summary, error) {
	query := `
        SELECT 
            c.id, c.cotizacion, c.fecha, c.cripto_id 
        FROM 
            cotizaciones c
        JOIN 
            monedas cm ON c.cripto_id = cm.id 
        WHERE 
            1=1`
	args := []interface{}{}

	appliedFilters := make(map[string]interface{})

	if filter.Nombre != nil {
		query += " AND cm.nombre LIKE ?"
		args = append(args, "%"+*filter.Nombre+"%")
		appliedFilters["Nombre"] = *filter.Nombre
	}
	if filter.MinCotizacion != nil {
		query += " AND c.cotizacion >= ?"
		args = append(args, *filter.MinCotizacion)
		appliedFilters["MinCotizacion"] = *filter.MinCotizacion
	}
	if filter.MaxCotizacion != nil {
		query += " AND c.cotizacion <= ?"
		args = append(args, *filter.MaxCotizacion)
		appliedFilters["MaxCotizacion"] = *filter.MaxCotizacion
	}
	if filter.StartDate != nil {
		query += " AND c.fecha >= ?"
		args = append(args, *filter.StartDate)
		appliedFilters["StartDate"] = *filter.StartDate
	}
	if filter.EndDate != nil {
		query += " AND c.fecha <= ?"
		args = append(args, *filter.EndDate)
		appliedFilters["EndDate"] = *filter.EndDate
	}

	// Add pagination
	query += " LIMIT ? OFFSET ?"
	args = append(args, filter.PageSize, filter.PageSize*(filter.PageNumber-1))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, criptomonedas.Summary{}, err
	}
	defer rows.Close()

	var cotizaciones []criptomonedas.Cotizacion
	for rows.Next() {
		var cotizacion criptomonedas.Cotizacion
		var cripto criptomonedas.CriptoMoneda
		var fechaString string
		if err := rows.Scan(&cotizacion.Id, &cotizacion.Cotizacion, &fechaString, &cripto.Id); err != nil {
			return nil, criptomonedas.Summary{}, err
		}
		cotizacion.Fecha, err = time.Parse("2006-01-02 15:04:05", fechaString)
		if err != nil {
			log.Println("Error al convertir fecha:", err)
			continue
		}
		//cotizacion.CriptoMoneda = cripto
		cotizaciones = append(cotizaciones, cotizacion)
	}

	// Calcular resumen
	summary := criptomonedas.Summary{
		TotalResults: len(cotizaciones),
		//Filters:      appliedFilters,
		PageNumber: filter.PageNumber,
		PageSize:   filter.PageSize,
	}

	return cotizaciones, summary, nil

}

// FindUltimaCotizacion retrieves the latest quotation for a given cryptocurrency name.
// @Summary Retrieve the latest quotation for a given cryptocurrency name
// @Description Retrieves the most recent quotation for a cryptocurrency identified by its name.
// @Tags cryptocurrencies
// @Accept json
// @Produce json
// @Param nombre query string true "Cryptocurrency name"
// @Success 200 {object} criptomonedas.Cotizacion "Success response with the latest quotation"
// @Failure 400 {object} map[string]string "error": "Bad Request"
// @Failure 404 {object} map[string]string "error": "Not Found"
// @Failure 500 {object} map[string]string "error": "Internal Server Error"
// @Router /cryptocurrencies/latest [get]
func (r *MySQLCryptoRepository) FindUltimaCotizacion(nombre string) (*criptomonedas.Cotizacion, error) {
	query := `
		SELECT
		c.id, c.cotizacion, c.fecha, c.cripto_id
	FROM
		cotizaciones c
	JOIN
		monedas cm ON c.cripto_id = cm.id
	WHERE
		cm.nombre = ?
	ORDER BY
		c.fecha DESC
	LIMIT 1
`

	row := r.db.QueryRow(query, nombre)
	cotizacion := criptomonedas.Cotizacion{}

	var fechaString string
	err := row.Scan(&cotizacion.Id, &cotizacion.Cotizacion, &fechaString, &cotizacion.CriptoMoneda_ID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no se encontro moneda con nombre %s", nombre)
			return nil, err
		}
		// Convertir fecha a time.Time
		cotizacion.Fecha, err = time.Parse("2006-01-02 15:04:05", fechaString)
		if err != nil {
			log.Println("Error al convertir fecha:", err)
		}
	}
	return &cotizacion, nil
}

func (r *MySQLCryptoRepository) FindAllByFilterForUser(filter criptomonedas.CriptoMonedaFilter, usuarioId int) ([]criptomonedas.Cotizacion, criptomonedas.Summary, error) {
	query := `
        SELECT 
            c.id, c.cotizacion, c.fecha, c.cripto_id,
            JSON_ARRAYAGG(c.cotizacion) AS cotizaciones_valores,
			JSON_ARRAYAGG(c.fecha) AS cotizaciones_fechas, 
			JSON_ARRAYAGG(cm.nombre) AS cripto_nombres  
        FROM 
            cotizaciones c
        JOIN 
            monedas cm ON c.cripto_id = cm.id 
        JOIN
            usuario_moneda um ON um.moneda_id = cm.id
        WHERE 
            um.usuario_id = ?`
	args := []interface{}{usuarioId}

	if filter.Nombre != nil {
		query += " AND cm.nombre LIKE ?"
		args = append(args, "%"+*filter.Nombre+"%")
	}
	if filter.MinCotizacion != nil {
		query += " AND c.cotizacion >= ?"
		args = append(args, *filter.MinCotizacion)
	}
	if filter.MaxCotizacion != nil {
		query += " AND c.cotizacion <= ?"
		args = append(args, *filter.MaxCotizacion)
	}
	if filter.StartDate != nil {
		query += " AND c.fecha >= ?"
		args = append(args, *filter.StartDate)
	}
	if filter.EndDate != nil {
		query += " AND c.fecha <= ?"
		args = append(args, *filter.EndDate)
	}

	// Add pagination
	query += " GROUP BY c.id, c.cotizacion, c.fecha, c.cripto_id"
	query += " LIMIT ? OFFSET ?"
	args = append(args, filter.PageSize, filter.PageSize*(filter.PageNumber-1))

	log.Println("Consulta SQL:", query)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, criptomonedas.Summary{}, err
	}
	defer rows.Close()

	var cotizaciones []criptomonedas.Cotizacion
	var summary criptomonedas.Summary
	for rows.Next() {
		var cotizacion criptomonedas.Cotizacion
		var fechaString string
		var cotizacionesValoresJSON, cotizacionesFechasJSON, criptoNombresJSON string

		if err := rows.Scan(&cotizacion.Id, &cotizacion.Cotizacion, &fechaString, &cotizacion.CriptoMoneda_ID, &cotizacionesValoresJSON, &cotizacionesFechasJSON, &criptoNombresJSON); err != nil {
			return nil, criptomonedas.Summary{}, err
		}

		cotizacion.Fecha, err = time.Parse(time.RFC3339, fechaString)
		if err != nil {
			log.Println("Error al convertir fecha:", err)
			continue
		}
		cotizaciones = append(cotizaciones, cotizacion)

		// Parse the summary JSON arrays for the first row only (since it's aggregated)
		if summary.TotalResults == 0 {
			var cotizacionesValores []float64
			var cotizacionesFechas []string
			var criptoNombres []string

			if err := json.Unmarshal([]byte(cotizacionesValoresJSON), &cotizacionesValores); err != nil {
				log.Println("Error al parsear cotizaciones valores JSON:", err)
				return nil, criptomonedas.Summary{}, err
			}
			if err := json.Unmarshal([]byte(cotizacionesFechasJSON), &cotizacionesFechas); err != nil {
				log.Println("Error al parsear cotizaciones fechas JSON:", err)
				return nil, criptomonedas.Summary{}, err
			}
			if err := json.Unmarshal([]byte(criptoNombresJSON), &criptoNombres); err != nil {
				log.Println("Error al parsear cripto nombres JSON:", err)
				return nil, criptomonedas.Summary{}, err
			}

			summary.CotizacionesValores = cotizacionesValores
			summary.CotizacionesFechas = cotizacionesFechas
			summary.CriptoNombres = criptoNombres
		}
	}

	summary.PageNumber = filter.PageNumber
	summary.PageSize = filter.PageSize
	summary.TotalResults = len(cotizaciones)

	return cotizaciones, summary, nil
}

func (r *MySQLCryptoRepository) BorrarCotizacionById(cotizacionId int) error {
	query := "DELETE FROM cotizaciones WHERE id = ?"
	args := []interface{}{cotizacionId}
	_, err := r.db.Exec(query, args)
	if err != nil {
		return fmt.Errorf("no se borro la cotizacion de id %v", cotizacionId)
	}
	return nil
}

func (r *MySQLCryptoRepository) GuardarCotizacionManual(usuarioId int, cotizacion criptomonedas.Cotizacion) (criptomonedas.Cotizacion, error) {
	var cotizacionCompleta criptomonedas.Cotizacion = criptomonedas.Cotizacion{
		Id:              cotizacion.Id,
		Cotizacion:      cotizacion.Cotizacion,
		Fecha:           cotizacion.Fecha,
		CriptoMoneda_ID: cotizacion.CriptoMoneda_ID,
		Manual:          true,
		UsuarioId:       &usuarioId,
	}

	// Inserta la cotización completa en la base de datos
	result, err := r.db.Exec(
		"INSERT INTO cotizaciones (cripto_id, cotizacion, fecha, manual, usuario_id) VALUES (?, ?, ?, TRUE, ?)",
		cotizacionCompleta.CriptoMoneda_ID,
		cotizacionCompleta.Cotizacion,
		cotizacionCompleta.Fecha,
		usuarioId,
	)
	if err != nil {
		log.Println("Error al guardar cripto:", err)
		return criptomonedas.Cotizacion{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return criptomonedas.Cotizacion{}, err
	}
	cotizacionCompleta.Id = int(id)

	return cotizacionCompleta, nil
}

func (r *MySQLCryptoRepository) ActualizarCotizacionManual(usuarioId int, cotizacion criptomonedas.Cotizacion) (criptomonedas.Cotizacion, error) {
	// Construye la consulta SQL
	query := "UPDATE cotizaciones SET cripto_id = ?, cotizacion = ?, fecha = ?, manual = TRUE, usuario_id = ? WHERE id = ?"

	// Imprime la consulta SQL con los parámetros
	fmt.Printf("Ejecutando consulta SQL: %s\n", query)
	fmt.Printf("Parámetros: cripto_id=%d, cotizacion=%f, fecha=%s, usuario_id=%d, id=%d\n",
		cotizacion.CriptoMoneda_ID,
		cotizacion.Cotizacion,
		cotizacion.Fecha.Format("2006-01-02 15:04:05"), // Formato de fecha según tu base de datos
		usuarioId,
		cotizacion.Id,
	)

	// Ejecuta la consulta SQL
	_, err := r.db.Exec(
		query,
		cotizacion.CriptoMoneda_ID,
		cotizacion.Cotizacion,
		cotizacion.Fecha,
		usuarioId,
		cotizacion.Id,
	)
	if err != nil {
		log.Println("Error al actualizar cotización:", err)
		return criptomonedas.Cotizacion{}, err
	}
	return cotizacion, nil
}

func (r *MySQLCryptoRepository) BorrarCotizacionManual(cotizacion criptomonedas.Cotizacion) error {
	// Primero, borrar la cotización
	_, err := r.db.Exec("DELETE FROM cotizaciones WHERE id = ?", cotizacion.Id)
	if err != nil {
		log.Println("Error al borrar la cotización:", err)
		return err
	}

	return nil
}
