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

import "github.com/ElPuig/netgraph/pkg/xml_loader"

type Node struct {
	Id       string
	Ip       string
	Model    string
	Location string
}

type Noder interface {
	GetLabel() string
	GetGroup() string
	GetNodeType() string
	GetShape() string
	GetSize() int
}

func (n Node) GetLabel() string {
	return n.Id + "\n" + n.Ip + "\nModel: " + n.Model + "\nLocation: " + n.Location
}

func (n Node) GetGroup() string {
	return "TODO"
}

func (n Node) GetNodeType() string {
	return "TODO"
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

func (n Node) GetSize() int {
	return 10
}

func GetNodeList(xml_data []xml_loader.RequestXMLData) map[string]Noder {
	res := make(map[string]Noder)
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
