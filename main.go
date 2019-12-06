package main

import (
	"os"

	"github.com/streadway/amqp"
)

func main() {
	// Get the connection string from the environment variable
	url := os.Getenv("AMQP_URL")

	// If it doesn't exist, use the default connection string.
	if url == "" {
		// Don't do this in production, this is for testing purpose only.
		url = "amqp://guest:guest@localhost:5673"
	}

	// Connect to rabbitmq instance
	conn, err := amqp.Dial(url)

	if err != nil {
		panic("Could not establish connection with rabbitmq:" + err.Error())
	}

	// Create a channel from the connection. We'll use channels to access the data in the
	// queue rather than the connection itself.
	channel, err := conn.Channel()

	if err != nil {
		panic("Could not open rabbitmq channel:" + err.Error())
	}

	// We create an exchange that will bind to the queue to send and receive messages.
	err = channel.ExchangeDeclare("events", "topic", true, true, true, true, nil)

	if err != nil {
		panic(err)
	}

	// We create a message to be sent to the queue.
	// It has to be an instance of the aqmp publishing struct
	message := amqp.Publishing{
		Body: []byte("Hello Suchaimi"),
	}

	// We publish the message to the exchange we created earlier
	err = channel.Publish("events", "random", false, false, message)

	if err != nil {
		panic("Error publishing a message to the queue.")
	}

	// We create a queue named test.
	_, err = channel.QueueDeclare("test", true, false, false, false, nil)

	if err != nil {
		panic("Error declaring the queue:" + err.Error())
	}

	// We bind the queue to the exchange to send and receive data from the queue.
	err = channel.QueueBind("test", "#", "events", false, nil)

	if err != nil {
		panic("Error binding to the queue:" + err.Error())
	}
}
