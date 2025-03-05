# Golang Yandex Intensive Autumn 2024

## Описание проекта
Проект представляет веб-сервер, на котором можно параллельно посчитать арифметическое выражение, состоящее из операций: +, -, *, /.\
Также поддерживаются скобки для выставления приоритета.

## Принцип работы
Проект состоит из двух сущностей: Оркестратор и Агент.
Запрос с арифметическим выражением поступает на Оркестратора, после чего преобразуется в RPN, по которой строится дерево, а также заполняется список задач. Каждая задача представляет собой следующую структуру: операция, аргумент1, аргумент2. Вторая сущность — Агент, постоянно опрашивает Оркестратора на наличие задач. После получения, задача отправляется в пул с воркерами, которые расчитывают результат и делаю POST запрос на Оркестратора. После чего Оркестратор обновляет дерево зависимостей и после обновления корня дерева, записывает в PostgreSQL результат выражения.

В данном проекте используется две базы данных:
1) Memgraph — графовая БД в оперативной памяти, которая служит для хранения узлов дерева.
2) PostgreSQL — реляционная БД, служащая для хранения всех выражений пользователя.

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
    Date: Wed, 05 Mar 2025 17:21:08 GMT
    Content-Length: 46

    {"id":"f2472d01-5b85-40bd-bc6e-a6c3f8bb2a86"}
    ```
2. Отправление запроса, метод которого отличный от POST:
    ```bash
    curl --location localhost:8080/api/v1/calculate/ -i
    ```
    Результат:
    ```
    HTTP/1.1 405 Method Not Allowed
    Content-Type: application/json
    Date: Wed, 05 Mar 2025 17:21:45 GMT
    Content-Length: 31

    {"error":"Method Not Allowed"}
    ```
3. Отправление выражения содержащее синтаксические ошибки (пропущена скобка, пропущена операция, неизвестный символ, пропущено число):
    + Пропущена скобка:
        ```bash
        curl --location localhost:8080/api/v1/calculate/ \
        --header 'Content-Type: application/json' \
        --data '{"expression": "(2.3 + 2) - 2)"}' -i
        ```
        Результат:
        ```bash
        HTTP/1.1 422 Unprocessable Entity
        Content-Type: application/json
        Date: Wed, 05 Mar 2025 17:23:00 GMT
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
        HTTP/1.1 422 Unprocessable Entity
        Content-Type: application/json
        Date: Wed, 05 Mar 2025 17:23:35 GMT
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
        HTTP/1.1 422 Unprocessable Entity
        Content-Type: application/json
        Date: Wed, 05 Mar 2025 17:24:03 GMT
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
        HTTP/1.1 422 Unprocessable Entity
        Content-Type: application/json
        Date: Wed, 05 Mar 2025 17:24:26 GMT
        Content-Length: 27

        {"error":"missing number"}
        ```
4. Просмотр всех выражений, которые ввел пользователь:
    ```bash
    curl --location localhost:8000/api/v1/expressions -i
    ```
    Результат:
    ```bash
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Wed, 05 Mar 2025 17:25:47 GMT
    Content-Length: 97

    {"expressions":[{"id":"f2472d01-5b85-40bd-bc6e-a6c3f8bb2a86","status":"completed","result":12}]}
    ```
5. Просмотр выражения, по id:
    ```bash
    curl --location localhost:8000/api/v1/expressions/f2472d01-5b85-40bd-bc6e-a6c3f8bb2a86 -i
    ```
    Результат:
    ```bash
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Wed, 05 Mar 2025 17:27:35 GMT
    Content-Length: 79

    {"id":"f2472d01-5b85-40bd-bc6e-a6c3f8bb2a86","status":"completed","result":12}
    ```
6. При введение выражения, в котором производится деление на 0, будет статус выражения failed.

## Структура проекта
```bash
├── cmd/
│   ├── agent/
│   │   └── main.go
│   └── orchestrator/
│       └── main.go
├── config/
├── internal/
│   ├── consts/
│   ├── services/
│   │   ├── agent/
│   │   │   ├── listener/
│   │   │   │   └── listener.go
│   │   │   ├── logger/
│   │   │   └── worker/
│   │   │       ├── pool.go
│   │   │       └── worker.go
│   │   └── orchestrator/
│   │       ├── app/
│   │       │   └── app.go
│   │       ├── calculation/
│   │       │   ├── parser.go
│   │       │   ├── task_manager.go
│   │       │   └── tree.go
│   │       ├── handlers/
│   │       │   ├── calculation.go
│   │       │   ├── expression.go
│   │       │   └── task.go
│   │       ├── logger/
│   │       └── storage/
│   │           ├── memgraph.go
│   │           └── postgres.go
│   └── shared/
│       ├── models/
│       └── utils/
├── pkg
│   └── utils
│       └── config.go
├── docker-compose.yml
└── README.md
``` 

## Запуск проекта
Перед запуском проекта необходимо создать .env файл. Дефолтные настройки есть в .env.example.

1) Запустить БД (можно с флагом -d):
```bash
docker compose up --build
```
2) Запустить Оркестратора (отдельный терминал)
```bash
go run ./cmd/orchestrator/main.go
```
3) Запустить Агента (отдельный терминал)
```bash
go run ./cmd/agent/main.go
```