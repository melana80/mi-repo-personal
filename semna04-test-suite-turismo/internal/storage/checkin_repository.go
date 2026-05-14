package storage

import (
	"strings"

	"github.com/uleam/awii/turismo/internal/errs"
	"github.com/uleam/awii/turismo/internal/models"
)

// CheckInRepository define las operaciones del repositorio de check-ins.
//
// IMPORTANTE: a diferencia de NegocioRepository y TuristaRepository, este
// repositorio depende de los OTROS dos repositorios para validar que el
// TuristaID y el NegocioID que se quieren guardar realmente existan. Esto
// se llama "validación cruzada" o "integridad referencial" — la versión
// en memoria del concepto que en bases de datos se conoce como FOREIGN KEY.
type CheckInRepository interface {
	Guardar(c models.CheckIn) error
	BuscarPorTurista(turistaID int) ([]models.CheckIn, error)
	Listar() []models.CheckIn
}

// CheckInMemoria es la implementación en memoria de CheckInRepository.
//
// Recibe los repositorios de Turista y Negocio en su constructor —
// inyección de dependencias manual, exactamente como se hará después con
// frameworks. Esto permite probarlo con repos reales o con repos en
// diferentes estados de población.
type CheckInMemoria struct {
	datos    map[int]models.CheckIn
	turistas TuristaRepository
	negocios NegocioRepository
}

// NewCheckInMemoria construye un repositorio vacío que valida sus IDs
// contra los repositorios pasados como parámetros.
//
// Si pasás nil en cualquiera de los dos parámetros, el repositorio se
// construye igual pero los Guardar siempre fallarán con ErrDatosInvalidos
// porque no hay forma de validar las referencias.
func NewCheckInMemoria(turistas TuristaRepository, negocios NegocioRepository) *CheckInMemoria {
	return &CheckInMemoria{
		datos:    make(map[int]models.CheckIn),
		turistas: turistas,
		negocios: negocios,
	}
}

// Guardar agrega un nuevo check-in al repositorio.
//
// Validaciones (en orden):
//  1. Fecha no puede estar vacía → ErrDatosInvalidos
//  2. Calificación debe estar entre 1 y 5 → ErrDatosInvalidos
//  3. TuristaID debe existir en el repositorio de turistas → ErrNoEncontrado
//  4. NegocioID debe existir en el repositorio de negocios → ErrNoEncontrado
//  5. ID no puede estar ya ocupado → ErrYaExiste
//
// La validación cruzada (puntos 3 y 4) se hace ANTES de chequear si el ID
// del check-in ya existe, porque tiene sentido reportar primero un error
// de datos inválidos antes que un conflicto de ID.
func (r *CheckInMemoria) Guardar(c models.CheckIn) error {
	if strings.TrimSpace(c.Fecha) == "" { 
		return errs.ErrDatosInvalidos 
	}
	if c.Calificacion < 1 || c.Calificacion > 5 {
		return errs.ErrDatosInvalidos
	}
	if r.turistas == nil || r.negocios == nil {
		return errs.ErrDatosInvalidos
	}
	if _, err := r.turistas.BuscarPorID(c.TuristaID); err != nil {
		return errs.ErrNoEncontrado
	}
	if _, err := r.negocios.BuscarPorID(c.NegocioID); err != nil {
		return errs.ErrNoEncontrado
	}
	if _, existe := r.datos[c.ID]; existe {
		return errs.ErrYaExiste
	}
	r.datos[c.ID] = c
	return nil
}

// BuscarPorTurista retorna todos los check-ins hechos por un turista.
//
// Validaciones:
//  1. TuristaID negativo → ErrDatosInvalidos
//  2. TuristaID no existe en el repositorio de turistas → ErrNoEncontrado
//
// Si el turista existe pero no tiene check-ins, retorna un slice vacío
// (no nil) sin error.
func (r *CheckInMemoria) BuscarPorTurista(turistaID int) ([]models.CheckIn, error) {
	if turistaID < 0 {
		return nil, errs.ErrDatosInvalidos
	}
	if r.turistas == nil {
		return nil, errs.ErrDatosInvalidos
	}
	if _, err := r.turistas.BuscarPorID(turistaID); err != nil {
		return nil, errs.ErrNoEncontrado
	}
	resultado := make([]models.CheckIn, 0)
	for _, c := range r.datos {
		if c.TuristaID == turistaID {
			resultado = append(resultado, c)
		}
	}
	return resultado, nil
}

// Listar retorna todos los check-ins. Si no hay ninguno retorna un slice
// vacío (no nil).
func (r *CheckInMemoria) Listar() []models.CheckIn {
	resultado := make([]models.CheckIn, 0, len(r.datos))
	for _, c := range r.datos {
		resultado = append(resultado, c)
	}
	return resultado
}
