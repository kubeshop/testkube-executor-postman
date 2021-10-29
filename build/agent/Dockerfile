# syntax=docker/dockerfile:1
FROM golang:1.17
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0 
ENV GOOS=linux

RUN cd cmd/agent;go build -o /app -mod mod -a .

FROM postman/newman
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app /bin/app
EXPOSE 8082
ENTRYPOINT ["/bin/app"]
