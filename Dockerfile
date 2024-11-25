ARG GO_VERSION=""
FROM golang:${GO_VERSION}alpine AS builder
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o ./bin/monitor ./cmd/monitor

FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /src/bin/monitor /app/monitor
CMD ["/app/monitor"]
