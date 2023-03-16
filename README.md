# Тестовое задание на позицию Middle/Senior Go developer в компанию Алиф

## Задание
Внедрите Rest API для финансового учреждения, где он предоставляет своим партнёрам услуги электронного кошелька. У него есть два типа учетных записей электронного кошелька: идентифицированные и неидентифицированные. API может поддерживать несколько клиентов, и следует использовать только методы http, post с json в качестве формата данных. Клиенты должны быть аутентифицированы через http параметр заголовок X-UserId и X-Digest. X-Digest — это hmac-sha1, хэш-сумма тела запроса. Должны быть предварительно записанные электронные кошельки, с разными балансами, а максимальный баланс составляет 10.000 сомони для неидентифицированных счетов и 100.000 сомони для идентифицированных счетов. Для хранения данных можете использовать по вашему выбору. API методы сервиса: 1. Проверить существует ли аккаунт электронного кошелька. 2. Пополнение электронного кошелька. 3. Получить общее количество и суммы операций пополнения за текущий месяц. 4. Получить баланс электронного кошелька. Во время разработки используйте git и Github и делайте значимые коммиты.

## Stack
Golang, PostgreSQL, Redis

## Инструкция запуска
Postgres и Redis запускаются в Docker, поэтому предполагается, что он установлен на машине.

Для сборки выполните команду `make`
``` bash
$ make
...
```

Для запуска выполните команду
```
./wallet
```

# Endpoints
## Аутентификация
`POST /login`

Отправить
```json
{
    "email": "email"
}
```

В случае успеха, отправляет 200 OK
```
HTTP/1.1 200 OK
Content-Length: 0
Content-Type: application/json; charset=utf-8
Date: Thu, 16 Mar 2023 08:52:21 GMT
X-Userid: <X-Userid>
```

## Кошелёк
### Проверить существует ли аккаунт электронного кошелька:
```
HEAD http://localhost:8080/alif/wallets/{id}
X-Userid: <X-Userid>
```

### Получить текущий баланс кошелька
```
GET http://localhost:8080/alif/wallets/{id}
X-Userid: <X-Userid>
```

Ответ
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 16 Mar 2023 09:00:04 GMT
Content-Length: 41
Connection: close

{
  "is_identified": bool,
  "balance": int (diram)
}
```

### Получить количество пополнений и общую сумму пополнений
```
GET http://localhost:8080/alif/wallets/{id}/replenishments/{month}
X-Userid: <X-Userid>
```

Ответ
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 16 Mar 2023 09:03:14 GMT
Content-Length: 49
Connection: close

{
  "total_replenishments": int,
  "total_amount": int (diram)
}
```

### Пополнить кошелёк
```
PUT http://localhost:8080/alif/wallets/{id}
X-Userid: <X-Userid>
X-Digest: <X-Digest>

{
    "balance": string
}
```

X-Digest — HMAC-SHA1 сумма тела запроса в шестнадцатеричной кодировке.