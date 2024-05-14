package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// printGraph returns a simplified string representation of this graph
func (server *HttpServer) printGraph(context *gin.Context) {
	context.Header(contentType, textPlain)
	context.String(http.StatusOK, server.g.Root.Stringify())
}
