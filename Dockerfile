FROM golang:1.20.0-bullseye

ENV PORT=${PORT}

RUN mkdir /app
RUN mkdir /app/bin

ADD . /app
WORKDIR /app

RUN go mod download

RUN go build -o bin/stocks-sync .

CMD ["/app/bin/stocks-sync"]
