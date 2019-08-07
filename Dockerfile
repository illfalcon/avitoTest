FROM golang:1.12
RUN mkdir -p /go/src/github.com/illfalcon/avitoTest
COPY . /go/src/github.com/illfalcon/avitoTest
WORKDIR /go/src/github.com/illfalcon/avitoTest
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go build cmd/main.go cmd/handlers.go
EXPOSE 9000
CMD ["./main"]