FROM docker.io/golang:1.25.5-alpine3.22 AS build

ARG PROXY_HOST="fire.lvnbb.de"
ARG PROXY_PORT="80"

RUN apk update --no-cache -U \
  && apk add --no-cache gcc g++ git make build-base
WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o apiwrapper ./cmd/server/

FROM alpine:3.22 AS prod

ARG PROXY_HOST="fire.lvnbb.de"
ARG PROXY_PORT="80"

RUN apk update --no-cache -U \
  && apk add --no-cache curl tzdata bash file openjdk21-jre
WORKDIR /app 
COPY --from=build /app/apiwrapper ./apiwrapper

COPY ./third /app/third
RUN echo "PROXY_HOST=$PROXY_HOST" >> /etc/environment \
    && echo "PROXY_PORT=$PROXY_PORT" >> /etc/environment \
    && sed -i \
    's|exec "\$JAVACMD" "\$@"|exec "\$JAVACMD" \
        -Dhttp.proxyHost='$PROXY_HOST' \
        -Dhttp.proxyPort='$PROXY_PORT' \
        -Dhttps.proxyHost='$PROXY_HOST' \
        -Dhttps.proxyPort='$PROXY_PORT' \
        "\$@"|' \
        /app/third/bin/ingestlist \
    && chmod +x /app/third/bin/ingestlist
CMD ["./apiwrapper"]

EXPOSE 8080