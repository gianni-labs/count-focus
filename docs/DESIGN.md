# Test Countdown - Diseño Técnico

## Objetivo

Construir una aplicación CLI/TUI llamada `countdown` que ejecute una cuenta regresiva desde la terminal, recibiendo una duración por argumento y mostrando el tiempo restante con una interfaz simple pero más moderna que un `print` tradicional.

El foco del MVP es que sea funcional, simple y fácil de entender.

## Comando esperado

```bash
countdown 10s
countdown 5m
countdown 1h
countdown 1h30m10s
countdown --help
countdown -h
```

## Alcance del MVP

Incluido:

- Ejecutar countdown desde terminal.
- Aceptar duración como argumento.
- Soportar combinaciones con horas, minutos y segundos:
  - `10s`
  - `5m`
  - `1h`
  - `1h30m`
  - `1h30m10s`
  - `2m15s`
- Mostrar tiempo restante de forma clara.
- Actualizar la interfaz mientras corre.
- Salir con:
  - `q`
  - `Esc`
  - `Ctrl+C`
- Mostrar ayuda básica con `--help` o `-h`.
- Al finalizar:
  - Mostrar mensaje `Done!` por defecto.
  - Mostrar un efecto visual básico tipo confeti en terminal.
- Agregar tests para lógica reusable.

No incluido en el MVP:

- Pausar/reanudar.
- Count-up.
- Títulos personalizados.
- Presets como pomodoro.
- Sonidos o notificaciones del sistema.
- Configuración avanzada de efectos.

## Tecnología propuesta

### Lenguaje

Go.

Razones:

- Muy bueno para CLIs pequeñas.
- Permite generar un binario simple llamado `countdown`.
- Tiene buen soporte estándar para tiempo, señales y tests.
- Mantiene el proyecto simple.

### Librería TUI

Se propone usar Bubble Tea:

- Permite manejar eventos de teclado con claridad.
- Facilita actualizar la pantalla periódicamente.
- Da una base moderna sin complejidad excesiva.
- Encaja bien con un modelo simple de countdown.

Opcionalmente se puede usar Lip Gloss para estilos visuales básicos, aunque el MVP puede mantenerse sobrio.

## Diseño de comportamiento

### Flujo normal

1. Usuario ejecuta:

   ```bash
   countdown 1m30s
   ```

2. La app parsea la duración.
3. Si la duración es válida, inicia la TUI.
4. La pantalla muestra el tiempo restante.
5. Cada segundo se actualiza el tiempo.
6. Al llegar a cero:
   - Se muestra `Done!`.
   - Se muestra una animación breve de confeti.
   - La app queda en pantalla final esperando que el usuario presione una tecla de salida.

### Salida manual

Durante la ejecución, el usuario puede salir con:

- `q`
- `Esc`
- `Ctrl+C`

En ese caso la app termina sin mostrar `Done!`.

### Ayuda

Si el usuario ejecuta:

```bash
countdown --help
countdown -h
```

Debe mostrar una ayuda simple:

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

### Errores

Si no se pasa duración:

```bash
countdown
```

Debe mostrar error y ayuda resumida.

Si la duración es inválida:

```bash
countdown abc
countdown 10x
countdown 1m2h
```

Debe mostrar un mensaje claro, por ejemplo:

```text
Invalid duration: abc
Expected format like 10s, 5m, 1h30m, or 1h30m10s.
```

## Formato de duración

El formato aceptado será una secuencia de bloques numéricos positivos seguidos por unidad:

```text
<number>h<number>m<number>s
```

Reglas propuestas:

- Unidades soportadas: `h`, `m`, `s`.
- Cada unidad puede aparecer como máximo una vez.
- El orden debe ser descendente: horas, luego minutos, luego segundos.
- Al menos una unidad debe estar presente.
- No se aceptan números negativos.
- No se aceptan espacios internos.
- No se aceptan valores vacíos como `h`, `m`, `s`.

Válidos:

```text
10s
5m
1h
1h30m
1h30m10s
30m10s
```

Inválidos:

```text
abc
10x
1m2h
1h1h
1.5h
-10s
```

Nota: aunque Go tiene `time.ParseDuration`, aceptar directamente ese parser permitiría formatos como `1.5h` o `100ms`. Para mantener el MVP alineado al requerimiento, conviene implementar un parser pequeño propio para `h`, `m`, `s`.

## Diseño visual inicial

Pantalla durante countdown:

```text

      COUNTDOWN

       01:30

   Press q to quit

```

Para duraciones con horas:

```text

      COUNTDOWN

      01:30:10

   Press q to quit

```

Pantalla final:

```text

       Done!

   *  .  *  .  *
     .  *  .  *
   *  .  *  .  *

```

El confeti del MVP será una animación breve usando caracteres como `*`, `.`, `+`, `•`, `✦`. No necesita ser físicamente realista ni compleja; basta con alternar algunos frames simples para dar sensación de celebración.

## Componentes internos propuestos

### `main.go`

Responsabilidades:

- Leer argumentos CLI.
- Mostrar ayuda.
- Validar input básico.
- Ejecutar la aplicación TUI.
- Manejar errores de alto nivel.

### Parser de duración

Responsabilidades:

- Convertir strings como `1h30m10s` a `time.Duration`.
- Rechazar formatos inválidos.
- Ser testeable sin depender de la TUI.

Funciones candidatas:

```go
func ParseDuration(input string) (time.Duration, error)
```

### Formateador de tiempo

Responsabilidades:

- Convertir una duración restante a texto legible.
- Mostrar `MM:SS` cuando no hay horas.
- Mostrar `HH:MM:SS` cuando hay horas.

Funciones candidatas:

```go
func FormatRemaining(d time.Duration) string
```

### Modelo TUI

Responsabilidades:

- Guardar duración inicial.
- Guardar instante de inicio o tiempo restante.
- Recibir ticks periódicos.
- Manejar teclas de salida.
- Detectar finalización.
- Renderizar la vista.

Estado candidato:

```go
type model struct {
    total         time.Duration
    remaining     time.Duration
    done          bool
    quitting      bool
    confettiFrame int
}
```

Cuando `done` sea `true`, el modelo seguirá aceptando ticks durante algunos frames para animar el confeti, pero no saldrá automáticamente. La salida seguirá dependiendo de `q`, `Esc` o `Ctrl+C`.

## Tests recomendados

### Parser

Casos válidos:

- `10s`
- `5m`
- `1h`
- `1h30m`
- `1h30m10s`
- `30m10s`

Casos inválidos:

- `` vacío
- `abc`
- `10x`
- `1m2h`
- `1h1h`
- `1.5h`
- `-10s`

### Formateo

- `10s` -> `00:10`
- `5m` -> `05:00`
- `1h` -> `01:00:00`
- `1h30m10s` -> `01:30:10`

## Decisiones confirmadas

- Al terminar, la app muestra `Done!` y queda esperando una tecla de salida.
- El confeti será una animación breve.
- Se usará Bubble Tea para la TUI. También se puede incluir Lip Gloss si aporta claridad visual sin complicar demasiado el código.

## Decisión técnica recomendada

Usar Bubble Tea + Lip Gloss:

- Bubble Tea: modelo, ticks y manejo de teclado.
- Lip Gloss: estilos simples para título, tiempo restante, ayuda y pantalla final.

La dependencia extra de Lip Gloss es razonable porque pertenece al mismo ecosistema Charm y mejora bastante la presentación sin introducir mucha complejidad.
