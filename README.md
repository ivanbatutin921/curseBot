# Telegram Bot для продажи курса

Бот для автоматизации продажи онлайн-курса с функциями проверки подписки на канал и приема платежей через Tinkoff API.

## Стек технологий

- Go 1.21+
- Fiber v2 (веб-фреймворк)
- PostgreSQL
- GORM (ORM для Go)
- godotenv
- Tinkoff Payment API SDK

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/course-sale-bot.git
cd course-sale-bot
```

2. Установите зависимости:
```bash
go mod init course-sale-bot
go mod tidy
```

3. Создайте файл .env на основе .env.example и заполните необходимые переменные окружения:
```bash
cp .env.example .env
```

4. Создайте базу данных PostgreSQL и укажите параметры подключения в .env

## Запуск

1. Запустите бот:
```bash
go run cmd/main.go
```

## Функциональность

- Автоматическая регистрация пользователей
- Проверка подписки на канал
- Информация о курсе
- Прием платежей (полная оплата/рассрочка)
- Система доступа к курсу
- Административные функции

## Структура проекта

```
cmd/
├── main.go            # Точка входа
internal/
├── api/              # HTTP обработчики
│   └── handlers/     # Обработчики команд
├── bot/              # Логика бота
│   ├── keyboards/    # Клавиатуры
│   └── middleware/   # Middleware
├── config/          # Конфигурация
├── models/          # Модели данных
├── repository/      # Работа с БД
└── service/         # Бизнес-логика
    └── payment/     # Интеграция с Tinkoff API
pkg/
└── utils/          # Вспомогательные функции
```

## Переменные окружения

- `BOT_TOKEN` - Токен Telegram бота
- `CHANNEL_ID` - ID канала для проверки подписки
- `DATABASE_URL` - URL подключения к PostgreSQL
- `TINKOFF_TERMINAL_KEY` - Ключ терминала Tinkoff
- `TINKOFF_TERMINAL_PASSWORD` - Пароль терминала Tinkoff
- `PRIVATE_CHAT_ID` - ID приватного чата с курсом

## Разработка

1. Создайте ветку для новой функциональности:
```bash
git checkout -b feature/your-feature-name
```

2. Внесите изменения и создайте коммит:
```bash
git add .
git commit -m "Add your feature description"
```

3. Отправьте изменения в репозиторий:
```bash
git push origin feature/your-feature-name
```

## Тестирование

Для запуска тестов используйте:
```bash
go test ./...
```

## Лицензия

MIT 