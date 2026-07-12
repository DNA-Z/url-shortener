// staticlint - статический анализатор для URL Shortener
//
// Использование:
// staticlint [flags] [packages]
//
// Флаги:
// -h, -help    Показать эту справку
//
// Примеры:
// staticlint .         # Проверка текущей директории
// staticlint ./...     # Проверка всех поддиректорий
// staticlint ./cmd/staticlint  # Проверка конкретной директории
//
// Состав анализаторов:
// 1. Стандартные анализаторы (30+ шт.)
// 2. Все анализаторы класса SA из staticcheck (50+ шт.)
// 3. Анализаторы классов S, ST, QF из staticcheck (минимум 1 каждого)
// 4. Дополнительные: errchkjson, bodyclose
package main
