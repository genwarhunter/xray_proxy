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
		}
		if time.Now().Sub(lastTimeUpdateConfigs) > 5*time.Minute {
			GetConfigs()
			lastTimeUpdateConfigs = time.Now()
		}
		time.Sleep(5 * time.Second)
	}
}

func runner() {
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
			response := httpGET(value.(infoPackageRow).Url, 1)
			for _, link := range strings.Split(response, "\n") {
				hash := GenerateConfig(link)
				if hash == "" {
					continue
				}
			}
		}
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
		var tmp = string(configData)
		tmp = strings.Replace(tmp, "\"port\":0,\"protocol\"", "\"port\":"+strconv.Itoa(int(port))+",\"protocol\"", 1)
		io.WriteString(stdin, tmp)
		stdin.Close()
		_, err = cmd.CombinedOutput()
		err = cmd.Start()
		if err != nil {
			fmt.Printf("Ошибка при запуске команды: %v\n", err)
			return
		}
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
