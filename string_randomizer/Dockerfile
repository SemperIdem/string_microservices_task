FROM golang:latest
WORKDIR /go/src/string_microservices_task/string_randomizer
COPY ./ ./
RUN go build -o main .
CMD ("./main")