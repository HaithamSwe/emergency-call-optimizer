package main

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
)

// ---------------------------
// Data Structures
// ---------------------------

// Point represents a 2D coordinate.
type Point struct {
	X, Y float64
}

// TelecomSite represents a telecom site with an ID and location.
type TelecomSite struct {
	ID       string
	Location Point
}

// Hub represents a hub site with an ID and location.
type Hub struct {
	ID       string
	Location Point
}

// EmergencyCenter represents an emergency center with an ID and location.
type EmergencyCenter struct {
	ID       string
	Location Point
}

// Validate validates the data for a TelecomSite.
func (site *TelecomSite) Validate() error {
	if site.ID == "" || site.Location == (Point{}) {
		return errors.New("invalid TelecomSite data")
	}
	return nil
}

// Validate validates the data for a Hub.
func (hub *Hub) Validate() error {
	if hub.ID == "" || hub.Location == (Point{}) {
		return errors.New("invalid Hub data")
	}
	return nil
}

// Validate validates the data for an EmergencyCenter.
func (ec *EmergencyCenter) Validate() error {
	if ec.ID == "" || ec.Location == (Point{}) {
		return errors.New("invalid EmergencyCenter data")
	}
	return nil
}

// CalculateDistance calculates the Euclidean distance between two points.
func CalculateDistance(a, b Point) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

// ---------------------------
// K-Nearest Neighbors (KNN)
// ---------------------------

// FindNearestEmergencyCenter finds the nearest emergency center from a given telecom site.
func FindNearestEmergencyCenter(site TelecomSite, centers []EmergencyCenter) (EmergencyCenter, error) {
	if err := site.Validate(); err != nil {
		return EmergencyCenter{}, err
	}

	var nearestCenter EmergencyCenter
	minDistance := math.MaxFloat64

	for _, center := range centers {
		dist := CalculateDistance(site.Location, center.Location)
		if dist < minDistance {
			minDistance = dist
			nearestCenter = center
		}
	}

	return nearestCenter, nil
}

// FindNearestHubs finds the nearest N hubs from a given point.
func FindNearestHubs(point Point, hubs []Hub, N int) ([]Hub, error) {
	if N <= 0 {
		return nil, errors.New("N must be greater than 0")
	}
	if len(hubs) < N {
		return nil, errors.New("not enough hubs to find the nearest N hubs")
	}

	// Simple linear search to find the nearest N hubs
	type hubDistance struct {
		hub      Hub
		distance float64
	}

	var distances []hubDistance
	for _, hub := range hubs {
		dist := CalculateDistance(point, hub.Location)
		distances = append(distances, hubDistance{hub, dist})
	}

	// Sort the hubs by distance
	for i := 0; i < len(distances)-1; i++ {
		for j := i + 1; j < len(distances); j++ {
			if distances[j].distance < distances[i].distance {
				distances[i], distances[j] = distances[j], distances[i]
			}
		}
	}

	nearestHubs := make([]Hub, N)
	for i := 0; i < N; i++ {
		nearestHubs[i] = distances[i].hub
	}

	return nearestHubs, nil
}

// ---------------------------
// Graph Structures and Dijkstra's Algorithm
// ---------------------------

// Graph represents a graph with nodes and edges.
type Graph struct {
	Nodes     map[string]bool
	Neighbors map[string]map[string]float64
}

// NewGraph initializes and returns a new Graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes:     make(map[string]bool),
		Neighbors: make(map[string]map[string]float64),
	}
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(node string) {
	if _, exists := g.Nodes[node]; !exists {
		g.Nodes[node] = true
		g.Neighbors[node] = make(map[string]float64)
	}
}

// AddEdge adds an edge to the graph with a given cost.
func (g *Graph) AddEdge(from, to string, cost float64) {
	g.AddNode(from)
	g.AddNode(to)
	g.Neighbors[from][to] = cost
	// Assuming undirected graph; add both connections
	g.Neighbors[to][from] = cost
}

// Item represents an item in the priority queue.
type Item struct {
	Value    string
	Priority float64
	Index    int
}

// MinPriorityQueue implements a priority queue for Dijkstra's Algorithm.
type MinPriorityQueue []*Item

func (pq MinPriorityQueue) Len() int { return len(pq) }

func (pq MinPriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq MinPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *MinPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *MinPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // Avoid memory leak
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// Dijkstra finds the shortest paths from a starting node to all other nodes in the graph.
func Dijkstra(graph *Graph, start string) (map[string]float64, map[string]string, error) {
	if graph == nil {
		return nil, nil, errors.New("graph is nil")
	}
	if _, exists := graph.Nodes[start]; !exists {
		return nil, nil, fmt.Errorf("start node %s does not exist in the graph", start)
	}

	distances := make(map[string]float64)
	previous := make(map[string]string)
	for node := range graph.Nodes {
		distances[node] = math.Inf(1)
	}
	distances[start] = 0

	priorityQueue := &MinPriorityQueue{}
	heap.Init(priorityQueue)
	heap.Push(priorityQueue, &Item{Value: start, Priority: 0})

	visited := make(map[string]bool)

	for priorityQueue.Len() > 0 {
		currentItem := heap.Pop(priorityQueue).(*Item)
		current := currentItem.Value

		if visited[current] {
			continue
		}
		visited[current] = true

		for neighbor, cost := range graph.Neighbors[current] {
			alt := distances[current] + cost
			if alt < distances[neighbor] {
				distances[neighbor] = alt
				previous[neighbor] = current
				heap.Push(priorityQueue, &Item{Value: neighbor, Priority: alt})
			}
		}
	}

	return distances, previous, nil
}

// ReconstructPath reconstructs the shortest path from start to target using the previous map.
func ReconstructPath(previous map[string]string, start, target string) ([]string, error) {
	path := []string{}
	current := target

	for current != start {
		path = append([]string{current}, path...)
		prev, exists := previous[current]
		if !exists {
			return nil, fmt.Errorf("no path found from %s to %s", start, target)
		}
		current = prev
	}
	path = append([]string{start}, path...)
	return path, nil
}

// ---------------------------
// Network Initialization
// ---------------------------

// InitializeSampleData initializes sample telecom sites, hubs, and emergency centers.
func InitializeSampleData() ([]TelecomSite, []Hub, []EmergencyCenter) {
	// Sample Telecom Sites
	telecomSites := []TelecomSite{
		{ID: "TS1", Location: Point{X: 10, Y: 20}},
		{ID: "TS2", Location: Point{X: 25, Y: 35}},
		{ID: "TS3", Location: Point{X: 40, Y: 50}},
	}

	// Sample Hubs
	hubs := []Hub{
		{ID: "H1", Location: Point{X: 12, Y: 22}},
		{ID: "H2", Location: Point{X: 28, Y: 38}},
		{ID: "H3", Location: Point{X: 45, Y: 55}},
		{ID: "H4", Location: Point{X: 50, Y: 60}},
		{ID: "H5", Location: Point{X: 5, Y: 15}},
	}

	// Sample Emergency Centers
	emergencyCenters := []EmergencyCenter{
		{ID: "EC1", Location: Point{X: 15, Y: 25}},
		{ID: "EC2", Location: Point{X: 30, Y: 40}},
		{ID: "EC3", Location: Point{X: 35, Y: 45}},
	}

	return telecomSites, hubs, emergencyCenters
}

// BuildNetworkGraph builds the network graph based on the infrastructure assumptions.
func BuildNetworkGraph(telecomSites []TelecomSite, hubs []Hub, emergencyCenters []EmergencyCenter) (*Graph, error) {
	graph := NewGraph()

	// Connect each Telecom Site to the nearest Hub
	for _, site := range telecomSites {
		nearestHubs, err := FindNearestHubs(site.Location, hubs, 1)
		if err != nil {
			return nil, fmt.Errorf("error finding nearest hub for %s: %v", site.ID, err)
		}
		for _, hub := range nearestHubs {
			cost := CalculateDistance(site.Location, hub.Location)
			graph.AddEdge(site.ID, hub.ID, cost)
		}
	}

	// Connect each Hub to its two nearest Hubs
	for _, hub := range hubs {
		nearestHubs, err := FindNearestHubs(hub.Location, hubs, 3) // Including itself
		if err != nil {
			return nil, fmt.Errorf("error finding nearest hubs for %s: %v", hub.ID, err)
		}
		// Skip the first one as it is the hub itself
		for i := 1; i < len(nearestHubs) && i <= 2; i++ {
			neighborHub := nearestHubs[i]
			cost := CalculateDistance(hub.Location, neighborHub.Location)
			graph.AddEdge(hub.ID, neighborHub.ID, cost)
		}
	}

	// Connect each Emergency Center to the five nearest Hubs
	for _, ec := range emergencyCenters {
		nearestHubs, err := FindNearestHubs(ec.Location, hubs, 5)
		if err != nil {
			return nil, fmt.Errorf("error finding nearest hubs for %s: %v", ec.ID, err)
		}
		for _, hub := range nearestHubs {
			cost := CalculateDistance(ec.Location, hub.Location)
			graph.AddEdge(ec.ID, hub.ID, cost)
		}
	}

	return graph, nil
}

// ---------------------------
// Main Function
// ---------------------------

func main() {
	// Initialize sample data
	telecomSites, hubs, emergencyCenters := InitializeSampleData()

	// Build the network graph
	graph, err := BuildNetworkGraph(telecomSites, hubs, emergencyCenters)
	if err != nil {
		fmt.Printf("Error building network graph: %v\n", err)
		return
	}

	// Display the graph connections
	fmt.Println("Network Graph Connections:")
	for node, neighbors := range graph.Neighbors {
		for neighbor, cost := range neighbors {
			fmt.Printf("%s <--> %s : %.2f\n", node, neighbor, cost)
		}
	}
	fmt.Println()

	// Find the nearest emergency center for each telecom site and find the shortest path
	for _, site := range telecomSites {
		nearestEC, err := FindNearestEmergencyCenter(site, emergencyCenters)
		if err != nil {
			fmt.Printf("Error finding nearest emergency center for %s: %v\n", site.ID, err)
			continue
		}
		fmt.Printf("Nearest Emergency Center for %s is %s\n", site.ID, nearestEC.ID)

		// Find shortest path from Telecom Site to Emergency Center using Dijkstra's Algorithm
		// Since the graph is undirected, we can start from the Telecom Site and find the path to EC
		distances, previous, err := Dijkstra(graph, site.ID)
		if err != nil {
			fmt.Printf("Error running Dijkstra's algorithm for %s: %v\n", site.ID, err)
			continue
		}

		// Check if there's a path to the nearest emergency center
		if distances[nearestEC.ID] == math.Inf(1) {
			fmt.Printf("No path found from %s to %s\n\n", site.ID, nearestEC.ID)
			continue
		}

		// Reconstruct the path
		path, err := ReconstructPath(previous, site.ID, nearestEC.ID)
		if err != nil {
			fmt.Printf("Error reconstructing path from %s to %s: %v\n", site.ID, nearestEC.ID, err)
			continue
		}

		// Display the shortest path and distance
		fmt.Printf("Shortest path from %s to %s: %v\n", site.ID, nearestEC.ID, path)
		fmt.Printf("Total distance: %.2f\n\n", distances[nearestEC.ID])
	}

	// Example Output for a Specific Telecom Site
	specificSite := telecomSites[0]
	nearestEC, err := FindNearestEmergencyCenter(specificSite, emergencyCenters)
	if err != nil {
		fmt.Printf("Error finding nearest emergency center for %s: %v\n", specificSite.ID, err)
		return
	}
	fmt.Printf("Example: Nearest Emergency Center for %s is %s\n", specificSite.ID, nearestEC.ID)

	// Find shortest path from specific telecom site to its nearest emergency center
	distances, previous, err := Dijkstra(graph, specificSite.ID)
	if err != nil {
		fmt.Printf("Error running Dijkstra's algorithm for %s: %v\n", specificSite.ID, err)
		return
	}

	if distances[nearestEC.ID] == math.Inf(1) {
		fmt.Printf("No path found from %s to %s\n", specificSite.ID, nearestEC.ID)
		return
	}

	path, err := ReconstructPath(previous, specificSite.ID, nearestEC.ID)
	if err != nil {
		fmt.Printf("Error reconstructing path from %s to %s: %v\n", specificSite.ID, nearestEC.ID, err)
		return
	}

	fmt.Printf("Shortest path from %s to %s: %v\n", specificSite.ID, nearestEC.ID, path)
	fmt.Printf("Total distance: %.2f\n", distances[nearestEC.ID])
}
