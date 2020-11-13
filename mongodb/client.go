package mongodb

import (
	"context"
	"e-learn/dbg"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type MongoClient struct {
	Client *mongo.Client
}

var DefaultDbConfig = MongoDb{
	Url:    "mongodb+srv://go-server:DfQ7i-hj9pnWD45@cluster0-r16y6.mongodb.net/test?retryWrites=true&w=majority",
	DbName: "e-learn",
}

var MyDb MongoDb

// Here we read mongo config file
func init() {
	dbg.SetDebug(true)
	defer dbg.MonitorFunc("mongodb config init")()

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dbg.ConsoleLog(err)
		MyDb = DefaultDbConfig
		return
	}

	jsonFile, err := os.Open(fmt.Sprintf("%s%cmongodb%cconfig.json", path, os.PathSeparator, os.PathSeparator))
	if err != nil {
		dbg.ConsoleLog(err)
		MyDb = DefaultDbConfig
		return
	}
	defer func() {
		if err := jsonFile.Close(); err != nil {
			dbg.ConsoleLog("bad config file close")
		}
	}()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		dbg.ConsoleLog(err)
		MyDb = DefaultDbConfig
		return
	}

	if err := json.Unmarshal(bytes, &MyDb); err != nil {
		// If we have trouble reading config file we use default config
		dbg.ConsoleLog(err)
		MyDb = DefaultDbConfig
		return
	}
}

func (mc *MongoClient) InitConn() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	dbg.ConsoleLog(MyDb.Url)
	clientOptions := options.Client().ApplyURI(MyDb.Url)
	var err error
	dbg.ConsoleLog(clientOptions)
	if mc.Client, err = mongo.Connect(ctx, clientOptions); err != nil {
		return err
	}

	if err := mc.Client.Ping(ctx, nil); err != nil {
		return err
	}

	return nil
}
