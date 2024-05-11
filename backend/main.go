package main

import (
	"backend/internal/graph"
	_ "embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	applicationJson = "application/json"
	contentType     = "Content-Type"
	maxMem          = 1 << 16
	uploadFailed    = "Upload failed [%s]"
)

// Graph contains the graph's root node and provide the methods used by the backend http server
type Graph struct {
	root *graph.Node
}

// NewGraph create a new graph
func NewGraph() *Graph {
	root, err := graph.NewLexeme("0", "ens", "")
	if err != nil {
		log.Fatalf("Failed to create the root node [%v]", err)
	}
	return &Graph{root: root}
}

// graph returns the graph
func (g *Graph) graph(context *gin.Context) {
	json, err := g.root.String()
	if err != nil {
		msg := fmt.Sprintf("Failed to generate the JSON string [%s]", err)
		log.Error(msg)
		handleFailedRequest(context, http.StatusInternalServerError, msg)
		return
	}
	context.Header(contentType, applicationJson)
	context.String(http.StatusOK, json)
}

// addChildToRootNode adds a child to the root node
func (g *Graph) addChildToRootNode(context *gin.Context) {
	node := &graph.Node{}
	err := context.BindJSON(node)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse the JSON payload [%s]", err)
		log.Error(msg)
		handleFailedRequest(context, http.StatusBadRequest, msg)
		return
	}
	root, err := g.root.AddNode("0", node)
	if err != nil {
		msg := fmt.Sprintf("Failed to add the node to the graph's root [%s]", err)
		log.Error(msg)
		handleFailedRequest(context, http.StatusBadRequest, msg)
		return
	}
	g.root = root
}

// deleteNode deletes a node
func (g *Graph) deleteNode(context *gin.Context) {
	parent := context.Param("parent")
	target := context.Param("node")
	root, err := g.root.RemoveNode(parent, target)
	if err != nil {
		msg := fmt.Sprintf("Failed to remove the node %q from its parent %q [%s]", target, parent, err)
		log.Error(msg)
		handleFailedRequest(context, http.StatusBadRequest, msg)
		return
	}
	g.root = root
}

// updateNode updates a node. The updated node may include a new child node
func (g *Graph) updateNode(context *gin.Context) {
	parent := context.Param("parent")
	node := &graph.Node{}
	err := context.BindJSON(node)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse the JSON payload [%s]", err)
		log.Error(msg)
		handleFailedRequest(context, http.StatusBadRequest, msg)
		return
	}
	root, err := g.root.UpdateNode(parent, node)
	if err != nil {
		msg := fmt.Sprintf("Failed to add the node %q to its parent %q [%s]", node.Name, parent, err)
		log.Error(msg)
		handleFailedRequest(context, http.StatusBadRequest, msg)
		return
	}
	g.root = root
}

// healthCheck returns a "200 OK" response to indicate that the backend service is available
func healthCheck(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}

// upload uploads a graph
func (g *Graph) upload(context *gin.Context) {
	fh, err := context.FormFile("file")
	if err != nil {
		msg := fmt.Sprintf(uploadFailed, err)
		log.Error(msg)
		handleFailedRequest(context, http.StatusBadRequest, msg)
		return
	}
	log.Debug("File uploaded")
	file, err := fh.Open()
	if err != nil {
		msg := fmt.Sprintf(uploadFailed, err)
		log.Error(msg)
		handleFailedRequest(context, http.StatusInternalServerError, msg)
		return
	}
	defer file.Close()

	bytes := make([]byte, fh.Size)
	n, err := file.Read(bytes)
	if err != nil {
		msg := fmt.Sprintf(uploadFailed, err)
		log.Error(msg)
		handleFailedRequest(context, http.StatusInternalServerError, msg)
		return
	}
	log.Debugf("Read %d bytes", n)
	root, err := g.root.Parse(bytes)
	if err != nil {
		msg := fmt.Sprintf(uploadFailed, err)
		log.Error(msg)
		handleFailedRequest(context, http.StatusInternalServerError, msg)
		return
	}
	g.root = root
	context.Status(http.StatusOK)
}

// handleFailedRequest writes a response with the given error code and message
func handleFailedRequest(context *gin.Context, statusCode int, message string) {
	context.JSON(statusCode, gin.H{
		"status":  fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		"message": message,
	})
}

// main starts the backend http server
func main() {
	g := NewGraph()

	router := gin.Default()
	router.MaxMultipartMemory = maxMem
	router.Use(cors.Default())

	router.GET("/apis", healthCheck)
	router.GET("/apis/graph", g.graph)
	router.GET("/apis/health", healthCheck)
	router.PUT("/apis/nodes", g.addChildToRootNode)
	router.PUT("/apis/nodes/:parent", g.updateNode)
	router.DELETE("/apis/nodes/:parent/:node", g.deleteNode)
	router.POST("/apis/upload", g.upload)

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
