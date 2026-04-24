package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Define	el	struct	Cliente con	los	4	campos	especificados	en	la	Sección	2.1	de
// este	handout.	Respeta	los	tipos	exactos.	Consulta:	Snippets	Sección	1.1
type Cliente struct {
	ID      int     // Identificador	único	del	cliente
	Nombre  string  // Nombre	completo
	Carrera string  // Carrera	a	la	que	pertenece	(ej:	“TI”,	“Civil”)
	Saldo   float64 // Saldo	disponible	en	dólares
}

// Define	el	struct	Producto con	los	5	campos	especificados	en	la	Sección	2.2.
type Producto struct {
	ID        int     // Identificador	único	del	producto
	Nombre    string  // Nombre	del	producto
	Precio    float64 // Precio	unitario	en	dólares
	Stock     int     // Cantidad	disponible
	Categoria string  // "bebida",	"snack" o	"almuerzo"
}

// Define	el	struct	Pedido con	los	6	campos	especificados	en	la	Sección	2.3.
type Pedido struct {
	ID         int     // Identificador único del pedido
	ClienteID  int     // ID del cliente que hizo la compra (referencia)
	ProductoID int     // ID del producto comprado (referencia)
	Cantidad   int     // Unidades compradas
	Total      float64 // Precio total (Precio × Cantidad)
	Fecha      string  // Fecha del pedido (ej: "2026-04-10")
}

var clientes []Cliente
var productos []Producto
var pedidos []Pedido

func leerLinea(lector *bufio.Reader) string {
	line, err := lector.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(line)
}

func leerEntero(lector *bufio.Reader, prompt string) int {
	for {
		fmt.Print(prompt)
		s := leerLinea(lector)
		n, err := strconv.Atoi(s)
		if err == nil {
			return n
		}
		fmt.Println("Entrada no válida, intente de nuevo.")
	}
}

func leerFloat(lector *bufio.Reader, prompt string) float64 {
	for {
		fmt.Print(prompt)
		s := leerLinea(lector)
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return f
		}
		fmt.Println("Entrada no válida, intente de nuevo.")
	}
}

func mostrarMenu() {
	fmt.Println()
	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║     CAFETERÍA — MENÚ PRINCIPAL     ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Println("║  1. Listar clientes                ║")
	fmt.Println("║  2. Listar productos               ║")
	fmt.Println("║  3. Agregar cliente                ║")
	fmt.Println("║  4. Agregar producto               ║")
	fmt.Println("║  5. Registrar pedido               ║")
	fmt.Println("║  6. Ver pedidos de un cliente      ║")
	fmt.Println("║  0. Salir                          ║")
	fmt.Println("╚════════════════════════════════════╝")
	fmt.Print("Opción: ")
}

func siguienteIDCliente() int {
	max := 0
	for _, c := range clientes {
		if c.ID > max {
			max = c.ID
		}
	}
	return max + 1
}

func siguienteIDProducto() int {
	max := 0
	for _, p := range productos {
		if p.ID > max {
			max = p.ID
		}
	}
	return max + 1
}

// AgregarCliente añade un cliente al slice global. Si c.ID es 0, asigna el siguiente ID disponible.
func AgregarCliente(c Cliente) {
	if c.ID == 0 {
		c.ID = siguienteIDCliente()
	}
	clientes = append(clientes, c)
}

// BuscarClientePorID devuelve el cliente y true si existe.
func BuscarClientePorID(id int) (Cliente, bool) {
	for _, c := range clientes {
		if c.ID == id {
			return c, true
		}
	}
	return Cliente{}, false
}

// ListarClientes muestra los clientes en formato tabla.
func ListarClientes(clientes []Cliente) {
	fmt.Printf("%-5s %-15s %-15s %-10s\n", "ID", "Nombre", "Carrera", "Saldo")
	fmt.Println(strings.Repeat("-", 50))

	for _, cliente := range clientes {
		fmt.Printf("%-5d %-15s %-15s %-10.2f\n",
			cliente.ID, cliente.Nombre, cliente.Carrera, cliente.Saldo)
	}
}

// EliminarCliente quita el cliente con el ID indicado. Devuelve true si se eliminó.
func EliminarCliente(id int) bool {
	for i, c := range clientes {
		if c.ID == id {
			clientes = append(clientes[:i], clientes[i+1:]...)
			return true
		}
	}
	return false
}

// AgregarProducto añade un producto al slice global. Si p.ID es 0, asigna el siguiente ID disponible.
func AgregarProducto(p Producto) {
	if p.ID == 0 {
		p.ID = siguienteIDProducto()
	}
	productos = append(productos, p)
}

// BuscarProductoPorID devuelve el producto y true si existe.
func BuscarProductoPorID(id int) (Producto, bool) {
	for _, p := range productos {
		if p.ID == id {
			return p, true
		}
	}
	return Producto{}, false
}

// ListarProductos muestra los productos en formato tabla.
func ListarProductos(productos []Producto) {
	fmt.Printf("%-5s %-20s %-10s %-8s %-12s\n", "ID", "Nombre", "Precio", "Stock", "Categoria")
	fmt.Println(strings.Repeat("-", 60))

	for _, p := range productos {
		fmt.Printf("%-5d %-20s %-10.2f %-8d %-12s\n",
			p.ID, p.Nombre, p.Precio, p.Stock, p.Categoria)
	}
}

// EliminarProducto quita el producto con el ID indicado. Devuelve true si se eliminó.
func EliminarProducto(id int) bool {
	for i, p := range productos {
		if p.ID == id {
			productos = append(productos[:i], productos[i+1:]...)
			return true
		}
	}
	return false
}

// DescontarSaldo resta monto del saldo del cliente apuntado. Exige monto no negativo y saldo suficiente.
func DescontarSaldo(cliente *Cliente, monto float64) error {
	if monto < 0 {
		return fmt.Errorf("el monto no puede ser negativo")
	}
	if cliente.Saldo < monto {
		return fmt.Errorf("saldo insuficiente: tiene %.2f y se requieren %.2f", cliente.Saldo, monto)
	}
	cliente.Saldo -= monto
	return nil
}

// DescontarStock resta cantidad del stock del producto apuntado. Exige cantidad positiva y stock suficiente.
func DescontarStock(producto *Producto, cantidad int) error {
	if cantidad <= 0 {
		return fmt.Errorf("la cantidad debe ser mayor que cero")
	}
	if producto.Stock < cantidad {
		return fmt.Errorf("stock insuficiente: hay %d unidades y se pidieron %d", producto.Stock, cantidad)
	}
	producto.Stock -= cantidad
	return nil
}

// RegistrarPedido valida cliente y producto, descuenta stock y saldo, y registra el pedido. Ante fallo no deja estado inconsistente (revierte stock si falla el saldo).
func RegistrarPedido(
	clientes []Cliente,
	productos []Producto,
	pedidos []Pedido,
	clienteID int,
	productoID int,
	cantidad int,
	fecha string,
) ([]Pedido, error) {
	idxC := -1
	for i := range clientes {
		if clientes[i].ID == clienteID {
			idxC = i
			break
		}
	}
	if idxC == -1 {
		return pedidos, fmt.Errorf("cliente no encontrado")
	}

	idxP := -1
	for i := range productos {
		if productos[i].ID == productoID {
			idxP = i
			break
		}
	}
	if idxP == -1 {
		return pedidos, fmt.Errorf("producto no encontrado")
	}

	total := productos[idxP].Precio * float64(cantidad)

	if err := DescontarStock(&productos[idxP], cantidad); err != nil {
		return pedidos, err
	}
	if err := DescontarSaldo(&clientes[idxC], total); err != nil {
		productos[idxP].Stock += cantidad
		return pedidos, err
	}

	nuevo := Pedido{
		ID:         len(pedidos) + 1,
		ClienteID:  clienteID,
		ProductoID: productoID,
		Cantidad:   cantidad,
		Total:      total,
		Fecha:      fecha,
	}
	pedidos = append(pedidos, nuevo)
	return pedidos, nil
}

// PedidosDeCliente imprime un reporte cruzando pedidos, clientes y productos para un cliente.
func PedidosDeCliente(
	pedidos []Pedido,
	clientes []Cliente,
	productos []Producto,
	clienteID int,
) {
	var nombreCliente string
	encontrado := false
	for _, c := range clientes {
		if c.ID == clienteID {
			nombreCliente = c.Nombre
			encontrado = true
			break
		}
	}
	if !encontrado {
		fmt.Printf("Error: no existe un cliente con ID %d.\n", clienteID)
		return
	}

	fmt.Printf("\n--- Pedidos de %s (ID %d) ---\n", nombreCliente, clienteID)

	nombreProducto := func(productoID int) string {
		for _, p := range productos {
			if p.ID == productoID {
				return p.Nombre
			}
		}
		return "(producto no encontrado)"
	}

	var acumulado float64
	tienePedidos := false
	for _, ped := range pedidos {
		if ped.ClienteID != clienteID {
			continue
		}
		if !tienePedidos {
			fmt.Printf("%-6s %-22s %-10s %-10s %-12s\n", "Pedido", "Producto", "Cantidad", "Total", "Fecha")
			fmt.Println(strings.Repeat("-", 70))
			tienePedidos = true
		}
		acumulado += ped.Total
		fmt.Printf("%-6d %-22s %-10d %-10.2f %-12s\n",
			ped.ID, nombreProducto(ped.ProductoID), ped.Cantidad, ped.Total, ped.Fecha)
	}

	if !tienePedidos {
		fmt.Println("Este cliente no tiene pedidos registrados.")
		return
	}

	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("Total acumulado gastado por el cliente: %.2f\n", acumulado)
}

func main() {
	clientes = []Cliente{
		{ID: 1, Nombre: "Juan", Carrera: "Ingenieria", Saldo: 100.0},
		{ID: 2, Nombre: "Maria", Carrera: "Medicina", Saldo: 50.0},
		{ID: 3, Nombre: "Pedro", Carrera: "Arquitectura", Saldo: 200.0},
		{ID: 4, Nombre: "Luis", Carrera: "Derecho", Saldo: 150.0},
	}
	productos = []Producto{
		{ID: 1, Nombre: "Papas fritas", Precio: 1.0, Stock: 10, Categoria: "snack"},
		{ID: 2, Nombre: "Doritos", Precio: 2.0, Stock: 5, Categoria: "snack"},
		{ID: 3, Nombre: "Coca cola", Precio: 3.0, Stock: 8, Categoria: "bebida"},
		{ID: 4, Nombre: "Arroz con pollo", Precio: 4.0, Stock: 7, Categoria: "almuerzo"},
	}
	pedidos = []Pedido{
		{ID: 1, ClienteID: 1, ProductoID: 1, Cantidad: 2, Total: 2.0, Fecha: "2026-04-10"},
		{ID: 2, ClienteID: 2, ProductoID: 2, Cantidad: 3, Total: 6.0, Fecha: "2026-04-11"},
		{ID: 3, ClienteID: 3, ProductoID: 3, Cantidad: 4, Total: 12.0, Fecha: "2026-04-12"},
		{ID: 4, ClienteID: 4, ProductoID: 4, Cantidad: 5, Total: 20.0, Fecha: "2026-04-13"},
	}

	lector := bufio.NewReader(os.Stdin)

	for {
		mostrarMenu()
		opcion := leerEntero(lector, "")

		switch opcion {
		case 1:
			ListarClientes(clientes)
		case 2:
			ListarProductos(productos)
		case 3:
			fmt.Println("\n--- Nuevo cliente ---")
			fmt.Print("Nombre: ")
			nombre := leerLinea(lector)
			fmt.Print("Carrera: ")
			carrera := leerLinea(lector)
			saldo := leerFloat(lector, "Saldo: ")
			AgregarCliente(Cliente{Nombre: nombre, Carrera: carrera, Saldo: saldo})
			fmt.Println("Cliente agregado correctamente.")
		case 4:
			fmt.Println("\n--- Nuevo producto ---")
			fmt.Print("Nombre: ")
			nombre := leerLinea(lector)
			precio := leerFloat(lector, "Precio: ")
			stock := leerEntero(lector, "Stock: ")
			fmt.Print("Categoría (bebida / snack / almuerzo): ")
			categoria := leerLinea(lector)
			AgregarProducto(Producto{Nombre: nombre, Precio: precio, Stock: stock, Categoria: categoria})
			fmt.Println("Producto agregado correctamente.")
		case 5:
			fmt.Println("\n--- Registrar pedido ---")
			clienteID := leerEntero(lector, "ID del cliente: ")
			productoID := leerEntero(lector, "ID del producto: ")
			cantidad := leerEntero(lector, "Cantidad: ")
			fmt.Print("Fecha (ej. 2026-04-16): ")
			fecha := leerLinea(lector)
			var err error
			pedidos, err = RegistrarPedido(clientes, productos, pedidos, clienteID, productoID, cantidad, fecha)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Pedido registrado correctamente.")
			}
		case 6:
			id := leerEntero(lector, "\nID del cliente para el reporte: ")
			PedidosDeCliente(pedidos, clientes, productos, id)
		case 0:
			fmt.Println("¡Hasta luego!")
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}