FROM golang:1.22 as build

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod download

COPY ../../cmd/officetime-api/main.go ./
COPY ../../app app

RUN CGO_ENABLED=0 GOOS=linux go build -o /officetime-api

FROM scratch

COPY --from=build /officetime-api .
COPY --from=ghcr.io/tarampampam/curl:8.6.0 /bin/curl /bin/curl

EXPOSE 8080

CMD ["/officetime-api", "--config", "config"]