package main

import (
	"context"
	// "encoding/json"
	"fmt"
	"time"

	// "fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// refer to https://www.mongodb.com/docs/drivers/go/current/quick-start/#std-label-golang-quickstart
func main() {

	uri := "mongodb://localhost:27017/ycsb?maxPoolSize=1000"
	// if uri == "" {
	// 	log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	// }
	const threadNum = 10
	const queryPerThread = 10000

	var clientArr [threadNum]*mongo.Client // 声明一个包含5个整数的数组

	for i := 0; i < threadNum; i++ {
		var err error
		clientArr[i], err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}

	}

	var wg sync.WaitGroup
	wg.Add(threadNum)
	start := time.Now()
	for i := 0; i < threadNum; i++ {
		go func(client *mongo.Client) {
			defer wg.Done()

			coll := client.Database("ycsb").Collection("usertable")
			for cnt := 0; cnt < queryPerThread; cnt++ {
				var result bson.M

				err := coll.FindOne(context.TODO(), bson.D{}).Decode(&result)
				if err == mongo.ErrNoDocuments {
					// fmt.Printf("No document was found with the title ")
					continue
				}

				if err != nil {
					fmt.Println(err)
				}
				// jsonData, err := json.MarshalIndent(result, "", "    ")
				// _, err = json.MarshalIndent(result, "", "    ")
				// if err != nil {
				// 	panic(err)
				// }
				// fmt.Printf("%s\n", jsonData)
			}

			if err := client.Disconnect(context.TODO()); err != nil {
				panic(err)
			}

		}(clientArr[i])
	}

	wg.Wait()
	elapsed := time.Since(start).Seconds()
	fmt.Println("Elapsed time:", elapsed)
	fmt.Println("qps: ", threadNum*queryPerThread/elapsed)
}
