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

        // TODO: Perform Token validation etc.

        // Set source and connection data into message
        // TODO: Possible to init this in model struct?
        inMsg.Source = "mcp.ws"
        inMsg.ConnectionType = "websocket"
        inMsg.ConnectionId = connectionId

        // Call handler
        go HandleRequestMessage(inMsg.BaseMessage)
    }
}

/**
 * Handles a received message from WS Queue, sending relevant info to connected client
 */
func HandleWsResponseMessage(msg models.WsMessage) {

    // Reference to original connection vars if it applies
    var connId uint
    var conn *websocket.Conn

    // TODO: Move these validations to some common place

    // Try to get target
    target := msg.Target
    if target == "" {
        fmt.Println("No target provided in queue message!");
        return
    }

    // Try to get code
    action := msg.Action
    if action == "" {
        fmt.Println("No action provided in queue message!");
        return
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

    // Logic to determine if broadcast is necessary
    msg.Broadcast = false
    switch msg.Target {
    case "data":
        switch msg.Action {
        case "post", "delete", "update":
            msg.Broadcast = true
        }
    }

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

    // Prepare message struct
    outMsg := models.WsMessage { BaseMessage: msg.BaseMessage }

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
