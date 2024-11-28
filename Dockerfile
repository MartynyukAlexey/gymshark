FROM golang:1.23-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /build ./cmd

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:3.20 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /build /build

ENTRYPOINT ["/build"]