package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// moveNode moves a node
func (server *HttpServer) moveNode(context *gin.Context) {
	parent := context.Param("parent")
	target := context.Param("node")
	newParent := context.Param("newParent")
	root, err := server.g.Root.MoveNode(parent, target, newParent)
	if err != nil {
		msg := fmt.Sprintf("Failed to move the node %q from %q to %q [%s]", target, parent, newParent, err)
		log.Error(msg)
		handleFailedRequest(context, err, msg)
		return
	}
	server.g.Root = root
	server.g.Save()
	server.getGraph(context)
}
