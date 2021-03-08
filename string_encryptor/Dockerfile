FROM golang:latest
WORKDIR /go/src/microservices/string_encryptor
COPY ./ ./
RUN go build -o main .
CMD ("./main")