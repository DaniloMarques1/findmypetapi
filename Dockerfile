FROM golang:1.16

WORKDIR /usr/home/findmypet

COPY . .

RUN go mod download

RUN go build main.go

CMD ["./main"]
