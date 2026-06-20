# Counters

Servicio de conteo de eventos polimórfico con soporte para múltiples tipos de items.

## Arquitectura

```
counters/
├── pkg/
│   └── client/        # Cliente Go para interactuar con el servicio
├── data/
│   ├── models/        # Modelos de datos D1
│   └── schemas/       # SQL schemas para D1
├── handlers/          # Handlers HTTP del servidor
├── singleton/         # Estado global del servidor
└── examples/          # Ejemplos de uso del cliente
```

## Servidor

El servidor se ejecuta en `:8080` y expone los siguientes endpoints:

### API v1

- `POST /api/v1/metrics/batch` - Registrar múltiples eventos en batch
- `GET /api/v1/:item_type/:item_id/:field` - Obtener valor de un campo
- `GET /api/v1/:item_type/:item_id` - Obtener métricas básicas

### Dashboard Endpoints

- `GET /api/v1/dashboard/historical/:item_type/:item_id` - Total histórico por item
- `GET /api/v1/dashboard/daily/:item_type/:item_id/:event_type` - Histograma diario
- `GET /api/v1/dashboard/hourly/:item_type/:item_id/:event_type` - Histograma horario
- `GET /api/v1/dashboard/versus/:item_type/:item_id/:event_type` - Comparativa top 10
- `GET /api/v1/dashboard/owner/:item_type/:item_id/:event_type` - Últimos usuarios
- `POST /api/v1/dashboard/cleanup` - Limpiar logs antiguos

## Base de Datos

### Tablas

1. **item_logs**
   - Registro detallado de cada interacción
   - Campos: id, item_id, item_type, user_id, user_name, event_type, created_at_unix

2. **item_snapshots**
   - Agregaciones por hora para consultas rápidas
   - Campos: item_id, item_type, event_type, period_hour_unix, total_count

### Índices

- `idx_logs_item`: (item_id, item_type, event_type)
- `idx_logs_cleanup`: (created_at_unix)
- Primary key: (item_id, item_type, event_type, period_hour_unix)

## Cliente

El cliente está disponible en `pkg/client` y puede ser instalado con:

```bash
go get github.com/pablodz/counters/pkg/client
```

Ver `pkg/client/README.md` para documentación completa.

## Desarrollo

### Compilar

```bash
go build -o counters .
```

### Ejecutar

```bash
go run main.go
```

Las variables de entorno requeridas son:
- `CF_ACCOUNT_ID`: Cloudflare account ID
- `CF_D1_DATABASE_ID`: D1 database ID
- `CF_D1_API_TOKEN`: D1 API token

### Correr ejemplo

```bash
cd examples
go run client_example.go
```

## Design Decisions

### Polimorfismo (item_id + item_type)

En lugar de tener tablas separadas para cada tipo (listings, leads, etc.), usamos un esquema polimórfico con:
- `item_id`: Identificador único del item
- `item_type`: Tipo de item (listing, lead, etc.)

Esto permite:
- Consultas uniformes independientes del tipo
- Análisis cruzado entre tipos
- Escalabilidad sin necesidad de migraciones

### Snapshots vs Logs

- **logs**: Registro completo, histórico detallado
- **snapshots**: Agregaciones por hora para performance

El cron job `CleanupOldLogs` elimina logs de más de 30 días automáticamente.

## License

MIT
