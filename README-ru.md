# Shifty's Facade

[![Go Version](https://img.shields.io/github/go-mod/go-version/ShiftyX1/Facade)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/ShiftyX1/Facade)](https://github.com/ShiftyX1/Facade/releases)
[![License](https://img.shields.io/github/license/ShiftyX1/Facade)](https://github.com/ShiftyX1/Facade/blob/master/LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/ShiftyX1/Facade/ci.yml?branch=master)](https://github.com/ShiftyX1/Facade/actions)
[![GitHub Issues](https://img.shields.io/github/issues/ShiftyX1/Facade)](https://github.com/ShiftyX1/Facade/issues)
[![GitHub Stars](https://img.shields.io/github/stars/ShiftyX1/Facade)](https://github.com/ShiftyX1/Facade/stargazers)
[![Last Commit](https://img.shields.io/github/last-commit/ShiftyX1/Facade)](https://github.com/ShiftyX1/Facade/commits/master)

Простой и гибкий mock HTTP сервер для разработчиков.

## Установка

```bash
make build
```

## Быстрый старт

1. Создайте конфигурационный файл `config.yaml`:

```yaml
server:
  port: 8080
  host: localhost

routes:
  - path: /users
    method: GET
    response:
      status: 200
      schema:
        type: array
        items:
          type: object
          properties:
            id: {type: integer, min: 1, max: 1000}
            name: {type: string, faker: name}
            email: {type: string, faker: email}
        count: 5

global:
  cors: true
  log_requests: true
```

2. Запустите сервер:

```bash
make run
```

Или собрав и запустив вручную:
```bash
make build
./build/facade -c configs/example.yaml
```

3. Тестируйте API:

```bash
curl http://localhost:8080/users
```

## Параметры командной строки

- `-c, --config` - путь к конфигурационному файлу (по умолчанию: configs/example.yaml)
- `-p, --port` - порт сервера (переопределяет конфиг)
- `--host` - хост сервера (переопределяет конфиг)
- `-v, --verbose` - подробное логирование

## Конфигурация

### Базовая структура

```yaml
server:
  port: 8080
  host: localhost

routes:
  - path: /api/endpoint
    method: GET
    response:
      status: 200
      body: '{"message": "Hello World"}'

global:
  cors: true
  delay: 100ms
  log_requests: true
```

### Типы ответов

#### Статичный ответ

```yaml
- path: /static
  method: GET
  response:
    status: 200
    body: '{"message": "Static response"}'
```

#### Шаблонизированный ответ

```yaml
- path: /users/{id}
  method: GET
  response:
    status: 200
    body: '{"id": {{.PathParams.id}}, "name": "User {{.PathParams.id}}"}'
```

#### Генерация по схеме

```yaml
- path: /users
  method: GET
  response:
    status: 200
    schema:
      type: array
      items:
        type: object
        properties:
          id: {type: integer, min: 1, max: 1000}
          name: {type: string, faker: name}
          email: {type: string, faker: email}
      count: 5
```

### Условные ответы

```yaml
- path: /data
  method: GET
  response:
    status: 200
    conditions:
      - if: '{{.QueryParams.format}} == "xml"'
        body: '<data>XML response</data>'
      - if: '{{.QueryParams.format}} == "json"'
        body: '{"data": "JSON response"}'
      - default: true
        body: '{"data": "Default response"}'
```

### Шаблонные функции

#### Доступ к данным запроса

- `{{.PathParams.id}}` - path параметры
- `{{.QueryParams.limit}}` - query параметры  
- `{{.Body.name}}` - данные из body запроса

#### Faker данные

- `{{fakerName}}` - случайное имя
- `{{fakerEmail}}` - случайный email
- `{{fakerPhone}}` - случайный телефон
- `{{fakerCompany}}` - случайная компания
- `{{fakerUUID}}` - UUID

#### Случайные данные

- `{{randomInt 1 100}}` - случайное число
- `{{randomString 10}}` - случайная строка
- `{{randomBool}}` - случайный boolean

#### Дата и время

- `{{now}}` - текущее время
- `{{now.Format "2006-01-02"}}` - форматированная дата

### Типы схем

#### Примитивные типы

```yaml
properties:
  name: {type: string, faker: name}
  age: {type: integer, min: 18, max: 65}
  price: {type: number, min: 0.1, max: 999.99}
  active: {type: boolean}
```

#### Массивы

```yaml
schema:
  type: array
  items:
    type: object
    properties:
      id: {type: integer}
      name: {type: string}
  count: 10
```

#### Объекты

```yaml
schema:
  type: object
  properties:
    user:
      type: object
      properties:
        id: {type: integer}
        profile:
          type: object
          properties:
            name: {type: string, faker: name}
```

### Faker типы

- `name` - имя человека
- `email` - email адрес
- `phone` - номер телефона
- `address` - адрес
- `company` - название компании
- `word` - случайное слово
- `sentence` - предложение
- `paragraph` - параграф
- `url` - URL
- `uuid` - UUID
- `date` - дата
- `timestamp` - timestamp

### Системные endpoints

Facade предоставляет системные endpoints для управления:

- `GET /_facade/health` - проверка состояния сервера
- `GET /_facade/state` - получение всего состояния
- `DELETE /_facade/state` - очистка состояния
- `GET /_facade/state/keys` - получение ключей состояния

### Stateful режим

Facade автоматически сохраняет данные из POST/PUT запросов и может использовать их в последующих GET запросах.

Пример сохранения пользователя:

```bash
# Создание пользователя
curl -X POST http://localhost:8080/users \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Получение состояния
curl http://localhost:8080/_facade/state
```

## Команды сборки и запуска

### Основные команды Makefile

```bash
make help          # Показать все доступные команды
make build         # Собрать приложение в build/facade
make run           # Собрать и запустить с примером конфигурации
make clean         # Очистить файлы сборки
make deps          # Установить/обновить зависимости
make build-all     # Кросс-компиляция для всех платформ
```

### Запуск с различными конфигурациями

```bash
# Запуск с примером конфигурации
make run

# Запуск с кастомной конфигурацией
./build/facade -c configs/frontend-dev.yaml -v

# Запуск с переопределением порта
./build/facade -c configs/testing.yaml -p 9000 -v
```

## Тестирование

### Базовое тестирование

После запуска сервера (`make run`) можно протестировать API:

```bash
# Проверка здоровья сервера
curl http://localhost:8080/_facade/health

# Получение списка пользователей (генерация по схеме)
curl http://localhost:8080/users | jq

# Получение конкретного пользователя (шаблонизация)
curl http://localhost:8080/users/123 | jq

# Создание пользователя (POST с данными)
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}' | jq

# Проверка состояния после POST
curl http://localhost:8080/_facade/state | jq

# Тестирование условной логики
curl "http://localhost:8080/products?category=electronics" | jq
curl "http://localhost:8080/products" | jq

# Тестирование query параметров
curl "http://localhost:8080/search?q=test&limit=10" | jq
```

### Тестирование конфигурации testing.yaml

```bash
# Запуск тестового сервера
./build/facade -c configs/testing.yaml -v

# Успешный платеж
curl -X POST http://localhost:8081/api/payments \
  -H "Content-Type: application/json" \
  -d '{"amount": "500", "card_number": "4111111111111111"}' | jq

# Платеж с верификацией (сумма > 1000)
curl -X POST http://localhost:8081/api/payments \
  -H "Content-Type: application/json" \
  -d '{"amount": "1500", "card_number": "4111111111111111"}' | jq

# Отклоненный платеж
curl -X POST http://localhost:8081/api/payments \
  -H "Content-Type: application/json" \
  -d '{"amount": "100", "card_number": "4000000000000002"}' | jq

# Статусы заказов
curl http://localhost:8081/api/orders/1/status | jq  # pending
curl http://localhost:8081/api/orders/2/status | jq  # shipped
curl http://localhost:8081/api/orders/3/status | jq  # delivered
curl http://localhost:8081/api/orders/999/status | jq # not found

# Симуляция ошибок
curl http://localhost:8081/api/error/500 | jq
curl http://localhost:8081/api/error/404 | jq
curl http://localhost:8081/api/error/401 | jq

# Медленный endpoint (5 секунд задержки)
time curl http://localhost:8081/api/slow | jq
```

## Примеры использования

### Frontend разработка

Настройте Facade как backend для разработки frontend приложения:

```yaml
server:
  port: 3001

routes:
  - path: /api/auth/login
    method: POST
    response:
      status: 200
      body: '{"token": "{{fakerUUID}}", "user": {"id": 1, "name": "{{.Body.username}}"}}'

  - path: /api/users
    method: GET
    response:
      status: 200
      schema:
        type: array
        items:
          type: object
          properties:
            id: {type: integer, min: 1, max: 1000}
            name: {type: string, faker: name}
            avatar: {type: string, faker: url}
        count: 20

global:
  cors: true
  delay: 200ms
```

### Тестирование интеграций

Симуляция различных сценариев для тестов:

```yaml
routes:
  - path: /api/payment
    method: POST
    response:
      status: 200
      conditions:
        - if: '{{.Body.amount}} > "1000"'
          body: '{"status": "pending", "id": "{{fakerUUID}}"}'
        - if: '{{.Body.card_number}} == "4111111111111111"'
          body: '{"status": "success", "id": "{{fakerUUID}}"}'
        - default: true
          status: 400
          body: '{"error": "Invalid card"}'
```

### Демонстрация API

Создание реалистичного API для демонстраций:

```yaml
routes:
  - path: /api/products
    method: GET
    response:
      status: 200
      schema:
        type: object
        properties:
          data:
            type: array
            items:
              type: object
              properties:
                id: {type: integer, min: 1, max: 1000}
                name: {type: string, faker: word}
                description: {type: string, faker: sentence}
                price: {type: number, min: 10, max: 500}
                image: {type: string, faker: url}
                category: {type: string, faker: word}
            count: 12
          total: {type: integer, min: 100, max: 1000}
          page: 1
          per_page: 12
```

## Сборка и развертывание

### Использование Makefile (рекомендуется)

```bash
# Показать все доступные команды
make help

# Быстрый старт - собрать и запустить
make run

# Только сборка
make build

# Очистка
make clean

# Кросс-компиляция для всех платформ
make build-all
```

### Ручная сборка

```bash
# Локальная сборка
mkdir -p build
go build -o build/facade cmd/facade/main.go

# Запуск
./build/facade -c configs/example.yaml
```

## Лицензия

MIT License - см. файл LICENSE для деталей.
