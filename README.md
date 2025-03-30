# Mattermost Voting Bot - Инструкция по установке и запуску

**О проекте**  
Бот для системы голосований внутри чатов Mattermost.

Позволяет:
- Создавать голосования
- Участвовать в голосованиях
- Просматривать результаты
- Управлять голосованиями


## Установка

### 1. Клонирование репозитория
```bash
git clone https://github.com/ваш-username/mattermost-voting-bot.git
cd mattermost-voting-bot
```
### 2. Настройка окружения
  - Используйте файл local.env увказав в нем подходящие параметры:
```
MATTERMOST_BOT_TOKEN=your_bot_token
MATTERMOST_URL=http://mattermost:8065
TARANTOOL_HOST=tarantool
```

## Запуск проекта

### 1. Сборка и запуск

```bash
docker-compose up -d --build
```
### 2. Проверка работы

```bash
docker-compose ps
```
  - Должны быть активны 3 сервиса:

    - **voting-bot**

    - **tarantool**

    - **mattermost**
## Использование бота

 - Доступные команды:
```bash
/vote create Опция1|Опция2|Опция3   # Создать голосование
/vote vote [ID] [вариант]           # Проголосовать
/vote results [ID]                 # Показать результаты
/vote end [ID]                     # Завершить голосование
/vote delete [ID]                  # Удалить голосование
```
## В этом задании я старался показать все свои скиллы в написании сервисов на GO
