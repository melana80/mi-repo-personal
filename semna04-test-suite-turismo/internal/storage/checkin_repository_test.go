// table-driven con los 7 casos: 
//1 CheckIn válido (turista y negocio existen) nil 
//2 Fecha vacía ErrDatosInvalidos
// 3 Calificación = 0 (fuera de 1-5) ErrDatosInvalidos
// 4 Calificación = 6 (fuera de 1-5) ErrDatosInvalidos 
// 5 TuristaID que NO existe en el repo de turistas ErrNoEncontrado 
// 6 NegocioID que NO existe en el repo de negocios ErrNoEncontrado 
// 7 ID de CheckIn ya usado ErrYaExiste 

package storage

import (
	"errors"
	"testing"

	"github.com/uleam/awii/turismo/internal/errs"
	"github.com/uleam/awii/turismo/internal/models"
)

// setupRepos crea los tres repositorios enlazados y siembra un turista y un
// negocio válidos (IDs 1) para que las pruebas de CheckIn puedan validar
// referencias cruzadas.
func setupRepos(t *testing.T) (TuristaRepository, NegocioRepository, *CheckInMemoria) {
	t.Helper()
	turistas := NewTuristaMemoria() // repositorio de turistas
	negocios := NewNegocioMemoria()
	checkins := NewCheckInMemoria(turistas, negocios)
	if err := turistas.Guardar(models.Turista{ // sembrar turista
		ID: 1, Nombre: "John", Nacionalidad: "USA", IdiomaPreferido: "en",
	}); err != nil {
		t.Fatalf("setup sembrar turista: %v", err)
	}
	if err := negocios.Guardar(models.Negocio{ 
		ID: 1, Nombre: "Café", Tipo: "restaurante",
		Ciudad: "Manta", IdiomasHablados: []string{"es", "en"}, Activo: true,
	}); err != nil {
		t.Fatalf("setup sembrar negocio: %v", err)
	}
	return turistas, negocios, checkins 
}

// TestCheckInMemoria_Guardar prueba la función Guardar de CheckInMemoria.
//agregamos 7 casos

func TestCheckInMemoria_Guardar(t *testing.T) { 
	casos := []struct {
		nombre      string
		preparar    func(t *testing.T, checkins *CheckInMemoria)
		entrada     models.CheckIn
		errEsperado error
	}{
		{
			nombre: "check-in válido (turista y negocio existen)",
			entrada: models.CheckIn{
				ID: 10, TuristaID: 1, NegocioID: 1,
				Fecha: "2026-04-10", Calificacion: 5,
			},
			errEsperado: nil,
		},
		{
			nombre: "fecha vacía → ErrDatosInvalidos",
			entrada: models.CheckIn{
				ID: 11, TuristaID: 1, NegocioID: 1,
				Fecha: "", Calificacion: 3,
			},
			errEsperado: errs.ErrDatosInvalidos,
		},
		{
			// Calificación = 0 (fuera de 1-5)
			nombre: "calificación 0 → ErrDatosInvalidos", 
			entrada: models.CheckIn{
				ID: 12, TuristaID: 1, NegocioID: 1,
				Fecha: "2026-04-10", Calificacion: 0,
			},
			errEsperado: errs.ErrDatosInvalidos,
		},
		{
			// Calificación = 6 (fuera de 1-5)
			nombre: "calificación 6 → ErrDatosInvalidos",
			entrada: models.CheckIn{
				ID: 13, TuristaID: 1, NegocioID: 1,
				Fecha: "2026-04-10", Calificacion: 6,
			},
			errEsperado: errs.ErrDatosInvalidos,
		},
		{
			//Turista existente
			nombre: "turista inexistente → ErrNoEncontrado",
			entrada: models.CheckIn{
				ID: 14, TuristaID: 999, NegocioID: 1,
				Fecha: "2026-04-10", Calificacion: 4,
			},
			errEsperado: errs.ErrNoEncontrado,
		},
		{ 
			//Negocio existente
			nombre: "negocio inexistente → ErrNoEncontrado",
			entrada: models.CheckIn{
				ID: 15, TuristaID: 1, NegocioID: 999,
				Fecha: "2026-04-10", Calificacion: 4,
			},
			errEsperado: errs.ErrNoEncontrado,
		},
		{
			//ID de check-in ya usado
			nombre: "ID de check-in ya usado → ErrYaExiste",
			preparar: func(t *testing.T, checkins *CheckInMemoria) {
				t.Helper()
				prev := models.CheckIn{
					ID: 20, TuristaID: 1, NegocioID: 1,
					Fecha: "2026-04-01", Calificacion: 2,
				}
				if err := checkins.Guardar(prev); err != nil {
					t.Fatalf("preparar Guardar: %v", err) 
				}
			},
			entrada: models.CheckIn{
				ID: 20, TuristaID: 1, NegocioID: 1,
				Fecha: "2026-04-11", Calificacion: 5,
			},
			errEsperado: errs.ErrYaExiste,
		},
	}
	// Casos 
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) { 
			_, _, checkins := setupRepos(t)
			if c.preparar != nil {
				c.preparar(t, checkins)
			}
			err := checkins.Guardar(c.entrada) // Guardar
			if !errors.Is(err, c.errEsperado) {// si no son iguales
				t.Fatalf("Guardar(): esperaba error=%v, obtuvo error=%v", c.errEsperado, err) 
			}
		})
	}
}

// BuscarPorTurista retorna todos los check-ins hechos por un turista. el turista 
// debe existir en el repositorio de turistas
//se retorna un error si el turistaID es negativo
//se retorna un error si el turistaID no existe en el repositorio de turistas


func TestCheckInMemoria_BuscarPorTurista(t *testing.T) {
	casos := []struct {
		nombre        string
		preparar      func(t *testing.T, checkins *CheckInMemoria)
		turistaID     int
		errEsperado   error
		esperadosLen  int
		esperadosIDs  map[int]struct{} // opcional: IDs esperados cuando len > 0
	}{
		{
			nombre: "turista con varios check-ins",
			preparar: func(t *testing.T, checkins *CheckInMemoria) { // preparar
				t.Helper()
				a := models.CheckIn{ID: 1, TuristaID: 1, NegocioID: 1, Fecha: "2026-04-10", Calificacion: 5} 
				b := models.CheckIn{ID: 2, TuristaID: 1, NegocioID: 1, Fecha: "2026-04-11", Calificacion: 4}
				if err := checkins.Guardar(a); err != nil {
					t.Fatalf("Guardar: %v", err)
				}
				if err := checkins.Guardar(b); err != nil {
					t.Fatalf("Guardar: %v", err)
				}
			},
			turistaID:    1,
			errEsperado:  nil,
			esperadosLen: 2,
			esperadosIDs: map[int]struct{}{1: {}, 2: {}},
		},
		{
			nombre:       "turista existe pero sin check-ins → slice vacío, sin error",
			preparar:     nil,
			turistaID:    1,
			errEsperado:  nil,
			esperadosLen: 0,
		},
		{
			nombre:       "turista inexistente → ErrNoEncontrado",
			preparar:     nil,
			turistaID:    999,
			errEsperado:  errs.ErrNoEncontrado,
			esperadosLen: -1,
		},
		{
			nombre:       "ID turista negativo → ErrDatosInvalidos",
			preparar:     nil,
			turistaID:    -1,
			errEsperado:  errs.ErrDatosInvalidos,
			esperadosLen: -1,
		},
	}
	
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, _, checkins := setupRepos(t)
			if c.preparar != nil {
				c.preparar(t, checkins)
			}
			visitas, err := checkins.BuscarPorTurista(c.turistaID) 
			if !errors.Is(err, c.errEsperado) {
				t.Fatalf("BuscarPorTurista(%d): esperaba error=%v, obtuvo error=%v",
					c.turistaID, c.errEsperado, err)
			}
			if c.errEsperado != nil {
				if visitas != nil {
					t.Errorf("con error esperaba visitas nil, obtuvo len=%d", len(visitas))
				}
				return
			}
			if visitas == nil {
				t.Fatalf("sin error esperaba slice no nil (puede estar vacío)")
			}
			if len(visitas) != c.esperadosLen {
				t.Fatalf("esperaba %d visitas, obtuvo %d", c.esperadosLen, len(visitas))
			}
			if c.esperadosIDs != nil {
				for _, v := range visitas {
					if _, ok := c.esperadosIDs[v.ID]; !ok {
						t.Errorf("visitas inesperada con ID=%d", v.ID)
					}
				}
			}
		})
	}
}

// Listar retorna todos los check-ins. Si no hay ninguno retorna un slice
//los tests dependen de esta garantía
//se retorna un slice vacío (no nil)
//se retorna un error si el turistaID es negativo
//se retorna un error si el turistaID no existe en el repositorio de turistas
//retorna un error si el negocioID es negativo
//se retorna un error si el negocioID no existe en el repositorio de negocios


func TestCheckInMemoria_Listar(t *testing.T) {
	t.Run("sin check-ins → slice vacío (no nil)", func(t *testing.T) {
		_, _, checkins := setupRepos(t)
		got := checkins.Listar()
		if got == nil {
			t.Fatal("Listar(): esperaba slice vacío (no nil), obtuvo nil")
		}
		if len(got) != 0 {
			t.Fatalf("Listar(): esperaba len=0, obtuvo len=%d", len(got))
		}
	})

	t.Run("varios check-ins listados (sin depender del orden)", func(t *testing.T) {
		_, _, checkins := setupRepos(t)
		c1 := models.CheckIn{ID: 1, TuristaID: 1, NegocioID: 1, Fecha: "2026-04-10", Calificacion: 5}
		c2 := models.CheckIn{ID: 2, TuristaID: 1, NegocioID: 1, Fecha: "2026-04-11", Calificacion: 3}
		if err := checkins.Guardar(c1); err != nil {
			t.Fatalf("Guardar: %v", err)
		}
		if err := checkins.Guardar(c2); err != nil {
			t.Fatalf("Guardar: %v", err)
		}
		got := checkins.Listar()
		if len(got) != 2 {
			t.Fatalf("Listar(): esperaba len=2, obtuvo len=%d", len(got))
		}
		porID := make(map[int]models.CheckIn, len(got))
		for _, c := range got {
			porID[c.ID] = c
		}
		if _, ok := porID[c1.ID]; !ok {
			t.Errorf("esperaba check-in ID=%d", c1.ID)
		}
		if _, ok := porID[c2.ID]; !ok {
			t.Errorf("esperaba check-in ID=%d", c2.ID)
		}
	})
}
//listar todos los checkins
