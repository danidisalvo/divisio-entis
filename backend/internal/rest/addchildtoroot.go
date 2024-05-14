package rest

import (
	"backend/internal/graph"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// addChildToRootNode adds a child to the root node
func (server *HttpServer) addChildToRootNode(context *gin.Context) {
	node := &graph.Node{}
	err := context.BindJSON(node)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse the JSON payload [%s]", err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	root, err := server.g.Root.AddNode("0", node)
	if err != nil {
		msg := fmt.Sprintf("Failed to add the node to the graph's root [%s]", err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	server.g.Root = root
	server.g.Save()
}
