// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package controller

import (
	"context"
	"strconv"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"

	"github.com/api7/ingress-controller/pkg/kube"
	sevenConf "github.com/api7/ingress-controller/pkg/seven/conf"
	"github.com/api7/ingress-controller/pkg/seven/state"
	apisixv1 "github.com/api7/ingress-controller/pkg/types/apisix/v1"
)

const (
	ADD           = "ADD"
	UPDATE        = "UPDATE"
	DELETE        = "DELETE"
	WatchFromKind = "watch"
)

func Watch() {
	c := &controller{
		queue: make(chan interface{}, 100),
	}
	kube.EndpointsInformer.Informer().AddEventHandler(&QueueEventHandler{c: c})
	go c.run()
}

func (c *controller) pop() interface{} {
	e := <-c.queue
	return e
}

func (c *controller) run() {
	for {
		ele := c.pop()
		c.process(ele)
	}
}

func (c *controller) process(obj interface{}) {
	qo, _ := obj.(*queueObj)
	ep, _ := qo.Obj.(*v1.Endpoints)
	if ep.Namespace != "kube-system" { // todo here is some ignore namespaces
		for _, s := range ep.Subsets {
			// if upstream need to watch
			// ips
			ips := make([]string, 0)
			for _, address := range s.Addresses {
				ips = append(ips, address.IP)
			}
			// ports
			for _, port := range s.Ports {
				upstreamName := ep.Namespace + "_" + ep.Name + "_" + strconv.Itoa(int(port.Port))
				// find upstreamName is in apisix
				// sync with all apisix group
				for _, cluster := range sevenConf.Client.ListClusters() {
					upstreams, err := cluster.Upstream().List(context.TODO())
					if err == nil {
						for _, upstream := range upstreams {
							if upstream.Name == upstreamName {
								nodes := make([]apisixv1.Node, 0)
								for _, ip := range ips {
									node := apisixv1.Node{IP: ip, Port: int(port.Port), Weight: 100}
									nodes = append(nodes, node)
								}
								upstream.Nodes = nodes
								// update upstream nodes
								// add to seven solver queue
								//apisix.UpdateUpstream(upstream)
								upstream.FromKind = WatchFromKind
								upstreams := []*apisixv1.Upstream{upstream}
								comb := state.ApisixCombination{Routes: nil, Services: nil, Upstreams: upstreams}
								if _, err = comb.Solver(); err != nil {
									glog.Errorf(err.Error())
								}
							}
						}
					}
				}
			}
		}
	}
}

type controller struct {
	queue chan interface{}
}

type queueObj struct {
	OpeType string      `json:"ope_type"`
	Obj     interface{} `json:"obj"`
}

type QueueEventHandler struct {
	c *controller
}

func (h *QueueEventHandler) OnAdd(obj interface{}) {
	h.c.queue <- &queueObj{ADD, obj}
}

func (h *QueueEventHandler) OnDelete(obj interface{}) {
	h.c.queue <- &queueObj{DELETE, obj}
}

func (h *QueueEventHandler) OnUpdate(old, update interface{}) {
	h.c.queue <- &queueObj{UPDATE, update}
}
