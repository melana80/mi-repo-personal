// Package storage define los repositorios del dominio de Turismo.
//
// Cada entidad (Negocio, Turista, CheckIn) tiene su propia interface y su
// implementación en memoria. Los tests del taller del Día B se escriben
// contra estas implementaciones.
package storage

import (
	"strings"

	"github.com/uleam/awii/turismo/internal/errs"
	"github.com/uleam/awii/turismo/internal/models"
)

// NegocioRepository define las operaciones del repositorio de negocios.
//
// La interface existe para que en el futuro podamos cambiar la
// implementación (a SQLite, PostgreSQL, etc.) sin tocar el código que la
// consume. Por ahora solo hay una implementación: NegocioMemoria.
type NegocioRepository interface {
	Guardar(n models.Negocio) error
	BuscarPorID(id int) (models.Negocio, error)
	Listar() []models.Negocio
	Eliminar(id int) error
}

// NegocioMemoria es la implementación en memoria de NegocioRepository.
// Los datos se guardan en un map indexado por ID. NO es seguro para uso
// concurrente — se asume que cada test crea su propio repo aislado.
type NegocioMemoria struct {
	datos map[int]models.Negocio
}

// NewNegocioMemoria construye un repositorio vacío listo para usar.
func NewNegocioMemoria() *NegocioMemoria {
	return &NegocioMemoria{
		datos: make(map[int]models.Negocio),
	}
}

// Guardar agrega un nuevo negocio al repositorio.
//
// Validaciones (en orden):
//  1. Nombre no puede estar vacío (después de TrimSpace) → ErrDatosInvalidos
//  2. Tipo debe estar en models.TiposValidos → ErrDatosInvalidos
//  3. IdiomasHablados no puede estar vacío → ErrDatosInvalidos
//  4. Cada idioma debe estar en models.IdiomasValidos → ErrDatosInvalidos
//  5. ID no puede estar ya ocupado → ErrYaExiste
func (r *NegocioMemoria) Guardar(n models.Negocio) error {
	if strings.TrimSpace(n.Nombre) == "" {
		return errs.ErrDatosInvalidos
	}
	if !models.TiposValidos[n.Tipo] {
		return errs.ErrDatosInvalidos
	}
	if len(n.IdiomasHablados) == 0 {
		return errs.ErrDatosInvalidos
	}
	for _, idioma := range n.IdiomasHablados {
		if !models.IdiomasValidos[idioma] {
			return errs.ErrDatosInvalidos
		}
	}
	if _, existe := r.datos[n.ID]; existe {
		return errs.ErrYaExiste
	}
	r.datos[n.ID] = n
	return nil
}
//


// BuscarPorID retorna el negocio con el ID solicitado.
//
// Validaciones:
//  1. ID negativo → ErrDatosInvalidos
//  2. ID no existe → ErrNoEncontrado
func (r *NegocioMemoria) BuscarPorID(id int) (models.Negocio, error) {
	if id < 0 {
		return models.Negocio{}, errs.ErrDatosInvalidos
	}
	n, existe := r.datos[id]
	if !existe {
		return models.Negocio{}, errs.ErrNoEncontrado
	}
	return n, nil
}

// Listar retorna todos los negocios. Si no hay ninguno retorna un slice
// vacío (no nil) — los tests dependen de esta garantía.
func (r *NegocioMemoria) Listar() []models.Negocio {
	resultado := make([]models.Negocio, 0, len(r.datos))
	for _, n := range r.datos {
		resultado = append(resultado, n)
	}
	return resultado
}

// Eliminar quita el negocio con el ID dado del repositorio.
//
// Validaciones:
//  1. ID no existe → ErrNoEncontrado
func (r *NegocioMemoria) Eliminar(id int) error {
	if _, existe := r.datos[id]; !existe {
		return errs.ErrNoEncontrado
	}
	delete(r.datos, id)
	return nil
}
