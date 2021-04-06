# builder image
FROM golang:1.16.3-alpine3.12 as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o mario .


# generate clean, final image for end users
FROM alpine:3.12
COPY --from=builder /build/mario .

EXPOSE 8080
# executable
ENTRYPOINT [ "./mario" ]