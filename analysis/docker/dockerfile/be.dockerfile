FROM golang:latest

RUN curl -fsSL https://deb.nodesource.com/setup_current.x | bash -
RUN apt-get update && apt-get install -y nodejs
WORKDIR /Documents/learning-english/analysis/
COPY ./analysis/handler/ ./handler/
COPY ./enviroment/.env /Documents/learning-english/enviroment/
COPY ./analysis/handler/go.mod ./handler/
COPY ./analysis/handler/go.sum ./handler/
COPY ./analysis/middleware/ /Documents/learning-english/analysis/middleware/
COPY ./analysis/handler/sourcegraph-cody/ ./handler/

WORKDIR /Documents/learning-english/analysis/handler/
RUN go mod tidy

RUN npm install -g @sourcegraph/cody
ENV SRC_ENDPOINT="https://sourcegraph.com"
ENV SRC_ACCESS_TOKEN="TOKEN"
RUN go build -ldflags="-s -w" -o start_api .
ENTRYPOINT ["./start_api"]