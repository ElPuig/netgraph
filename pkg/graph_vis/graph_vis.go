/*
Copyright 2025 David García De Mercado.

This file is part of Netgraph.

   Netgraph is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   Netgraph is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with Netgraph.  If not, see <https://www.gnu.org/licenses/>.
*/

package graph_vis

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ElPuig/netgraph/pkg/xml_loader"
)

var regex_switch string = `^R.*-.*`
var regex_ap string = `^\d{3}`

type Node struct {
	Id       string
	Ip       string
	Model    string
	Location string
}

type Noder interface {
	GetName() string
	GetLabel() string
	GetGroup() string
	GetNodeType() string
	GetShape() string
	GetSize() string
	GetUrl() string
}

type NodeMap map[string]Noder

func (n Node) GetName() string {
	return n.Id
}

func (n Node) GetLabel() string {
	return n.Id + "\\n" + n.Ip + "\\nModel: " + n.Model + "\\nLocation: " + n.Location
}

func (n Node) GetGroup() string {
	match, _ := regexp.MatchString(regex_switch, n.Id)
	if match {
		return strings.Split(n.Id, "-")[0][1:]
	}
	match, _ = regexp.MatchString(regex_ap, n.Id)
	if match {
		return "AP"
	}
	return "UNKNOWN"
}

func (n Node) GetNodeType() string {
	match, _ := regexp.MatchString(regex_switch, n.Id)
	if match {
		return "SWITCH"
	}
	match, _ = regexp.MatchString(regex_ap, n.Id)
	if match {
		return "AP"
	}
	return "UNDEFINED"
}

func (n Node) GetShape() string {
	switch n.GetNodeType() {
	case "SWITCH":
		return "box"
	case "AP":
		return "dot"
	default:
		return "triangle"
	}
}

func (n Node) GetSize() string {
	return "10"
}

func (n Node) GetUrl() string {
	return "/" + n.Ip
}

func GetNodeMap(xml_data []xml_loader.RequestXMLData) NodeMap {
	res := make(NodeMap)
	for _, xml := range xml_data {
		res[xml.IP] = Node{
			Id:       xml.Device.Info.Name,
			Ip:       xml.IP,
			Model:    xml.Device.Info.Model,
			Location: xml.Device.Info.Location,
		}
	}
	return res
}

func (nm NodeMap) ToVisJson() json.RawMessage {
	res := ``
	res += `[`
	for _, n := range nm {
		res += `{ "id": "` + n.GetName() + `", `
		res += `  "label": "` + n.GetLabel() + `", `
		res += `  "shape": "` + n.GetShape() + `", `
		res += `  "group": "` + n.GetGroup() + `", `
		res += `  "size": ` + string(n.GetSize()) + `, `
		res += `  "url": "` + n.GetUrl() + `"},`
	}
	res = res[:len(res)-1]
	res += `]`
	return json.RawMessage(res)
}

type Edge struct {
	From string
	To   string
}

type Edger interface {
	GetFrom() string
	GetTo() string
	GetLabel() string
	GetLength() string
	GetArrowType() string
}

type EdgeMap map[string]Edger

func (e Edge) GetFrom() string {
	return e.From
}

func (e Edge) GetTo() string {
	return e.To
}

func (e Edge) GetLabel() string {
	return "TODO"
}

func (e Edge) GetLength() string {
	return "TODO"
}

func (e Edge) GetArrowType() string {
	return "TODO"
}

func GetEdgeMap(xml_data []xml_loader.RequestXMLData) EdgeMap {
	res := make(EdgeMap)
	for _, xml := range xml_data {
		for _, p := range xml.Device.Ports {
			for _, c := range p.Connections.Connection {
				if c.SysName != "" {
					var name0, name1 string
					if xml.Device.Info.Name < c.SysName {
						name0 = xml.Device.Info.Name
						name1 = c.SysName
					} else {
						name0 = c.SysName
						name1 = xml.Device.Info.Name
					}
					edge_key := name0 + name1
					if value, exists := res[edge_key]; exists {
						fmt.Println(value)
					} else {
						res[edge_key] = Edge{
							From: name0,
							To:   name1,
						}
					}
				}
			}
		}
	}
	return res
}

func (em EdgeMap) ToVisJson() json.RawMessage {
	res := ``
	res += `[`
	for _, e := range em {
		res += `{ "from": "` + e.GetFrom() + `", `
		res += `  "to": "` + e.GetTo() + `"},`
	}
	res = res[:len(res)-1]
	res += `]`
	return json.RawMessage(res)
}
