package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func main() {
	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("go-service")
	conn, err := amqp.DialConfig(os.Getenv("RABBIT_CONNECTION_STRING"), config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer amqpChannel.Close()
	queue, err := amqpChannel.QueueDeclare("exchange", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	router := httprouter.New()
	router.POST("/api", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		var t map[string]interface{}
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		newJson := map[string]interface{}{
			"data": t,
			"pattern": queue.Name,
		}

		fmt.Println(newJson)

		newJsonBytes, err := json.Marshal(newJson)

		if err != nil {
			log.Fatal(err)
		}

		err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
			Body: newJsonBytes,
		})

		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)

	})
	log.Fatal(http.ListenAndServe(":8080", router))
}
