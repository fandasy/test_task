Test Task: Music Library API
---

### Environment file

Чтобы запустить проект, используйте шаблонный файл config.env

(имя файла должно быть «config»)

Не забудьте изменить данные для подключения к базе данных и YOUR_API_HOST

```
# local - level: Debug, Type: Text
# dev   - level: Debug, Type: Json
# prod  - level: Info,  Type: Json
SLOG=local

# address where the service is launched
ADDR=localhost:8080

# database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=yourDB
DB_USER=yourUser
DB_PASSWORD=yourPassword

# your api host (не нужно ставить префикс http)
# example: localhost:1234 or yourApiHost.com/api/v5 or yourApiHost.com
YOUR_API_HOST=example.com
```

### Swagger

Файлы Swagger находятся в каталоге ./docs.

```
  /docs
    docs.go
    swagger.json
    swagger.yaml
```

Swagger GUI path: [GET] localhost:8080/swagger/index.html

Если вы используете другой графический интерфейс Swagger, обязательно измените адрес и порт на свои значения

### Important points

1. Реализация запроса к стороннему API, расположенному по адресу /internal/clients/your-api
2. Реализация handlers находится по пути ./internal/http-server/handlers
3. Реализация всей логики базы данных находится по пути ./internal/storage/psql
4. При запросе обновления [PATCH] проверьте поля Link и Release date, так как они проходят валидацию, время должно быть в формате DD.MM.YYYY, а ссылка должна быть действительной

