package server

import (
	"context"
	"e-learn/dbg"
	"e-learn/mongodb"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type Server struct {
	MongoClient *mongodb.MongoClient
	Router      *mux.Router
	Config      Config
}

type Config struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

var DefaultConfig = Config{
	Address: "0.0.0.0",
	Port:    "9999",
}

func (s *Server) getConfig() {
	defer dbg.MonitorFunc("server config init")()

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dbg.ConsoleLog(err)
		s.Config = DefaultConfig
		return
	}

	jsonFile, err := os.Open(fmt.Sprintf("%s%cserver%cconfig.json", path, os.PathSeparator, os.PathSeparator))
	if err != nil {
		dbg.ConsoleLog(err)
		s.Config = DefaultConfig
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
		s.Config = DefaultConfig
		return
	}

	if err := json.Unmarshal(bytes, &s.Config); err != nil {
		// If we have trouble reading config file we use default config
		dbg.ConsoleLog(err)
		s.Config = DefaultConfig
		return
	}
	dbg.ConsoleLog(s.Config)
}

func (s *Server) Init() error {
	s.getConfig()
	s.Router = mux.NewRouter()
	s.MongoClient = &mongodb.MongoClient{}
	return s.MongoClient.InitConn()
}

func (s *Server) BuildHandlers() {
	s.Router.HandleFunc("/register", s.Register).Methods("POST")
	s.Router.HandleFunc("/login", s.Login).Methods("POST")
}

func (s *Server) Run() {
	mc := &mongodb.MongoClient{}
	if err := mc.InitConn(); err != nil {
		log.Fatal(err)
	} else {

	}

	if err := s.Init(); err != nil {
		log.Fatal(err)
	} else {
		dbg.ConsoleLog("successful rest server init")
	}

	s.BuildHandlers()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    s.Config.Address + ":" + s.Config.Port,
		Handler: s.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		time.Sleep(2 * time.Second)
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
