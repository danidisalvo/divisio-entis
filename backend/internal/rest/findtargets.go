package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// findTargets returns the nodes to which the given node can be moved
func (server *HttpServer) findTargets(context *gin.Context) {
	node := context.Param("node")
	nodes, err := server.g.Root.FindTargetNodes(node)
	if err != nil {
		msg := fmt.Sprintf("Failed to find the target nodes of node %q [%s]", node, err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	bytes, err := json.Marshal(nodes)
	if err != nil {
		msg := fmt.Sprintf("Failed to serialize the nodes [%s]", err)
		log.Error(msg)
	}
	context.Header(contentType, applicationJson)
	context.String(http.StatusOK, string(bytes))
}
