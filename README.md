# GoConfig

This is an hierarchial config manager using redis. In case of website handling multiple country, we might need to maintain a config for payment username. But this key will vary for each country as shown below
payment_username.sg = "hhhhh"
payment_username.au = "yyyyy"
all other countries will use:
payment_username = "aaaaa"

In this case, if we feed in all these to our GoConfig
payment_username = "aaaaa"
payment_username.sg = "hhhhh"
payment_username.au = "yyyyy"

and
If we ask for
payment_username.kr, it will return back "aaaaa"
payment_username.sg, it will return back "hhhhh" accordingly.
## Installation and Play

* go get github.com/karuppaiah/goconfig
* cd $GOPATH/src/github.com/karuppaiah/goconfig
* install dep(https://github.com/golang/dep)
* dep ensure
* cd examples
* go run main.go 


## Usage
```
    // Handle errors accordingly
	redisClient, err := config.DefaultRedis("", "", "", 0)
	if err != nil {
		log.Errorln(err)
	}
	goConfig := config.NewGoConfig(redisClient, ".")
	err = goConfig.Set("a.b.c", "a")
	if err != nil {
		log.Errorln(err)
	}
    val, err := goConfig.Get("a.b.c")
	if err != nil {
		log.Errorln(err)
	}
```


# TODO :
- [ ] Instead of storing just string values. See possibility to store any type of value.
- [ ] Option to Load Data on startup and keep it in memory from DB
- [ ] Docker image and publish in hub.docker.com
- [ ] Write unit testing
- [ ] Flutter/web assembly frontend to playaround (++)