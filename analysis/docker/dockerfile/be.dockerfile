FROM golang:latest

WORKDIR /Documents/learning-english/analysis
COPY ./analysis/handler .

WORKDIR /Documents/learning-english/analysis/handler
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

ENTRYPOINT ["./main"]