package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
)

var Imfile imfile

type imfile struct {
}

func (i *imfile) ImportFile(c *gin.Context) {
	text, err := service.Imfile.ImportFile()
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, text)
}
