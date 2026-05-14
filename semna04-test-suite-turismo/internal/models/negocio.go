// Package models define las entidades del dominio de Turismo.
package models

// IdiomasValidos define el set cerrado de idiomas que un Negocio puede
// declarar como hablados. Mantenerlo como variable exportada permite que
// los repositorios lo usen para validar y que los tests lo referencien sin
// duplicar la lista.
//
// Nota: en código de producción real esto vendría probablemente de
// configuración o base de datos. Para esta semana, la lista hardcodeada es
// suficiente.
var IdiomasValidos = map[string]bool{
	"es": true, // español
	"en": true, // inglés
	"fr": true, // francés
	"de": true, // alemán
	"it": true, // italiano
}

// TiposValidos define el set cerrado de tipos de negocio aceptados por la
// plataforma de turismo.
var TiposValidos = map[string]bool{
	"restaurante": true,
	"hotel":       true,
	"tour":        true,
	"artesania":   true,
}

// Negocio representa un establecimiento turístico registrado en la
// plataforma. Los campos IdiomasHablados y Tipo se validan contra los sets
// cerrados IdiomasValidos y TiposValidos.
type Negocio struct {
	ID              int
	Nombre          string
	Tipo            string
	Ciudad          string
	IdiomasHablados []string
	Activo          bool
}
