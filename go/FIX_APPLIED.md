# ✅ Исправление применено

## Проблема
В шаблоне `endpoint_manager.html` использовалось несуществующее поле `.Auth` вместо `.AuthRequired`.

## Исправление
Исправлена строка 571 в файле `go/internal/http/handlers/endpoint_manager.html`:

**Было:**
```html
<tr data-method="{{.Method}}" data-auth="{{.Auth}}" ...>
```

**Стало:**
```html
<tr data-method="{{.Method}}" data-auth="{{if .AuthRequired}}required{{else}}none{{end}}" ...>
```

## Что нужно сделать

**Перезапустите сервер**, чтобы изменения вступили в силу:

```powershell
# 1. Остановите текущий сервер (Ctrl+C)

# 2. Запустите заново:
cd go
.\start-local.ps1
```

## После перезапуска

- ✅ Endpoint Manager будет работать без ошибок
- ✅ Таблица endpoints отобразится правильно
- ✅ Swagger UI будет доступен на `/swagger`
- ✅ OpenAPI spec будет доступен на `/api/openapi.json`

