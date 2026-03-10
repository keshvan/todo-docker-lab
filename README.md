## Запуск
1. Создать `.env` в корне проекта:
```env
POSTGRES_USER=todo
POSTGRES_PASSWORD=todo123
POSTGRES_DB=todo_db
POSTGRES_HOST=db # DNS имя в сети Docker
POSTGRES_PORT=5432
SERVER_PORT=8080
```
2. Поднять контейнеры:
```bash
docker-compose up -d
```
3. API доступно по:
### http://localhost:8080/api/v1/todos