# build stage
# FROM golang:alpine AS build
FROM golang:1.19-alpine AS build
ARG PASSWORD=secret
WORKDIR /go/src/app
COPY . .

# RUN go mod init eclaim-api
RUN apk add build-base
RUN go mod tidy
# RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN swag init
# RUN apk add --no-cache git
RUN GOOS=linux go build -tags musl -ldflags="-s -w  -X 'mcash-finance-console-core/configs.passphrase=$PASSWORD'" -o ./bin/api ./main.go

# final stage
FROM alpine:latest
# RUN apk add --no-cache git

WORKDIR /data

COPY --from=build /go/src/app/bin /go/bin
EXPOSE 3030
ENTRYPOINT ENV=DEV /go/bin/api
