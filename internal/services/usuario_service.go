package services

import (
	"errors"
	"fmt"
	repositories "primerProjecto/internal/adapters/repositories"
	"primerProjecto/internal/entities/criptomonedas"
)

/*type UsuarioService struct {
	repoUsuario *repositories.MySQLUsuarioRepository
	repoCripto  *repositories.MySQLCryptoRepository
}

func NewUsuarioService(repo *repositories.MySQLUsuarioRepository, repo2 *repositories.MySQLCryptoRepository) *UsuarioService {
	s := &UsuarioService{repoUsuario: repo, repoCripto: repo2}
	return s
}*/

type UsuarioService struct {
	repoUsuario repositories.UsuarioRepository
	repoCripto  repositories.CryptoRepository
}

func NewUsuarioService(repoUsuario repositories.UsuarioRepository, repoCripto repositories.CryptoRepository) *UsuarioService {
	s := &UsuarioService{repoUsuario: repoUsuario, repoCripto: repoCripto}
	return s
}

func (s *UsuarioService) CreateUsuario(usuario criptomonedas.Usuario, monedasFavoritas []string) error {
	id, err := s.repoUsuario.SaveUsuario(usuario)
	if err != nil {
		return err
	}

	for _, monedaCodigo := range monedasFavoritas {
		moneda, err := s.repoCripto.FindCryptoByCode(monedaCodigo)
		if err != nil {
			return err
		}
		_, err = s.repoUsuario.AgregarMonedaFavorita(id, moneda.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UsuarioService) UpdateUsuarioById(id int, usuario criptomonedas.Usuario) error {
	return s.repoUsuario.UpdateUsuarioById(id, usuario)
}

func (s *UsuarioService) FindUsuarioByID(id int) (*criptomonedas.Usuario, error) {
	return s.repoUsuario.FindUsuarioById(id)
}

func (s *UsuarioService) FindMonedasByUsuarioID(id int) ([]int, error) {
	return s.repoUsuario.FindMonedasByUsuarioID(id)
}

func (s *UsuarioService) FindUsuariosByMonedaID(id int) ([]int, error) {
	return s.repoUsuario.FindUsuariosByMonedaID(id)
}

func (s *UsuarioService) PatchUsuarioByID(id int, updates map[string]interface{}, monedas []string, eliminarMonedas bool) error {
	if len(updates) == 0 && !eliminarMonedas && len(monedas) == 0 {
		return errors.New("no hay actualizaciones para realizar")
	}

	if err := s.repoUsuario.PatchUsuarioByID(id, updates); err != nil {
		return err
	}
	var ids []int
	for _, codigo := range monedas {
		moneda, err := s.repoCripto.FindCryptoByCode(codigo)
		if err != nil {
			return nil
		}
		ids = append(ids, moneda.Id)
	}

	if eliminarMonedas {
		if err := s.repoUsuario.DeleteMonedasDeInteres(id); err != nil {
			return err
		}
	} else if len(monedas) > 0 {
		if err := s.repoUsuario.UpdateMonedasDeInteres(id, ids); err != nil {
			return err
		}
	}

	return nil
}

func (s *UsuarioService) UpdateMonedasDeInteres(id int, monedas []string) error {
	var ids []int
	for _, codigo := range monedas {
		moneda, err := s.repoCripto.FindCryptoByCode(codigo)
		if err != nil {
			return nil
		}
		ids = append(ids, moneda.Id)
	}
	if err := s.repoUsuario.UpdateMonedasDeInteres(id, ids); err != nil {
		return err
	}
	return nil
}

func (s *UsuarioService) GuardarMonedaFavorita(nombreMoneda string, UsuarioId int) error {

	// Buscar la criptomoneda por nombre
	cripto, err := s.repoCripto.FindCryptoByName(nombreMoneda)
	if err != nil {
		return err
	}

	if cripto != nil {
		_, err = s.repoUsuario.AgregarMonedaFavorita(UsuarioId, cripto.Id)
		if err != nil {
			return err
		}
		return nil
	}
	// Guardar la nueva criptomoneda
	cripto = &criptomonedas.CriptoMoneda{Nombre: nombreMoneda}
	if err := s.repoCripto.SaveMoneda(*cripto); err != nil {
		return fmt.Errorf("la criptomoneda %s no se pudo guardar", nombreMoneda)
	}
	_, err = s.repoUsuario.AgregarMonedaFavorita(UsuarioId, cripto.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UsuarioService) GuardarCotizacionManual(usuarioId int, cotizacion criptomonedas.Cotizacion) (criptomonedas.Cotizacion, error) {
	cotizacionCreada, err := s.repoCripto.GuardarCotizacionManual(usuarioId, cotizacion)
	if err != nil {
		return cotizacionCreada, fmt.Errorf("la cotizacion no se pudo guardar")
	}
	s.repoUsuario.RegistrarAuditoria(usuarioId, cotizacionCreada.Id, "cotizacion Creada")
	return cotizacionCreada, nil
}

func (s *UsuarioService) ActualizarCotizacionManual(usuarioId int, cotizacion criptomonedas.Cotizacion) (criptomonedas.Cotizacion, error) {
	cotizacionActualizada, err := s.repoCripto.ActualizarCotizacionManual(usuarioId, cotizacion)
	if err != nil {
		return cotizacionActualizada, fmt.Errorf("la cotizacion no se pudo actualizar")
	}
	s.repoUsuario.RegistrarAuditoria(usuarioId, cotizacionActualizada.Id, "cotizacion Actualizada")
	return cotizacionActualizada, nil
}

func (s *UsuarioService) BorrarCotizacionManual(cotizacionId int) (*criptomonedas.Cotizacion, error) {
	cotizacion, err := s.repoCripto.FindByCotizacionID(cotizacionId)
	if err != nil {
		return cotizacion, fmt.Errorf("no se encontro cotizacion de id %v", cotizacionId)
	}

	if !cotizacion.Manual {
		return cotizacion, fmt.Errorf("la cotizacion no es manual, no se puede borrar")
	}
	//s.repoUsuario.RegistrarAuditoria(usuarioId, cotizacionActualizada.Id, "cotizacion Actualizada")
	s.repoCripto.BorrarCotizacionById(cotizacion.Id)

	return cotizacion, nil
}
