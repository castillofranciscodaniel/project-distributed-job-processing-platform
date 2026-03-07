package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting Distributed Job Worker...")

	// Aquí irá la lógica de conexión a la cola de mensajes (SQS/RabbitMQ/etc)
	// y el procesamiento de contratos.

	// Por ahora, simulamos un proceso que se mantiene vivo
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			log.Println("Worker waiting for jobs...")
			time.Sleep(30 * time.Second)
		}
	}()

	log.Println("Worker is running. Press Ctrl+C to stop.")
	<-stop
	log.Println("Shutting down worker...")
}
