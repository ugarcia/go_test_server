package server

import (
    "fmt"
    "encoding/json"
    "golang.org/x/net/websocket"

    "github.com/ugarcia/go_test_common/models"
)

// Global connections collection variables
var connectionsCount uint = 0
var wsConnections = make(map[uint]*websocket.Conn)


/**
 * Gets websockets handler for use in endpoints
 */
func GetWsHandler() websocket.Handler {
    return websocket.Handler(wsHandler)
}

/**
 * Internal handler for websockets connections
 */
func wsHandler(ws *websocket.Conn) {

    // Increment connections count, set the local one and add to collection map
    connectionsCount++
    connectionId := connectionsCount
    wsConnections[connectionId] = ws

    // Endless loop for WS messages listener
    for {

        // Fetch message (looks blocking thing while not messages)
        var in string
        if err := websocket.Message.Receive(ws, &in); err != nil {
            fmt.Println(err.Error())
            break
        }
        fmt.Printf("Received: %s\n", in)

        // Unwrap/parse bytes into our message type
        inMsg := models.WsMessage{}
        if err := json.Unmarshal([]byte(in), &inMsg); err != nil {
            fmt.Println(err.Error())
            break
        }

        // Set sender and connection id into message
        inMsg.Sender = "ws"
        inMsg.ConnectionId = connectionId

        // Check message target and deliver to proper handler
        switch target := inMsg.Target; target {
            case "db":
                handleDbRequest(inMsg)
            default:
                fmt.Printf("Unknown target: %s\n", target)
        }
    }
}

/**
 * Handles a request to DB service
 */
func handleDbRequest(inMsg models.WsMessage) {

    // Create DB Queue message from original
    msg := models.DbQueueMessage {
        BaseMessage:inMsg.BaseMessage,
        Entity: inMsg.Code,
    }

    // Parse message to bytes
    data, err := json.Marshal(msg)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // Call sender handler
    SendDbQueueMessage(data)
}

/**
 * Handles a received message from WS Queue, sending relevant info to connected client
 */
func HandleWsQueueMessage(msg models.WsQueueMessage) {

    // Reference to original connection vars if it applies
    var connId uint
    var conn *websocket.Conn

    // Check original connection if not broadcasting
    if !msg.Broadcast {

        // Try to get original connection identifier
        connId = msg.ConnectionId
        if connId == 0 {
            fmt.Println("No connection specified in queue message!");
            return
        }

        // Try to get original connection (if not broadcasting)
        conn = wsConnections[connId]
        if conn == nil {
            fmt.Printf("No connection found for index %d\n", connId);
            return
        }
    }

    // Try to get code
    code := msg.Code
    if code == "" {
        fmt.Println("No code provided in queue message!");
        return
    }

    // Try to get delivered data
    inData := msg.Data
    if inData == nil {
        fmt.Println("No data provided in queue message!");
        return
    }

    // Prepare message struct
    outMsg := models.WsMessage {
        BaseMessage: msg.BaseMessage,
        Code: msg.Code,
    }

    // Wrap into bytes
    resp, err := json.Marshal(outMsg)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // Get string version
    response := string(resp)

    // Check if we need to broadcast it
    if msg.Broadcast {
        fmt.Println("Broadcasting message thru WS!");
        for _, v := range wsConnections {
            websocket.Message.Send(v, response)
        }

    // No broadcast, send to single connection
    } else {
        websocket.Message.Send(conn, response)
    }
}