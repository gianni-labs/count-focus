# Mejoras Futuras

Este documento lista mejoras posibles después del MVP, con una breve explicación, prioridad sugerida y estado.

## Resumen de estado

| # | Mejora | Prioridad | Estado |
|---|--------|-----------|--------|
| 1 | Layout responsive / ajuste a pantalla | Alta | Implementado |
| 2 | Tiempo con letras grandes estilo letrero | Alta | Implementado base |
| 3 | Barra de progreso | Media | Pendiente |
| 4 | Pausar y reanudar | Media | Pendiente |
| 5 | Título personalizado | Media | Pendiente |
| 6 | Presets | Media | Pendiente |
| 7 | Modo count-up | Baja | Pendiente |
| 8 | Sonido o notificación al terminar | Baja | Pendiente |
| 9 | Temas visuales | Baja | Pendiente |
| 10 | Más efectos finales | Baja | Pendiente |

Estados posibles:

- `Pendiente`: todavía no implementado.
- `En progreso`: implementación iniciada pero no cerrada.
- `Implementado base`: existe una primera versión funcional, pero puede refinarse.
- `Implementado`: completo para el alcance definido.

## Prioridad Alta

### 1. Layout responsive / ajuste a pantalla

Estado: implementado en el MVP inicial.

La TUI se adapta automáticamente al tamaño actual de la terminal.

Incluye:

- Centrado vertical y horizontal.
- Manejo de cambios de tamaño de ventana.
- Uso de `tea.WindowSizeMsg` para guardar `width` y `height`.
- Uso de `lipgloss.Place(width, height, ...)` para renderizar el contenido.

Mejoras futuras posibles sobre este punto:

- Ajustar tamaños de fuente/arte según ancho disponible.
- Ocultar ayuda o confeti en terminales extremadamente pequeñas.
- Agregar layouts alternativos para pantallas muy anchas.

### 2. Tiempo con letras grandes estilo letrero

Estado: implementado base.

Mostrar el tiempo restante con números grandes, tipo cartel o letrero de terminal, cuando el tamaño de la terminal lo permite. Si la terminal es pequeña, vuelve automáticamente al formato simple `MM:SS` / `HH:MM:SS`.

Esta primera versión usa una fuente tipo banner basada en caracteres de bloque y esquinas. Más adelante se podrá convertir en un set configurable de fuentes visuales.

Objetivo:

- Que el countdown tenga una identidad visual más fuerte.
- Mejorar legibilidad desde lejos.
- Lograr una sensación más moderna y divertida.

Ejemplo aproximado:

```text
  ██████╗  ██████╗       ██╗ ██████╗ 
 ██╔═████╗██╔═████╗     ███║██╔═████╗
 ██║██╔██║██║██╔██║     ╚██║██║██╔██║
 ████╔╝██║████╔╝██║      ██║████╔╝██║
 ╚██████╔╝╚██████╔╝      ██║╚██████╔╝
  ╚═════╝  ╚═════╝       ╚═╝ ╚═════╝ 
```

Notas técnicas:

- Se implementó una fuente propia tipo banner para dígitos `0-9` y `:`.
- El render grande se activa solo si hay suficiente ancho y alto.
- Más adelante se podrá agregar un set configurable de fuentes.

## Prioridad Media

### 3. Barra de progreso

Agregar una barra que muestre visualmente cuánto falta y cuánto ya transcurrió.

Objetivo:

- Dar contexto visual además del tiempo numérico.
- Hacer más clara la progresión del countdown.

Ejemplo:

```text
[████████████░░░░░░░░] 60%
```

### 4. Pausar y reanudar

Permitir pausar y continuar el countdown.

Objetivo:

- Mejorar usabilidad.
- Útil para sesiones largas.

Tecla sugerida:

- `Space` para pausar/reanudar.

### 5. Título personalizado

Permitir pasar un título al countdown.

Ejemplo:

```bash
countdown 25m --title "Focus"
```

Objetivo:

- Dar contexto al countdown.
- Preparar el camino para presets como pomodoro.

### 6. Presets

Agregar presets útiles.

Ejemplos:

```bash
countdown pomodoro
countdown short-break
countdown long-break
```

O alternativamente:

```bash
countdown --preset pomodoro
```

Objetivo:

- Hacer más cómodo el uso repetitivo.

## Prioridad Baja

### 7. Modo count-up

Agregar un modo que cuente hacia arriba desde cero.

Ejemplo:

```bash
countdown --up
```

Objetivo:

- Usar la app como cronómetro simple.

### 8. Sonido o notificación al terminar

Emitir un sonido o mostrar una notificación del sistema cuando termine.

Objetivo:

- Avisar al usuario aunque no esté mirando la terminal.

Notas:

- Puede ser dependiente del sistema operativo.
- Conviene dejarlo para después de consolidar el MVP visual.

### 9. Temas visuales

Permitir elegir estilos de color.

Ejemplo:

```bash
countdown 10m --theme neon
```

Objetivo:

- Personalizar la apariencia.
- Separar estilos de la lógica.

### 10. Más efectos finales

Agregar efectos alternativos al confeti.

Ejemplos:

- Fuegos artificiales.
- Flash de pantalla.
- Animación de partículas.
- Mensajes aleatorios de finalización.

Objetivo:

- Hacer más entretenida la experiencia final.

## Orden recomendado de implementación

1. Layout responsive / ajuste a pantalla. Estado: Implementado.
2. Letras grandes estilo letrero para el tiempo. Estado: Implementado base.
3. Barra de progreso. Estado: Pendiente.
4. Pausar y reanudar.
5. Título personalizado.
6. Presets.
7. Sonido/notificación.
8. Temas visuales.
9. Count-up.
10. Más efectos finales.

## Próximo paso sugerido

La siguiente mejora debería ser combinar:

1. Refinar la fuente grande del tiempo para que se vea más como un letrero.
2. Agregar barra de progreso responsive.
3. Ocultar o simplificar elementos secundarios en terminales pequeñas.
