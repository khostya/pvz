FROM golang:1.24.0 AS builder
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app

FROM scratch
COPY --from=builder /bin/app /app
CMD ["/app"]