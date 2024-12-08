# Тестовое задание для MEDODS

**Задание:**

[Написать часть сервиса аутентифика](https://medods.notion.site/Test-task-BackDev-623508ed85474f48a721e43ab00e9916)

Для сборки и запуска:

1. Изменить .env и ./postgres/.env:
   .env:\
   
   ```dotenv
   #tokens parameters
   ACCESS_TTL=120 #in minutes
   REFRESH_TTL=30 #in days
   SECRET_KEY=secret
   
   #database parameters
   DB_USER=service
   DB_PASS=password
   DB_NAME=service
   ```
   
   ./postgres/.env:\
   
   ```dotenv
   POSTGRES_USER=postgres
   POSTGRES_DB=postgres
   POSTGRES_PASSWORD=postgres
   
   #database parameters
   DB_USER=service
   DB_PASS=password
   DB_NAME=service
   ```
2. Выполнить
   
   ```bash
   docker compose -f "docker-compose.yml" up -d --build
   ```

Сервис поддерживает 2 маршрута:
/generate/{guid} - get запрос возвращает json с парой токенов\
/refresh - post запрос, принимает json с парой токенов и возвращает такой же json с новыми токенами

