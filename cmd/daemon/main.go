package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Iniciando Sync Agent...")
	go func() {

	}()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Println("\n🛑 Señal de interrupción recibida. Apagando Sync Agent de forma segura...")
}
