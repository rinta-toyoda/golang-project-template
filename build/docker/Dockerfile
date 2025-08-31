FROM golang:1.24.1

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz \
  -o migrate.tar.gz && \
  tar -xzf migrate.tar.gz && \
  mv migrate /usr/local/bin/migrate && \
  chmod +x /usr/local/bin/migrate && \
  rm migrate.tar.gz

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
