package common

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(stop func()) {
	log.Printf("Starting server ... \n")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting Down menu ... \n")

	if stop != nil {
		stop()
	}
}
