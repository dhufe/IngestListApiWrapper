FROM docker.io/golang:1.24.6-alpine3.22 AS build
RUN apk update --no-cache -U \
  && apk add --no-cache gcc g++ git make build-base
WORKDIR /app 
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -o apiwrapper ./cmd/server/

FROM alpine:3.21 AS prod
RUN apk update --no-cache -U \
  && apk add --no-cache curl tzdata openjdk21-jre bash file
WORKDIR /app 
COPY --from=build /app/apiwrapper ./apiwrapper
COPY ./third /app/third
RUN chmod +x /app/third/bin/ingestlist
CMD ["./apiwrapper"]

EXPOSE 8080
