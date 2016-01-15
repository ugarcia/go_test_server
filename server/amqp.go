package server

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/streadway/amqp"
    "github.com/ugarcia/go_test_common/models"
)

// Global variables for DB queue/channel
var dbQueue amqp.Queue
var dbChannel *amqp.Channel

/**
 * Helper function for handling errors
 */
func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

/**
 * Initialize Queues, Channels and Consumer Loops
 */
func InitAMQP() {

    // Create connection for DB
    dbConn, err := amqp.Dial("amqp://guest:guest@mq.gamewheel.local:5672/")
    failOnError(err, "Failed to connect to RabbitMQ DB")
    defer dbConn.Close()

    // Create DB Channel
    dbChannel, err = dbConn.Channel()
    failOnError(err, "Failed to open DB channel")
    defer dbChannel.Close()

    // Create DB Queue
    dbQueue, err = dbChannel.QueueDeclare(
        "db", // name
        true,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare DB queue")

    // Create connection for WS
    wsConn, err := amqp.Dial("amqp://guest:guest@mq.gamewheel.local:5672/")
    failOnError(err, "Failed to connect to RabbitMQ WS")
    defer wsConn.Close()

    // Create WS Channel
    wsChannel, err := wsConn.Channel()
    failOnError(err, "Failed to open WS channel")
    defer wsChannel.Close()

    // Create WS Queue
    wsQueue, err := wsChannel.QueueDeclare(
        "ws", // name
        true,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare WS queue")

    err = wsChannel.Qos(
        1,     // prefetch count
        0,     // prefetch size
        false, // global
    )
    failOnError(err, "Failed to set QoS")

    // Consume WS Channel messages
    msgs, err := wsChannel.Consume(
        wsQueue.Name, // queue
        "", // consumer
        false, // auto-ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil, // args
    )
    failOnError(err, "Failed to register WS consumer")

    // Loop forever for messages, and call handler for each
    forever := make(chan bool)
    go func() {
        for d := range msgs {
            log.Printf("Received a message: %s", d.Body)
            d.Ack(false)
            receiveWsQueueMessage(d)
        }
    }()
    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}

/**
 * Receives a message from WS Queue and calls WS handler
 */
func receiveWsQueueMessage(msg amqp.Delivery) {
    // TODO: Check delivery parameters first?
    data := models.WsQueueMessage{}
    if err := json.Unmarshal([]byte(msg.Body), &data); err != nil {
        fmt.Println(err.Error())
        return
    }
    HandleWsQueueMessage(data)
}

/**
 * Sends a message to DB Queue
 */
func SendDbQueueMessage(msg []byte) {
    sendQueueMessage(msg, dbChannel, dbQueue.Name)
}

/**
 * Sends a message to a Queue
 */
func sendQueueMessage(msg []byte, ch *amqp.Channel, queue string) {
    err := ch.Publish(
        "",     // exchange
        queue, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
            DeliveryMode: amqp.Persistent,
            ContentType: "application/json",
            Body:        []byte(msg),
        })
    failOnError(err, "Failed to publish a message")
    log.Printf(" [x] Sent %s", string(msg))
}
