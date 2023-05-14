package proxyhdl

import (
	"io/ioutil"
	"net/http"
	"strings"
	"web-app-test/ports"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	proxyService ports.ProxyService
}

func NewHTTPHandler(proxyService ports.ProxyService) *HTTPHandler {
	return &HTTPHandler{
		proxyService: proxyService,
	}
}

func (hdl *HTTPHandler) Get(ctx *gin.Context) {
	resp, err := hdl.proxyService.GetModified(ctx.Request.URL.Path)
	// handle errors
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "my custom error",
		})
		return
	}

	// handle non JSON responses
	if !strings.Contains(resp.Header.Get("Content-Type"), gin.MIMEJSON) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "this is not JSON",
		})
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ctx.Data(http.StatusOK, gin.MIMEJSON, body)
}
