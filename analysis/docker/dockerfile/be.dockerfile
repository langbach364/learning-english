FROM golang:latest

WORKDIR /Documents/learning-english/analysis
COPY ./handler .

WORKDIR /Documents/learning-english/analysis/handler
RUN go mod tidy

RUN go build -o main .

ENTRYPOINT ["./main"]