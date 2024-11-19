package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"picmask/imageutils"

	"github.com/gin-gonic/gin"
)

type Img struct {
	filename string
}

func (img Img) GetImg(c *gin.Context) {
	// 获取当前脚本的绝对路径
	currentFilePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 假设项目根目录是当前脚本的父目录的父目录
	projectRoot := filepath.Dir(currentFilePath)
	projectRoot = filepath.Dir(projectRoot)

	id := c.Param("id")
	operate := c.Param("operate")
	passwd := c.Param("passwd")
	newpsd, _ := strconv.ParseFloat(passwd, 64)
	fmt.Println(id, operate, newpsd, passwd)
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form file err: %s", err.Error()))
		return
	}
	// 获取文件名称
	fileName := file.Filename
	dst := filepath.Join(projectRoot, fileName)
	fmt.Println(dst)
	// 创建目录
	// os.MkdirAll("./uploads", os.ModePerm)

	// 保存文件到指定路径
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}
	output := filepath.Join(projectRoot, id+"_"+fileName)
	if operate == "encrypt" {
		// 加密
		if imageutils.ProcessImage(dst, newpsd, true, output) != nil {
			fmt.Println("error occur in encrpyt")
		}
	} else {
		// 解密
		if imageutils.ProcessImage(dst, newpsd, false, output) != nil {
			fmt.Println("error occur in decrpyt")
		}
	}
	fmt.Println(output)
	// c.FileAttachment(output, id+"_"+fileName)
	// c.File(output)
	// imageutils.ProcessImage(dst,operate)
	// c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", fileName))
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success dealed the image",
		"file":    fileName,
	})
}

func (img Img) DownImg(c *gin.Context) {
	// 获取当前脚本的绝对路径
	currentFilePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 假设项目根目录是当前脚本的父目录的父目录
	projectRoot := filepath.Dir(currentFilePath)
	projectRoot = filepath.Dir(projectRoot)
	id := c.Param("id")
	filename := c.Param("filename")
	filePath := filepath.Join(projectRoot, id+"_"+filename)
	fmt.Println(filePath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Name", filePath)
	c.File(filePath)
}
