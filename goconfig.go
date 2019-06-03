package goconfig

import (
	"strings"
	"sync"

	"container/list"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var initGoConfigCtx sync.Once
var instanceGoConfig *GoConfig

// GoConfig is the representation of the heirarchial storage and retreival system
type GoConfig struct {
	RedisClient     *redis.Client
	ConfigList      *list.List
	DelimiterForKey string
}

// DefaultRedis allows us to create a redis client by passing respective params
// If you want to default to localhost, 6370 and db 0 with password "", leave all empty and
// pass dbNum as 0.
func DefaultRedis(host string, password string, port string, dbNum int) (*redis.Client, error) {
	if dbNum <= 0 {
		dbNum = 0
	}
	if strings.TrimSpace(host) == "" {
		host = "localhost"
	}
	if strings.TrimSpace(port) == "" {
		port = "6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr:     host + `:` + port,
		Password: password, // no password set
		DB:       dbNum,    // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		return nil,
			errors.Wrapf(err,
				"defaultRedis: Error Pinging redis at %v:%v with password ##### and db Number: %v",
				host,
				port,
				dbNum)
	}
	log.Debugln(pong)
	return client, nil
}

// NewGoConfig return the singleton GoConfig
func NewGoConfig(redisClient *redis.Client, delimiterForKey string) *GoConfig {
	initGoConfigCtx.Do(func() {
		if redisClient == nil {
			var err error
			redisClient, err = DefaultRedis("", "", "", 0)
			if err != nil {
				log.Errorln("Invalid redisClient", err)
				redisClient = nil
			}
		}
		if delimiterForKey == "" {
			delimiterForKey = "."
		}
		instanceGoConfig = &GoConfig{RedisClient: redisClient, ConfigList: list.New(), DelimiterForKey: delimiterForKey}
	})
	return instanceGoConfig
}
func parseKey(keys string, delimiter string) []string {
	return strings.Split(keys, delimiter)
}

// Set will store the value in the respective key. If key
// doesnt exist it will store the key and put the value inside it
// Please note that key will be like a.b which says under a node, b will
// have value. Here separator in key is "."
func (g *GoConfig) Set(key string, value string) error {
	keys := parseKey(key, g.DelimiterForKey)
	if len(keys) == 0 {
		return errors.Errorf("Key is empty")
	}
	err := g.RedisClient.Set(key, value, 0).Err()
	if err != nil {
		return errors.Wrapf(err, "Error during Set key value:%v:%v", key, value)
	}
	return nil
}

// Get will retrieve the value in the respective key. If key
// doesnt exist it will recursively check if parent has the value
// and return the value whoever has. Always child with value gets high preference
// in Get.
// Please note that key will be like a.b which says under a node, b will
// might or might not have value. Here separator in key is "."
// If b exists, then its value is returned, else if a exists , then its value is returned
// else error returned
func (g *GoConfig) Get(keyStr string) (string, error) {
	keys := parseKey(keyStr, g.DelimiterForKey)
	currentConstructedKey := keyStr
	for i := 0; i < len(keys); i++ {
		val, err := g.RedisClient.Get(currentConstructedKey).Result()
		if err != nil {
			// go to next loop to find the value
			currentConstructedKey = lastItemTrim(currentConstructedKey, g.DelimiterForKey)
			continue
		}
		return val, nil
	}
	return "", errors.Errorf("No value available")
}

// Delete will Delete the key in redis store
func (g *GoConfig) Delete(key string) error {
	err := g.RedisClient.Del(key).Err()
	if err != nil {
		return errors.Wrapf(err, "Error during Delete key value:%v", key)
	}
	return nil
}

// GetAll will get all list of key values as map[string]string
func (g *GoConfig) GetAll(keyStr string) (map[string]string, error) {
	keys := parseKey(keyStr, g.DelimiterForKey)
	currentConstructedKey := keyStr
	var returnMap = make(map[string]string)
	for i := 0; i < len(keys); i++ {
		val, err := g.RedisClient.Get(currentConstructedKey).Result()
		if err != nil {
			// go to next loop to find the value
			currentConstructedKey = lastItemTrim(currentConstructedKey, g.DelimiterForKey)
			continue
		}
		returnMap[currentConstructedKey] = val
		currentConstructedKey = lastItemTrim(currentConstructedKey, g.DelimiterForKey)
	}
	return returnMap, nil
}

func lastItemTrim(currentKey string, delimiter string) string {
	keys := parseKey(currentKey, delimiter)
	return strings.Join(keys[0:len(keys)-1], ".")
}
