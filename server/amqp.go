package server

import (
    "fmt"

    "github.com/streadway/amqp"

    "github.com/ugarcia/go_test_common/models"
    "github.com/ugarcia/go_test_common/mq"
)

// Global variables for queue/channel etc.
var q *mq.AMQP

// Constants for this module
const MQ_URL = "amqp://guest:guest@mq.gamewheel.local:5672/"
const QUEUE = "mcp_q"
const ROUTE = "mcp.*"
const ID = "mcp.*"
const EXCHANGE = "mcp"
var EXCHANGES = []string{"mcp", "modules"}


/**
 * Initialize Queues, Channels and Consumer Loops
 */
func InitAMQP() {

    // Create struct
    q = new(mq.AMQP)

    // Init it
    q.Init(MQ_URL)

    // Defer closing
    defer q.Close()

    // Register exchanges
    q.RegisterExchanges(EXCHANGES)

    // Register queues
    q.RegisterQueues([]string{QUEUE})

    // Bind queues
    q.BindQueuesToExchange([]string{QUEUE}, EXCHANGE, ROUTE)

    // Start consuming
    q.Consume(QUEUE, receiveQueueMessage)
}

/**
 * Receives a message from MCP Queue and calls handler
 */
func receiveQueueMessage(msg models.QueueMessage, d amqp.Delivery) {

    // Lookup original sender from message
    switch msg.Source {

        // From websockets origin, send response back
        case "mcp.ws":
            outMsg := models.WsMessage{ BaseMessage: msg.BaseMessage }
            HandleWsResponseMessage(outMsg)

        // TODO: Handle malformed/other messages here
    }

    // TODO: Pass this delivery object along ans send ACK only after finishing everything???
    d.Ack(false)
}

/**
 * Handles a request from WS/HTTP
 */
func HandleRequestMessage(inMsg models.BaseMessage) {

    // Create Queue message from original
    msg := models.QueueMessage {
        BaseMessage: inMsg,
        Sender: inMsg.Source,
    }

    // Lookup message type
    switch(msg.Target) {

        // Request for data, direct to basic module
        case "data":
            msg.Receiver = "modules.basic"

        // TODO: Add more options here

        // Unknown target
        default:
            fmt.Printf("Unknown message target: %s", msg.Target)
            return
    }

    // Call sender handler
    q.SendMessage(msg)
}