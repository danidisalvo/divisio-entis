package rest

import (
	"backend/internal/graph"
	graphErrors "backend/internal/graph/errors"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	address         = ":8080"
	applicationJson = "application/json"
	contentType     = "Content-Type"
	maxMem          = 1 << 16
	textPlain       = "text/plain"
	uploadFailed    = "Upload failed [%s]"
)

type HttpServer struct {
	g *graph.Graph
}

func NewHttpServer(g *graph.Graph) *HttpServer {
	return &HttpServer{g: g}
}

func (server *HttpServer) StartHttpServer() error {
	router := gin.Default()
	router.HandleMethodNotAllowed = true
	router.MaxMultipartMemory = maxMem
	router.Use(cors.Default())

	router.GET("/apis", healthCheck)
	router.GET("/apis/health", healthCheck)

	router.DELETE("/apis/graph", server.deleteGraph)
	router.GET("/apis/graph", server.getGraph)
	router.GET("/apis/graph/print", server.printGraph)
	router.PUT("/apis/nodes", server.addChildToRootNode)
	router.PUT("/apis/nodes/:parent", server.updateNode)
	router.DELETE("/apis/nodes/:parent/:node", server.deleteNode)
	router.POST("/apis/nodes/:parent/:node/:newParent", server.moveNode)
	router.POST("/apis/upload", server.upload)

	err := router.Run(address)
	if err != nil {
		return err
	}
	return nil
}

// healthCheck returns a "200 OK" response to indicate that the backend service is available
func healthCheck(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}

// handleFailedRequest writes a response with the given error code and message
func handleFailedRequest(context *gin.Context, err error, message string) {

	var duplicatedNodeError *graphErrors.DuplicatedNodeError
	var illegalArgumentError *graphErrors.IllegalArgumentError
	var nodeNotFoundError *graphErrors.NodeNotFoundError

	var statusCode int
	if errors.As(err, &duplicatedNodeError) || errors.As(err, &illegalArgumentError) {
		statusCode = http.StatusBadRequest
	} else if errors.As(err, &nodeNotFoundError) {
		statusCode = http.StatusNotFound
	} else {
		statusCode = http.StatusInternalServerError
	}

	context.JSON(statusCode, gin.H{
		"status":  fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		"message": message,
	})
}
