# ---- build ----
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata git

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY internal/ ./internal
COPY migrations/ ./migrations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pocketbase ./main.go

# ---- runtime ----
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

# Copiar binario y migraciones
COPY --from=builder /app/pocketbase /app/pocketbase
COPY ./migrations/ /app/pb_migrations

# Copiar cred.json (clave de Firebase)
COPY internal/utils/cred.json /app/internal/utils/cred.json

# Crear carpeta de datos y ejecutar PocketBase
ENTRYPOINT ["/bin/sh","-lc","mkdir -p /data && ./pocketbase migrate up --dir /data && exec ./pocketbase serve --http 0.0.0.0:${PORT:-8080} --dir /data"]
