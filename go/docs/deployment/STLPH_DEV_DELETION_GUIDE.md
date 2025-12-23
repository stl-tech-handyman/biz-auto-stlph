# Удаление проекта stlph-dev

## 📋 Анализ проекта

**Проект:** `stlph-dev` (BizOps360-Dev)  
**Статус:** ACTIVE (но пустой)

## ✅ Проверка активности

### 1. Проверка Cloud Run сервисов:
```bash
gcloud run services list --project=stlph-dev
```

### 2. Проверка Cloud Functions:
```bash
gcloud functions list --project=stlph-dev
```

### 3. Проверка Compute Instances:
```bash
gcloud compute instances list --project=stlph-dev
```

### 4. Проверка использования в коде:
```bash
# В активном коде (go/) нет упоминаний stlph-dev
grep -r "stlph-dev" go/
# Результат: No matches found ✅
```

## 📊 Текущее состояние

### Активные проекты (используются):
- ✅ **bizops360-dev** - Основной Go API проект
- ✅ **bizops360-prod** - Production проект
- ✅ **bizops360-email-dev** - Email сервис
- ✅ **bizops360-maps** - Maps API keys

### Неактивный проект:
- ⚠️ **stlph-dev** - Пустой проект (все сервисы удалены)

## 🗑️ История очистки

Согласно `archive/oldcode/stlph-cloud/CLEANUP_STLPH_DEV_COMPLETE.md`:

**Удалено из stlph-dev:**
1. ✅ `geo` - Legacy Cloud Function
2. ✅ `health` - Legacy Cloud Function  
3. ✅ `healthz` - Legacy Cloud Function
4. ✅ `leads` - Cloud Function
5. ✅ `stlph-api` - Old JavaScript API
6. ✅ `stlph-api-go-dev` - Duplicate Go API
7. ✅ `stlph-email-api` - Legacy email service
8. ✅ `stlph-geo-api` - Legacy geo service
9. ✅ `stlph-health-api` - Legacy health service

**Результат:** Проект пустой, можно удалить ✅

## ⚠️ Перед удалением - финальная проверка

### 1. Проверьте, нет ли активных ресурсов:
```bash
# Cloud Run
gcloud run services list --project=stlph-dev

# Cloud Functions
gcloud functions list --project=stlph-dev

# Compute
gcloud compute instances list --project=stlph-dev

# Storage
gsutil ls -p stlph-dev

# Secrets (если нужны - экспортируйте перед удалением)
gcloud secrets list --project=stlph-dev
```

### 2. Проверьте billing:
```bash
gcloud billing projects describe stlph-dev
```

### 3. Проверьте, нет ли внешних ссылок:
- Проверьте документацию на ссылки на `stlph-dev` URLs
- Проверьте интеграции (Zapier, webhooks, etc.)
- Проверьте DNS записи

## 🗑️ Удаление проекта

### Вариант 1: Удалить проект полностью (РЕКОМЕНДУЕТСЯ)

```bash
# 1. Убедитесь, что проект пустой
gcloud run services list --project=stlph-dev
gcloud functions list --project=stlph-dev

# 2. Удалите проект (WARNING: Это необратимо!)
gcloud projects delete stlph-dev

# 3. Подтвердите удаление (введите project ID)
# Введите: stlph-dev
```

### Вариант 2: Отключить billing (если хотите сохранить проект)

```bash
# Отключить billing (проект останется, но не будет работать)
gcloud billing projects unlink stlph-dev
```

## ✅ После удаления

### Проверьте активные проекты:
```bash
gcloud projects list --filter="projectId:bizops360*"
```

### Ожидаемый результат:
- ✅ `bizops360-dev` - активен
- ✅ `bizops360-prod` - активен
- ✅ `bizops360-email-dev` - активен
- ✅ `bizops360-maps` - активен
- ❌ `stlph-dev` - удален

## 📝 Примечания

1. **Все активные deployment используют `bizops360-dev`** ✅
2. **В коде нет ссылок на `stlph-dev`** ✅
3. **Проект был очищен ранее** ✅
4. **Безопасно удалить** ✅

## ⚠️ ВАЖНО

- Удаление проекта **необратимо**
- Убедитесь, что нет активных ресурсов
- Экспортируйте secrets, если они нужны
- Проверьте billing перед удалением

---

**Дата создания:** 2025-01-XX  
**Статус:** Готово к удалению ✅



