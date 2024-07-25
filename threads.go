package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func threads() {
	go runner()
	var lastTimeUpdateQueue = time.Now()
	var lastTimeUpdateConfigs = time.Now()
	var lastTimeUpdatePackages = time.Now()
	for {
		if time.Now().Sub(lastTimeUpdateQueue) > 5*time.Minute {
			loadHashes()
			updateQueue()
			lastTimeUpdateQueue = time.Now()
		}
		if time.Now().Sub(lastTimeUpdatePackages) > 5*time.Minute {
			updatePackages()
			lastTimeUpdatePackages = time.Now()
		}
		if time.Now().Sub(lastTimeUpdateConfigs) > 5*time.Minute {
			GetConfigs()
			lastTimeUpdateConfigs = time.Now()
		}
		time.Sleep(5 * time.Second)
	}
}

func runner() {
	log.Println("Run Runner")
	for {
		var t = uint16(threadsNow.Load())
		if t < AppConfig.RangePort {
			for i := uint16(0); i < AppConfig.RangePort-t; i++ {
				threadsNow.Add(1)
				go runXray()
			}
		}

		if t > AppConfig.RangePort {
			for i := uint16(0); i < t-AppConfig.RangePort; i++ {
				go func() {
					quit <- struct{}{}
				}()
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func GetConfigs() {
	Package.Range(func(key, value any) bool {
		if value.(infoPackageRow).Use {
			if value.(infoPackageRow).Id == 2 {
				GenerateConfig("vmess://eyJhZGQiOiJsaW5kZTA2LmluZGlhdmlkZW8uc2JzIiwiYWlkIjoiMCIsImhvc3QiOiJsaW5kZTA2LmluZGlhdmlkZW8uc2JzIiwiaWQiOiJlZGJiMTA1OS0xNjMzLTQyNzEtYjY2ZS1lZDRmYmE0N2ExYmYiLCJuZXQiOiJ3cyIsInBhdGgiOiIvbGlua3dzIiwicG9ydCI6IjQ0MyIsInBzIjoi8J+HuvCfh7hVUyB8IPCfn6IgfCB2bWVzcyB8IEBEZWFtTmV0X1Byb3h5IHwgMCIsInNjeSI6ImF1dG8iLCJzbmkiOiIiLCJ0bHMiOiJ0bHMiLCJ0eXBlIjoiIiwidiI6IjIifQ==")
			}
			response := httpGET(value.(infoPackageRow).Url, 1)
			for _, link := range strings.Split(response, "\n") {
				hash := GenerateConfig(link)
				if hash == "" {
					continue
				}
			}
		}
		log.Println("Config OK!")
		return true
	})
}

func runXray() {
	select {
	case <-quit:
		return
	default:
		var hash = <-queue
		path := AppConfig.PathToConfDir + hash
		var port, err = freePorts.ExtractMin()
		if err != nil {
			return
		}
		configData, err := os.ReadFile(path)
		PortConfMap.Store(port, hash)
		var cmd = exec.Command("xray")
		stdin, err := cmd.StdinPipe()

		err = cmd.Start()
		if err != nil {
			fmt.Printf("Ошибка при запуске команды: %v\n", err)
			return
		}
		var tmp = string(configData)
		tmp = strings.Replace(tmp, "\"port\":0,\"protocol\"", "\"port\":"+strconv.Itoa(int(port))+",\"protocol\"", 1)
		io.WriteString(stdin, tmp)
		stdin.Close()
		_, err = cmd.CombinedOutput()
		pid := cmd.Process.Pid
		log.Println(pid)
		PortPidMap.Store(port, pid)
		err = cmd.Wait()
		defer func() {
			hashMap[hash] = false
			PortPidMap.Delete(port)
			PortConfMap.Delete(port)
			freePorts.Insert(port)
			threadsNow.Add(-1)
		}()
	}
}
