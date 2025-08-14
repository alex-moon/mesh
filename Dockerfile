FROM golang:1.23-alpine

# Install git for go modules that might need it
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod* go.sum* ./

# Download dependencies
RUN go mod download

# Install air and templ for local
RUN go install github.com/air-verse/air@v1.61.7 \
    && go install github.com/a-h/templ/cmd/templ@v0.3.898

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Install ca-certificates for HTTPS requests and basic tools
RUN apk --no-cache add ca-certificates curl

# Create app directory
WORKDIR /app

# Expose port
EXPOSE 8000

# Run the application
CMD ["./main"]
