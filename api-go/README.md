# API Go — Factorización QR

Recibe una matriz y devuelve su factorización QR (`A = Q × R`). También tiene un endpoint de rotación adicional.

## Cómo Correr

```
cd api-go
go run .
```

Levanta en `:8080` (o `$PORT`). Necesita un `.env` con `JWT_SECRET` (ejemplo en `.env.example`).

## Autenticación

Las rutas bajo `/api/v1/*` requieren un JWT (`/health` queda pública). Para generar uno:

```
go run ./cmd/mint-token
```

Se utiliza como `Authorization: Bearer <token>`.

## Cómo Probar

```
curl http://localhost:8080/health

curl -X POST http://localhost:8080/api/v1/matrix/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d "{\"matrix\": [[1,2,3],[4,5,6],[7,8,10]]}"

curl -X POST http://localhost:8080/api/v1/matrix/rotate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d "{\"matrix\": [[1,2,3],[4,5,6]]}"
```

Tests: `go test ./...`

## Contrato

`POST /api/v1/matrix/qr`
```json
{ "matrix": [[1, 2, 3], [4, 5, 6], [7, 8, 10]] }
```

Respuesta:
```json
{
  "qr": { "q": [[...]], "r": [[...]] },
  "stats": { "max": 0, "min": 0, "average": 0, "sum": 0, "diagonalCheck": { "q": false, "r": false, "anyDiagonal": false } }
}
```

`POST /api/v1/matrix/rotate`
```json
{ "matrix": [[1, 2, 3], [4, 5, 6]] }
```
Respuesta: `{ "matrix": [[...]] }`

Errores:
```json
{ "code": "VALIDATION_ERROR", "message": "...", "detail": "..." }
```

## Estructura

```
main.go
cmd/mint-token/   genera tokens JWT
internal/
  handler/   HTTP y validación
  service/   lógica de QR y rotación
  client/    llamada a la API Node
  auth/      middleware de JWT
  model/     DTOs y errores
```
