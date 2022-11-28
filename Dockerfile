FROM golang:1.19-alpine

WORKDIR /app

COPY ./ ./

RUN go mod download

RUN apk add \
    openssh \
    sshpass \
    ansible \
    python3 

RUN go build -o /action-runner

ENV ANSIBLE_HOST_KEY_CHECKING false

EXPOSE 7032

CMD [ "/action-runner" ]