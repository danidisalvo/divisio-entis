package rest

import (
	"backend/internal/graph"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// updateNode updates a node. The updated node may include a new child node
func (server *HttpServer) updateNode(context *gin.Context) {
	parent := context.Param("parent")
	node := &graph.Node{}
	err := context.BindJSON(node)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse the JSON payload [%s]", err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	root, err := server.g.Root.UpdateNode(parent, node)
	if err != nil {
		msg := fmt.Sprintf("Failed to add the node %q to its parent %q [%s]", node.Name, parent, err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	server.g.Root = root
	server.g.Save()
}
