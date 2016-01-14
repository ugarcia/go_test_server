# go_test_server
Testing Go Capabilities

06.01.2016:
-----------
- An HTTP server which does:
    - Serve static content
    - API endpoint, communicating directly with DB Worker process
    - Handle Websocket connections from/to clients
    - Publish job messages into DB Worker Queue (as a result of received websocket ones)
    - Consume job responses from Workers into WS Queue (and delivers/broadcasts results to client(s) thru websockets)
    Pending:
        - Use Queues for API too, instead of directly use the DB Worker package (coupling ...)
        - Pub/Sub? Only if relevant for our use cases ...
        - Reorganize code in a better way, using structs simulated OO

Easiest way to test it is running the process and browsing the html examples

    go run src/github.com/ugarcia/go_test_server/main.go
    http://localhost:8080  (and click examples - right now only one)

So far, after some load tests using both intervals for websockets/queues and siege for API/static, it looks pretty fast compared to others.

Dependencies:
-------------
- RabbitMq >= 3.5.1
- Golang >= 1.5.x (and proper system config for $GOPATH and $GOROOT)
- Package github.co/ugarcia/go_test_db_common
