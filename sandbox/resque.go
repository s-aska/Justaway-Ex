package main

import (
	"github.com/kavu/go-resque"
	_ "github.com/kavu/go-resque/godis" // Use godis driver
	"github.com/simonz05/godis/redis"
)

func main() {
	client := redis.New("tcp:127.0.0.1:6379", 0, "")
	enqueuer := resque.NewRedisEnqueuer("godis", client)
	enqueuer.Enqueue("resque:queue:default", "CreateActivity", "69886580", "@bookma_org_testさんがリツイートしました")
}
