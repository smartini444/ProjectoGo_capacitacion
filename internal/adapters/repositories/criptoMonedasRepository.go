package repositories

//go:generate echo $GOPACKAGE/$GOFILE
//go:generate mockgen -source=./$GOFILE -destination=./mock/$GOFILE -package mock

import (
	"database/sql"
	"log"
	"primerProjecto/internal/entities/criptomonedas"
)

type MySQLCryptoRepository struct {
	db *sql.DB
}

func NewMySQLCryptoRepository(db *sql.DB) *MySQLCryptoRepository {
	return &MySQLCryptoRepository{db: db}
}

type CryptoRepository interface {
	SaveMoneda(cripto criptomonedas.CriptoMoneda) error
	FindAllMonedas() ([]*criptomonedas.CriptoMoneda, error)
	FindCryptoByName(name string) (*criptomonedas.CriptoMoneda, error)
	FindCryptoByCode(codigo string) (*criptomonedas.CriptoMoneda, error)
	FindByMonedaID(id int) (*criptomonedas.CriptoMoneda, error)
	UpdateMoneda(id int, moneda criptomonedas.CriptoMoneda) error

	//cotizaciones
	SaveCotizacion(cripto criptomonedas.Cotizacion) error
	FindByCotizacionID(id int) (*criptomonedas.Cotizacion, error)
	FindAllCotizaciones() ([]*criptomonedas.Cotizacion, error)
	UpdateCotizacion(id int, cotizacion criptomonedas.Cotizacion) error
	FindAllByFilter(filter criptomonedas.CriptoMonedaFilter) ([]criptomonedas.Cotizacion, criptomonedas.Summary, error)
	FindAllByFilterForUser(filter criptomonedas.CriptoMonedaFilter, usuarioId int) ([]criptomonedas.Cotizacion, criptomonedas.Summary, error)
	FindUltimaCotizacion(nombre string) (*criptomonedas.Cotizacion, error)
	BorrarCotizacionManual(cotizacion criptomonedas.Cotizacion) error
	GuardarCotizacionManual(usuarioId int, cotizacion criptomonedas.Cotizacion) (criptomonedas.Cotizacion, error)
	ActualizarCotizacionManual(usuarioId int, cotizacion criptomonedas.Cotizacion) (criptomonedas.Cotizacion, error)
	BorrarCotizacionById(id int) error
}

func (r *MySQLCryptoRepository) SaveMoneda(cripto criptomonedas.CriptoMoneda) error {
	_, err := r.db.Exec("INSERT INTO monedas (Id,Nombre,Codigo) VALUES (?, ?, ?)", cripto.Id, cripto.Nombre, cripto.Codigo)
	if err != nil {
		log.Println("Error al guardar cripto:", err)
		return err
	}
	return nil
}

func (r *MySQLCryptoRepository) FindByMonedaID(id int) (*criptomonedas.CriptoMoneda, error) {
	query := "SELECT id, nombre,codigo FROM monedas WHERE id = ?"
	row := r.db.QueryRow(query, id)
	moneda := criptomonedas.CriptoMoneda{}

	err := row.Scan(&moneda.Id, &moneda.Nombre, &moneda.Codigo)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no se encontro moneda con id %d", id)
			return nil, err
		}
	}
	return &moneda, nil
}

func (r *MySQLCryptoRepository) FindAllMonedas() ([]*criptomonedas.CriptoMoneda, error) {
	query := "SELECT * FROM monedas"
	moneda := criptomonedas.CriptoMoneda{}
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("no se encotraron lineas")
		return nil, err
	}
	defer rows.Close()

	var criptomonedas []*criptomonedas.CriptoMoneda

	for rows.Next() {
		moneda1 := moneda
		err := rows.Scan(&moneda1.Nombre, &moneda1.Id, &moneda.Codigo)
		if err != nil {
			log.Println("Error al escanear fila:", err)
			continue
		}
		criptomonedas = append(criptomonedas, &moneda)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error al iterar sobre las filas:", err)
		return nil, err
	}

	return criptomonedas, nil

}

func (r *MySQLCryptoRepository) UpdateMoneda(id int, moneda criptomonedas.CriptoMoneda) error {
	query := "UPDATE monedas SET nombre = ? WHERE id = ?"
	_, err := r.db.Exec(query, moneda.Nombre, id)
	if err != nil {
		log.Println("Error al actualizar la moneda:", err)
		return err
	}
	return nil
}

func (r *MySQLCryptoRepository) FindCryptoByName(name string) (*criptomonedas.CriptoMoneda, error) {
	query := "SELECT id, nombre, codigo FROM monedas WHERE nombre = ?"
	var cripto criptomonedas.CriptoMoneda
	err := r.db.QueryRow(query, name).Scan(&cripto.Id, &cripto.Nombre, &cripto.Codigo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró la criptomoneda
		}
		return nil, err
	}
	return &cripto, nil
}

func (r *MySQLCryptoRepository) FindCryptoByCode(codigo string) (*criptomonedas.CriptoMoneda, error) {
	query := "SELECT id, nombre, codigo FROM monedas WHERE codigo = ?"
	var cripto criptomonedas.CriptoMoneda
	err := r.db.QueryRow(query, codigo).Scan(&cripto.Id, &cripto.Nombre, &cripto.Codigo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró la criptomoneda
		}
		return nil, err
	}
	return &cripto, nil
}
