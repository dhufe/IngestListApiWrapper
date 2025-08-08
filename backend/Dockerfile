FROM golang:1.24.6-alpine3.22 AS build
RUN apk update --no-cache -U 
WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o apiwrapper ./cmd/server/

FROM alpine:3.21 AS prod
RUN apk update --no-cache -U \
  && apk add --no-cache curl tzdata
WORKDIR /app 
COPY --from=build /app/apiwrapper ./apiwrapper
COPY ./third /app/
# COPY ./config/config.yaml /app
CMD ["./apiwrapper"]

EXPOSE 8080
