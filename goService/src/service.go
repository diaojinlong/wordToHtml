package main

import (
	"archive/zip"
	"fmt"
	"github.com/gofrs/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/convert", convertHandler)
	log.Fatal(http.ListenAndServe(":8083", nil))
}

// 打包成zip文件
func Zip(src_dir string, zip_file_name string) {

	// 预防：旧文件无法覆盖
	os.RemoveAll(zip_file_name)

	// 创建：zip文件
	zipfile, _ := os.Create(zip_file_name)
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历路径信息
	filepath.Walk(src_dir, func(path string, info os.FileInfo, _ error) error {

		// 如果是源路径，提前进行下一个遍历
		if path == src_dir {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src_dir+`/`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	// 实现请求处理
	// 从请求中获取文件参数的值
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	daystr := time.Now().Format("20060102")
	id, _ := uuid.NewV4()
	os.Mkdir(daystr, os.ModePerm)
	upload_path := daystr + "/" + strings.ReplaceAll(id.String(), "-", "")

	//创建上传目录
	os.Mkdir(upload_path, os.ModePerm)
	os.Mkdir(upload_path+"/html", os.ModePerm)
	file_ext := strings.ToLower(path.Ext(handler.Filename))
	if file_ext != `.doc` && file_ext != `.docx` {
		http.Error(w, `文档格式不正确`, http.StatusInternalServerError)
		return
	}
	upload_file := upload_path + "/index" + file_ext
	//创建上传文件
	f, err := os.Create(upload_file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	// 执行命令行操作
	cmd := exec.Command("/opt/libreoffice7.5/program/soffice", "--headless", "--convert-to", "html", upload_file, "--outdir", upload_path+"/html")
	stdout, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%s", stdout)
	// 目录压缩为 zip 文件
	Zip(upload_path+"/html", upload_path+"/html.zip")

	// 打包后提供下载
	// 如果使用 stdlib 包下载压缩包，示例代码如下：
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=html.zip")
	zipFile, err := os.Open(upload_path + "/html.zip")
	defer zipFile.Close()
	io.Copy(w, zipFile)
}

//docker exec -it golang /bin/bash
//cd /usr/src/myapp/src && go run service.go
//go run service.go
//cd /usr/src/myapp/src && go build -o wordToHtml
