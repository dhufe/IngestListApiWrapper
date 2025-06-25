FROM golang:alpine3.21 AS build
WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o apiwrapper 

FROM alpine:3.21 AS prod
WORKDIR /borg
COPY --from=build /app/apiwrapper ./apiwrapper
CMD ["./apiwrapper"]
