FROM golang:1.18 AS builder
WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /work/solaredge cmd/*.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /work/solaredge /
ENTRYPOINT ["/solaredge"]