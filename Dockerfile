# ------------ ÉTAPE 1 : BUILD ------------
FROM golang:1.22-alpine AS builder

# Variables d'environnement pour Go
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Crée un répertoire de travail
WORKDIR /app

# Copie les fichiers go.mod et go.sum pour cacher le cache de modules
COPY go.mod go.sum ./

# Télécharge les dépendances
RUN go mod download

# Copie tout le code source
COPY . .

# Compile le binaire
RUN go build -o chess-bot ./main.go

# ------------ ÉTAPE 2 : IMAGE FINALE ------------
FROM alpine:latest

# Ajoute un utilisateur non-root pour plus de sécurité
RUN adduser -D -g '' botuser

WORKDIR /app

# Copie uniquement le binaire et les fichiers nécessaires
COPY --from=builder /app/chess-bot /app/chess-bot
COPY config.json /app/config.json

# Assure-toi que les permissions sont correctes
RUN chown -R botuser:botuser /app
USER botuser

# Commande par défaut
CMD ["/app/chess-bot"]
