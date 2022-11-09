package nodectrl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"regexp"
)

type Node struct {
	Id          string `json:"nodeID"`
	IpAddrStr   string `json:"IP"`
	Username    string `json:"username"`
	KeyFile     string `json:"keyfile"`
	NodeCommand string `json:"nodecommand"`
	sshClient   ssh.ClientConfig
	pid         string
}

type NodeList struct {
	Nodes        []Node `json:"nodes"`
	StartCommand string `json:"startcommand"`
}

// ParseNodeList read the JSON formatted file for the cluster information
func ParseNodeList(jsonSource string) error {
	bytes, _ := os.ReadFile(jsonSource)

	var templist NodeList
	err := json.Unmarshal(bytes, &templist)
	globalNodeList = templist

	for i, node := range globalNodeList.Nodes {
		var hostKey ssh.PublicKey
		key, err := os.ReadFile(node.KeyFile)
		if err != nil {
			return err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return err
		}
		globalNodeList.Nodes[i].sshClient = ssh.ClientConfig{
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

// startNode executes the start command for the referenced node
func (n *Node) startNode() error {
	var startCmd string
	if n.NodeCommand != "" {
		startCmd = n.NodeCommand
	} else {
		startCmd = globalNodeList.StartCommand
	}

	client, err := ssh.Dial("tcp", n.IpAddrStr, &n.sshClient)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	cmdStr := fmt.Sprintf("%v \\'%v %v\\'", "sh -c", "echo $$; exec", startCmd)
	session.Run(cmdStr)
	ok, err := regexp.Match("[0-9]+", b.Bytes())
	if ok {
		n.pid = b.String()
	} else {
		return errors.New(fmt.Sprintf("Error parsing pid for node id %v", n.Id))
	}
	return nil
}

// StartNodes starts all the nodes specified by the cluster file
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

// StartNodeById starts the node specified by the node id
func StartNodeById(nodeId string) error {
	for _, node := range globalNodeList.Nodes {
		if node.Id == nodeId {
			return node.startNode()
		}
	}
	return errors.New(fmt.Sprintf("Node id [%v] to start not found", nodeId))
}

// stopNode connects to the node and calls the kill command with the process id stored
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

// StopNodes stops all nodes specified by the cluster file
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

// StopNodeById stops the node specified by the node id
func StopNodeById(nodeId string) error {
	for _, node := range globalNodeList.Nodes {
		if node.Id == nodeId {
			return node.stopNode()
		}
	}
	return errors.New(fmt.Sprintf("Node id [%v] to stop not found", nodeId))
}

var globalNodeList NodeList
