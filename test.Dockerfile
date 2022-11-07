FROM golang:1.19-alpine AS build

# Fetch dependencies
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Build
COPY . .
ENV CGO_ENABLED=0
CMD ["go", "test", "-count=1", "-v", "./..."]
