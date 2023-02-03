FROM golang:1.18.10 as builder

RUN mkdir /build
ADD *.go *.mod *.sum /build/
WORKDIR /build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o stocks-sync .

FROM alpine:3.11.13

RUN mkdir /app
RUN mkdir /app/bin

COPY --from=builder /build/stocks-sync /app/bin/
WORKDIR /app
ENV PORT=${PORT}
EXPOSE ${PORT}
CMD ["/app/bin/stocks-sync"]
