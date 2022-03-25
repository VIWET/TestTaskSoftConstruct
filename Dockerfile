FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN mkdir /build

RUN go build -o build/main ./cmd/main/main.go 
CMD ./build/main