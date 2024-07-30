package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
		if time.Now().Sub(lastTimeKiller) > AppConfig.DelayKill*time.Minute {
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
			err = process.Signal(syscall.Signal(0))
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
	var wg sync.WaitGroup
	Package.Range(func(key, value any) bool {
		if value.(infoPackageRow).Use {
			if value.(infoPackageRow).Id == 2 {
				wg.Add(1)
				go GenerateConfig("ss://Y2hhY2hhMjAtaWV0Zi1wb2x5MTMwNTpGNGZDaFdHb2pJY214dU1ST3VWUUlRS0dLUmRnMDl2S2ZHcDJTVHpTcmR6clVTd1o=@45.87.219.110:51348#🔒 SS-TCP-NA 🇧🇾 BY-45.87.219.110:51348", &wg)
			}
			response := httpGET(value.(infoPackageRow).Url, 1)
			for _, link := range strings.Split(response, "\n") {
				wg.Add(1)
				go GenerateConfig(link, &wg)

			}
		}
		log.Println("Config OK!")
		return true
	})
	wg.Wait()
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
	_, _ = io.WriteString(stdin, tmp)
	_ = stdin.Close()
	_, err = cmd.CombinedOutput()
	pid := cmd.Process.Pid
	log.Println("PID: ", pid, "	PORT: ", port, "	STATUS: Started")
	PortPidMap.Store(port, pid)
}

func killer() {
	var response = httpGET(AppConfig.CheckerUrl, 1)
	for i := AppConfig.StartPort; i < AppConfig.StartPort+AppConfig.RangePort; i++ {
		ok := strings.Contains(response, AppConfig.Ip+":"+strconv.Itoa(int(i)))
		pid, ok2 := PortPidMap.Load(i)
		if !ok2 {
			continue
		}
		if !ok {
			process, _ := os.FindProcess(pid.(int))
			_ = process.Signal(os.Kill)
			_, _ = process.Wait()
			log.Println("PID: ", pid, "	PORT: ", i, "	STATUS: Killed")
			deletePortInfo(i)
		} else {
			log.Println("PID: ", pid, "	PORT: ", i, "	STATUS: Live")
		}
	}
}
