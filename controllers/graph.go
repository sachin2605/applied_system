package controllers

import (
	"net/http"

	"applied_system/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Edge struct {
	V1 string `json:"v1"`
	V2 string `json:"v2"`
}

type GraphController struct {
	GraphMap map[string]models.Graph `json:"graph_map"`
}

func NewGraphController() *GraphController {
	return &GraphController{
		GraphMap: make(map[string]models.Graph),
	}
}

func (g *GraphController) PostGraphHandler(c *gin.Context) {
	var edges []Edge
	if err := c.ShouldBindJSON(&edges); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	graph := models.NewGraph()
	for _, edge := range edges {
		graph.AddEdge(edge.V1, edge.V2)
	}

	id := uuid.New().String()
	g.GraphMap[id] = *graph

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (g *GraphController) GetShortestPathHandler(c *gin.Context) {
	id := c.Param("id")

	graph, exist := g.GraphMap[id]
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{"error": "graph not found"})
		return
	}

	start := c.Query("start")
	end := c.Query("end")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end vertices required"})
		return
	}

	path, err := graph.ShortestPath(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"path": path})
}

func (g *GraphController) DeleteGraphHandler(c *gin.Context) {
	id := c.Param("id")

	_, exist := g.GraphMap[id]
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{"error": "graph not found"})
		return
	}
	delete(g.GraphMap, id)
	c.JSON(http.StatusOK, gin.H{"message": "graph deleted"})
}
