FROM golang:1.24-alpine AS builder

# Variables d'environnement pour Go compilées pour ARMv7 (Raspberry Pi 2)
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm \
    GOARM=7

WORKDIR /app

# Copie les fichiers go.mod et go.sum pour gérer le cache des modules
COPY go.mod go.sum ./

RUN go mod download

# Copie tout le code source
COPY . .

# Compile le binaire pour ARMv7
RUN go build -o chess-bot ./main.go

# ------------ ÉTAPE 2 : IMAGE FINALE ------------
FROM alpine:latest

# Ajout d'un utilisateur non-root
RUN adduser -D -g '' botuser

WORKDIR /app

# Copie du binaire et du fichier JSON
COPY --from=builder /app/chess-bot /app/chess-bot
COPY tracked_users.json /app/tracked_users.json

RUN chown -R botuser:botuser /app
USER botuser

CMD ["/app/chess-bot"]

