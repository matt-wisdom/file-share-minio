FROM golang:1.20

WORKDIR /app

COPY . .

COPY ../db /app/db

RUN go build -o main .

CMD ["./main"]
