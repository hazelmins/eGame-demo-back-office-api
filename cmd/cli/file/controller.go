// go run .\cmd\ginadmin\ file controller -p=shop -c=shopController -t=admin
package file

import (
	"eGame-demo-back-office-api/configs"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	gstrings "eGame-demo-back-office-api/pkg/utils/strings"

	"github.com/spf13/cobra"
)

var controllerStr string = `package {{.PackageName}}

import (
	"eGame-demo-back-office-api/internal/controllers/{{.TypeName}}"

	"github.com/gin-gonic/gin"
)

type {{.ClassName}} struct {
	{{.TypeName}}.BaseController
}

func New{{.UpClassName}}() {{.ClassName}}  {
	return {{.ClassName}}{}
}

func (con {{.ClassName}}) Routes(rg *gin.RouterGroup) {
	{{range $kk,$vv := .Methods}}
	rg.{{$vv}}("/{{$kk}}", con.{{$kk}})
	{{end}}
}

{{range $k,$v := .Methods}}
func (con {{$.ClassName}}) {{$k}}(c *gin.Context) {
}
{{end}}
`

var cmdController = &cobra.Command{
	Use:   "controller [-p pagename -c controllerName -m methods]",
	Short: "create controller file",
	Run:   controllerFunc,
}

var (
	pagename       string
	controllerName string
	methods        string
	typename       string
)

func init() {
	cmdController.Flags().StringVarP(&typename, "typename", "t", "", "input typename api or admin")
	cmdController.Flags().StringVarP(&pagename, "pagename", "p", "", "input pagename eg: setting")
	cmdController.Flags().StringVarP(&controllerName, "controllerName", "c", "", "input controller name eg: AdminController")
	cmdController.Flags().StringVarP(&methods, "methods", "m", "list:get,add:get,save:post,edit:get,del:get", "input methods eg: index:get,add:get")
	cmdController.MarkFlagRequired("typename")
}

func controllerFunc(cmd *cobra.Command, args []string) {

	//判断 typename 类型
	if typename != "admin" && typename != "api" {
		fmt.Println("typename: api or admin")
		os.Exit(1)
	}

	err := configs.Init("")
	if err != nil {

		fmt.Printf("start fail:[Config Init] %s", err.Error())
		os.Exit(1)
	}

	if len(pagename) == 0 || len(controllerName) == 0 {
		cmd.Help()
		return
	}

	pageSlice := strings.Split(pagename, "\\")
	packageName := pageSlice[len(pageSlice)-1]

	upName, _, _ := gstrings.StrFirstToUpper(packageName)

	err = writeController(upName, packageName)
	if err != nil {
		fmt.Printf("[error] %s", err.Error())
		os.Exit(1)
	}
}

func writeController(upName string, packageName string) error {
	// 定義一個名為 parms 的結構，用於儲存控制器生成所需的參數
	parms := struct {
		ClassName   string
		Pagename    string
		PackageName string
		UpClassName string
		TypeName    string
		Methods     map[string]string
	}{
		ClassName:   controllerName, // 控制器名稱
		Pagename:    pagename,       // 頁面名稱
		PackageName: packageName,    // 包名稱
		TypeName:    typename,       // 類型名稱
		UpClassName: upName,         // 大寫的類型名稱
	}

	// 將傳遞進來的 methods 字串拆分為方法名和HTTP方法的映射
	//methodMap 是一個 map（映射），用於將 HTTP 請求方法映射到相應的大寫字母的字符串。
	methods := strings.Split(methods, ",")
	methodMap := make(map[string]string)
	for _, v := range methods {
		methodSlice := strings.Split(v, ":")
		methodMap[methodSlice[0]] = strings.ToUpper(methodSlice[1])
		//methodSlice[0] 是切片中的第一個元素，即 "GET"，代表 HTTP 請求方法。

		//methodSlice[1] 是切片中的第二個元素，即 "read"，代表與該方法相關的操作或名稱。
	}
	parms.Methods = methodMap

	// 設置控制器文件的基本路徑
	basePath := configs.RootPath + "internal" + string(filepath.Separator) + "controllers" + string(filepath.Separator) + typename + string(filepath.Separator) + pagename

	// 檢查路徑是否存在，如果不存在則創建之
	_, err := os.Lstat(basePath)
	if err != nil {
		os.Mkdir(basePath, os.ModeDir)
	}

	// 設置新控制器文件的路徑
	newPath := basePath + string(filepath.Separator) + controllerName + ".go"

	// 檢查是否已經存在相同名稱的控制器文件，如果存在則返回錯誤
	_, err = os.Lstat(newPath)
	if err == nil {
		return err
	}

	// 創建新的控制器文件
	file, err := os.Create(newPath)
	if err != nil {
		cobra.CompError(err.Error())
		return err
	}
	defer file.Close()

	// 使用模板將參數填充到控制器文件中
	tem, _ := template.New("controller_file").Parse(controllerStr)
	tem.ExecuteTemplate(file, "controller_file", parms)

	return nil
}

/*
這個代碼段的目的是生成一個 Gin 框架的控制器（controller）文件。它接受一些參數，包括控制器名稱、頁面名稱、包名稱、類型名稱、大寫的類型名稱和一組方法映射。然後，它根據這些參數創建控制器文件，將其放置在指定的路徑下。

具體步驟包括：

創建一個 parms 結構，用於存儲控制器生成所需的參數。
拆分傳遞進來的 methods 字串，並將方法名和HTTP方法映射存儲在 methodMap 中。
設置控制器文件的基本路徑，檢查路徑是否存在，如果不存在則創建之。
設置新控制器文件的路徑，檢查是否已經存在相同名稱的控制器文件，如果存在則返回錯誤。
創建新的控制器文件，並使用模板將參數填充到文件中。
這個函數的主要目的是生成控制器文件，以便在 Gin 框架中使用。
*/
