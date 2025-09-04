# ============
# Base stage
# ============
FROM golang:1.23-bullseye AS base
WORKDIR /app

# Install git, curl, certs for Go + Node
RUN apt-get update && apt-get install -y \
    git \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js & npm
ENV NODE_VERSION=22
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g npm@latest \
    && rm -rf /var/lib/apt/lists/*

# Copy and download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install air & templ for dev
RUN go install github.com/air-verse/air@v1.61.7 \
    && go install github.com/a-h/templ/cmd/templ@v0.3.898

# ============
# Dev stage
# ============
FROM base AS dev
WORKDIR /app
COPY . .
EXPOSE 8000

# Use air in dev for live reload
CMD ["air", "-c", ".air.toml"]

# ============
# Build stage (for production)
# ============
FROM base AS build
WORKDIR /app

COPY . .

# Build frontend
RUN npm ci && npm run build

# Build backend
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ============
# Production stage
# ============
FROM debian:bullseye-slim AS prod
WORKDIR /app

# Copy binary and static assets from build stage
COPY --from=build /app/main .
COPY --from=build /app/index.html .
COPY --from=build /app/static ./static

# HTTPS certs for outbound requests
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

EXPOSE 8000
CMD ["./main"]
