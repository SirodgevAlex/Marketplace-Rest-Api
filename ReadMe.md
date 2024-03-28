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

4. Запрос для создания пользователя

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

5. Запрос для аутентификации

```bash
curl -i -X POST http://localhost:8080/login \
-H 'Content-Type: application/json' \
-d '{"Email": "sirodgev@yandex.ru", "Password": "Sneeeir1_"}'
```

получим в ответ такой результат

```bash
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 28 Mar 2024 19:05:49 GMT
Content-Length: 158

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InNpcm9kZ2V2dkB5YW5kZXgucnUiLCJleHAiOjE3MTE2NTMwNDl9.Fq53ihagDLfCUhvrGOMVYR1bnMfVx7UHD9wlAVZLnh8"}
```

6. Запрос для создания объявления

```bash
curl -i -X POST http://localhost:8080/create-ad \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InNpcm9kZ2VAeWFuZGV4LnJ1IiwiZXhwIjoxNzExNjU0Mjc0fQ.6ytapN2fThQq9QrLXa47G8hw725PvOiI9wRuSRW8PmI' \
-d '{"Title": "bmx", "Text": "Really nice bmx", "ImageURL" : "https://example.com/image.jpg", "price" : 1000}'
```
