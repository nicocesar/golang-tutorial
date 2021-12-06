FROM golang:1.17.4-bullseye

COPY . /app
WORKDIR /app
RUN go mod tidy

CMD ./golang-tutorial
