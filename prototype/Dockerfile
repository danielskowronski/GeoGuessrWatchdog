FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN CGO_ENABLED=0 go build -o /out/fetcher ./main.go

FROM alpine:3.21

COPY --from=build /out/fetcher /usr/local/bin/fetcher

ENTRYPOINT ["/usr/local/bin/fetcher"]