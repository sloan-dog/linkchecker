package structures

import (
	"fmt"
	"strings"
	"sync"
)

// Node arbitrary graph node structure
type Node struct {
	name  string
	links string
	sync.Mutex
	edges []*Node
}

// NewNode Creates new node from @param name, links, edges
func NewNode(name string, links string, edges []*Node) *Node {
	return &Node{
		name,
		links,
		sync.Mutex{},
		edges,
	}
}

// GetLinks returns node's links prop
func (n *Node) GetLinks() []string {
	return strings.Split(n.links, "")
}

// AddEdge attaches an edge node to the target node's edges prop
func (n *Node) AddEdge(edge *Node) {
	n.edges = append(n.edges, edge)
}

// NodeConfig config strucutre for defining a node
type NodeConfig struct {
	name  string
	links string
	edges []string
}

// GraphConfig container structure for all node configs and entry points
type GraphConfig struct {
	NodeConfigs map[string]NodeConfig
	EntryPoints []string
}

// Graph wrapper around node pointers comprising a graph
type Graph struct {
	Nodes       []*Node
	EntryPoints []*Node
	Created     map[string]*Node
}

// ConfigGraphDefault creates an arbitrary default graph structure
func ConfigGraphDefault() Graph {
	config := GraphConfig{
		NodeConfigs: map[string]NodeConfig{
			"one": NodeConfig{
				name:  "one",
				links: "ABC",
				edges: []string{"two", "three"},
			},
			"two": NodeConfig{
				name:  "two",
				links: "BBA",
				edges: []string{"four"},
			},
			"three": NodeConfig{
				name:  "three",
				links: "DEF",
				edges: []string{"five", "six", "seven"},
			},
			"four": NodeConfig{
				name:  "four",
				links: "FHHG",
				edges: []string{},
			},
			"five": NodeConfig{
				name:  "five",
				links: "AAA",
				edges: []string{"seven", "eight"},
			},
			"six": NodeConfig{
				name:  "six",
				links: "AAA",
				edges: []string{"eight", "nine"},
			},
			"seven": NodeConfig{
				name:  "seven",
				links: "JBK",
				edges: []string{},
			},
			"eight": NodeConfig{
				name:  "eight",
				links: "BBB",
				edges: []string{"two"},
			},
			"nine": NodeConfig{
				name:  "nine",
				links: "KVOSD",
				edges: []string{"three"},
			},
		},
		EntryPoints: []string{"one", "nine"},
	}
	return ConfigGraph(&config)
}

// ConfigGraph generates a graph struct and creates all constituent nodes from a graphconfig
func ConfigGraph(config *GraphConfig) Graph {
	g := Graph{
		Nodes:       make([]*Node, 0, 20),
		EntryPoints: []*Node{},
		Created:     map[string]*Node{},
	}
	for _, nodeConfig := range config.NodeConfigs {
		CreateNodeFromConfig(nodeConfig.name, config, &g)
	}
	// create entry point refs
	for _, name := range config.EntryPoints {
		n := g.Created[name]
		g.EntryPoints = append(g.EntryPoints, n)
	}
	return g
}

// CreateNodeFromConfig creates a node structure from a config, recursively generating its dependency nodes
func CreateNodeFromConfig(name string, config *GraphConfig, grph *Graph) *Node {
	// if it doesn't exist create it
	var n *Node
	var edgeNode *Node
	if _, exists := grph.Created[name]; !exists {
		cfg, _ := config.NodeConfigs[name]
		n = NewNode(name, cfg.links, []*Node{})
		grph.Created[cfg.name] = n
		for _, edgeName := range cfg.edges {
			// create its edge nodes
			edgeNode = CreateNodeFromConfig(edgeName, config, grph)
			n.AddEdge(edgeNode)
		}
	}
	// always return the node even if it has already been created
	grph.Nodes = append(grph.Nodes, grph.Created[name])
	return grph.Created[name]
}

// Queue a thread safe fifo queue for holding node pointers
type Queue struct {
	sync.Mutex
	queue []*Node
}

// NewQueue returns a pointer to a new queue
func NewQueue() *Queue {
	return &Queue{
		queue: make([]*Node, 0, 20),
	}
}

// IsEmpty returns whether queue is empty (bool)
func (q *Queue) IsEmpty() bool {
	return len(q.queue) < 1
}

// Add puts a node pointer on the queue
func (q *Queue) Add(n *Node) {
	q.queue = append(q.queue, n)
}

// Pop removes the first node p from the queue (does not check for empty)
func (q *Queue) Pop() *Node {
	element := q.queue[0]
	q.queue = q.queue[1:]
	return element
}

// Counter thread safe tracker of link counts
type Counter struct {
	sync.Mutex
	counter map[string]int
}

func NewCounter() *Counter {
	return &Counter{
		counter: map[string]int{},
	}
}

// AddCountForLink increments count for a particular link
func (c *Counter) AddCountForLink(link string) {
	c.Lock()
	defer c.Unlock()
	_, exists := c.counter[link]
	if !exists {
		c.counter[link] = 0
	}
	c.counter[link]++
}

// GetCounts returns an array of strings representing a link and its count
func (c *Counter) GetCounts() []string {
	c.Lock()
	defer c.Unlock()
	arr := make([]string, 0, len(c.counter))
	for k, v := range c.counter {
		arr = append(arr, fmt.Sprintf("%v:%v", k, v))
	}
	return arr
}

// VisitedTracker tracks which nodes have been visited already
type VisitedTracker struct {
	sync.Mutex
	visited map[*Node]bool
}

func NewTracker() *VisitedTracker {
	return &VisitedTracker{
		visited: map[*Node]bool{},
	}
}

// isVisited is non thread safe checker for node having been visited
func (vt *VisitedTracker) isVisited(node *Node) bool {
	_, exists := vt.visited[node]
	return exists
}

// IsVisited is thread safe method for checking if node exists
func (vt *VisitedTracker) IsVisited(node *Node) bool {
	vt.Lock()
	defer vt.Unlock()
	return vt.isVisited(node)
}

// Visit is thread safe method which tags a node in the tracker as visited
func (vt *VisitedTracker) Visit(node *Node) {
	vt.Lock()
	defer vt.Unlock()
	vt.visited[node] = true
}

func CountLinks(grph *Graph, counter *Counter, queue *Queue, tracker *VisitedTracker, entryPoint *Node, name int) {
	queue.Add(entryPoint)
	for !queue.IsEmpty() {
		n := queue.Pop()
		if !tracker.IsVisited(n) {
			tracker.Visit(n)
			for _, link := range n.GetLinks() {
				counter.AddCountForLink(link)
			}
			for _, edge := range n.edges {
				queue.Add(edge)
			}
		}
	}
}
