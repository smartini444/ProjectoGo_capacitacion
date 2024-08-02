package repositories

//go:generate echo $GOPACKAGE/$GOFILE
//go:generate mockgen -source=./$GOFILE -destination=./mock/$GOFILE -package mock

import (
	"database/sql"
	"fmt"
	"log"
	"primerProjecto/internal/entities/criptomonedas"
	"strings"
)

type MySQLUsuarioRepository struct {
	db *sql.DB
}

func NewMySQLUsuarioRepository(db *sql.DB) *MySQLUsuarioRepository {
	return &MySQLUsuarioRepository{db: db}
}

type UsuarioRepository interface {
	SaveUsuario(usuario criptomonedas.Usuario) (int, error)
	UpdateUsuarioById(id int, usuario criptomonedas.Usuario) error
	FindUsuarioById(id int) (*criptomonedas.Usuario, error)
	FindMonedasByUsuarioID(id int) ([]int, error)
	FindUsuariosByMonedaID(id int) ([]int, error)
	PatchUsuarioByID(id int, updates map[string]interface{}) error
	AgregarMonedaFavorita(idUsuario, idMoneda int) ([]int, error)
	UpdateMonedasDeInteres(usuarioId int, monedas []int) error
	DeleteMonedasDeInteres(usuarioId int) error
	RegistrarAuditoria(usuarioId, cotizacionID int, logOperacion string) error
}

func (r *MySQLUsuarioRepository) SaveUsuario(usuario criptomonedas.Usuario) (int, error) {
	query := `INSERT INTO usuarios (nombre, apellidos, fecha_nacimiento, codigo_usuario, email, tipo_documento, fecha_registro, esta_activo)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, usuario.Nombre, usuario.Apellidos, usuario.Fecha_Nacimiento, usuario.CodigoUsuario, usuario.Email, usuario.TipoDocumento, usuario.Fecha_registro, usuario.Esta_activo)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *MySQLUsuarioRepository) UpdateUsuarioById(id int, usuario criptomonedas.Usuario) error {
	query := `
		UPDATE usuarios
		SET nombre = ?, apellidos = ?, fecha_nacimiento = ?, codigo_usuario = ?, email = ?, tipo_documento = ?, fecha_registro = ?, esta_activo = ?
		WHERE id = ?`
	_, err := r.db.Exec(query,
		usuario.Nombre, usuario.Apellidos, usuario.Fecha_Nacimiento, usuario.CodigoUsuario,
		usuario.Email, usuario.TipoDocumento, usuario.Fecha_registro, usuario.Esta_activo, id,
	)
	if err != nil {
		log.Println("Error al actualizar el Usuario:", err)
		return err
	}
	return nil
}

func (r *MySQLUsuarioRepository) FindUsuarioById(id int) (*criptomonedas.Usuario, error) {
	query := `
		SELECT id, nombre, apellidos, fecha_nacimiento, codigo_usuario, email, tipo_documento, fecha_registro, esta_activo
		FROM usuarios
		WHERE id = ?`
	var usuario criptomonedas.Usuario
	err := r.db.QueryRow(query, id).Scan(
		&usuario.Id, &usuario.Nombre, &usuario.Apellidos, &usuario.Fecha_Nacimiento,
		&usuario.CodigoUsuario, &usuario.Email, &usuario.TipoDocumento,
		&usuario.Fecha_registro, &usuario.Esta_activo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró el usuario
		}
		return nil, err
	}
	return &usuario, nil
}

func (r *MySQLUsuarioRepository) FindMonedasByUsuarioID(id int) ([]int, error) {
	query := "SELECT moneda_id FROM usuario_moneda WHERE usuario_id = ?"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monedasId []int
	for rows.Next() {
		var monedaId int
		if err := rows.Scan(&monedaId); err != nil {
			return nil, err
		}
		monedasId = append(monedasId, monedaId)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return monedasId, nil
}

func (r *MySQLUsuarioRepository) FindUsuariosByMonedaID(id int) ([]int, error) {
	query := "SELECT usuario_id FROM usuario_moneda WHERE moneda_id = ?"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuariosId []int
	for rows.Next() {
		var usuarioId int
		if err := rows.Scan(&usuariosId); err != nil {
			return nil, err
		}
		usuariosId = append(usuariosId, usuarioId)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usuariosId, nil
}

func (r *MySQLUsuarioRepository) PatchUsuarioByID(id int, updates map[string]interface{}) error {
	setParts := []string{}
	args := []interface{}{}
	for key, value := range updates {
		setParts = append(setParts, key+" = ?")
		args = append(args, value)
	}
	args = append(args, id)
	query := fmt.Sprintf("UPDATE usuarios SET %s WHERE id = ?", strings.Join(setParts, ", "))
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *MySQLUsuarioRepository) AgregarMonedaFavorita(idUsuario, idMoneda int) ([]int, error) {
	query := "INSERT INTO usuario_moneda (usuario_id, moneda_id) VALUES (?, ?)"
	_, err := r.db.Exec(query, idUsuario, idMoneda)
	if err != nil {
		log.Printf("Error al asociar usuario con moneda: %s", err)
		return []int{}, err
	}
	return []int{idUsuario, idMoneda}, nil
}

func (r *MySQLUsuarioRepository) UpdateMonedasDeInteres(usuarioId int, monedas []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM usuario_moneda WHERE usuario_id = ?", usuarioId)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, monedaId := range monedas {
		_, err := tx.Exec("INSERT INTO usuario_moneda (usuario_id, moneda_id) VALUES (?, ?)", usuarioId, monedaId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *MySQLUsuarioRepository) DeleteMonedasDeInteres(usuarioId int) error {
	_, err := r.db.Exec("DELETE FROM usuario_moneda WHERE usuario_id = ?", usuarioId)
	return err
}

func (r *MySQLUsuarioRepository) RegistrarAuditoria(usuarioId, cotizacionID int, logOperacion string) error {
	query := "INSERT INTO auditoria_cotizacion (usuario_id, cotizacion_id, log) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, usuarioId, cotizacionID, logOperacion)
	if err != nil {
		log.Println("Error al registrar auditoría:", err)
		return err
	}
	return nil
}
