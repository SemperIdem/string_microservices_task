FROM golang:latest
WORKDIR /go/src/string_microservices_task/string_encryptor
COPY ./ ./
RUN go build -o main .
CMD ("./main")