/**
 * Copyright (c) 2015 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package graph

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/trustedanalytics/go-cf-lib/api"
	"github.com/trustedanalytics/go-cf-lib/types"
	"github.com/twmb/algoimpl/go/graph"
	"net/url"
	"strings"
)

type DependencyGraph struct {
	cf    *api.CfAPI
	nodes map[string]graph.Node
}

func NewDependencyGraph() *DependencyGraph {
	toReturn := new(DependencyGraph)
	toReturn.cf = api.NewCfAPI()
	toReturn.nodes = make(map[string]graph.Node)
	return toReturn
}

func (dg *DependencyGraph) NewNode(g *graph.Graph, guid, name string, typ types.ComponentType,
	dependent *graph.Node, clone bool) graph.Node {

	if node, ok := dg.nodes[guid]; ok {
		dg.appendDependency(&node, dependent)
		return node
	}
	node := g.MakeNode()
	*node.Value = types.Component{
		GUID:         guid,
		Name:         name,
		Type:         typ,
		DependencyOf: []string{},
		Clone:        clone,
	}
	dg.appendDependency(&node, dependent)
	dg.nodes[guid] = node
	return node
}

func (dg *DependencyGraph) appendDependency(node *graph.Node, dependent *graph.Node) {
	if dependent != nil {
		dependencyOf := (*node.Value).(types.Component).DependencyOf
		dependencyOf = append(dependencyOf, (*dependent.Value).(types.Component).GUID)
		value := (*node.Value).(types.Component)
		value.DependencyOf = dependencyOf
		*node.Value = value
	}
}

func (dg *DependencyGraph) addDependenciesToGraph(g *graph.Graph, parent graph.Node, sourceAppGUID string) error {
	log.Infof("addDependenciesToGraph for parent %v", *parent.Value)
	sourceAppSummary, err := dg.cf.GetAppSummary(sourceAppGUID)
	if err != nil {
		return err
	}
	for _, svc := range sourceAppSummary.Services {
		if dg.isNormalService(svc) {
			node := dg.NewNode(g, svc.GUID, svc.Name, types.ComponentService, &parent, true)
			g.MakeEdgeWeight(parent, node, 1)
		} else {
			node := dg.NewNode(g, svc.GUID, svc.Name, types.ComponentUPS, &parent, true)
			g.MakeEdgeWeight(parent, node, 1)
			// Retrieve UPS
			response, err := dg.cf.GetUserProvidedService(svc.GUID)
			if err != nil {
				return err
			}
			val, ok := response.Entity.Credentials["url"]
			if !ok {
				continue
			}
			urlStr, ok := val.(string)
			if !ok {
				continue
			}
			appID, appName, err := dg.getAppIdAndNameFromSpaceByUrl(sourceAppSummary.SpaceGUID, urlStr)
			if err != nil {
				return err
			}
			if len(appID) > 0 {
				log.Infof("Application %v is bound using %v", appID, svc.Name)
				node2 := dg.NewNode(g, appID, appName, types.ComponentApp, &node, true)
				g.MakeEdgeWeight(node, node2, 1)
				if dg.graphHasCycles(g) {
					log.Errorf("Graph got cycle. Stopping graph traversing...")
					break
				}
				_ = dg.addDependenciesToGraph(g, node2, appID)
			}
		}
	}
	return nil
}

func (dg *DependencyGraph) graphHasCycles(g *graph.Graph) bool {
	components := g.StronglyConnectedComponents()
	for _, comp := range components {
		if len(comp) > 1 {
			log.Warnf("Cycle of length %v", len(comp))
			return true
		}
	}
	return false
}

func (dg *DependencyGraph) isNormalService(svc types.CfAppSummaryService) bool {
	// Normal services require plan.
	// User provided services does not support Plans so this field is empty then.
	return len(svc.Plan.Service.Label) > 0
}

func (dg *DependencyGraph) getAppIdAndNameFromSpaceByUrl(spaceGUID, urlStr string) (string, string, error) {
	appURL, err := url.Parse(urlStr)
	if err != nil {
		log.Infof("[%v] is not a correct URL. Parsing failed.", urlStr)
		return "", "", err
	}
	log.Infof("URL Host %v", appURL.Host)
	routes, err := dg.cf.GetSpaceRoutesForHostname(spaceGUID, strings.Split(appURL.Host, ".")[0])
	if err != nil {
		return "", "", err
	}
	if routes.Count == 0 {
		log.Infof("No routes found for host: %v", appURL.Host)
		return "", "", nil
	}
	log.Infof("%v route(s) retrieved for host %v", routes.Count, appURL.Host)
	routeGUID := routes.Resources[0].Meta.GUID
	apps, err := dg.cf.GetAppsFromRoute(routeGUID)
	if err != nil {
		return "", "", err
	}
	if apps.Count == 0 {
		log.Infof("No apps bound to route: [%v]", routeGUID)
		return "", "", nil
	}
	app := apps.Resources[0]
	log.Debugf("App %+v", app)
	isSearched, err := dg.doesUrlMatchApplication(urlStr, app.Meta.GUID)
	if err != nil {
		return "", "", err
	}
	if !isSearched {
		log.Infof("url of found app does not match url in user provided service")
		return "", "", nil

	}
	log.Infof("Found app match url in user provided service")
	return app.Meta.GUID, app.Entity.Name, nil
}

func (dg *DependencyGraph) doesUrlMatchApplication(appUrlStr, appID string) (bool, error) {
	appURL, err := url.Parse(appUrlStr)
	if err != nil {
		return false, err
	}
	appSummary, err := dg.cf.GetAppSummary(appID)
	log.Debugf("App summary retrieved is [%+v]", appSummary)
	if err != nil {
		return false, err
	}
	for i := range appSummary.Routes {
		route := appSummary.Routes[i]
		if appURL.Host == fmt.Sprintf("%v.%v", route.Host, route.Domain.Name) {
			return true, nil
		}
	}
	return false, nil
}
