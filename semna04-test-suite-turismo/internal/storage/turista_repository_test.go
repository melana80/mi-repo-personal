package storage

import (
	"errors"
	"testing"

	"github.com/uleam/awii/turismo/internal/errs"
	"github.com/uleam/awii/turismo/internal/models"
)

func TestTuristaMemoria_Guardar_TableDriven(t *testing.T) {
	repo := NewTuristaMemoria()

	// Pre-condición: sembramos un turista para poder probar "ID duplicado".
	base := models.Turista{
		ID: 1, Nombre: "Alice", Nacionalidad: "Francia", IdiomaPreferido: "fr",
	}
	if err := repo.Guardar(base); err != nil {
		t.Fatalf("setup falló: %v", err)
	}

	casos := []struct {
		nombre      string
		entrada     models.Turista
		errEsperado error
	}{
		{
			nombre: "caso feliz - turista válido",   //   el caso 1 tiene los datos correctos, aqui no debe haber error.
			entrada: models.Turista{
				ID: 100, Nombre: "Bob", Nacionalidad: "Alemania", IdiomaPreferido: "de",
			},
			errEsperado: nil,
		},
		{
			nombre: "nombre vacío falla",  //Aqui esta el caso 2, donde tiene el nombre vacio el cual da error de validacion.
			entrada: models.Turista{
				ID: 101, Nombre: "", Nacionalidad: "Italia", IdiomaPreferido: "it",
			},
			errEsperado: errs.ErrDatosInvalidos,
		},
		{
			nombre: "nacionalidad vacía falla",
			entrada: models.Turista{
				ID: 102, Nombre: "Carla", Nacionalidad: "", IdiomaPreferido: "es",
			},
			errEsperado: errs.ErrDatosInvalidos,
		},
		{
			nombre: "idioma preferido no válido falla",
			entrada: models.Turista{
				ID: 103, Nombre: "Diego", Nacionalidad: "Ecuador", IdiomaPreferido: "ja",
			},
			errEsperado: errs.ErrDatosInvalidos,
		},
		{
			nombre: "ID duplicado falla",
			entrada: models.Turista{
				ID: 1, Nombre: "Eve", Nacionalidad: "EEUU", IdiomaPreferido: "en",
			},
			errEsperado: errs.ErrYaExiste,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			err := repo.Guardar(c.entrada)
			if !errors.Is(err, c.errEsperado) {
				t.Errorf("esperaba error=%v, obtuvo error=%v", c.errEsperado, err)  				  //Reporta los errores de manera continua
			}
		})
	}
}

func TestTuristaMemoria_BuscarPorID_TableDriven(t *testing.T) {
	casos := []struct {
		nombre      string
		idBuscar    int
		errEsperado error
	}{
		{
			nombre:      "ID existente retorna turista",
			idBuscar:    1,
			errEsperado: nil,
		},
		{
			nombre:      "ID negativo retorna ErrDatosInvalidos",
			idBuscar:    -1,
			errEsperado: errs.ErrDatosInvalidos,
		},
		{
			nombre:      "ID inexistente retorna ErrNoEncontrado",
			idBuscar:    999,
			errEsperado: errs.ErrNoEncontrado,
		},
	}

	for _, c := range casos { // Itera sobre el slice de casos y ejecuta el test para cada uno
		t.Run(c.nombre, func(t *testing.T) { // Define el nombre del test basado en el nombre del caso
			repo := NewTuristaMemoria() // Crea una instancia de la clase TuristaMemoria

			// guardamos un turista para el caso feliz.
			esperado := models.Turista{
				ID: 1, Nombre: "Alice", Nacionalidad: "Francia", IdiomaPreferido: "fr",
			}
			if err := repo.Guardar(esperado); err != nil { 
				t.Fatalf("setup falló: %v", err) 
			}

			got, err := repo.BuscarPorID(c.idBuscar) 

			if !errors.Is(err, c.errEsperado) { 
				t.Fatalf("esperaba error=%v, obtuvo error=%v", c.errEsperado, err)
			}
			if c.errEsperado == nil { 
				if got.ID != esperado.ID {
					t.Errorf("ID: esperaba %d, obtuvo %d", esperado.ID, got.ID)
				}
				if got.Nombre != esperado.Nombre {
					t.Errorf("Nombre: esperaba %q, obtuvo %q", esperado.Nombre, got.Nombre)
				}
			} else {
				if got != (models.Turista{}) {
					t.Errorf("esperaba Turista vacío en error, obtuvo %+v", got)
				}
			}
		})
	}
}

func TestTuristaMemoria_Listar(t *testing.T) {
	t.Run("repo vacío retorna slice vacío (no nil)", func(t *testing.T) {
		repo := NewTuristaMemoria()

		got := repo.Listar()

		if got == nil {
			t.Fatalf("Listar(): esperaba slice vacío (no nil), obtuvo nil")
		}
		if len(got) != 0 {
			t.Fatalf("Listar(): esperaba len=0, obtuvo len=%d", len(got))
		}
	})

	t.Run("repo con turistas retorna todos (sin depender del orden)", func(t *testing.T) {
		repo := NewTuristaMemoria()

		t1 := models.Turista{ID: 1, Nombre: "A", Nacionalidad: "X", IdiomaPreferido: "es"}
		t2 := models.Turista{ID: 2, Nombre: "B", Nacionalidad: "Y", IdiomaPreferido: "en"}
		if err := repo.Guardar(t1); err != nil {
			t.Fatalf("setup falló: %v", err)
		}
		if err := repo.Guardar(t2); err != nil {
			t.Fatalf("setup falló: %v", err)
		}

		got := repo.Listar()

		if len(got) != 2 {
			t.Fatalf("Listar(): esperaba len=2, obtuvo len=%d", len(got))
		}

		porID := make(map[int]models.Turista, len(got))
		for _, tt := range got {
			porID[tt.ID] = tt
		}
		if _, ok := porID[t1.ID]; !ok {
			t.Errorf("Listar(): esperaba turista con ID=%d", t1.ID)
		}
		if _, ok := porID[t2.ID]; !ok {
			t.Errorf("Listar(): esperaba turista con ID=%d", t2.ID)
		}
	})
}