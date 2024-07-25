package models

import (
	"testing"
)

func TestSameEdge(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "A")
	if !g.vertices["A"]["A"] {
		t.Fatal("expected edge between A and B")
	}
}

func TestAddEdge(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	if !g.vertices["A"]["B"] || !g.vertices["B"]["A"] {
		t.Fatal("expected edge between A and B")
	}
	if !g.vertices["B"]["C"] || !g.vertices["C"]["B"] {
		t.Fatal("expected edge between B and C")
	}
	if g.vertices["A"]["C"] || g.vertices["C"]["A"] {
		t.Fatal("did not expected edge between A and C")
	}
}

func TestShortestPath(t *testing.T) {
	g := NewGraph()
	// A-B-C-D
	//  \   /
	//    E
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "D")
	g.AddEdge("D", "E")
	g.AddEdge("A", "E")

	//path1
	path, err := g.ShortestPath("A", "E")
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	expectedPath := []string{"A", "E"}
	if !equalPaths(path, expectedPath) {
		t.Fatalf("expected path %v, got %v", expectedPath, path)
	}

	//path2
	path, err = g.ShortestPath("A", "D")
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}

	expectedPath = []string{"A", "E", "D"}
	if !equalPaths(path, expectedPath) {
		t.Fatalf("expected path %v, got %v", expectedPath, path)
	}

	//no path
	path, err = g.ShortestPath("A", "Z")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if path != nil {
		t.Fatalf("expected nil path, got %v", path)
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
