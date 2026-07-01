# Test Countdown - Plan de Implementación

## Objetivo

Implementar el MVP de `countdown`, una aplicación CLI/TUI en Go que recibe una duración, muestra una cuenta regresiva en terminal y finaliza con mensaje `Done!` más una animación breve de confeti.

## Fase 1 - Inicializar proyecto Go

Tareas:

- Crear `go.mod`.
- Definir el módulo del proyecto.
- Crear estructura inicial de archivos.

Estructura propuesta:

```text
.
├── go.mod
├── main.go
├── duration.go
├── duration_test.go
├── format.go
├── format_test.go
├── tui.go
├── REQUIREMENTS.md
├── DESIGN.md
└── PLAN.md
```

Notas:

- Mantener todo en `package main` para simplicidad del MVP.
- Más adelante, si crece, se puede separar en paquetes internos.

## Fase 2 - Parser de duración

Tareas:

- Implementar:

```go
func ParseDuration(input string) (time.Duration, error)
```

Reglas:

- Aceptar unidades `h`, `m`, `s`.
- Permitir combinaciones como:
  - `10s`
  - `5m`
  - `1h`
  - `1h30m`
  - `1h30m10s`
  - `30m10s`
- Exigir orden descendente: `h`, luego `m`, luego `s`.
- Rechazar formatos inválidos como:
  - vacío
  - `abc`
  - `10x`
  - `1m2h`
  - `1h1h`
  - `1.5h`
  - `-10s`

Resultado esperado:

- Devuelve `time.Duration` válido.
- Devuelve error claro si el formato es inválido.

## Fase 3 - Tests del parser

Tareas:

- Crear tests unitarios para casos válidos e inválidos.
- Verificar que las duraciones parseadas coincidan con el valor esperado.

Ejemplos:

```go
ParseDuration("1h30m10s") == time.Hour + 30*time.Minute + 10*time.Second
ParseDuration("10s") == 10*time.Second
```

## Fase 4 - Formateador de tiempo

Tareas:

- Implementar:

```go
func FormatRemaining(d time.Duration) string
```

Reglas:

- Si la duración tiene 1 hora o más, usar `HH:MM:SS`.
- Si la duración es menor a 1 hora, usar `MM:SS`.
- No mostrar valores negativos; si llega negativo, mostrar cero.

Ejemplos:

```text
10s        -> 00:10
5m         -> 05:00
1h         -> 01:00:00
1h30m10s   -> 01:30:10
```

## Fase 5 - Tests del formateador

Tareas:

- Crear tests para distintos formatos.
- Incluir caso de duración negativa.

## Fase 6 - CLI básica y ayuda

Tareas:

- Implementar lectura de argumentos en `main.go`.
- Soportar:

```bash
countdown --help
countdown -h
countdown <duration>
```

- Mostrar ayuda si se usa `--help` o `-h`.
- Mostrar error si falta duración.
- Mostrar error si hay demasiados argumentos.
- Mostrar error si la duración es inválida.

Ayuda propuesta:

```text
Usage:
  countdown <duration>

Examples:
  countdown 10s
  countdown 5m
  countdown 1h30m
  countdown 1h30m10s

Keys:
  q, Esc, Ctrl+C   Quit
```

## Fase 7 - Agregar dependencias TUI

Tareas:

- Agregar Bubble Tea.
- Agregar Lip Gloss.

Comandos esperados:

```bash
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
```

Uso:

- Bubble Tea para modelo, ticks y teclado.
- Lip Gloss para estilos visuales básicos.

## Fase 8 - Modelo TUI del countdown

Tareas:

- Crear `tui.go`.
- Implementar modelo Bubble Tea.
- Manejar estado:

```go
type model struct {
    total         time.Duration
    remaining     time.Duration
    startedAt     time.Time
    done          bool
    confettiFrame int
}
```

- Crear tick periódico, idealmente cada 200ms o 1s.
- Calcular el tiempo restante según `startedAt` para evitar drift.
- Al llegar a cero:
  - Marcar `done = true`.
  - Mantener pantalla final.
  - Seguir animando confeti por frames.

## Fase 9 - Manejo de teclado

Tareas:

- Permitir salir con:
  - `q`
  - `Esc`
  - `Ctrl+C`

Comportamiento:

- Si el usuario sale antes de terminar, no mostrar `Done!`.
- Si ya terminó, salir desde la pantalla final.

## Fase 10 - Render visual

Tareas:

- Pantalla durante countdown:

```text
      COUNTDOWN

       01:30

   Press q to quit
```

- Pantalla final:

```text
       Done!

   *  .  *  .  *
     .  *  .  *
   *  .  *  .  *

   Press q to quit
```

- Usar Lip Gloss para:
  - título
  - tiempo restante
  - texto de ayuda
  - mensaje final

## Fase 11 - Animación de confeti

Tareas:

- Definir una lista de frames estáticos.
- Alternar frames con cada tick cuando `done = true`.
- Mantener animación activa mientras la pantalla final esté abierta.

Ejemplo:

```go
var confettiFrames = []string{
    "*  .  +  .  *",
    ".  +  *  +  .",
    "+  *  .  *  +",
}
```

## Fase 12 - Pruebas manuales

Comandos a probar:

```bash
go test ./...
go run . --help
go run . -h
go run .
go run . abc
go run . 10s
go run . 1m10s
go run . 1h30m10s
```

Verificar:

- Countdown inicia correctamente.
- Tiempo restante se actualiza.
- Teclas de salida funcionan.
- Al terminar aparece `Done!`.
- Confeti se anima brevemente o continuamente en pantalla final.
- La app queda esperando salida manual.

## Fase 13 - Build local

Tareas:

- Compilar binario:

```bash
go build -o countdown .
```

- Probar:

```bash
./countdown 10s
```

Opcional:

- Agregar instrucciones de instalación en `README.md` después del MVP.

## Criterios de aceptación del MVP

El MVP se considera completo si:

- `countdown 10s` ejecuta una cuenta regresiva funcional.
- `countdown 1h30m10s` es aceptado correctamente.
- `countdown --help` y `countdown -h` muestran ayuda.
- Formatos inválidos muestran error claro.
- La UI se actualiza mientras corre.
- `q`, `Esc` y `Ctrl+C` salen de la app.
- Al llegar a cero se muestra `Done!` con confeti animado.
- La app no se cierra automáticamente al terminar.
- Existen tests para parser y formateador.
- `go test ./...` pasa correctamente.
