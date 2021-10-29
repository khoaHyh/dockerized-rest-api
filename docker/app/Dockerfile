# Build environment
# -----------------
FROM golang:1.17-alpine as build-env
WORKDIR /dockerized-rest-api

RUN apk update && apk add --no-cache gcc musl-dev git 

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/app ./cmd/app \
    && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate

# Deployment environment
# ----------------------
FROM alpine
RUN apk update && apk add --no-cache bash mysql-client

COPY --from=build-env /dockerized-rest-api/bin/app /dockerized-rest-api/
COPY --from=build-env /dockerized-rest-api/bin/migrate /dockerized-rest-api/
COPY --from=build-env /dockerized-rest-api/migrations /dockerized-rest-api/migrations

COPY --from=build-env /dockerized-rest-api/docker/app/bin /usr/local/bin/dockerized-rest-api
RUN chmod +x /usr/local/bin/dockerized-rest-api/*

EXPOSE 8080
CMD ["/dockerized-rest-api/app"]
