FROM golang:1.17
RUN apt update && apt install -y nodejs npm
RUN npm install -g newman
WORKDIR /
ADD . /
RUN go build -o /usr/local/bin/agent agents/postman-agent/main.go 
