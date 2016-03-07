package graph

import (
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/trustedanalytics/go-cf-lib/api"
	"github.com/trustedanalytics/go-cf-lib/types"
	"github.com/twmb/algoimpl/go/graph"
)

type GraphAPI struct {
	w *api.CfAPI
}

func NewGraphAPI() *GraphAPI {
	toReturn := new(GraphAPI)
	toReturn.w = api.NewCfAPI()
	return toReturn
}

// Returns a list of services and apps in application stack in reversed topological order
func (gr *GraphAPI) Discover(sourceAppGUID string) ([]types.Component, error) {
	sourceAppSummary, err := gr.w.GetAppSummary(sourceAppGUID)
	if err != nil {
		return nil, err
	}

	g := graph.New(graph.Directed)
	dg := NewDependencyGraph()
	root := dg.NewNode(g, sourceAppGUID, sourceAppSummary.Name, types.ComponentApp, nil, true)
	_ = dg.addDependenciesToGraph(g, root, sourceAppGUID)
	if dg.graphHasCycles(g) {
		return nil, errors.New("Graph has cycles and stack cannot be copied")
	} else {
		log.Infof("Graph has no cycles")
	}

	// Calculate topological order
	sorted := g.TopologicalSort()
	log.Infof("Topological Order:\n")
	ret := make([]types.Component, len(sorted))
	// Reverse order
	for i, node := range sorted {
		log.Infof(gr.showNodeWithNeighbours(g, &node))
		ret[len(sorted)-1-i] = (*node.Value).(types.Component)
	}

	return ret, nil
}

func (gr *GraphAPI) showNodeWithNeighbours(g *graph.Graph, node *graph.Node) string {
	text := ""
	for _, n := range g.Neighbors(*node) {
		text += fmt.Sprint((*n.Value).(types.Component).Name) + ", "
	}
	return fmt.Sprintf("%v [%v]", (*node.Value).(types.Component).Name, text)
}
