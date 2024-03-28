Привет!
Это инструкция

1) Запускаем db.sql

```pgsql
psql db.sql
```

у себя я делал не так, ибо у меня там pgAdmin

```pgsql
psql -U postgres -d db -f db.sql -W 
```

postgres - название пользователя в pgAdmin, у меня потом пароль попросило от пользователя

2. Запускаем main.go

```go
go run main.go
```

3. Все, мы все запустили, можно делать сами запросы

4.  Запрос для создания пользователя

```bash
curl -i -X POST http://localhost:8080/register \
-H 'Content-Type: application/json' \
-d '{"Email": "sirodgev@yandex.ru", "Password": "Sneeeir1_"}'
```
получим в ответ такой результат

```bash
HTTP/1.1 201 Created
Date: Thu, 28 Mar 2024 11:44:53 GMT
Content-Length: 60
Content-Type: text/plain; charset=utf-8

{"Id":4,"Email":"sirodgev@yandex.ru","Password":"Sneeeir1_"}
```


curl -i -X POST http://localhost:8080/login \
-H 'Content-Type: application/json' \
-d '{"Email": "sirodgev@yandex.ru", "Password": "Sneeeir1_"}'