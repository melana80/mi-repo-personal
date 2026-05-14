package storage

import (
	"strings"

	"github.com/uleam/awii/turismo/internal/errs"
	"github.com/uleam/awii/turismo/internal/models"
)

// TuristaRepository define las operaciones del repositorio de turistas.
type TuristaRepository interface {
	Guardar(t models.Turista) error
	BuscarPorID(id int) (models.Turista, error)
	Listar() []models.Turista
}

// TuristaMemoria es la implementación en memoria de TuristaRepository.
type TuristaMemoria struct {
	datos map[int]models.Turista
}

// NewTuristaMemoria construye un repositorio vacío listo para usar.
func NewTuristaMemoria() *TuristaMemoria {
	return &TuristaMemoria{
		datos: make(map[int]models.Turista),
	}
}

// Guardar agrega un nuevo turista al repositorio.
//
// Validaciones (en orden):
//  1. Nombre no puede estar vacío → ErrDatosInvalidos
//  2. Nacionalidad no puede estar vacía → ErrDatosInvalidos
//  3. IdiomaPreferido debe estar en models.IdiomasValidos → ErrDatosInvalidos
//  4. ID no puede estar ya ocupado → ErrYaExiste
func (r *TuristaMemoria) Guardar(t models.Turista) error {
	if strings.TrimSpace(t.Nombre) == "" {
		return errs.ErrDatosInvalidos
	}
	if strings.TrimSpace(t.Nacionalidad) == "" {
		return errs.ErrDatosInvalidos
	}
	if !models.IdiomasValidos[t.IdiomaPreferido] {
		return errs.ErrDatosInvalidos
	}
	if _, existe := r.datos[t.ID]; existe {
		return errs.ErrYaExiste
	}
	r.datos[t.ID] = t
	return nil
}




// BuscarPorID retorna el turista con el ID solicitado.
//
// Validaciones:
//  1. ID negativo → ErrDatosInvalidos
//  2. ID no existe → ErrNoEncontrado
func (r *TuristaMemoria) BuscarPorID(id int) (models.Turista, error) {
	if id < 0 {
		return models.Turista{}, errs.ErrDatosInvalidos
	}
	t, existe := r.datos[id]
	if !existe {
		return models.Turista{}, errs.ErrNoEncontrado
	}
	return t, nil
}

// Listar retorna todos los turistas. Si no hay ninguno retorna un slice
// vacío (no nil).
func (r *TuristaMemoria) Listar() []models.Turista {
	resultado := make([]models.Turista, 0, len(r.datos))
	for _, t := range r.datos {
		resultado = append(resultado, t)
	}
	return resultado
}
