# Golang Yandex Intensive Autumn 2024

## Описание проекта
Проект представляет веб-сервер, на котором можно посчитать арифметическое выражение, состоящее из операций: +, -, *, /.\
Также поддерживаются скобки для выставления приоритета.

## Варианты использования
1. Введено валидное выражение:
    ```bash
    curl --location localhost:8080/api/v1/calculate/ \
    --header 'Content-Type: application/json' \
    --data '{"expression": "(-2+4.4)*5"}' -i
    ```
    Результат:
    ```
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Sun, 22 Dec 2024 16:09:53 GMT
    Content-Length: 30

    {"result":12.000000000000002}
    ```
2. Отправление запроса, метод которого отличный от POST:
    ```bash
    curl --location localhost:8080/api/v1/calculate/ -i
    ```
    Результат:
    ```
    HTTP/1.1 405 Method Not Allowed
    Content-Type: application/json
    Date: Sun, 22 Dec 2024 15:48:57 GMT
    Content-Length: 31

    {"error":"Method Not Allowed"}
    ```
3. Отправление невалидного сообщения (деление на ноль, пропущена скобка, пропущена операция, неизвестный символ, пропущено число):
    + Деление на ноль:
        ```bash
        curl --location localhost:8080/api/v1/calculate/ \
        --header 'Content-Type: application/json' \
        --data '{"expression": "2.3 + 2 / 0"}' -i
        ```
        Результат:
        ```
        HTTP/1.1 400 Bad Request
        Content-Type: application/json
        Date: Sun, 22 Dec 2024 15:53:52 GMT
        Content-Length: 29

        {"error":"division by zero"}
        ```
    + Пропущена скобка:
        ```bash
        curl --location localhost:8080/api/v1/calculate/ \
        --header 'Content-Type: application/json' \
        --data '{"expression": "(2.3 + 2) - 2)"}' -i
        ```
        Результат:
        ```bash
        HTTP/1.1 400 Bad Request
        Content-Type: application/json
        Date: Sun, 22 Dec 2024 15:57:29 GMT
        Content-Length: 77

        {"error":"there is no opening parenthesis corresponding to the closing one"}
        ```
    + Пропущена операция:
        ```bash
        curl --location localhost:8080/api/v1/calculate/ \
        --header 'Content-Type: application/json' \
        --data '{"expression": "(2.3 + 2.5) 2"}' -i
        ```
        Результат:
        ```bash
        HTTP/1.1 400 Bad Request
        Content-Type: application/json
        Date: Sun, 22 Dec 2024 15:59:32 GMT
        Content-Length: 30

        {"error":"missing operation"}
        ```
    + Неизвестный символ:
        ```bash
        curl --location localhost:8080/api/v1/calculate/ \
        --header 'Content-Type: application/json' \
        --data '{"expression": "(2.3 ? 2.5) + 2"}' -i
        ```
        Результат:
        ```bash
        HTTP/1.1 400 Bad Request
        Content-Type: application/json
        Date: Sun, 22 Dec 2024 16:00:29 GMT
        Content-Length: 27

        {"error":"unknown symbol"}
        ```
    + Пропущено число:
        ```bash
        curl --location localhost:8080/api/v1/calculate/ \
        --header 'Content-Type: application/json' \
        --data '{"expression": "(2.3 + 2.5) + "}' -i
        ```
        Результат:
        ```bash
        HTTP/1.1 400 Bad Request
        Content-Type: application/json
        Date: Sun, 22 Dec 2024 16:01:24 GMT
        Content-Length: 27

        {"error":"missing number"}
        ```
4. В остальных случаях выдается ошибка 500 (Internal Server Error)

## Структура проекта
```bash
├── cmd
│   └── main.go
├── go.mod
├── internal
│   └── application
│       └── application.go
├── pkg
│   ├── calculation
│   │   ├── calculation.go
│   │   ├── calculation_test.go
│   │   └── errors.go
│   └── server
│       ├── consts.go
│       ├── server.go
│       └── server_test.go
└── README.md
``` 

## Запуск проекта
```bash
go run ./cmd
```
Также можно поменять порт, на котором запущен веб-сервер:
```bash
PORT=8000 go run ./cmd
```