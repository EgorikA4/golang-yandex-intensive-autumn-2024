package application

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/pkg/calculation"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/pkg/server"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) Run() error {
	for {
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
		}
		text = strings.TrimSpace(text)
		if text == "exit" {
			log.Println("aplication was successfully closed")
			return nil
		}
		result, err := calculation.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed wit error: ", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate/", server.CalcHandler)
	return http.ListenAndServe(":" + a.config.Addr, nil)
}
