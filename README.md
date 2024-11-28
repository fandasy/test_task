Test Task: Music Library API
---

### Environment file

To start the project, use the template config.env file

(the file name must be "config")

Don't forget to change the data for connecting to the database and YOUR_API_HOST

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

# your api host (http prefix must not be used)
# example: localhost:1234 or yourApiHost.com/api/v5 or yourApiHost.com
YOUR_API_HOST=example.com
```

### Swagger

Swagger files are located in the ./docs directory

```
  /docs
    docs.go
    swagger.json
    swagger.yaml
```

Swagger GUI path: localhost:8080/swagger/index.html

### Important points

1. Implementation of the request to the third-party API, located at /internal/clients/your-api
2. The implementation of handlers is located at the path ./internal/http-server/handlers
3. Implementation of all database logic is located at the path ./internal/storage/psql
4. When requesting an [PATCH] update, check the Link and Release date fields, as they pass validation, the time should be in DD.MM.YYYY format and the link should be valid

