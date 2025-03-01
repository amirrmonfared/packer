FROM golang:1.23.4 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o packer ./cmd/packer/main.go

FROM alpine:3.18
COPY --from=builder /app/packer /usr/bin/packer

ENTRYPOINT ["/usr/bin/packer"]
