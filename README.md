# CinemaBooking

В данном проекте реализована Back-end часть веб-приложения по покупке билетов в кино
## 🌐 Возможности системы

- Регистрация и авторизация пользователей
- Просмотр афиш, описания фильмов, их сеансов их показа
- Пополнение баланса
- Покупка билетов
- Бонусная система оплаты билетов 
- Возмность оставить отзыв на просмотренные фильмы
---

## 🛠 Backend (Golang + Gin)

### 📦 Технологии

- **Go (Golang)** — основной язык
- **Gin** — фреймворк для создания REST API
- **GORM** — ORM для взаимодействия с PostgreSQL
- **PostgreSQL** — база данных
- **JWT** — для авторизации пользователей
- **Docker** — контейнеризация
- **Swagger** — встроенная документация API
- **GitHub Actions** — автоматическая сборка и деплой

### 📁 Структура проекта

BookingKart-Platform/  
├── .github/workflows # CI-процессы  
├── cmd/ # main.go и запуск сервера  
├── config/ # настройки, инициализация  
├── docs/ # Swagger-документация  
├── nginx/conf.d/ # конфигурация nginx  
├── pkg/ # бизнес-логика, роутеры, сервисы  
├── Dockerfile # билд сервиса  
├── docker-compose.yml # сборка и запуск всего приложения  
├── go.mod, go.sum # зависимости  

### 🚀 Быстрый старт (локально)

1. **Клонировать репозиторий:**  
--bash  
git clone https://github.com/n1tees/CinemaBooking
cd CinemaBooking 

2. **Создать .env файл с переменными окружения:**  
--dotenv  
DB_HOST=localhost  
DB_PORT=5432  
DB_USER=postgres  
DB_PASSWORD=your_password  
DB_NAME=bookingkart  
JWT_SECRET=your_jwt_secret  

3. Запуск в Docker:  
--bash  
docker-compose up --build  

4. Swagger UI доступен по адресу:  
--bash  
http://<IP>:8080/swagger/index.html

4. Others
go run ./cmd/main.go
swag init -g cmd/main.go -o ./docs --parseDependency --parseInternal  

