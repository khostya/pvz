#!/usr/bin/env bash

set -eumxo pipefail

DEFAULT_PG_URL="postgres://user:password@localhost:5432"
PG_URL=$DEFAULT_PG_URL

go test ./internal/... ./pkg/... -coverprofile unit.coverage.txt
TEST_DATABASE_URL=$PG_URL/pvz go test ./tests/postgres/... -tags=integration -coverprofile integration.coverage.txt -coverpkg=./...
go test ./tests/http/... -tags=integration.http -coverprofile integration.http.coverage.txt -coverpkg=./...
go test ./tests/grpc/... -tags=integration.grpc -coverprofile integration.grpc.coverage.txt -coverpkg=./...
echo "mode: set" > coverage.txt && cat *.coverage.txt | grep -v mode: | sort -r | awk '{if($1 != last) {print $0;last=$1}}' >> coverage.txt