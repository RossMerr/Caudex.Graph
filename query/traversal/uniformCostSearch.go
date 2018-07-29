package traversal

import (
	"sort"

	"github.com/RossMerr/Caudex.Graph/query"
)

func UniformCostSearch(storage query.Storage, predicates []query.Predicate, frontier *query.Frontier) bool {
	if frontier.Len() > 0 {
		queue := frontier.Pop()
		depth := len(queue.Parts)
		searchDepth := len(predicates)
		part := queue.Parts[depth-1]

		if _, ok := frontier.Explored[part.UUID]; !ok {
			frontier.Explored[part.UUID] = true
			pv := predicates[depth-1]
			if variable, p, _ := pv(part.UUID, nil, depth-1); p == query.Matched {
				queue.Parts[depth-1].Variable = variable
				frontier.AppendKeyValue(queue, part.UUID, part.Variable)
				sort.Sort(frontier)
				return depth == searchDepth
			}

		}

		if pe := predicates[depth]; pe != nil {
			iterator := storage.Edges(part.UUID)
			for kv, hasEdges := iterator(); hasEdges; kv, hasEdges = iterator() {
				if _, ok := frontier.Explored[kv]; !ok {
					if variable, p, weight := pe(part.UUID, kv, depth); p == query.Visiting || p == query.Matching {
						frontier.AppendEdgeKeyValue(queue, kv, variable, weight)
					}
				}
			}
		}

		sort.Sort(frontier)
	}
	return false
}
