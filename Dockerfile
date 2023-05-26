FROM golang:1.20-alpine

RUN mkdir /app
WORKDIR /app

COPY . /app/
RUN apk update && apk add --no-cache git

RUN go build -o main .

CMD ["go", "run", "."]

EXPOSE 9090
