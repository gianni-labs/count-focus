# Count Focus

Timer de foco para la terminal.

![Demo de Count Focus](assets/count-focus-demo.gif)

## Instalación

Con Homebrew:

```bash
brew tap gianni-labs/tap
brew trust --formula gianni-labs/tap/count-focus
brew install count-focus
```

Esto instala el comando:

```bash
count-focus
```

## Uso

```bash
count-focus <duration>
```

Ejemplos:

```bash
count-focus 10s
count-focus 5m
count-focus 1h
count-focus 1h30m
count-focus 1h30m10s
```

### Título

Podés ponerle un título al timer para saber qué estás haciendo:

```bash
count-focus 25m --title "Escribir informe"
count-focus 1h -t "Deep work"
```

### Hasta una hora específica

En vez de una duración, podés apuntar a una hora del reloj (formato 24h). Cuenta regresiva hasta esa hora de hoy:

```bash
count-focus --until 15:00      # hasta las 15:00
count-focus -u 15:30:30        # con segundos
```

Si la hora ya pasó hoy, muestra un error.

### Teclas

Mientras corre el timer:

- `Space` — pausar / reanudar
- `q`, `Esc`, `Ctrl+C` — salir

## Presets

En vez de una duración, podés usar un preset con nombre:

```bash
count-focus --preset pomodoro     # 25m
count-focus -p short-break        # 5m
count-focus -p long-break         # 15m
```

> Nota: `pomodoro` acá es solo un atajo a una duración de 25 minutos, no el ciclo completo de trabajo/descanso.

### Presets personalizados

Podés definir tus propios presets (o cambiar los que vienen por defecto) creando este archivo:

```
~/.config/count-focus/presets.conf
```

Con una línea por preset, en formato `nombre = duración`:

```conf
pomodoro = 30m
deep-work = 90m
review = 45m
```

Los presets del archivo sobreescriben o extienden los built-in (`pomodoro`, `short-break`, `long-break`). Hay un ejemplo completo en [`examples/presets.conf`](examples/presets.conf).

## Versión

```bash
count-focus --version
```

## Licencia

MIT

## Actualizar

```bash
brew update
brew upgrade count-focus
```

## Desinstalar

```bash
brew uninstall count-focus
```
