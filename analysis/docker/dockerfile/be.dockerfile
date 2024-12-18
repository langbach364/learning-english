FROM golang:latest

WORKDIR /Documents/learning-english/analysis/
COPY ./analysis/handler/ ./handler/

WORKDIR /Documents/learning-english/enviroment/
COPY ./enviroment/.env ./

WORKDIR /Documents/learning-english/analysis/handler/
COPY ./analysis/handler/go.mod ./
COPY ./analysis/handler/go.sum ./

WORKDIR /Documents/learning-english/analysis/middleware/
COPY ./analysis/middleware/ ./


WORKDIR /Documents/learning-english/analysis/handler/
RUN go mod tidy
RUN go build -ldflags="-s -w" -o ./start_api .

ENTRYPOINT ["./start_api"]