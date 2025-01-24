# Сервис по работе с пользователем и его кошельком

## Структура таблиц
### User

| Field    | Type    | Description                           |
|----------|---------|---------------------------------------|
| Id       | int64   | Primary key, уникальный идентификатор |
| Username | string  | Уникальное имя пользователя           |
| Password | string  | Пароль пользователя                   |
| Email    | string  | Уникальный адрес электронной почты    |


### Wallet

| Field   | Type    | Description                                                 |
|---------|---------|-------------------------------------------------------------|
| Id      | int     | Primary key, уникальный идентификатор                       |
| Usd     | float64 | Количество средств в долларах                               |
| Rub     | float64 | Количество средств в рублях                                 |
| Eur     | float64 | Количество средств в евро                                   |
| UserId  | int64   | Идентификатор пользователя                                  |
| User    | User    | Связь с пользователем (с каскадным обновлением и удалением) |


---
## Api имеет 3 основных деления

Основные ссылки API: \
http://localhost:5000/ - возвращает точку, означает что сервер жив \
http://localhost:5000/api/v1/swagger/index.html - swagger 

User:
 - Register - http://localhost:5000/api/v1/user/register - регистрация нового пользователя 
 - Login - http://localhost:5000/api/v1/user/login - Авторизация пользователя 

Balance:
 - Get Balance - http://localhost:5000/api/v1/user/balance - Баланс на кошельке пользователя 
 - Deposit - http://localhost:5000/api/v1/user/balance/deposit - Добавление валюты на кошелёк пользователя 
 - Withdraw - http://localhost:5000/api/v1/user/balance/withdraw - Списание валюты с кошелка пользователя

Exchange:
 - Get exchange rates - http://localhost:5000/api/v1/user/exchange/rates - Получение текущих валютных курсов
  - Exchange currencies - http://localhost:5000/api/v1/user/exchange - Обмен двух типов валют в кошельке пользователя по курсу

P.S. (Примеры запросов-ответов см. в реализации Swagger)

---
## Используемый стек:
Go:
- [go-chi](https://github.com/go-chi/chi) - RestAPI
- [eko-gocache](https://github.com/eko/gocache) - Библиотека для реализации кеша, имеет много возможностей работы с разными реализациями
- [go-redis](https://github.com/redis/go-redis) - Библиотека для взаимодействия с Redis
- [godotenv](https://github.com/joho/godotenv) - Библиотека для подгрузки .env
- [swaggo](https://github.com/swaggo/swag) - Библиотека для упрощённый генерации swagger через аннотации. Все аннотации написаны в пакете handler
- [gorm](https://gorm.io/) - Библиотека для работы с базой данных
- [gRPC](https://pkg.go.dev/google.golang.org/grpc) - Библиотека для получения данных из второго сервиса через gRPC

PostgresSQL - Основная бд для хранения данных о пользователе \
Redis - Используется в качестве быстро-доступного кеша \
[Taskfile](https://taskfile.dev/) - Инструмент для упрощения написания скриптов позволяя создавать несколько команд в одном yml. 
Используется для быстрого создания независимых контейнеров docker для PostgresSQL, Redis, а также генерация swagger файлов
---
Для локального запуска приложения, в .env нужно заменить адреса с имен докер контейнеров, на localhost и обратно  