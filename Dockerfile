FROM golang:1.21

WORKDIR /app
COPY aats ./
COPY .env ./


CMD ["/app/aats"]
