# Countdown

Una pequeña aplicación de countdown para terminal, escrita en Go, con una TUI simple usando Bubble Tea y Lip Gloss.

## Uso

```bash
countdown <duration>
```

Ejemplos:

```bash
countdown 10s
countdown 5m
countdown 1h
countdown 1h30m
countdown 1h30m10s
```

Ayuda:

```bash
countdown --help
countdown -h
```

## Formato de duración

El formato acepta combinaciones de horas, minutos y segundos:

- `h` para horas
- `m` para minutos
- `s` para segundos

Ejemplos válidos:

```text
10s
5m
1h
1h30m
30m10s
1h30m10s
```

El orden debe ser siempre: horas, minutos, segundos.

Ejemplos inválidos:

```text
abc
10x
1m2h
1h1h
1.5h
-10s
```

## Teclas

Durante la ejecución puedes salir con:

- `q`
- `Esc`
- `Ctrl+C`

Cuando el countdown llega a cero, muestra `Done!` con una animación simple de confeti y queda esperando una tecla de salida.

## Desarrollo

Instalar dependencias:

```bash
go mod tidy
```

Ejecutar tests:

```bash
go test ./...
```

Ejecutar localmente:

```bash
go run . 10s
```

Compilar binario local:

```bash
go build -o countdown .
./countdown 10s
```

Instalar el comando en tu `GOPATH/bin`:

```bash
go install .
```

Luego, si `$(go env GOPATH)/bin` está en tu `PATH`, puedes ejecutar:

```bash
countdown 10s
```

## Documentación

La documentación del proyecto está en:

```text
docs/
├── REQUIREMENTS.md
├── DESIGN.md
├── PLAN.md
└── FUTURE_IMPROVEMENTS.md
```
