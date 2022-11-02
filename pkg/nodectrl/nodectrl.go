package nodectrl

import (
	"encoding/json"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
)

type Node struct {
	Id        string `json:"nodeID"`
	IpAddrStr string `json:"IP"`
	Username  string `json:"username"`
	keyFile   string `json:"keyfile"`
	cert      ssh.Certificate
	IpAddr    net.IPAddr
}

type NodeList struct {
	Nodes []Node `json:"nodes"`
}

func ParseNodeList(jsonSource string) error {
	bytes, _ := os.ReadFile(jsonSource)

	var templist NodeList
	err := json.Unmarshal(bytes, &templist)

	return err
}

func StartNodes() error {
	return nil
}

func StartNodeById(nodeId string) error {
	return nil
}

func StopNodes() error {
	return nil
}

func StopNodeById(nodeId string) error {
	return nil
}

var globalNodeList NodeList
