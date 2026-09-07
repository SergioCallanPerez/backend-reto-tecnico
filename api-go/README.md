# API Go — Factorización QR

Api que recibe una matriz rectangular y devuelve su factorización QR (`A = Q × R`), calculada con Gram-Schmidt modificado.

## Cómo iniciar

Requiere Go instalado.

```
cd api-go
go run .
```

Levanta en `http://localhost:{PORT}`.

Compilar un binario en vez de correrlo directo:
```
go build ./...
```

## Cómo probar si el servicio está UP

Health check:
```
curl http://localhost:8080/health
```

Para Factorización QR:
```
curl -X POST http://localhost:8080/api/v1/matrix/qr \
  -H "Content-Type: application/json" \
  -d "{\"matrix\": [[1,2,3],[4,5,6],[7,8,10]]}"
```

Esquema de respuesta:
```json
{
  "qr": {
    "q": [[...]],
    "r": [[...]]
  }
}
```

## Contrato

**Esquema de Request** — `POST /api/v1/matrix/qr`
```json
{ "matrix": [[1, 2, 3], [4, 5, 6], [7, 8, 10]] }
```
- `matrix`: rectangular, m×n, con m ≥ n, dimensiones entre 1 y 500.

**Esquema de Response (caso exitoso)**
```json
{ "qr": { "q": [[...]], "r": [[...]] } }
```
Los valores están redondeados a 6 decimales (solo resultado, no cálculo).

**Esquema de Response (caso error)**
```json
{ "code": "VALIDATION_ERROR", "message": "...", "detail": "..." }
```
Códigos: `VALIDATION_ERROR` (400), `UPSTREAM_ERROR` (502), `INTERNAL_ERROR` (500).
