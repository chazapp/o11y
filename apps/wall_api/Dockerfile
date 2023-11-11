FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o wall_api

FROM alpine:3.18.4
WORKDIR /app
COPY --from=builder /app /app
RUN adduser -D -H -h /app appuser
RUN chown -R appuser:appuser /app
EXPOSE 8080
USER appuser
ENTRYPOINT ["/app/wall_api"]

