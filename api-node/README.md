# API Node — Estadísticas

Recibe las matrices QR que calculó las apis y devuelve máximo, mínimo, promedio, suma total, y si alguna de las dos es diagonal.

## Cómo Correr

```
cd api-node
npm install
npm run dev
```

## Cómo Probar

```
curl http://localhost:3000/health

curl -X POST http://localhost:3000/api/v1/matrix/stats \
  -H "Content-Type: application/json" \
  -d "{\"q\": [[1,0],[0,1]], \"r\": [[2,3],[0,4]]}"
```

Tests: `npm test`

## Contrato

`POST /api/v1/matrix/stats`
```json
{ "q": [[1, 0], [0, 1]], "r": [[2, 3], [0, 4]] }
```

Respuesta:
```json
{ "max": 4, "min": 0, "average": 1.375, "sum": 11, "diagonalCheck": { "q": true, "r": false, "anyDiagonal": true } }
```

Errores: mismo formato que la API Go: `{ "code": "...", "message": "...", "detail": "..." }`.

## Estructura

```
src/
  server.ts    punto de entrada
  app.ts       App Express
  logger.ts    logging en JSON
  controller/  HTTP y validación
  service/     cálculo de estadísticas
  model/       DTOs y errores
```
