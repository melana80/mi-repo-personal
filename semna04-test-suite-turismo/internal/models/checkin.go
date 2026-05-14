package models

// CheckIn representa un registro de visita: un Turista visitó un Negocio
// en una fecha específica y dejó una calificación.
//
// La relación con Turista y Negocio es POR ID (no anidada). Esto sigue el
// patrón profesional usado en bases de datos relacionales: en memoria, el
// repositorio es responsable de validar que TuristaID y NegocioID existan
// antes de guardar.
//
// Nota sobre Fecha: usamos string en formato "YYYY-MM-DD" en lugar de
// time.Time porque time.Time se introduce en Semana 7. Por ahora basta con
// que el campo no esté vacío.
type CheckIn struct {
	ID           int
	TuristaID    int
	NegocioID    int
	Fecha        string // formato "YYYY-MM-DD"
	Calificacion int    // 1-5
}
