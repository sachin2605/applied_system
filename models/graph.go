package models

import (
	"errors"
	"sync"
)

type Graph struct {
	vertices map[string]map[string]bool
	mu       sync.Mutex
}

func NewGraph() *Graph {
	return &Graph{
		vertices: make(map[string]map[string]bool),
	}
}

func (g *Graph) AddEdge(v1, v2 string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.vertices[v1] == nil {
		g.vertices[v1] = make(map[string]bool)
	}
	if g.vertices[v2] == nil {
		g.vertices[v2] = make(map[string]bool)
	}

	g.vertices[v1][v2] = true
	g.vertices[v2][v1] = true
}

func (g *Graph) NumVertices() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.vertices)
}

func (g *Graph) NumEdges() int {
	g.mu.Lock()
	defer g.mu.Unlock()

	count := 0
	for _, neighbors := range g.vertices {
		count += len(neighbors)
	}
	return count / 2 // Each edge is counted twice
}

func (g *Graph) ShortestPath(start, end string) ([]string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.vertices[start] == nil || g.vertices[end] == nil {
		return nil, errors.New("one or both vertices not found")
	}

	visited := make(map[string]bool)
	queue := [][]string{{start}}
	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		node := path[len(path)-1]
		if node == end {
			return path, nil
		}
		if !visited[node] {
			visited[node] = true
			for neighbor := range g.vertices[node] {
				newPath := append([]string{}, path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	return nil, errors.New("no path found")
}
