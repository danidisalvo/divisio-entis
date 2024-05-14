package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// deleteGraph clear, i.e., resets, the graph
func (server *HttpServer) deleteGraph(context *gin.Context) {
	server.g.Clear()
	context.Writer.WriteHeader(http.StatusNoContent)
}
