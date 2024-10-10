FROM golang:1.22 as build

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod download

COPY ../../cmd/tgtime-api/api.go ./
COPY ../../api api
COPY ../../app app
COPY ../../pkg pkg
COPY ../../internal internal
COPY ../../.env /

RUN CGO_ENABLED=0 GOOS=linux go build -o /tgtime-api

FROM scratch

COPY --from=build /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
COPY --from=build /tgtime-api .
COPY --from=ghcr.io/tarampampam/curl:8.6.0 /bin/curl /bin/curl
COPY --from=build /.env .

ENV TZ Europe/Moscow

ENV EXPOSE_PORT=$EXPOSE_PORT

EXPOSE $EXPOSE_PORT

CMD ["/tgtime-api"]