# Monitoring Guide — Prometheus + Loki + Grafana

เอกสารนี้อธิบายการติดตั้ง การใช้งาน และการดู logs/metrics ของระบบ Product CRUD ที่รันผ่าน Docker

## Architecture

```
┌─────────────┐     /api/*      ┌─────────────┐
│  Frontend   │ ──────────────► │   Backend   │◄──── scrape /metrics
│  (Nginx)    │                 │  (Go Fiber) │
│  :5173      │                 │   :8080     │
└─────────────┘                 └──────┬──────┘
                                       │ stdout logs
                                       ▼
                                ┌─────────────┐
                                │  Promtail   │──────► ┌─────────┐
                                │  (agent)    │        │  Loki   │
                                └─────────────┘        └────┬────┘
                                                            │
┌─────────────┐     query metrics    ┌─────────────┐       │
│ Prometheus  │ ◄────────────────────│   Grafana   │◄──────┘
│   :9090     │                      │   :3000     │  query logs
└─────────────┘                      └─────────────┘
```

## Components

| Component | Port | หน้าที่ |
|-----------|------|---------|
| **Backend** | 8080 | API + `/metrics` endpoint สำหรับ Prometheus |
| **Frontend** | 5173 | React app (Nginx) proxy `/api` ไป backend |
| **Prometheus** | 9099 | เก็บ metrics (request rate, latency, Go runtime) |
| **Loki** | 3100 | เก็บ logs จาก containers |
| **Promtail** | 9080 (internal) | ดึง logs จาก Docker containers → ส่ง Loki |
| **Grafana** | 3000 | Dashboard รวม metrics + logs |

## Prerequisites

- Docker Desktop หรือ Docker Engine + Docker Compose v2
- RAM ว่างอย่างน้อย **4 GB**
- Port ว่าง: `3000`, `5173`, `8080`, `9099`, `3100`

## การติดตั้งและรัน

### 1. รันทั้งระบบ (App + Monitoring)

จาก root ของโปรเจกต์:

```bash
docker compose up -d --build
```

### 2. ตรวจสอบสถานะ containers

```bash
docker compose ps
```

ทุก service ควรอยู่ในสถานะ `running`:

| Container | Service |
|-----------|---------|
| product-crud-backend | backend |
| product-crud-frontend | frontend |
| product-crud-prometheus | prometheus |
| product-crud-loki | loki |
| product-crud-promtail | promtail |
| product-crud-grafana | grafana |

### 3. เข้าใช้งาน

| URL | คำอธิบาย |
|-----|----------|
| http://localhost:5173 | Frontend (Product CRUD UI) |
| http://localhost:8080/api/products | Backend API |
| http://localhost:8080/metrics | Prometheus metrics (raw) |
| http://localhost:9099 | Prometheus UI |
| http://localhost:3000 | Grafana (user: `admin`, password: `admin`) |

### 4. หยุดระบบ

```bash
docker compose down
```

ลบ volumes ด้วย (metrics/logs ที่เก็บไว้):

```bash
docker compose down -v
```

## Grafana Dashboard

1. เปิด http://localhost:3000
2. Login: `admin` / `admin`
3. ไปที่ **Dashboards** → folder **Product CRUD** → **Product CRUD Monitoring**

Dashboard แสดง:

- **Request Rate** — จำนวน request ต่อวินาที
- **Error Rate (4xx/5xx)** — อัตรา error
- **Response Latency (p95)** — latency เปอร์เซนไทล์ 95
- **Go Runtime** — goroutines และ memory
- **Backend Logs** — log stream จาก backend container

## การดู Logs (Loki)

### ผ่าน Grafana Explore

1. เปิด Grafana → **Explore**
2. เลือก datasource **Loki**
3. ใส่ LogQL query

### ตัวอย่าง LogQL

```logql
# Log ทั้งหมดของ backend
{service="backend"}

# Log ของ frontend (Nginx)
{service="frontend"}

# Filter HTTP GET requests
{service="backend"} |= "GET"

# Filter errors
{service="backend"} |= "error"

# Filter โดย container name
{container="product-crud-backend"}
```

### ผ่าน Docker CLI (realtime)

```bash
# Backend logs
docker compose logs -f backend

# Frontend logs
docker compose logs -f frontend

# Promtail logs (debug log shipping)
docker compose logs -f promtail

# ทุก services
docker compose logs -f
```

## การดู Metrics (Prometheus)

### ผ่าน Prometheus UI

1. เปิด http://localhost:9099
2. ไปที่ **Status → Targets** — ตรวจว่า `backend` เป็น **UP**
3. ไปที่ **Graph** แล้วลอง query

### ตัวอย่าง PromQL

```promql
# Request rate
sum(rate(http_requests_total{job="backend"}[1m]))

# Error rate
sum(rate(http_requests_total{job="backend",status_code=~"4..|5.."}[5m]))
/ sum(rate(http_requests_total{job="backend"}[5m]))

# Latency p95
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket{job="backend"}[5m])) by (le)
)

# Go goroutines
go_goroutines{job="backend"}

# Memory usage
process_resident_memory_bytes{job="backend"}
```

### ตรวจ metrics โดยตรง

```bash
curl http://localhost:8080/metrics
```

## การทำงานของแต่ละส่วน

### Backend (Go Fiber)

- ใช้ `fiberprometheus` middleware เก็บ metrics ทุก HTTP request
- Expose endpoint `/metrics` ให้ Prometheus scrape
- Log ออก stdout ผ่าน Fiber logger middleware
- Format log: `time | status | latency | ip | method | path | error`

### Frontend (Nginx)

- Serve React static files
- Proxy `/api/*` ไป `backend:8080` (same-origin, ไม่ต้อง CORS)
- Access log ออก stdout → Promtail เก็บไป Loki

### Prometheus

- Scrape `backend:8080/metrics` ทุก 15 วินาที
- Config: `monitoring/prometheus/prometheus.yml`

### Loki + Promtail

- Promtail ใช้ Docker socket อ่าน logs จาก containers
- Label `service` มาจาก Docker Compose service name
- Loki เก็บ logs retention 7 วัน (168h)
- Config: `monitoring/loki/loki.yml`, `monitoring/promtail/promtail.yml`

### Grafana

- Auto-provision datasources (Prometheus + Loki)
- Auto-load dashboard จาก `monitoring/grafana/dashboards/`
- Config: `monitoring/grafana/provisioning/`

## Troubleshooting

| ปัญหา | วิธีแก้ |
|-------|---------|
| Grafana ไม่มี metrics | ตรวจ Prometheus → Targets → backend UP หรือไม่ |
| ไม่มี logs ใน Loki | ตรวจ `docker compose logs promtail` และ Promtail มี access Docker socket |
| Port 8080 ถูกใช้แล้ว | หยุด process/container ที่ใช้ port 8080 ก่อน |
| Frontend เรียก API ไม่ได้ | ตรวจ backend running และ nginx proxy config |
| Dashboard ว่าง | สร้าง traffic ก่อน (เปิด frontend, refresh หรือ CRUD สินค้า) |

### คำสั่ง debug

```bash
# ตรวจ backend health
curl http://localhost:8080/api/products

# ตรวจ metrics endpoint
curl -s http://localhost:8080/metrics | head -20

# Restart service เดียว
docker compose restart backend

# Rebuild หลังแก้ code
docker compose up -d --build backend
```

## โครงสร้างไฟล์ Monitoring

```
monitoring/
├── prometheus/
│   └── prometheus.yml          # Scrape config
├── loki/
│   └── loki.yml                # Loki storage config
├── promtail/
│   └── promtail.yml            # Docker log collection
└── grafana/
    ├── provisioning/
    │   ├── datasources/
    │   │   └── datasources.yml # Auto-config Prometheus + Loki
    │   └── dashboards/
    │       └── dashboards.yml  # Dashboard loader
    └── dashboards/
        └── product-crud.json   # Product CRUD dashboard

docker/
├── backend/Dockerfile
├── frontend/Dockerfile
└── nginx/nginx.conf

docker-compose.yml              # รวม app + monitoring stack
```

## รันแบบ Local (ไม่ใช้ Docker)

Monitoring stack ออกแบบสำหรับ Docker แต่ยังรัน app แบบ local ได้ตามเดิม:

```bash
# Terminal 1 — Backend
cd backend && go run .

# Terminal 2 — Frontend
cd frontend && npm run dev
```

ในโหมด local, metrics อยู่ที่ http://localhost:8080/metrics แต่ Prometheus/Loki/Grafana ต้องรันผ่าน Docker
