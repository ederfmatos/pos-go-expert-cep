FROM golang:1.23-alpine AS builder
WORKDIR /build
COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -ldflags '-s -w -extldflags' -o app .

FROM scratch
COPY --from=builder /build/app ./app
EXPOSE 8080
ENTRYPOINT ["/app"]