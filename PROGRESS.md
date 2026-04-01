## 2026-03-26
- Добавил структуры Storage, TenantTask, Task
- Реализовал findOrCreateTaskTenant
- Застрял на newApp

## 2026-04-01
- Реализовал newApp: загрузка storage, выбор тенанта из os.Args с fallback на defaultTenant
- Перенёс инкремент Counter внутрь createTask, передаём *TenantTask вместо int
- Убрал локальную переменную tasks, работаем напрямую через app.Current.Tasks
- Исправил main: os.Exit(1) вместо panic при ошибке инициализации
- Подключил saveFile к exit, сохраняем весь Storage