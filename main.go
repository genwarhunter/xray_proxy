package main

import (
	"log"
	"path/filepath"
	"sync"
	"sync/atomic"
)

var PortPidMap sync.Map
var PortConfMap sync.Map
var HashProtocolMap sync.Map

var Package sync.Map

var hashMap sync.Map

var freePorts *MinHeap
var threadsNow atomic.Int32

var queue chan string

//var servConnection map[server]string

func init() {
	queue = make(chan string, 100000)
	AppConfig.PathToConfDir, _ = filepath.Abs(AppConfig.PathToConfDir)
	GetConfig()
	freePorts = NewMinHeap()
	CreateConnectMysql()
	updatePackages()
	GetConfigs()
	loadHashes()
	updateQueue()
	for port := AppConfig.StartPort; port < AppConfig.StartPort+AppConfig.RangePort; port++ {
		freePorts.Insert(port)
	}
	log.Println("Initialization Complete")
}

func main() {
	log.Println("Run threads")
	threads()
}
