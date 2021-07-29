package main

import (
    "context"
    "net/http"
    "os"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v4"
)

func getHi(c *gin.Context) {
    response := make(map[string]string)
    response["message"] = "hi"

    c.IndentedJSON(http.StatusOK, response)
}

func getHealth(c *gin.Context) {

    dbConnectionUrl := os.Getenv("DB_CONNECTION_URL")

    conn, err := pgx.Connect(context.Background(), dbConnectionUrl)
    if err != nil {
        log.Println("Unable to connect to database: %v", err)
        c.String(http.StatusInternalServerError, "db connection error")
        return
    }

    defer conn.Close(context.Background())

    var one int64
    err = conn.QueryRow(context.Background(), "select 1").Scan(&one)
    if err != nil {
        log.Println("QueryRow failed: %v", err)
        c.String(http.StatusInternalServerError, "db query failed")
        return
    }

    response := make(map[string]string)
    response["healthy"] = "true"

    c.IndentedJSON(http.StatusOK, response)
}

func main() {
    port := os.Getenv("PORT")

    router := gin.Default()
    router.GET("/hi", getHi)
    router.GET("/health", getHealth)

    router.Run(":" + port)
}
