FROM golang:1.23@sha256:8c10f21bec412f08f73aa7b97ca5ac5f28a39d8a88030ad8a339fd0a781d72b4 as buildgo

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build ./cmd/beansd

FROM alpine@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099

RUN mkdir /app
COPY --from=buildgo /app/beansd /app/

ENV BEANS_HTTP_PORT=8080

CMD ["/app/beansd"]
