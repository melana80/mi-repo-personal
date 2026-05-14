package models

// Turista representa a un visitante extranjero registrado en la plataforma.
//
// El campo IdiomaPreferido se valida contra IdiomasValidos (definido en
// negocio.go) — la plataforma solo soporta los idiomas listados ahí.
type Turista struct {
	ID              int
	Nombre          string
	Nacionalidad    string
	IdiomaPreferido string
}
