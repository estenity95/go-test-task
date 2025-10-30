FROM golang:1.25-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Сборка бинаря. Флаги -w -s уменьшают размер.
# У тебя сейчас только cmd/api. Если добавишь cmd/migrate — просто добавь второй go build (см. ниже).
ENV CGO_ENABLED=0
RUN go build -ldflags "-w -s" -o ./bin/api     ./cmd/api \
 && go build -ldflags "-w -s" -o ./bin/migrate ./cmd/migrate

EXPOSE 8080
CMD ["./bin/api"]