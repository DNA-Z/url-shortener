// Package handler содержит HTTP-обработчики для API сервиса сокращения URL.
//
// Обработчики реализуют следующие эндпоинты:
//   - POST / - создание короткого URL (text/plain)
//   - POST /api/shorten - создание короткого URL (application/json)
//   - POST /api/shorten/batch - пакетное создание коротких URL
//   - GET /{id} - перенаправление по короткому URL
//   - GET /api/user/urls - получение всех URL пользователя
//   - DELETE /api/user/urls - удаление URL пользователя
//   - GET /ping - проверка подключения к БД
package handler
