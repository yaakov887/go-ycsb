package nodectrl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"regexp"
)

type Node struct {
	Id        string `json:"nodeID"`
	IpAddrStr string `json:"IP"`
	Username  string `json:"username"`
	keyFile   string `json:"keyfile"`
	IpAddr    net.IPAddr
	sshClient ssh.ClientConfig
	pid       string
}

type NodeList struct {
	Nodes []Node `json:"nodes"`
}

func ParseNodeList(jsonSource string) error {
	bytes, _ := os.ReadFile(jsonSource)

	var templist NodeList
	err := json.Unmarshal(bytes, &templist)
	globalNodeList = templist

	for i, node := range globalNodeList.Nodes {
		globalNodeList.Nodes[i].IpAddr.IP = net.ParseIP(node.IpAddrStr)
		var hostKey ssh.PublicKey
		key, err := os.ReadFile(node.keyFile)
		if err != nil {
			return err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return err
		}
		node.sshClient = ssh.ClientConfig{
			Config: ssh.Config{},
			User:   node.Username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.FixedHostKey(hostKey),
		}
	}

	return err
}

func (n *Node) startNode() error {
	client, err := ssh.Dial("tcp", n.IpAddrStr, &n.sshClient)
	defer client.Close()
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	defer session.Close()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	session.Stdout = &b
	cmdStr := fmt.Sprintf("%v \\'%v %v\\'", "sh -c", "echo $$; exec", "")
	session.Run(cmdStr)
	ok, err := regexp.Match("[0-9]+", b.Bytes())
	if ok {
		n.pid = b.String()
	} else {
		return errors.New(fmt.Sprintf("Error parsing pid for node id %v", n.Id))
	}
	return nil
}

func StartNodes() error {
	var errMap map[string]error
	for _, node := range globalNodeList.Nodes {
		err := node.startNode()
		if err != nil {
			errMap[node.Id] = err
		}
	}
	if len(errMap) > 0 {
		return errors.New(fmt.Sprintf("Error starting nodes: %v", errMap))
	} else {
		return nil
	}
}

func StartNodeById(nodeId string) error {
	for _, node := range globalNodeList.Nodes {
		if node.Id == nodeId {
			return node.startNode()
		}
	}
	return errors.New(fmt.Sprintf("Node id [%v] to start not found", nodeId))
}

func (n *Node) stopNode() error {
	client, err := ssh.Dial("tcp", n.IpAddrStr, &n.sshClient)
	defer client.Close()
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	defer session.Close()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	session.Stdout = &b
	cmdStr := fmt.Sprintf("kill %v", n.pid)
	session.Run(cmdStr)
	if len(b.Bytes()) > 0 {
		return errors.New(fmt.Sprintf("Error stopping node %v: %v", n.Id, b.String()))
	}

	return nil
}

func StopNodes() error {
	var errMap map[string]error
	for _, node := range globalNodeList.Nodes {
		err := node.stopNode()
		if err != nil {
			errMap[node.Id] = err
		}
	}
	if len(errMap) > 0 {
		return errors.New(fmt.Sprintf("Error stopping nodes: %v", errMap))
	} else {
		return nil
	}
}

func StopNodeById(nodeId string) error {
	for _, node := range globalNodeList.Nodes {
		if node.Id == nodeId {
			return node.stopNode()
		}
	}
	return errors.New(fmt.Sprintf("Node id [%v] to stop not found", nodeId))
}

var globalNodeList NodeList
