// Programa demo del proyecto base del Día B.
//
// Crea los 3 repositorios, siembra datos realistas de Manta, registra
// algunos check-ins y los imprime. Sirve para verificar que el proyecto
// arranca antes de empezar a escribir tests.
//
// Ejecutar:
//   go run ./cmd/turismo
package main

import (
	"fmt"

	"github.com/uleam/awii/turismo/internal/models"
	"github.com/uleam/awii/turismo/internal/storage"
)

func main() {
	// 1. Crear los 3 repositorios. CheckInMemoria recibe los otros dos
	//    como dependencias para hacer validación cruzada.
	negocios := storage.NewNegocioMemoria()
	turistas := storage.NewTuristaMemoria()
	checkins := storage.NewCheckInMemoria(turistas, negocios)

	// 2. Sembrar negocios — datos reales de Manta.
	for _, n := range datosNegocios() {
		if err := negocios.Guardar(n); err != nil {
			fmt.Printf("ERROR sembrando negocio %d: %v\n", n.ID, err)
		}
	}

	// 3. Sembrar turistas — visitantes de cruceros.
	for _, t := range datosTuristas() {
		if err := turistas.Guardar(t); err != nil {
			fmt.Printf("ERROR sembrando turista %d: %v\n", t.ID, err)
		}
	}

	// 4. Registrar algunos check-ins de muestra.
	for _, c := range datosCheckIns() {
		if err := checkins.Guardar(c); err != nil {
			fmt.Printf("ERROR sembrando check-in %d: %v\n", c.ID, err)
		}
	}

	// 5. Imprimir el estado de los 3 repositorios.
	fmt.Println("== Negocios registrados ==")
	for _, n := range negocios.Listar() {
		fmt.Printf("  [%d] %-25s %-12s idiomas=%v\n",
			n.ID, n.Nombre, n.Tipo, n.IdiomasHablados)
	}

	fmt.Println("\n== Turistas registrados ==")
	for _, t := range turistas.Listar() {
		fmt.Printf("  [%d] %-20s %-12s idioma=%s\n",
			t.ID, t.Nombre, t.Nacionalidad, t.IdiomaPreferido)
	}

	fmt.Println("\n== Check-ins registrados ==")
	for _, c := range checkins.Listar() {
		fmt.Printf("  [%d] turista=%d negocio=%d fecha=%s ★%d\n",
			c.ID, c.TuristaID, c.NegocioID, c.Fecha, c.Calificacion)
	}

	// 6. Demostración de búsqueda por turista (relación).
	fmt.Println("\n== Check-ins del turista 1 (John Smith) ==")
	visitas, err := checkins.BuscarPorTurista(1)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	for _, c := range visitas {
		negocio, _ := negocios.BuscarPorID(c.NegocioID)
		fmt.Printf("  %s en %s (★%d)\n",
			c.Fecha, negocio.Nombre, c.Calificacion)
	}
}

func datosNegocios() []models.Negocio {
	return []models.Negocio{
		{
			ID: 1, Nombre: "Café del Mar", Tipo: "restaurante",
			Ciudad: "Manta", IdiomasHablados: []string{"es", "en"}, Activo: true,
		},
		{
			ID: 2, Nombre: "Hotel Oro Verde", Tipo: "hotel",
			Ciudad: "Manta", IdiomasHablados: []string{"es", "en", "fr"}, Activo: true,
		},
		{
			ID: 3, Nombre: "Tour Operadora Manta", Tipo: "tour",
			Ciudad: "Manta", IdiomasHablados: []string{"es", "en"}, Activo: true,
		},
		{
			ID: 4, Nombre: "Manabita Crafts", Tipo: "artesania",
			Ciudad: "Manta", IdiomasHablados: []string{"es"}, Activo: true,
		},
		{
			ID: 5, Nombre: "Trattoria Manabita", Tipo: "restaurante",
			Ciudad: "Manta", IdiomasHablados: []string{"es", "en", "it"}, Activo: true,
		},
	}
}

func datosTuristas() []models.Turista {
	return []models.Turista{
		{ID: 1, Nombre: "John Smith", Nacionalidad: "USA", IdiomaPreferido: "en"},
		{ID: 2, Nombre: "Marie Dubois", Nacionalidad: "Francia", IdiomaPreferido: "fr"},
		{ID: 3, Nombre: "Klaus Müller", Nacionalidad: "Alemania", IdiomaPreferido: "de"},
		{ID: 4, Nombre: "Mario Rossi", Nacionalidad: "Italia", IdiomaPreferido: "it"},
	}
}

func datosCheckIns() []models.CheckIn {
	return []models.CheckIn{
		{ID: 1, TuristaID: 1, NegocioID: 1, Fecha: "2026-04-10", Calificacion: 5},
		{ID: 2, TuristaID: 1, NegocioID: 3, Fecha: "2026-04-11", Calificacion: 4},
		{ID: 3, TuristaID: 2, NegocioID: 2, Fecha: "2026-04-10", Calificacion: 5},
		{ID: 4, TuristaID: 4, NegocioID: 5, Fecha: "2026-04-12", Calificacion: 5},
	}
}
