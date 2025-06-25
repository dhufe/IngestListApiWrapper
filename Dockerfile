FROM golang:alpine3.21 AS build
apk update --no-cache -U \
  && apk add --no-cache curl
WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o apiwrapper 

FROM alpine:3.21 AS prod
apk update --no-cache -U \
  && apk add --no-cache curl
WORKDIR /app 
COPY --from=build /app/apiwrapper ./apiwrapper
CMD ["./apiwrapper"]
