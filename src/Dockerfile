FROM golang:alpine

WORKDIR /go/src/app
#WORKDIR /app
COPY . /go/src/app

RUN apk add --no-cache \
        libc6-compat curl

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

CMD ["app"]