package main

import (
	"fmt"
	//"os"

	"libraries/utils"

	"libraries/toml"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	//file, _ := os.Create("test.log")
	//log.SetOutput(file)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	clientBlocks := make(map[uint64]map[uint64]bool)
	blockStruct := make(map[uint64]bool)
	clientBlocks[1001] = blockStruct
	blockStruct[1002] = true
	log.Info(clientBlocks)
	log.Info(clientBlocks[1001][1002])

	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	/**
	log.WithFields(log.Fields{
		"omg":    true,
		"number": 100,
	}).Fatal("The ice breaks!")
	*/

	arr := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s := arr[0:4]
	fmt.Println(s)

	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")

	blogArticleViews := map[string]int{
		"unix":         0,
		"python":       1,
		"go":           2,
		"javascript":   3,
		"testing":      4,
		"philosophy":   5,
		"startups":     6,
		"productivity": 7,
		"hn":           8,
		"reddit":       9,
		"C++":          10,
	}
	for key, views := range blogArticleViews {
		fmt.Println("There are", views, "views for", key)
	}

	tomlConfig, err := toml.LoadTomlConfig("etc/config.toml")
	if err != nil {
		panic(err)
	}

	etcdClient := utils.EtcdDail(tomlConfig.Etcd)
	_, err = etcdClient.Put(context.Background(), "worker1", "115.29.164.163:18600")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := etcdClient.Get(context.Background(), "worker1")
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}
