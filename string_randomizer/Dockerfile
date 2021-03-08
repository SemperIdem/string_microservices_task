FROM golang:latest
WORKDIR /go/src/microservices/string_randomizer
COPY ./ ./
RUN go build -o main .
CMD ("./main")