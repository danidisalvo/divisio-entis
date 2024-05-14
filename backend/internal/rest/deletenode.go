package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// deleteNode deletes a node
func (server *HttpServer) deleteNode(context *gin.Context) {
	parent := context.Param("parent")
	target := context.Param("node")
	root, err := server.g.Root.RemoveNode(parent, target)
	if err != nil {
		msg := fmt.Sprintf("Failed to remove the node %q from its parent %q [%s]", target, parent, err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	server.g.Root = root
	server.g.Save()
}
