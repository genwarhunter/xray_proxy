package main

import (
	"log"
	"path/filepath"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

var PortPidMap sync.Map
var PortConfMap sync.Map
var Pid2Kill sync.Map

var Package sync.Map

var hashMap sync.Map

var freePorts *MinHeap
var threadsNow atomic.Int32

var queue chan string

//var servConnection map[server]string

func init() {
	queue = make(chan string, 100000)
	debug.SetMaxThreads(10 * 1e6)
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
}

func main() {
	log.Println("Run threads")
	threads()
	//var response = httpGET("https://raw.githubusercontent.com/Epodonios/v2ray-configs/main/All_Configs_Sub.txt", 1)
	//for _, link := range strings.Split("response, "\n") {
	//	hash := GenerateConfig(link)
	//	if hash == "" {
	//		continue
	//	}
	//	go runXray()
	//	// fmt.Println(hash)
	//}
}
