FROM golang:latest

WORKDIR /Documents/learning-english/analysis
COPY ./analysis/handler .

WORKDIR /Documents/learning-english/analysis/handler
RUN go mod tidy

ENTRYPOINT ["./start_api"]