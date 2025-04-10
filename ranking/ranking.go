package ranking

import (
	"fmt"
	"strings"

	"github.com/benbjohnson/immutable"

	"github.com/Snawoot/deprank/graph"
)

type Ranking struct {
	ranks *immutable.List[graph.NodeSet]
}

func NewRanking() Ranking {
	return Ranking{
		ranks: immutable.NewList[graph.NodeSet](),
	}
}

func (r Ranking) Append(n *graph.Node) Ranking {
	return Ranking{
		ranks: r.ranks.Append(graph.NewNodeSet(n)),
	}
}

func (r Ranking) Merge(o Ranking) Ranking {
	// assume r is greater or swap otherwise
	if r.ranks.Len() < o.ranks.Len() {
		return o.Merge(r)
	}
	commonLength := min(r.ranks.Len(), o.ranks.Len())
	totalLength := max(r.ranks.Len(), o.ranks.Len())
	builder := immutable.NewListBuilder[graph.NodeSet]()
	for i := 0; i < commonLength; i++ {
		builder.Append(graph.MergeNodeSets(r.ranks.Get(i), o.ranks.Get(i)))
	}
	for i := commonLength; i < totalLength; i++ {
		builder.Append(r.ranks.Get(i))
	}
	return Ranking{
		ranks: builder.List(),
	}
}

func (r Ranking) String() string {
	b := new(strings.Builder)
	fmt.Fprintf(b, "Ranking<len=%d>:\n", r.ranks.Len())
	itr := r.ranks.Iterator()
	for !itr.Done() {
		rank, nodeSet := itr.Next()
		fmt.Fprintf(b, "\tRank %d:\n", rank)
		setItr := nodeSet.Iterator()
		for !setItr.Done() {
			elem, _ := setItr.Next()
			fmt.Fprintf(b, "\t\t%s\n", elem.Name)
		}
	}
	return b.String()
}

func rankGraph(root *graph.Node, visited graph.NodeSet) (Ranking, error) {
	if visited.Has(root) {
		return Ranking{}, fmt.Errorf("loop detected: %v", root)
	}
	if root == nil {
		return NewRanking(), nil
	}
	childrenRanking := NewRanking()
	visited = visited.Add(root)
	for _, child := range root.Children {
		r, err := rankGraph(child, visited)
		if err != nil {
			return Ranking{}, err
		}
		childrenRanking = childrenRanking.Merge(r)
	}
	return childrenRanking.Append(root), nil
}

func RankGraph(root *graph.Node) (Ranking, error) {
	visited := graph.NewNodeSet()
	return rankGraph(root, visited)
}
