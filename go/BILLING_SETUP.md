# Привязка Billing Account

## Billing Account ID
`01C379-C9A8C1-3ED059`

## Способ 1: Через веб-интерфейс (рекомендуется)

1. Откройте страницу linked projects:
   https://console.cloud.google.com/billing/01C379-C9A8C1-3ED059/linkedprojects

2. Нажмите кнопку **"Link a project"** или **"Link project"**

3. Выберите проект **bizops360-dev**

4. Подтвердите привязку

## Способ 2: Через CLI (требует beta компонент)

```bash
gcloud beta billing projects link bizops360-dev --billing-account=01C379-C9A8C1-3ED059
```

Если beta компонент не установлен:
```bash
gcloud components install beta
```

## После привязки

Запустите полную настройку:
```bash
cd go
bash scripts/complete-setup-dev.sh
```

Скрипт автоматически:
- ✅ Включит все необходимые API
- ✅ Создаст Artifact Registry
- ✅ Настроит Docker auth
- ✅ Создаст secrets
- ✅ Задеплоит сервис

