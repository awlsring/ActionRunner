FROM golang:1.19-alpine

WORKDIR /app

COPY ./ ./

RUN go mod download

RUN go build -o /action-runner

EXPOSE 7032

CMD [ "/action-runner" ]