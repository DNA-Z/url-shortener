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
// 1. Анализаторы класса SA из staticcheck
// 2. Анализаторы класса S (simple)
// 3. Анализаторы класса ST (stylecheck)
// 4. Анализаторы класса QF (quickfix)
// 5. Дополнительные: errchkjson, bodyclose
// 6. Собственный анализатор OsExitAnalyzer
package main
