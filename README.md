# Coding Challenge Interseguro — QR + Estadísticas

Son dos APIs que se comunican entre sí: una en Go calcula la factorización QR de una matriz, y otra en Node calcula estadísticas sobre el resultado. El cliente solo le habla a la API Go.

## Flujo

Cliente → API Go (calcula QR) → API Node (calcula estadísticas) → API Go → Cliente.

Node nunca se expone al cliente final, es un servicio interno.

## Correr localmente

Sin Docker:
```
cd api-node && npm install && npm run dev
cd api-go && go run .   # necesita .env con NODE_API_URL y JWT_SECRET
```

Con Docker Compose:
```
docker compose up --build
```

Cada API tiene su propio README con el detalle de cómo correrla y probarla.

## Tests

```
cd api-go && go test ./...
cd api-node && npm test
```

## Despliegue

Ambas APIs corren en Cloud Run (`us-central1`):

- API Go: https://api-go-sert5qga2q-uc.a.run.app
- API Node: https://api-node-sert5qga2q-uc.a.run.app (solo la llama Go, no está pensada para uso directo)

Las rutas bajo `/api/v1` requieren JWT (`Authorization: Bearer <token>`) — se genera con `go run ./cmd/mint-token` dentro de `api-go`.

```
curl -X POST https://api-go-sert5qga2q-uc.a.run.app/api/v1/matrix/qr -H "Content-Type: application/json" -H "Authorization: Bearer <token>" -d "{\"matrix\": [[1,2,3],[4,5,6],[7,8,10]]}"
```

## CI/CD

Cloud Build tiene dos triggers (uno por servicio), disparados por push a `main` y filtrados por carpeta. Cada API se redespliega sola sin afectar a la otra, aunque compartan repo.

## Frontend

Hay un frontend básico en un repo aparte, desplegado en Firebase Hosting.

## Estructura

```
api-go/       API Go (Fiber)
api-node/     API Node (Express + TypeScript)
```
