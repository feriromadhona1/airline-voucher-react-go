# Airline Voucher Seat Assignment (Go)

React frontend + Go (Echo) backend with SQLite.

```
airline-voucher-react-go/
├── frontend/          # React + Vite
├── backend/           # Go (Echo) + SQLite
├── README.md
└── docker-compose.yml
```

## 1. Prerequisites

- Go 1.24+
- Node.js 18+ with npm
- Docker (optional)

## 2. Install dependencies

### Backend

```bash
cd backend
go mod tidy
```

### Frontend

```bash
cd frontend
npm install
```

## 3. Environment setup

### Frontend

```bash
cd frontend
cp .env.example .env
```

`.env`:

```env
VITE_API_URL=http://127.0.0.1:8000
```

### Backend

No `.env` required. Optional environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8000` | HTTP port |
| `DB_PATH` | `./vouchers.db` | SQLite file path |

SQLite file is created automatically as `vouchers.db` and includes the `vouchers` table.

## 4. Run the app

### Backend

```bash
cd backend
go run .
```

API: http://127.0.0.1:8000

- `POST /api/check`
- `POST /api/generate`

### Frontend

```bash
cd frontend
npm run dev
```

Open http://127.0.0.1:5173

## 5. Docker (optional)

From the project root:

```bash
docker compose up --build
```

- Backend: http://127.0.0.1:8000
- Frontend: http://127.0.0.1:5173

Stop:

```bash
docker compose down
```

## Sample requests

### Check voucher

```http
POST /api/check
Content-Type: application/json

{
  "flightNumber": "GA102",
  "date": "2025-07-12"
}
```

### Generate voucher

```http
POST /api/generate
Content-Type: application/json

{
  "name": "Sarah",
  "id": "98123",
  "flightNumber": "ID102",
  "date": "2025-07-12",
  "aircraft": "Airbus 320"
}
```

Aircraft types: `ATR`, `Airbus 320`, `Boeing 737 Max`

## Backend layout

```
backend/
├── main.go
├── vouchers.db          # created at runtime
└── internal/
    ├── database/        # SQLite access (parameterized queries)
    ├── handlers/        # HTTP handlers
    ├── models/          # request/response types
    └── seats/           # seat map + random generation
```
