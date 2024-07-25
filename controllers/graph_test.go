package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostGraph(t *testing.T) {

	graphCtrl := NewGraphController()
	router := gin.Default()
	router.POST("/graphs", graphCtrl.PostGraphHandler)

	// Define the graph edges to post
	edges := []map[string]string{
		{"v1": "A", "v2": "B"},
		{"v1": "B", "v2": "C"},
	}

	// Convert edges to JSON
	jsonEdges, _ := json.Marshal(edges)

	req, _ := http.NewRequest("POST", "/graphs", bytes.NewBuffer(jsonEdges))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}

	if _, ok := response["id"]; !ok {
		t.Fatal("expected response to contain 'id'")
	}
}

func BenchmarkPostGraph(b *testing.B) {
	graphCtrl := NewGraphController()
	router := gin.Default()
	router.POST("/graphs", graphCtrl.PostGraphHandler)

	// Define the graph edges to post
	edges := []map[string]string{
		{"v1": "A", "v2": "B"},
		{"v1": "B", "v2": "C"},
		{"v1": "C", "v2": "D"},
		{"v1": "D", "v2": "E"},
		{"v1": "E", "v2": "A"},
	}
	jsonEdges, _ := json.Marshal(edges)

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/graphs", bytes.NewBuffer(jsonEdges))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			b.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
		}
	}
}

func TestGetShortestPath(t *testing.T) {

	graphCtrl := NewGraphController()
	router := gin.Default()
	router.POST("/graphs", graphCtrl.PostGraphHandler)
	router.GET("/graphs/:id/shortest-path", graphCtrl.GetShortestPathHandler)

	// A-b-c-D
	edges := []map[string]string{
		{"v1": "A", "v2": "B"},
		{"v1": "B", "v2": "C"},
		{"v1": "C", "v2": "D"},
	}
	jsonEdges, _ := json.Marshal(edges)
	req, _ := http.NewRequest("POST", "/graphs", bytes.NewBuffer(jsonEdges))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var response map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}
	graphID := response["id"]

	req, _ = http.NewRequest("GET", "/graphs/"+graphID+"/shortest-path?start=A&end=D", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var pathResponse map[string][]string
	if err := json.Unmarshal(resp.Body.Bytes(), &pathResponse); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}

	expectedPath := []string{"A", "B", "C", "D"}
	if !equalPaths(pathResponse["path"], expectedPath) {
		t.Fatalf("expected path %v, got %v", expectedPath, pathResponse["path"])
	}

	//Test invalid Id
	graphID = "random-test-id"
	req, _ = http.NewRequest("GET", "/graphs/"+graphID+"/shortest-path?start=A&end=D", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

}

func TestDeleteGraph(t *testing.T) {

	graphCtrl := NewGraphController()
	router := gin.Default()
	router.POST("/graphs", graphCtrl.PostGraphHandler)
	router.GET("/graphs/:id/shortest-path", graphCtrl.GetShortestPathHandler)
	router.DELETE("/graphs/:id", graphCtrl.DeleteGraphHandler)

	edges := []map[string]string{
		{"v1": "A", "v2": "B"},
		{"v1": "B", "v2": "C"},
	}
	jsonEdges, _ := json.Marshal(edges)
	req, _ := http.NewRequest("POST", "/graphs", bytes.NewBuffer(jsonEdges))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var response map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}
	graphID := response["id"]

	// Now test deletion
	req, _ = http.NewRequest("DELETE", "/graphs/"+graphID, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	var deleteResponse map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &deleteResponse); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}

	// Verify that the graph is deleted by attempting to get the shortest path
	req, _ = http.NewRequest("GET", "/graphs/"+graphID+"/shortest-path?start=A&end=B", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.Code)
	}

	var errorResponse map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &errorResponse); err != nil {
		t.Fatalf("expected valid JSON response, got error: %v", err)
	}

	if errorResponse["error"] != "graph not found" {
		t.Fatalf("expected error 'graph not found', got %v", errorResponse["error"])
	}
}

func equalPaths(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
