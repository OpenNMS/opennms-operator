#Builder
FROM golang:1.17 as builder

WORKDIR /app

COPY Makefile Makefile

COPY go.mod go.mod
COPY go.sum go.sum

RUN make dependencies

COPY cmd/ cmd/
COPY charts/dependencies charts/dependencies
COPY api/ api/
COPY config/ config/
COPY internal/ internal/

RUN make alpine-build

#Runner
FROM alpine:3.15

WORKDIR /app

COPY --from=builder /app/operator .
COPY charts/ charts/
COPY config/ config/

ENTRYPOINT ["/app/operator"]