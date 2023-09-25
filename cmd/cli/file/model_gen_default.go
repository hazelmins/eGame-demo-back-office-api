package file

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"eGame-demo-back-office-api/configs"
)

func modifyDefault(fileName string) {
	// 構建預設文件的路徑，根據配置的根路徑和文件名
	var filePath = configs.RootPath + "internal" + string(filepath.Separator) + "models" + string(filepath.Separator) + "default.go"

	// 創建文件集合 fset 用於存儲文件
	fset := token.NewFileSet()

	// 解析預設文件，並使用 ParseComments 解析註解
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Println(err)
		return
	}

	// 使用自定義的 Visitor 來遍歷抽象語法樹 (AST)
	ast.Walk(&Visitor{
		fset: fset,
		name: fileName,
	}, f)

	// 將修改後的 AST 寫回到原文件
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0766)
	if err != nil {
		log.Fatalf("open err %s", err.Error())
	}

	// 格式化並寫回文件
	err = format.Node(file, fset, f)
	if err != nil {
		log.Fatal(err)
	}
}

// Visitor 結構實現了 ast.Visitor 接口，用於遍歷 AST
type Visitor struct {
	fset *token.FileSet
	name string
}

// Visit 方法是 ast.Visitor 接口的實現，用於處理不同類型的 AST 节点
func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.FuncDecl: // 如果是函數聲明節點
		demo := node.(*ast.FuncDecl)
		if demo.Name.Name == "GetModels" { // 如果函數名稱是 "GetModels"

			// 取得函數體的第一個語句，預期是 return 聲明
			returnStm, ok := demo.Body.List[0].(*ast.ReturnStmt)
			if !ok {
				return v
			}

			// 取得 return 聲明的返回值，預期是一個複合字面值
			comp, ok := returnStm.Results[0].(*ast.CompositeLit)
			if !ok {
				return v
			}

			// 向複合字面值的元素列表中添加一個新的元素，這是一個對新型別的指針
			comp.Elts = append(comp.Elts, &ast.UnaryExpr{
				Op: token.AND,
				X: &ast.CompositeLit{
					Type: &ast.Ident{
						Name: v.name, // 使用指定的文件名
					},
				},
			})

		}
	}

	return v
}

/*
這段代碼的目的是修改一個 Go 文件（通常是 default.go）中的 AST，具體來說，是尋找名稱為 "GetModels" 的函數，然後向它的返回值中添加一個新的元素，這個元素是一個對指定名稱文件的新型別的指針。

詳細步驟包括：

構建文件路徑：根據配置的根路徑和文件名，構建預設文件的路徑。

解析預設文件：使用 parser.ParseFile 函數解析預設文件，並使用 parser.ParseComments 選項解析註解。

使用自定義的 Visitor 來遍歷 AST：使用 ast.Walk 函數遍歷 AST，Visitor 結構實現了 ast.Visitor 接口，並在訪問 GetModels 函數時進行修改。

修改 AST：在 Visitor 的 Visit 方法中，當訪問到 GetModels 函數時，它會尋找函數體中的 return 聲明，然後向返回值中的複合字面值添加一個新的元素，這個元素是對指定名稱文件的新型別的指針。

將修改後的 AST 寫回文件：打開文件，將修改後的 AST 以適當的格式寫回文件中。

總之，這段代碼的目的是在指定的 Go 文件中進行 AST 修改，以實現特定的功能或行為。
*/
