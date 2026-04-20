FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SIGNING_TOKEN
ARG JWT_SECRET
ARG API_KEY
ENV SIGNING_TOKEN=${SIGNING_TOKEN}
ENV JWT_SECRET=${JWT_SECRET}
ENV API_KEY=${API_KEY}

RUN go build -o workshop ./cmd

EXPOSE 8080

CMD ["./workshop"]
