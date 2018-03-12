package channelstructures

type ChanNode struct {
	name  string
	links string
	edges []*ChanNode
}

type Queue struct {
	enqChannel chan *ChanNode
	deqChannel chan *ChanNode
	queue      []*ChanNode
}

type VisitMessage struct {
	op   string
	node *ChanNode
}

type VisitedTracker struct {
	visited          map[*ChanNode]bool
	isVisitedChannel chan *ChanNode
	visitChannel     chan *ChanNode
}

func (vt *VisitedTracker) Visit(n *ChanNode) {
	vt.visitChannel <- n
}

func (vt *VisitedTracker) IsVisited(n *ChanNode) {
	vt.isVisitedChannel <- n
}

func (vt *VisitedTracker) Listen() {
	var n *ChanNode
	for {
		select {
		case vt.isVisitedChannel <- n:

		}
	}
}
