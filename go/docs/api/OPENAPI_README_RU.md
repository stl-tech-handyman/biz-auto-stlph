# OpenAPI Спецификация - BizOps360 Go API

## 📋 Описание

Полная спецификация OpenAPI 3.0.3 для BizOps360 Go API на русском языке.

**Файл:** `openapi-ru.yaml`

## 🚀 Использование

### Просмотр в Swagger UI

1. Откройте [Swagger Editor](https://editor.swagger.io/)
2. Загрузите файл `openapi-ru.yaml`
3. Просматривайте и тестируйте API

### Генерация клиентов

Используйте инструменты для генерации клиентов из OpenAPI спецификации:

```bash
# Используя openapi-generator
openapi-generator generate -i openapi-ru.yaml -g go -o ./clients/go

# Используя swagger-codegen
swagger-codegen generate -i openapi-ru.yaml -l go -o ./clients/go
```

### Валидация

```bash
# Используя Python
python -c "import yaml; yaml.safe_load(open('openapi-ru.yaml', encoding='utf-8'))"

# Используя swagger-cli
npx @apidevtools/swagger-cli validate openapi-ru.yaml
```

## 📊 Статистика API

- **Всего эндпоинтов:** 15
- **Требуют аутентификацию:** 8
- **Без аутентификации:** 7

### По категориям:

- **Здоровье и информация:** 4 эндпоинта
- **Stripe:** 4 эндпоинта (требуют API ключ)
- **Расчет стоимости:** 2 эндпоинта (требуют API ключ)
- **Email:** 2 эндпоинта (требуют API ключ)
- **Календарь:** 1 эндпоинт
- **Обработка лидов:** 2 эндпоинта
- **V1 Pipeline:** 2 эндпоинта

## 🔐 Аутентификация

### API Key Authentication

Для эндпоинтов, требующих аутентификацию, используйте заголовок:

```
X-Api-Key: your-api-key-here
```

### Получение API ключа

```bash
gcloud secrets versions access latest \
  --secret="svc-api-key-dev" \
  --project="bizops360-dev"
```

## 📝 Примеры использования

### Пример запроса (cURL)

```bash
# Расчет стоимости
curl -X POST "https://bizops360-api-go-dev-gqqr4r256q-uc.a.run.app/api/estimate" \
  -H "X-Api-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "eventDate": "2025-06-15",
    "durationHours": 4.0,
    "numHelpers": 2
  }'
```

### Пример запроса (JavaScript)

```javascript
const response = await fetch('https://bizops360-api-go-dev-gqqr4r256q-uc.a.run.app/api/estimate', {
  method: 'POST',
  headers: {
    'X-Api-Key': 'your-api-key',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    eventDate: '2025-06-15',
    durationHours: 4.0,
    numHelpers: 2
  })
});

const data = await response.json();
```

## 🔄 Обновление спецификации

При добавлении новых эндпоинтов:

1. Обновите `openapi-ru.yaml`
2. Добавьте схемы запросов/ответов
3. Обновите примеры
4. Проверьте валидность YAML
5. Обновите счетчик эндпоинтов

## 📚 Дополнительная документация

- [ENDPOINTS_SUMMARY.md](../postman/ENDPOINTS_SUMMARY.md) - Краткое описание всех эндпоинтов
- [POSTMAN_COLLECTION.md](POSTMAN_COLLECTION.md) - Информация о Postman коллекции
- [API_KEY_AUTHENTICATION_GUIDE.md](API_KEY_AUTHENTICATION_GUIDE.md) - Руководство по аутентификации

---

**Последнее обновление:** 2025-01-XX  
**Версия спецификации:** 1.0.0



