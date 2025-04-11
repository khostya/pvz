FROM golang:1.24.0 AS modules
COPY ./../go.mod ./../go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:1.24.0 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY .. /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/docs

EXPOSE 2000
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /app/api /api
COPY --from=builder /bin/app /app
CMD ["/app"]