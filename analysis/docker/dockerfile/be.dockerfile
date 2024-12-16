FROM golang:latest

WORKDIR /Documents/learning-english/analysis
COPY ./analysis/handler .

WORKDIR /Documents/learning-english/analysis/handler
RUN go mod tidy

RUN go build -ldflags="-s -w" -o ./start_api *.go

ENTRYPOINT ["./start_api"]