FROM golang:1.26-alpine AS build

LABEL org.opencontainers.image.source=https://github.com/danielskowronski/GeoGuessrWatchdog

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o /out/worker ./cmd/ggwd
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

FROM alpine:3.21

COPY --from=build /out/worker /usr/local/bin/ggwd-worker
COPY --from=build /out/api /usr/local/bin/ggwd-api

ENTRYPOINT ["/usr/local/bin/ggwd-worker"]
