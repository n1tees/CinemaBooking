# Минимальный рантайм-образ
FROM debian:bookworm-slim

WORKDIR /app

# Копируем бинарник и .env (если нужно)
COPY bookingkart-platform .

EXPOSE 8080

ENTRYPOINT ["./cinema-booking"]
