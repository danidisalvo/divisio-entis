package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// graph returns the graph
func (server *HttpServer) getGraph(context *gin.Context) {
	json, err := server.g.Root.String()
	if err != nil {
		msg := fmt.Sprintf("Failed to generate the JSON string [%s]", err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	context.Header(contentType, applicationJson)
	context.String(http.StatusOK, json)
}
