package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// upload uploads a graph
func (server *HttpServer) upload(context *gin.Context) {
	fh, err := context.FormFile("file")
	if err != nil {
		msg := fmt.Sprintf(uploadFailed, err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	log.Debug("File uploaded")
	file, err := fh.Open()
	if err != nil {
		msg := fmt.Sprintf(uploadFailed, err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	defer file.Close()

	bytes := make([]byte, fh.Size)
	n, err := file.Read(bytes)
	if err != nil {
		msg := fmt.Sprintf(uploadFailed, err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	log.Debugf("Read %d bytes", n)
	server.g.Root, err = server.g.Root.Parse(bytes)
	if err != nil {
		msg := fmt.Sprintf(uploadFailed, err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	server.g.Save()
	context.Status(http.StatusOK)
}
