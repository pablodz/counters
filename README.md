# counters

Microservicio transversal de contadores sobre **Cloudflare D1** (acceso directo por HTTP REST, sin driver SQL ni ORM). Pensado para ser compartido por varios proyectos: cualquier app puede registrar y consultar contadores (views, likes, shares) de cualquier contenido identificado por `content_type` + `content_id`.

## Cómo conecta con D1

No usa `database/sql`. Cada operación es un `POST` a la REST API de D1 (`/client/v4/accounts/.../d1/database/.../query`), igual que `aldepa-backend`. Ver `d1.go`.

El schema se aplica automáticamente al arrancar (`migrate.go` con `CREATE ... IF NOT EXISTS`), así que el servicio se auto-configura la primera vez. El SQL canónico está en `schemas/post_metrics.sql`.

## Endpoints

Base: `/api/v1`. Todos requieren la API key (`X-API-Key` o `Authorization: Bearer <key>`) **solo si** `COUNTERS_API_KEY` está configurada; si está vacía el API queda abierto (útil para dev).

| Método | Ruta | Descripción |
|---|---|---|
| `POST` | `/metrics/{content_type}/{content_id}/view` | +1 view (sin body) |
| `POST` | `/metrics/{content_type}/{content_id}/like` | +1 like (sin body) |
| `POST` | `/metrics/{content_type}/{content_id}/share` | +1 share (sin body) |
| `POST` | `/metrics/{content_type}/{content_id}/increment` | suma genérica, body `{"field":"views_count","amount":1}` |
| `POST` | `/metrics/{content_type}/{content_id}/reset` | pone los contadores en 0 |
| `POST` | `/metrics/batch` | varios increments a la vez, body `{"events":[...]}` |
| `GET`  | `/metrics/{content_type}/{content_id}` | devuelve todos los contadores |
| `GET`  | `/metrics/{content_type}/{content_id}/{field}` | un solo contador |
| `GET`  | `/trending/{content_type}?limit=10` | top contenidos por views recientes |
| `GET`  | `/healthz` | probe de k8s (sin auth) |

`field` ∈ `views_count` | `likes_count` | `shares_count`.

### Ejemplo de respuesta

```json
{
  "content_id": "uuid-del-listing",
  "content_type": "listing",
  "views_count": 42,
  "likes_count": 5,
  "shares_count": 2,
  "updated_at": "2026-06-19 20:13:06"
}
```

### Ejemplo de uso (batch)

```json
POST /api/v1/metrics/batch
{
  "events": [
    {"content_type": "listing", "content_id": "abc", "field": "views_count"},
    {"content_type": "listing", "content_id": "def", "field": "views_count"}
  ]
}
```

## Configuración (env vars)

| Variable | Requerida | Descripción |
|---|---|---|
| `CF_ACCOUNT_ID` | sí | Account ID de Cloudflare |
| `CF_D1_DATABASE_ID` | sí | UUID de la base D1 de contadores |
| `CF_D1_API_TOKEN` | sí | Token con permisos D1 |
| `COUNTERS_API_KEY` | no | Si se setea, protege todo `/api` |
| `PORT` | no | Default `8080` |

## Desarrollo local

```bash
make run      # carga .env y arranca con go run
make vet      # go vet
make fmt      # go fmt
```

`.env` (gitignored) ya viene con las credenciales de la DB de contadores para dev.

## Cómo lo usan otros proyectos

Dentro del cluster k3s el servicio se alcanza como `http://counters:8080`. Ejemplo desde `aldepa-backend` para contar una view de un listing:

```go
// en un handler al abrir la página de un listing
req, _ := http.NewRequest("POST",
    "http://counters:8080/api/v1/metrics/listing/"+listingID+"/view", nil)
req.Header.Set("X-API-Key", os.Getenv("COUNTERS_API_KEY"))
client.Do(req)
```

Para leer (ej. mostrar views en la UI), el frontend del proyecto llama a su propio backend, que a su vez consulta:

```
GET http://counters:8080/api/v1/metrics/listing/{id}
```

## Deploy

```bash
make deploy
```

Eso build-ea la imagen `pablogod/counters`, la pushea al registry y ejecuta `make deploy-project p=counters` en `vps-deploy` (aplica `apps/counters/deployment.yaml` en k3s y reinicia el rollout).

## Estructura

```
counters/
├── main.go          # routing (stdlib net/http, Go 1.22+) + arranque
├── config.go        # env vars + validación
├── d1.go            # cliente HTTP REST de D1
├── store.go         # upsert atómico, get, trending, reset
├── handlers.go      # HTTP handlers
├── middleware.go    # auth (API key) + logging
├── migrate.go       # auto-migración del schema (embed)
├── models.go        # struct Metrics + whitelist de campos
├── schemas/
│   └── post_metrics.sql   # schema canónico (idempotente)
├── Dockerfile       # build CGO_ENABLED=0, debian-slim
└── Makefile         # deploy / run / fmt / vet
```

Sin dependencias externas: puro stdlib de Go. Imagen mínima, build rápido.
