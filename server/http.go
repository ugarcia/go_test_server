package server

import (
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/contrib/static"

    "github.com/ugarcia/go_test_db_worker/models"
)

/**
    Server initialization
 */
func InitHTTP() {

    // Get new router instance
    router := gin.Default()

    // Add API-like endpoint handler
    router.GET("/games", allGames)

    // Add websockets endpoint handler
    router.GET("/ws", func(c *gin.Context) {
        // Check auth headers here??
        handler := GetWsHandler()
        handler.ServeHTTP(c.Writer, c.Request)
    })

    // Serve static content
    router.Use(static.Serve("/", static.LocalFile("./res/html", true)))

    // Listen and log error if any
    if err := router.Run(":8080"); err != nil {
        log.Fatal(err.Error())
        os.Exit(2)
    }
}


/**
    Games Index Endpoint
 */
func allGames(c *gin.Context) {

    // TODO: Use AMQP instead!!

    // Init Games collection
    games := models.Games{}.All()

    // Send OK JSON response
    c.JSON(http.StatusOK, games)
}
