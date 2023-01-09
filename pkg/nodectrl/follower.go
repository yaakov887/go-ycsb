package nodectrl

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"strings"
)

type Follower struct {
	Id        string `json:"followerID"`
	IpAddrStr string `json:"IP"`
	Username  string `json:"username"`
	KeyFile   string `json:"keyfile"`
	Started   bool
	sshClient *ssh.ClientConfig
}

type FollowerList struct {
	Followers []Follower
}

const (
	StartFollowerFmtStr = "go-ycsb run %v -F %v -P %v"
)

func FollowersStarted() bool {
	anyFollowers := false
	for _, f := range globalFollowerList.Followers {
		anyFollowers = anyFollowers || f.Started
	}
	return anyFollowers
}

func ParseFollowerList(jsonSource string) error {
	bytes, _ := os.ReadFile(jsonSource)

	var templist FollowerList
	err := json.Unmarshal(bytes, &templist)
	globalFollowerList = templist

	return err
}

func (f *Follower) runFollowerCommand(command string) error {
	if f.sshClient == nil {
		sshClient, err := GenerateSSHClientConfig(f.Username, f.KeyFile)
		if err != nil {
			return err
		}
		f.sshClient = sshClient
	}

	client, err := ssh.Dial("tcp", f.IpAddrStr, f.sshClient)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Run(command)
	if err != nil {
		return err
	}
	return nil
}

func StartFollowers(dbName, workload string) {
	for _, follower := range globalFollowerList.Followers {
		startcmd := fmt.Sprintf(StartFollowerFmtStr, dbName, follower.Id, workload)
		err := follower.runFollowerCommand(startcmd)
		follower.Started = err == nil
	}
}

// Download file from sftp server
func downloadFile(sc sftp.Client, remoteFile, localFile string) (err error) {

	srcFile, err := sc.OpenFile(remoteFile, (os.O_RDONLY))
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localFile)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	if bytes == 0 {
		log.Printf("Copied %v bytes from follower file %v...", bytes, remoteFile)
	}

	return nil
}

func (f *Follower) getFollowerFiles() error {
	conn, err := ssh.Dial("tcp", f.IpAddrStr, f.sshClient)
	if err != nil {
		return err
	}
	defer conn.Close()

	sc, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer sc.Close()

	files, err := sc.ReadDir(".")
	if err != nil {
		return err
	}
	for _, fi := range files {
		if strings.Contains(fi.Name(), f.Id) {
			downloadFile(*sc, fi.Name(), "./"+fi.Name())
		}
	}
	return nil
}

func GetFollowersFiles() {
	for _, f := range globalFollowerList.Followers {
		if f.Started {
			err := f.getFollowerFiles()
			if err != nil {
				log.Printf("Error downloading files from follower %v.", f.Id)
			}
		} else {
			log.Printf("Follower %v was never started and has no files.", f.Id)
		}
	}
}

var globalFollowerList FollowerList
