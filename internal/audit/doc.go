// Package audit предоставляет механизм аудита запросов.
//
// Реализован паттерн Наблюдатель (Observer) для асинхронной отправки событий аудита в различные приемники.
//
// Поддерживаемые приемники:
//   - File - запись в файл
//   - HTTP - отправка на удаленный сервер
//
// # Использование
//
//	publisher := audit.NewPublisher()
//	publisher.Register(audit.NewFileObserver("/var/log/audit.log"))
//	publisher.Register(audit.NewHTTPObserver("http://audit-server:8080"))
//	publisher.Publish(audit.NewEvent("shorten", "user123", "https://example.com"))
package audit
