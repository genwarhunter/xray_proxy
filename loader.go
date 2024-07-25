package main

import (
	"log"
	"os"
)

func updatePackages() {
	pack, ok := selectFromPackage()
	if ok {
		pkgIds := make(map[interface{}]bool)
		for _, k := range pack {
			pkgIds[k.Id] = true
			Package.Store(k.Id, k)
		}
		removeOtherKeys(pkgIds, &Package)
	}
}

func loadHashes() bool {
	files, err := os.ReadDir(AppConfig.PathToConfDir)
	if err != nil {
		log.Println(err)
		return false
	}
	for _, file := range files {
		_, ok := hashMap[file.Name()]
		if !ok {
			hashMap[file.Name()] = false
		}
	}
	log.Println("loadHashes OK!")
	return true
}

func updateQueue() {
	for k, v := range hashMap {
		if !v {
			queue <- k
			hashMap[k] = false
		}
	}
}
