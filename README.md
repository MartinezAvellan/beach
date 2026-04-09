# Cameras API

Backend em Go para cameras de praia ao vivo em Portugal com dados meteorologicos e oceanograficos do IPMA, expondo uma API REST e frontend integrado.

Catalogo estatico, streams HLS via URL direta, condicoes via APIs JSON publicas.

---

## Stack

| Componente | Tecnologia |
|---|---|
| Linguagem | Go 1.25 |
| Router | chi v5 |
| Cache | Memoria com TTL |
| Frontend | HTML/CSS/JS + HLS.js (embebido no binario) |
| Container | Docker multi-stage (Alpine) |
| Orquestracao | Kubernetes (HPA) |

**Unica dependencia externa:** `github.com/go-chi/chi/v5`

---

## Arquitetura

```
┌──────────────────────────────────────────────────────┐
│                     Frontend                         │
│              (embed go:embed, HLS.js)                │
│         http://localhost:8080/                        │
└──────────────┬───────────────────────────────────────┘
               │
┌──────────────▼───────────────────────────────────────┐
│                    API REST                           │
│              (chi router + handlers)                  │
├──────────────────────────────────────────────────────┤
│                Application Layer                     │
│         CameraService (ListCameras,                  │
│         GetCameraByID, ResolveCameraStream,           │
│         ResolveCameraConditions)                      │
├──────────────────────────────────────────────────────┤
│                   Adapters                            │
│  ┌─────────┐  ┌──────────┐  ┌───────┐  ┌─────────┐ │
│  │ Catalog  │  │  Stream  │  │ IPMA  │  │  Cache  │ │
│  │ (embed)  │  │(URL dir.)│  │(APIs) │  │ (mem)   │ │
│  └─────────┘  └──────────┘  └───────┘  └─────────┘ │
├──────────────────────────────────────────────────────┤
│                    Domain                             │
│     Camera, CameraStream, CameraConditions           │
│     CameraRepository, StreamResolver,                │
│     ConditionsResolver (interfaces)                   │
└──────────────────────────────────────────────────────┘
```

### Camadas

| Camada | Path | Responsabilidade |
|---|---|---|
| Domain | `internal/domain/` | Entidades e interfaces |
| Application | `internal/application/` | Logica de negocio |
| Adapters/HTTP | `internal/adapters/http/` | Router, handlers, frontend |
| Adapters/Catalog | `internal/adapters/catalog/` | Catalogo estatico (187 cameras) |
| Adapters/Stream | `internal/adapters/stream/` | Monta URL HLS direto |
| Adapters/IPMA | `internal/adapters/ipma/` | Client APIs IPMA |
| Adapters/Cache | `internal/adapters/cache/` | Cache em memoria com TTL |
| Config | `internal/config/` | Variaveis de ambiente |

---

## Endpoints

| Metodo | Rota | Descricao | Fonte | Latencia |
|---|---|---|---|---|
| GET | `/` | Frontend web (SPA) | embed | ~50us |
| GET | `/health` | Healthcheck | hardcode | ~50us |
| GET | `/api/v1/cameras` | Lista 187 cameras | embed | ~100us |
| GET | `/api/v1/cameras/{id}` | Detalhe de camera | embed | ~20us |
| GET | `/api/v1/cameras/{id}/stream` | URL HLS (.m3u8) | URL direta | ~20us |
| GET | `/api/v1/cameras/{id}/conditions` | Condicoes meteo/ocean | APIs IPMA | ~100us (cache) |

### Exemplos de resposta

**GET /api/v1/cameras/{id}/stream**
```json
{
  "camera_id": "penichesupertubos",
  "stream_url": "https://video-auth1.iol.pt/beachcam/supertubos/playlist.m3u8",
  "status": "online",
  "expires_at": "2026-04-09T16:12:28Z"
}
```

**GET /api/v1/cameras/{id}/conditions**
```json
{
  "camera_id": "penichesupertubos",
  "wave_height": "2.8 m",
  "wave_period": "8.9 s",
  "wave_direction": "NW",
  "wind_speed": "9.4 km/h",
  "wind_direction": "Nordeste",
  "water_temp": "15.1 ºC",
  "air_temp": "21.4 ºC",
  "uv_index": "4.5",
  "weather": "Ceu parcialmente nublado",
  "humidity": "50%",
  "fetched_at": "2026-04-09T14:40:27Z"
}
```

### Erros

| Codigo | Corpo | Quando |
|---|---|---|
| 404 | `{"error": "camera_not_found"}` | ID nao existe no catalogo |
| 503 | `{"error": "stream_unavailable"}` | Camera sem stream proprio (38 de 187) |
| 503 | `{"error": "conditions_unavailable"}` | IPMA indisponivel |

---

## Fontes de dados

| Dado | Fonte | Tipo |
|---|---|---|
| Catalogo (187 cameras) | `cameras.json` embebido | Estatico (go:embed) |
| Stream HLS (149 cameras) | `video-auth1.iol.pt/beachcam/{slug}/playlist.m3u8` | URL construida |
| Ondulacao, periodo, temp. agua | IPMA Sea Forecast API | JSON publico |
| Temp. ar, vento, tipo de tempo | IPMA Weather Forecast API | JSON publico |
| UV | IPMA UV API | JSON publico |
| Vento real, temp. real, humidade | IPMA Observations API | JSON publico |

### APIs IPMA consumidas

| API | Endpoint | Cache |
|---|---|---|
| Sea Forecast | `api.ipma.pt/open-data/forecast/oceanography/daily/hp-daily-sea-forecast-day0.json` | 15min |
| Weather Forecast | `api.ipma.pt/open-data/forecast/meteorology/cities/daily/hp-daily-forecast-day0.json` | 15min |
| UV Index | `api.ipma.pt/open-data/forecast/meteorology/uv/uv.json` | 15min |
| Observations | `api.ipma.pt/open-data/observation/meteorology/stations/observations.json` | 15min |
| Weather Types | `api.ipma.pt/open-data/weather-type-classe.json` | Permanente |

Mapeamento camera-estacao IPMA feito por proximidade geografica (haversine).

---

## Frontend

SPA responsiva embebida no binario Go, acessivel em `http://localhost:8080/`.

- Lista de 187 cameras com busca por nome/localizacao
- Player HLS ao vivo (HLS.js - funciona em qualquer browser)
- Cards de condicoes meteorologicas e oceanograficas
- Auto-refresh a cada 5 minutos (stream + condicoes)
- Countdown visivel ate proxima atualizacao
- Dark mode nativo
- Responsivo: desktop (sidebar) e mobile (drawer)
- Video streaming vai direto do browser para `video-auth1.iol.pt` (nao passa pelo backend)

---

## Como correr

### Local

```bash
go run ./cmd/api
# abrir http://localhost:8080
```

### Build

```bash
go build -o bin/api ./cmd/api
./bin/api
```

### Docker

```bash
docker compose up --build
# http://localhost:8080
```

### Variaveis de ambiente

| Variavel | Default | Descricao |
|---|---|---|
| `APP_NAME` | cameras-api | Nome da aplicacao |
| `APP_ENV` | local | Ambiente (local/production) |
| `APP_PORT` | 8080 | Porta HTTP |
| `HTTP_CLIENT_TIMEOUT` | 10s | Timeout para calls IPMA |
| `CATALOG_CACHE_TTL` | 15m | TTL do cache de condicoes |
| `STREAM_CACHE_TTL` | 1m | TTL do cache de streams |
| `LOG_LEVEL` | info | Nivel de log |

---

## Estrutura do projeto

```
beach/
├── cmd/api/main.go                          # Entrypoint
├── internal/
│   ├── config/config.go                     # Env vars
│   ├── domain/
│   │   ├── camera.go                        # Camera, CameraStream, CameraConditions
│   │   └── repository.go                    # Interfaces
│   ├── application/service.go               # CameraService
│   └── adapters/
│       ├── http/
│       │   ├── router.go                    # Chi router + frontend
│       │   ├── handlers.go                  # REST handlers
│       │   └── static/index.html            # Frontend SPA
│       ├── catalog/
│       │   ├── catalog.go                   # StaticCatalog (go:embed)
│       │   └── cameras.json                 # 187 cameras
│       ├── stream/resolver.go               # URL direta HLS
│       ├── ipma/
│       │   ├── client.go                    # Client IPMA (4 APIs)
│       │   └── resolver.go                  # Mapeamento geo + resolver
│       └── cache/memory.go                  # Cache em memoria TTL
├── k8s/deployment.yaml                      # K8s config completa
├── Dockerfile                               # Multi-stage Alpine
├── docker-compose.yml
├── cameras-api.postman_collection.json      # Postman collection
├── go.mod
└── go.sum
```

---

## Performance e escala

### Perfil da aplicacao

Esta API e extremamente leve porque:
- **Zero I/O** nos endpoints principais (tudo servido de memoria)
- **Zero base de dados** (catalogo embebido, cache in-memory)
- **Zero I/O** nos endpoints principais (tudo servido de memoria)
- **Video streaming nao passa pelo backend** (browser -> video-auth1.iol.pt direto)
- Unico I/O externo: 4 calls IPMA a cada 15 minutos (cache global, nao por user)

### Benchmarks por endpoint

| Endpoint | Latencia (p99) | Throughput (1 pod) |
|---|---|---|
| `GET /` | ~50us | 100.000+ req/s |
| `GET /cameras` | ~100us | 80.000+ req/s |
| `GET /cameras/{id}` | ~20us | 150.000+ req/s |
| `GET /cameras/{id}/stream` | ~20us | 150.000+ req/s |
| `GET /cameras/{id}/conditions` (cache hit) | ~100us | 100.000+ req/s |
| `GET /cameras/{id}/conditions` (cache miss) | ~300ms | N/A (1x cada 15min) |

### Estimativa de carga por users simultaneos

```
Users online → Requests por segundo (steady state)

 1.000 users →   ~3 req/s   (conditions refresh cada 5min)
10.000 users →  ~33 req/s
30.000 users → ~100 req/s
100.000 users → ~333 req/s
```

Pico no page load (burst): ~10x o steady state durante segundos.

### Capacidade por pods

| Users simultaneos | Pods necessarios | CPU total | RAM total |
|---|---|---|---|
| 1.000 | 1 | 50m | 32Mi |
| 10.000 | 2 | 100m | 64Mi |
| 30.000 | 3 | 150m | 96Mi |
| 100.000 | 4-6 | 300m | 192Mi |

**Um unico pod Go aguenta +50.000 req/s** neste perfil. Os pods extras sao para redundancia e disponibilidade, nao para capacidade.

### Kubernetes (30k-100k users)

Configuracao completa em `k8s/deployment.yaml`:

| Recurso | Valor |
|---|---|
| Deployment replicas | 3 |
| HPA min/max | 2-6 pods |
| CPU request/limit | 50m / 200m por pod |
| Memory request/limit | 32Mi / 64Mi por pod |
| Health probes | `/health` (readiness 2s, liveness 5s) |
| Strategy | RollingUpdate (zero downtime) |
| Service | ClusterIP (porta 80 -> 8080) |
| Ingress | NGINX |

### Custo estimado em cloud

| Cenario | Pods | Custo mensal estimado |
|---|---|---|
| Ate 10k users | 2 | ~$5 |
| 30k users | 3 | ~$8 |
| 100k users | 6 | ~$15 |

O binario final tem ~15MB. A imagem Docker (Alpine) tem ~25MB.

### Bottlenecks potenciais

| Componente | Limite | Mitigacao |
|---|---|---|
| Cache em memoria | ~10MB (187 cameras + condicoes) | Mais que suficiente ate 1M users |
| IPMA APIs | Rate limit desconhecido | Cache 15min = max 4 calls/hora |
| video-auth1.iol.pt | Servidor IOL (nao controlamos) | Browser acede direto, nao passa pelo backend |
| Ingress/LB | Depende do provider | CDN para frontend statico |

### Escalabilidade futura

Para ir alem de 100k users simultaneos:
1. **CDN** na frente (Cloudflare/CloudFront) para cachear `/`, `/cameras`, `/conditions`
2. **Redis** para cache partilhado entre pods (substitui in-memory)
3. **Postgres** para persistencia e analytics
4. **WebSocket** para push de conditions (elimina polling de 5min)

---

## Postman

Importar `cameras-api.postman_collection.json` no Postman.

Variaveis:
- `base_url`: `http://localhost:8080`
- `camera_id`: `penichesupertubos` (ou qualquer outro ID)

---

## Roadmap

### MVP (concluido)
- [x] API REST completa
- [x] Healthcheck
- [x] Catalogo estatico de 187 cameras
- [x] Stream HLS via URL direta (149 cameras)
- [x] Condicoes meteo/ocean via IPMA (APIs JSON publicas)
- [x] Cache em memoria com TTL
- [x] Frontend responsivo com player HLS
- [x] Auto-refresh a cada 5 minutos
- [x] Dockerfile multi-stage
- [x] docker-compose
- [x] Kubernetes deployment + HPA
- [x] Postman collection

### Fase 2
- [ ] Redis para cache partilhado
- [ ] Postgres para persistencia
- [ ] Filtros por regiao/localizacao
- [ ] Paginacao
- [ ] Mares (fonte alternativa ao IPMA)

### Fase 3
- [ ] Autenticacao (JWT)
- [ ] Favoritos
- [ ] Notificacoes (condicoes ideais de surf)
- [ ] App Flutter (consumir esta API)
- [ ] WebSocket para push de dados
- [ ] CDN para escala 1M+ users
