package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func threads() {
	go runner()
	var lastTimeUpdateQueue = time.Now()
	var lastTimeUpdateConfigs = time.Now()
	var lastTimeUpdatePackages = time.Now()
	var lastTimeKiller = time.Now()
	for {
		if time.Now().Sub(lastTimeUpdateQueue) > 30*time.Minute {
			loadHashes()
			updateQueue()
			lastTimeUpdateQueue = time.Now()
		}
		if time.Now().Sub(lastTimeUpdatePackages) > 30*time.Minute {
			updatePackages()
			lastTimeUpdatePackages = time.Now()
		}
		if time.Now().Sub(lastTimeUpdateConfigs) > 30*time.Minute {
			GetConfigs()
			lastTimeUpdateConfigs = time.Now()
		}
		if time.Now().Sub(lastTimeKiller) > 5*time.Minute {
			killer()
			lastTimeKiller = time.Now()
		}
		time.Sleep(5 * time.Second)
	}
}

func runner() {
	defer func() {
		PortPidMap.Range(func(key, value any) bool {
			process, _ := os.FindProcess(value.(int))
			_ = process.Signal(os.Interrupt)
			return true
		})
	}()
	log.Println("Run Runner")
	for {
		var t = uint16(threadsNow.Load())
		if t < AppConfig.RangePort {
			for i := uint16(0); i < AppConfig.RangePort-t; i++ {
				threadsNow.Add(1)
				runXray()
			}
		}

		PortPidMap.Range(func(key, value any) bool {
			process, err := os.FindProcess(value.(int))
			if err != nil {
				return false
			}
			err = process.Signal(os.Signal(nil))
			if err != nil {
				log.Println(value.(int), " Умер")
				deletePortInfo(key.(uint16))
			}
			return true
		})
		time.Sleep(5 * time.Second)
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
}

func killer() {
	var response = httpGET(AppConfig.CheckerUrl, 1)
	var re = regexp.MustCompile(`(?si)port":(\d+)`)
	for _, match := range re.FindAllString(str, -1) {

	}

	Pid2Kill.Range(func(key, value any) bool {
		process, _ := os.FindProcess(value.(int))
		_ = process.Signal(os.Interrupt)
		return true
	})
}
