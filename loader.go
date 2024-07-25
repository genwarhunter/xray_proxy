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
		_, ok := hashMap.Load(file.Name())
		if !ok {
			hashMap.Store(file.Name(), false)
		}
	}
	log.Println("loadHashes OK!")
	return true
}

func updateQueue() {
	hashMap.Range(func(k, v any) bool {
		if !v.(bool) {
			queue <- k.(string)
			hashMap.Store(k, false)
		}
		return true
	})
}
