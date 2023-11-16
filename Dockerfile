FROM golang:1.21.4-alpine AS builder

WORKDIR /go/src/propmon

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o propmon /go/src/propmon/cmd/main.go

FROM alpine:3.18

COPY --from=builder /go/src/propmon/propmon /usr/bin/propmon

EXPOSE 9500

ENTRYPOINT ["/usr/bin/propmon"]