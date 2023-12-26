FROM golang:1.21 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main main.go

FROM alpine:3.15.0
WORKDIR /

COPY --from=builder /build/main /main

CMD [ "/main" ]
