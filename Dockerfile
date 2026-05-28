ARG GO_VERSION="1.26"
FROM golang:${GO_VERSION} AS builder
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . /src
RUN go build -o ./bin/monitor ./cmd/monitor

FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /src/bin/monitor /app/monitor
CMD ["/app/monitor"]
