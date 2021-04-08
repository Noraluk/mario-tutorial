# builder image
FROM golang:1.16.3-alpine3.12 as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o server .


# generate clean, final image for end users
FROM alpine:3.12

WORKDIR /app

COPY --from=builder /build/server /app/server
COPY static static

EXPOSE 8080
# executable
ENTRYPOINT [ "./server" ]