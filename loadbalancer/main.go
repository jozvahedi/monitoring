package main

import (
	Server "github.com/jozvahedi/loadbalancer/loadbalancer/cmd/httpServer"
	Config "github.com/jozvahedi/loadbalancer/loadbalancer/config"
)

func main() {
	Config.ReadYamlFileOrPanic()
	Config.ReadJsonFileOrPanic()
	Server.HttpServer(Config.ConfFile.HTTPServer.HTTPServerServer, Config.ConfFile.HTTPServer.HTTPServerPort)
}
