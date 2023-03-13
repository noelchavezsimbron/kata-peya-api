FROM golang:1.18 as builder
ENV GIT_TERMINAL_PROMPT=1
ENV CGO_ENABLED=0
ARG SSH_PRIVATE_KEY
RUN mkdir /root/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa && chmod 400 /root/.ssh/id_rsa
RUN touch /root/.ssh/known_hosts

WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN go build -o ./build/server cmd/server/main.go


FROM alpine:latest
RUN apk update && apk add ca-certificates tzdata && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /app/build .
CMD ["./server"]