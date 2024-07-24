package main

import (
	"path/filepath"
	"sync"
	"sync/atomic"
)

var PortPidMap sync.Map
var PortConfMap sync.Map

var Package sync.Map

var hashMap map[string]bool

var freePorts *MinHeap
var threadsNow atomic.Int32

var quit chan struct{}
var queue chan string

//var servConnection map[server]string

func init() {
	// s := GenerateConfig("vmess://eyJ2IjoiMiIsImFkZCI6ImxpbmRlMDYuaW5kaWF2aWRlby5zYnMiLCJwb3J0IjoiNDQzIiwiaWQiOiJlZGJiMTA1OS0xNjMzLTQyNzEtYjY2ZS1lZDRmYmE0N2ExYmYiLCJhaWQiOjAsInNjeSI6ImF1dG8iLCJuZXQiOiJ3cyIsImhvc3QiOiJsaW5kZTA2LmluZGlhdmlkZW8uc2JzIiwicGF0aCI6IlwvbGlua3dzIiwidGxzIjoidGxzIiwicHMiOiJcdWQ4M2NcdWRkZmFcdWQ4M2NcdWRkZjhVUyB8IFx1ZDgzZFx1ZGZlMiB8IHZtZXNzIHwgQERlYW1OZXRfUHJveHkgfCAwIn0=")
	// fmt.Println(s)

	//os.Exit(33)
	hashMap = make(map[string]bool)
	queue = make(chan string, 10000)
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
