# Getting Started

This guide walks you through running hada-commerce locally from scratch.

---

## Prerequisites

| Tool | Minimum version | Notes |
|------|----------------|-------|
| Go | 1.24+ | [golang.org/dl](https://golang.org/dl/) |
| Node.js | 20+ | Required for the admin panel |
| npm | 10+ | Bundled with Node 20 |
| Docker & Docker Compose | latest | Used for PostgreSQL and Redis |
| `psql` (optional) | 16 | For running migrations manually |

---

## 1. Clone the repository

```bash
git clone https://github.com/Abraxas-365/hada-commerce.git
cd hada-commerce
```

---

## 2. Start infrastructure services

Docker Compose starts PostgreSQL 16 and Redis 7 in the background.

```bash
docker-compose up -d
```

Verify the services are healthy:

```bash
docker-compose ps
```

Expected output:

```
NAME            IMAGE           STATUS
hada-postgres   postgres:16     Up (healthy)
hada-redis      redis:7         Up (healthy)
```

**Default connection details** (set in `docker-compose.yml`):

| Service | Host | Port | User | Password | Database |
|---------|------|------|------|----------|----------|
| PostgreSQL | localhost | 5433 | hada | hada | hada |
| Redis | localhost | 6379 | — | — | — |

---

## 3. Configure environment variables

Copy the example environment file and edit it:

```bash
cp backend/.env.example backend/.env
```

Minimum required variables:

```dotenv
# Database
DATABASE_URL=postgres://hada:hada@localhost:5433/hada?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# Server
PORT=8080
ENV=development

# File storage (local by default)
STORAGE_DRIVER=local
STORAGE_LOCAL_PATH=./uploads
```

---

## 4. Run database migrations

Migrations live in `backend/migrations/` and are numbered sequentially (e.g. `001_init.up.sql`). Run them in order with `psql`:

```bash
for f in backend/migrations/*.up.sql; do
  echo "Applying $f..."
  psql "$DATABASE_URL" -f "$f"
done
```

Or use the Makefile shortcut if available:

```bash
make migrate
```

After applying all migrations you should have these tables (among others):

```
products         orders          customers
categories       collections     pages
sections         blocks          block_types
themes           plugins         plugin_versions
plugin_installations  store_settings  media
promos           vendors
```

---

## 5. Start the backend

```bash
cd backend
go run ./cmd/...
```

The server starts on port `8080` by default. You should see:

```
[INFO] hada-commerce starting on :8080
[INFO] tenant extraction middleware active
[INFO] event bus started
```

Verify the API is reachable:

```bash
curl http://localhost:8080/api/v1/settings \
  -H "X-Tenant-ID: my-store"
```

---

## 6. Start the admin panel

In a second terminal:

```bash
npm install
npm run dev
```

The React admin panel is available at `http://localhost:3000`.

---

## 7. Create your first tenant

hada-commerce is multi-tenant. Every resource belongs to a tenant. The tenant ID is passed in the `X-Tenant-ID` request header.

To bootstrap a tenant, choose a slug (e.g. `my-store`) and create initial settings:

```bash
curl -X PUT http://localhost:8080/api/v1/settings \
  -H "X-Tenant-ID: my-store" \
  -H "Content-Type: application/json" \
  -d '{
    "store_name": "My Store",
    "store_email": "hello@mystore.com",
    "currency": "USD",
    "timezone": "UTC"
  }'
```

Response:

```json
{
  "tenant_id": "my-store",
  "store_name": "My Store",
  "store_email": "hello@mystore.com",
  "currency": "USD",
  "timezone": "UTC",
  "checkout_config": {
    "guest_checkout": true,
    "require_phone": false
  },
  "updated_at": "2024-01-15T10:00:00Z"
}
```

---

## 8. Create your first product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "X-Tenant-ID: my-store" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Classic T-Shirt",
    "description": "A comfortable everyday tee.",
    "price_amount": 2999,
    "currency": "USD",
    "sku": "TSHIRT-001",
    "stock": 100,
    "status": "active"
  }'
```

`price_amount` is always in the **smallest currency unit** (cents for USD). `2999` = $29.99.

---

## 9. Explore the admin panel

Open `http://localhost:3000` in your browser. The admin panel lets you:

- **Products** — create, edit, archive products; manage stock
- **Orders** — view and fulfil orders; transition status
- **Customers** — browse customer profiles and order history
- **Catalog** — manage categories and collections
- **Pages** — build pages with the block editor
- **Themes** — customize design tokens
- **Plugins** — install and configure plugins
- **Settings** — store-wide configuration

---

## Makefile targets

| Target | Description |
|--------|-------------|
| `make dev` | Start backend + frontend in watch mode |
| `make build` | Build backend binary to `bin/hada-commerce` |
| `make test` | Run all Go tests |
| `make migrate` | Apply pending SQL migrations |
| `make lint` | Run `go vet` and staticcheck |

---

## Troubleshooting

**Port 5433 already in use**
Another PostgreSQL instance may be running. Change the host port in `docker-compose.yml`:
```yaml
ports:
  - "5434:5432"  # use 5434 on the host
```
Update `DATABASE_URL` accordingly.

**`go: module not found` errors**
Make sure you are inside the `backend/` directory and dependencies are downloaded:
```bash
cd backend && go mod download
```

**Frontend blank screen**
Check that the backend is running and the API URL in the frontend `.env` matches `http://localhost:8080`.
