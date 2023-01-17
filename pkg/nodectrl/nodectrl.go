package nodectrl

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pingcap/go-ycsb/pkg/util"
	"golang.org/x/crypto/ssh"
	"os"
	"regexp"
	"time"
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

var pidHeader = []string{"nodeid", "pid"}

const (
	pidFileName = "pid_output.csv"
)

// updateNodePid updates the node, based on the Node ID, with the PID passed
func updateNodePid(nodeid, pid string) {
	for _, node := range globalNodeList.Nodes {
		if node.Id == nodeid {
			node.pid = pid
		}
	}
}

// writeNodePids Write the Node Process IDs to a CSV for eventual reading and stopping of the nodes
func writeNodePids() error {
	err := os.Remove(pidFileName)
	file, err := os.Create(pidFileName)
	if err != nil {
		return err
	}

	var values [][]string
	for _, node := range globalNodeList.Nodes {
		var tempValues []string
		tempValues = append(tempValues, node.Id)
		tempValues = append(tempValues, node.pid)
		values = append(values)
	}

	util.RenderCSV(pidHeader, values, file)

	return file.Close()
}

// readNodePids read the Node Process IDs from the CSV file
func readNodePids() error {
	file, err := os.Open(pidFileName)
	if err != nil {
		return err
	}

	r, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return err
	}

	for i, record := range r {
		if i == 0 {
			continue
		}

		updateNodePid(record[0], record[1])
	}

	return nil
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

// NodesParsed returns true if the globalNodeList is not empty
func NodesParsed() bool {
	return len(globalNodeList.Nodes) > 0
}

// NodesStarted returns true if any of the nodes have a PID set
func NodesStarted() bool {
	for _, node := range globalNodeList.Nodes {
		if node.pid != "" {
			return true
		}
	}
	return false
}

// getNodeById returns the node based on the ID passed
func getNodeById(nodeId string) (*Node, error) {
	for _, node := range globalNodeList.Nodes {
		if node.Id == nodeId {
			return &node, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Node id [%v] not found", nodeId))
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
		writeNodePids()
		return nil
	}
}

// StartNodeById starts the node specified by the node id
func StartNodeById(nodeId string) error {
	node, err := getNodeById(nodeId)
	if err != nil {
		return err
	}
	return node.startNode()
}

// stopNode connects to the node and calls the kill command with the process id stored
func (n *Node) stopNode() error {
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
	cmdStr := fmt.Sprintf("kill %v", n.pid)
	session.Run(cmdStr)
	if len(b.Bytes()) > 0 {
		return errors.New(fmt.Sprintf("Error stopping node %v: %v", n.Id, b.String()))
	}

	return nil
}

// StopNodes stops all nodes specified by the cluster file
func StopNodes() error {
	if !NodesStarted() {
		readNodePids()
	}

	var errMap map[string]error
	for _, node := range globalNodeList.Nodes {
		err := node.stopNode()
		if err != nil {
			if errMap == nil {
				errMap = make(map[string]error)
			}
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
	node, err := getNodeById(nodeId)
	if err != nil {
		return err
	}
	return node.stopNode()
}

// runNodeCmd executes the command passed on the referenced node
func (n *Node) runNodeCmd(command string) error {
	if &n.sshClient == nil {
		sshClient, err := GenerateSSHClientConfig(n.Username, n.KeyFile)
		if err != nil {
			return err
		}
		n.sshClient = *sshClient
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

	return session.Run(command)

}

// RunNodeCommand executes the command passed on the node specified by id—
func RunNodeCommand(nodeId, command string) error {
	node, err := getNodeById(nodeId)
	if err != nil {
		return err
	}
	return node.runNodeCmd(command)
}

// RunSSHCommand run a command over ssh connection
func RunSSHCommand(ipAddr, userName, keyFile, command string) error {
	if ipAddr == "" {
		return errors.New("[RunSSHCommand] IP Address required")
	}
	if userName == "" {
		return errors.New("[RunSSHCommand] Username required")
	}
	if keyFile == "" {
		return errors.New("[RunSSHCommand] Key file required")
	}
	if command == "" {
		return errors.New("[RunSSHCommand] Command required")
	}

	sshClient, err := GenerateSSHClientConfig(userName, keyFile)
	if err != nil {
		return err
	}

	client, err := ssh.Dial("tcp", ipAddr, sshClient)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(command)
}

// generateSSHClientConfig create the ssh client config from the username and keyfile provided
func GenerateSSHClientConfig(userName, keyFile string) (*ssh.ClientConfig, error) {
	var hostKey ssh.PublicKey
	key, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return &ssh.ClientConfig{
		Config: ssh.Config{},
		User:   userName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
		Timeout:         2 * time.Second,
	}, nil
}

var globalNodeList NodeList
