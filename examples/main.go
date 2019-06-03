package main

import (
	config "github.com/karuppaiah/goconfig"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Handle errors accordingly
	// Use default redis config as localhost, port, password and db as 0
	redisClient, err := config.DefaultRedis("", "", "", 0)
	if err != nil {
		log.Errorln(err)
	}
	goConfig := config.NewGoConfig(redisClient, ".")
	goConfig.DelimiterForKey = "."
	err = goConfig.Set("a", "a")
	if err != nil {
		log.Errorln(err)
	}
	err = goConfig.Set("a.b", "ab")
	if err != nil {
		log.Errorln(err)
	}
	err = goConfig.Set("a.b.c", "abc")
	if err != nil {
		log.Errorln(err)
	}
	err = goConfig.Set("a.b.d", "abd")
	if err != nil {
		log.Errorln(err)
	}
	val, err := goConfig.Get("a.b.c")
	if err != nil {
		log.Errorln(err)
	}
	log.Println(val)
	val, err = goConfig.Get("a.b.d")
	if err != nil {
		log.Errorln(err)
	}
	log.Println(val)
	val, err = goConfig.Get("a.b.c.d")
	if err != nil {
		log.Errorln(err)
	}
	log.Println(val)
	val, err = goConfig.Get("a.b.d.d")
	if err != nil {
		log.Errorln(err)
	}
	log.Println(val)
	val, err = goConfig.Get("v")
	if err != nil {
		log.Errorln(err)
	}
	log.Println(val)
	mapVal, err := goConfig.GetAll("a.b.d.d")
	if err != nil {
		log.Errorln(err)
	}
	log.Println(mapVal)
	mapVal, err = goConfig.GetAll("v")
	if err != nil {
		log.Errorln(err)
	}
	log.Println(mapVal)
}
