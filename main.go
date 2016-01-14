package main

import (
    "github.com/ugarcia/go_test_server/server"
)

/**
    Main execution
 */
func main() {

    // Init amqp connection
    go server.InitAMQP()

    // Init HTTP server
    server.InitHTTP()
}
