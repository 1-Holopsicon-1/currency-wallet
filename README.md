## Мини проект по созданию Кошелька пользователя и обмен валют

Система состоит из трех подсервисов:
- Сервис управления кошельками и взаимодействия с пользователем через REST API. 
- Сервис получения актуальных курсов валют с использованием gRPC.
- Общий gRPC интерфейс с описанием в .proto-файле и сгенерированными методами для взаимодействия между сервисами.

О каждом сервисе см. подробнее README внутри других папок

---
### Диаграмма взаимодействия

[Пользователь] -- REST API --> [wallet-service] -- gRPC --> [currency-service]

---
### Предварительные требования
- Docker (для запуска контейнеров) \
Проект собирается с помощью docker-compose командой

``` bash
docker-compose up --build 
```

### Для нативного запуска
- Go версии >= 1.20
- Протокол-буфер (`protoc`) и плагины для Go (`protoc-gen-go`)
- Swaggo для генерации свагера
- Taskfile (желательно но не обязательно, т.к. просто упрощает взаимодействие)

P.S. (Обо всех зависимостях также см. README каждого сервиса)

---

### Установка и запуск

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/your-repo/currency-wallet.git
   cd currency-wallet
   ```
2. Сгенерируйте gRPC-код из .proto-файлов:
    ```bash
   cd protos
   task gen 
   ```
3. Соберите и запустите сервисы:
- Запуск wallet-service
    ```bash
   go run gw-currency-wallet/cmd/main.go -migrate -start
   ```
- Запуск currency-rate-service
   ```bash
   go run gw-exchanger/cmd/main.go -migrate -startGrpc
   ```
