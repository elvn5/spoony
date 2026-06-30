# ============================================================
# Stage 1: Build Go backend
# ============================================================
FROM golang:1.23-alpine AS go-builder

RUN apk add --no-cache git
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# ============================================================
# Stage 2: Build Vue.js frontend
# ============================================================
FROM node:20-alpine AS node-builder

WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# ============================================================
# Stage 3: All-in-one runtime
# ============================================================
FROM alpine:3.19

RUN apk add --no-cache \
    postgresql15 \
    nginx \
    supervisor \
    ca-certificates \
    tzdata \
    bash \
    su-exec

# Built artifacts
COPY --from=go-builder  /app/server           /app/server
COPY --from=node-builder /app/dist            /usr/share/nginx/html

# Configs
COPY nginx.conf          /etc/nginx/http.d/default.conf
COPY supervisord.conf    /etc/supervisord.conf
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# PostgreSQL directories
RUN mkdir -p /var/lib/postgresql/data /run/postgresql /var/log/supervisor \
    && chown -R postgres:postgres /var/lib/postgresql /run/postgresql

EXPOSE 80

ENTRYPOINT ["/docker-entrypoint.sh"]
