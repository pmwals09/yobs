FROM golang:1.21
WORKDIR /app
COPY aats ./
ENTRYPOINT /app/aats
