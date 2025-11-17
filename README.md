## Сервис на Go, который обращается к "OpenExchangeRates" и "Giphy", сравнивает курс валюты к рублю за сегодня/вчера и отдаёт случайный GIF
- если курс вырос - GIF по запросу "rich"
- если курс упал или не вырос - GIF по запросу "broke"

## Технологии 
- Gin
- курсы валют (OpenExchangeRates)
- GIF (Giphy API)
- Docker

## Запуск в Docker 
- сборка образа: docker build -t uplatform:alpine .
- запуск: docker run --rm -p 8080:8080 --env-file .env uplatform:alpine
После запуска сервер доступен по адресу: http://localhost:8080
