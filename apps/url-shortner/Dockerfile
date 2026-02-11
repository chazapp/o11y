FROM golang:1.26 AS builder
ARG VERSION
RUN test -n "$VERSION"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -ldflags="-X main.Version=$VERSION" -o url-shortner cmd/main.go

FROM alpine:3.22.1
WORKDIR /app
COPY --from=builder /app/url-shortner /app/url-shortner
RUN adduser -D -H -h /app appuser
RUN chown -R appuser:appuser /app
EXPOSE 8080
USER appuser
ENTRYPOINT ["/app/url-shortner"]

