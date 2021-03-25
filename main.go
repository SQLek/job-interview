package main

import (
	"log"

	"github.com/JeremyLoy/config"
	"github.com/SQLek/wp-interview/model"
	"github.com/SQLek/wp-interview/server"
	"github.com/SQLek/wp-interview/worker"
)

type Config struct {
	S server.Config
	W worker.Config
	M model.Config
}

func main() {
	c := Config{}
	log.Println("App is sterting...")

	// TODO
	// i also do not like doubling code because of borked library
	// maybe this one? https://github.com/ilyakaznacheev/cleanenv
	err := config.FromEnv().To(&c.S)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = config.FromEnv().To(&c.W)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = config.FromEnv().To(&c.M)
	if err != nil {
		log.Fatalln(err.Error())
	}

	model, err := model.InitModel(c.M)
	if err != nil {
		log.Fatalln(err.Error())
	}

	sche := worker.MakeSheduler(c.W, model)

	server := server.MakeServer(c.S, model, sche)

	log.Println("App is sterted.")
	log.Fatal(server.ListenAndServe())
	log.Println("App shut down.")
}
