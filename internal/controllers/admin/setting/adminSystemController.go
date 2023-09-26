/*
 * @Description:系统管理 修改redis目录结构
 */

package setting

import (
	"bufio"
	"context"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"eGame-demo-back-office-api/configs"
	"eGame-demo-back-office-api/internal/controllers/admin"
	"eGame-demo-back-office-api/pkg/loggers"
	"eGame-demo-back-office-api/pkg/redisx"
	"eGame-demo-back-office-api/pkg/utils/filesystem"
	gstrings "eGame-demo-back-office-api/pkg/utils/strings"

	"github.com/gin-gonic/gin"
)

type adminSystemController struct {
	admin.BaseController
}

func NewAdminSystemController() adminSystemController {
	return adminSystemController{}
}

func (con adminSystemController) Routes(rg *gin.RouterGroup) {
	rg.GET("/index", con.index)
	rg.GET("/getdir", con.getDir)
	rg.GET("/view", con.view)
	rg.GET("/index_redis", con.indexRedis)
	rg.GET("/getdir_redis", con.getDirRedis)
	rg.GET("/view_redis", con.viewRedis)
}

/**
日志目录页面
*/
func (con adminSystemController) index(c *gin.Context) {
	// 定義變數
	var (
		path     string
		err      error
		log_path string
	)

	// 設定日誌路徑為根路徑下的 "logs" 目錄
	path = gstrings.JoinStr(configs.RootPath, string(filepath.Separator), "logs")

	// 讀取指定目錄下的所有文件
	files, err := ioutil.ReadDir(path)

	// 設定日誌路徑為單純的 "logs"，不包括根路徑
	log_path = gstrings.JoinStr(string(filepath.Separator), "logs")

	// 從 Gin 的上下文中獲取一個名為 "ctx" 的變數
	ctx, _ := c.Get("ctx")

	// 如果讀取目錄時發生錯誤，記錄錯誤日誌，並回傳錯誤頁面
	if err != nil {
		loggers.LogError(ctx.(context.Context), "admin", "讀取目錄失敗", map[string]string{"error": err.Error()})
		con.ErrorHtml(c, err)
		return
	}

	// 傳遞 HTML 響應，狀態碼為 200 OK，並傳遞以下變數到模板中：
	// - "log_path": 日誌路徑
	// - "files": 目錄中的文件列表
	// - "line": 文件路徑的分隔符
	con.Html(c, http.StatusOK, "setting/systemlog.html", gin.H{
		"log_path": log_path,
		"files":    files,
		"line":     string(filepath.Separator),
	})
}

/**
获取目录用於在網頁中列出指定目錄下的文件和子目錄，並返回 JSON 格式的結果
*/
func (con adminSystemController) getDir(c *gin.Context) {
	// 定義 FileNode 類型的結構，用於表示文件或目錄的信息
	type FileNode struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Type string `json:"type"`
	}

	var (
		path      string
		err       error
		fileSlice []FileNode
		files     []fs.FileInfo
	)

	// 初始化一個空的 FileNode 切片
	fileSlice = make([]FileNode, 0)

	// 從請求參數中獲取路徑，並過濾路徑，確保安全性
	path, err = filesystem.FilterPath(configs.RootPath+"logs", c.Query("path"))
	if err != nil {
		con.Error(c, err.Error())
		return
	}

	// 讀取指定路徑下的所有文件和目錄
	files, err = ioutil.ReadDir(path)
	if err != nil {
		con.Error(c, "獲取目錄失敗")
		return
	}

	// 遍歷文件列表，為每個文件或目錄創建對應的 FileNode 並添加到切片中
	for _, v := range files {
		var fileType string
		if v.IsDir() {
			fileType = "dir"
		} else {
			fileType = "file"
		}
		fileSlice = append(fileSlice, FileNode{
			Name: v.Name(),
			Path: gstrings.JoinStr(c.Query("path"), string(filepath.Separator), v.Name()),
			Type: fileType,
		})
	}

	// 將結果以 JSON 格式回傳給客戶端，狀態碼為 200 OK
	c.JSON(http.StatusOK, gin.H{
		"data": fileSlice,
	})
}

/**
获取日志详情用於查看指定文件的部分內容，並在 Web 頁面上顯示
*/
func (con adminSystemController) view(c *gin.Context) {
	// 定義變數
	var (
		err       error
		startLine int
		endLine   int
		scanner   *bufio.Scanner
		line      int
	)

	// 從 URL 查詢參數中獲取起始行數，如果無法獲取或轉換成整數，則處理錯誤
	startLine, err = strconv.Atoi(c.DefaultQuery("start_line", "1"))
	if err != nil {
		con.ErrorHtml(c, err)
		return
	}

	// 從 URL 查詢參數中獲取結束行數，如果無法獲取或轉換成整數，則處理錯誤
	endLine, err = strconv.Atoi(c.DefaultQuery("end_line", "20"))
	if err != nil {
		con.ErrorHtml(c, err)
		return
	}

	// 創建一個存儲文件內容的字符串切片
	var filecontents []string

	// 獲取文件路徑並過濾路徑以確保安全性
	filePath, err := filesystem.FilterPath(configs.RootPath+"logs", c.Query("path"))
	if err != nil {
		con.ErrorHtml(c, err)
		return
	}

	// 打開文件以進行讀取
	fi, err := os.Open(filePath)
	if err != nil {
		con.ErrorHtml(c, err)
		return
	}
	defer fi.Close()

	// 使用 bufio.Scanner 逐行讀取文件內容
	scanner = bufio.NewScanner(fi)
	for scanner.Scan() {
		line++
		if line >= startLine && line <= endLine {
			// 在指定行數範圍內取得數據
			filecontents = append(filecontents, scanner.Text())
		} else {
			continue
		}
	}

	// 傳遞 HTML 響應，狀態碼為 200 OK，並傳遞以下變數到模板中：
	// - "file_path": 文件路徑
	// - "filecontents": 文件內容的字符串切片
	// - "start_line": 起始行數
	// - "end_line": 結束行數
	// - "line": 文件的總行數
	con.Html(c, http.StatusOK, "setting/systemlog_view.html", gin.H{
		"file_path":    c.Query("path"),
		"filecontents": filecontents,
		"start_line":   startLine,
		"end_line":     endLine,
		"line":         line,
	})
}

/**
日志目录页面用於從 Redis 中獲取特定路徑的日誌鍵提取日期部分，並在 Web 頁面上顯示不重複的日期部分列表
*/
func (con adminSystemController) indexRedis(c *gin.Context) {
	// 定義日誌路徑
	path := "logs"

	// 從 Redis 中獲取所有鍵以匹配指定路徑的日誌
	dateSlice, err := redisx.GetRedisClient().Keys("logs:*").Result()

	// 從 Gin 上下文中獲取一個名為 "ctx" 的變數
	ctx, _ := c.Get("ctx")

	// 如果發生錯誤，記錄錯誤日誌，並回傳錯誤頁面
	if err != nil {
		loggers.LogError(ctx.(context.Context), "admin", "讀取目錄失敗", map[string]string{"error": err.Error()})
		con.ErrorHtml(c, err)
		return
	}

	// 創建一個名為 "dates" 的 map，用於存儲不重複的日期部分
	dates := make(map[string]struct{})

	// 遍歷日期切片，提取日期部分並存入 "dates" map 中
	for _, v := range dateSlice {
		temp := strings.Split(v, ":")

		if _, ok := dates[temp[1]]; !ok {
			dates[temp[1]] = struct{}{}
		}
	}

	// 傳遞 HTML 響應，狀態碼為 200 OK，並傳遞以下變數到模板中：
	// - "log_path": 日誌路徑
	// - "files": 存儲不重複日期部分的 map
	con.Html(c, http.StatusOK, "setting/systemlog_redis.html", gin.H{
		"log_path": path,
		"files":    dates,
	})
}

/**
获取目录用於從 Redis 中獲取特定路徑的日誌鍵，提取相關信息，然後返回 JSON 格式的結果
*/
func (con adminSystemController) getDirRedis(c *gin.Context) {
	// 從 URL 查詢參數中獲取路徑
	path := c.Query("path")

	// 定義 FileNode 類型的結構，用於表示文件或目錄的信息
	type FileNode struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Type string `json:"type"`
	}

	// 將路徑按下劃線 "_" 分割成切片
	pathSlice := strings.Split(path, "_")

	// 構建 Redis 鍵的模式，用於查找相關日誌鍵
	pattern := pathSlice[0] + ":*"

	// 從 Redis 中獲取所有符合模式的鍵
	dateSlice, err := redisx.GetRedisClient().Keys(pattern).Result()

	// 從 Gin 上下文中獲取一個名為 "ctx" 的變數
	ctx, _ := c.Get("ctx")

	// 如果發生錯誤，記錄錯誤日誌，並回傳錯誤頁面
	if err != nil {
		loggers.LogError(ctx.(context.Context), "admin", "讀取目錄失敗", map[string]string{"error": err.Error()})
		con.ErrorHtml(c, err)
		return
	}

	// 創建一個空的 FileNode 切片，用於存儲文件和目錄的信息
	fileSlice := make([]FileNode, 0)

	// 創建一個臨時的 map，用於檢查重複的條目
	tempMap := make(map[string]struct{})

	// 遍歷日期切片，提取相關信息並填充到 fileSlice 中
	for _, v := range dateSlice {
		temp := strings.Split(v, ":")
		index, _ := strconv.Atoi(pathSlice[1])
		var fileType string

		// 確定條目是文件還是目錄
		if index+2 == len(temp) {
			fileType = "file"
		} else {
			fileType = "dir"
		}

		// 檢查 tempMap 以避免重複條目
		if _, ok := tempMap[temp[index+1]]; ok {
			continue
		} else {
			tempMap[temp[index+1]] = struct{}{}
		}

		// 將條目信息添加到 fileSlice 中
		fileSlice = append(fileSlice, FileNode{
			Name: temp[index+1],
			Path: pathSlice[0] + ":" + temp[index+1] + "_" + strconv.Itoa(index+1),
			Type: fileType,
		})
	}

	// 將結果以 JSON 格式回傳給客戶端，狀態碼為 200 OK
	c.JSON(http.StatusOK, gin.H{
		"data": fileSlice,
	})
}

/**
获取日志详情用於從 Redis 中獲取指定範圍的日誌內容，然後在 Web 頁面上顯示
*/
func (con adminSystemController) viewRedis(c *gin.Context) {
	// 從 URL 查詢參數中獲取起始行數，如果無法獲取或轉換成整數，使用默認值 1
	startLine, _ := strconv.Atoi(c.DefaultQuery("start_line", "1"))

	// 從 URL 查詢參數中獲取結束行數，如果無法獲取或轉換成整數，使用默認值 20
	endLine, _ := strconv.Atoi(c.DefaultQuery("end_line", "20"))

	// 從 URL 查詢參數中獲取文件路徑
	filePath := c.Query("path")

	// 將文件路徑按下劃線 "_" 分割成切片
	pathSlice := strings.Split(filePath, "_")

	// 從 Redis 中獲取指定範圍的日誌內容
	filecontents, _ := redisx.GetRedisClient().LRange(pathSlice[0], int64(startLine-1), int64(endLine-1)).Result()

	// 從 Redis 中獲取指定鍵的列表長度，即文件的總行數
	line, _ := redisx.GetRedisClient().LLen(pathSlice[0]).Result()

	// 傳遞 HTML 響應，狀態碼為 200 OK，並傳遞以下變數到模板中：
	// - "file_path": 文件路徑
	// - "filecontents": 文件內容的字符串切片
	// - "start_line": 起始行數
	// - "end_line": 結束行數
	// - "line": 文件的總行數
	con.Html(c, http.StatusOK, "setting/systemlog_viewredis.html", gin.H{
		"file_path":    filePath,
		"filecontents": filecontents,
		"start_line":   startLine,
		"end_line":     endLine,
		"line":         line,
	})
}
