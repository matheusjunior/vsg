FROM golang:1.19-buster AS build

# Fetch dependencies
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 go build -o /app/app ./cmd/main.go

# Copy app
FROM gcr.io/distroless/static-debian11
COPY --from=build /app/app /
ENV AWS_ACCESS_KEY_ID=testUser
ENV AWS_SECRET_ACCESS_KEY=testAccessKey
ENV AWS_DEFAULT_REGION=us-east-1
CMD ["/app"]