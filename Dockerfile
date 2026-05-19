FROM golang:1.26-alpine AS build

LABEL org.opencontainers.image.source=https://github.com/danielskowronski/GeoGuessrWatchdog

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o /out/app ./cmd/ggwd

FROM alpine:3.21

COPY --from=build /out/app /usr/local/bin/app

ENTRYPOINT ["/usr/local/bin/app"]
