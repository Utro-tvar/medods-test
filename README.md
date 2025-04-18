# Тестовое задание для MEDODS

**Задание:**

[Написать часть сервиса аутентифика](https://medods.yonote.ru/share/a74f6d8d-1489-4b54-bd82-81af5bf50a03/doc/test-task-backdev-sCBrYs5n6e)

Для сборки и запуска:

1. Изменить .env и ./postgres/.env:<br />   .env:
   
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
   
   ./postgres/.env:
   
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

Сервис поддерживает 2 маршрута:<br />/generate?GUID={guid} - get запрос возвращает json с парой токенов<br />/refresh - post запрос, принимает json с парой токенов и возвращает такой же json с новыми токенами

