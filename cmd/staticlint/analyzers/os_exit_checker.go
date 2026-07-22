package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osExitChecker",
	Doc:  `Проверяет, что os.Exit не вызывается напрямую в функции main пакета main`,
	Run:  runOsExitCheck,
}

func runOsExitCheck(pass *analysis.Pass) (interface{}, error) {
	// Проходим по всем файлам пакета
	for _, file := range pass.Files {
		// Проверяем, что это файл пакета main
		if file.Name.Name != "main" {
			continue
		}

		// Инспектируем AST дерево
		ast.Inspect(file, func(n ast.Node) bool {
			// Ищем функцию main
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// Проверяем, что это функция main
			if funcDecl.Name.Name != "main" {
				return true
			}

			// Проверяем тело функции
			if funcDecl.Body == nil {
				return true
			}

			// Инспектируем тело функции
			ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
				// Ищем вызов функции
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				// Проверяем, что это вызов функции из пакета os
				selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				// Проверяем, что это os.Exit
				if pkg, ok := selExpr.X.(*ast.Ident); ok {
					if pkg.Name == "os" && selExpr.Sel.Name == "Exit" {
						// Получаем позицию в файле
						pos := pass.Fset.Position(callExpr.Pos())

						// Сообщаем об ошибке
						pass.Reportf(
							callExpr.Pos(),
							"прямой вызов os.Exit в функции main запрещен. "+
								"Используйте log.Fatal или возврат ошибки. (файл: %s, строка: %d)",
							pos.Filename,
							pos.Line,
						)
					}
				}
				return true
			})
			return true
		})
	}
	return nil, nil
}
