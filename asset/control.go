package asset

import (
	"Peregrine/stru"
	"log"
)

var List map[string]string

func Init(a []stru.Asset) {
	List = make(map[string]string)
	if len(a) <= 0 {
		log.Fatalln("资产数为0。。")
		return
	}
	for _, v := range a {
		List[v.Name] = v.Host
	}
}

func GetAsset(name string) string {
	return List[name]
}
