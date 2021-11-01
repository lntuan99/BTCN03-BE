package main

import (
    "busmap.vn/librarycore/api/routers"
    "busmap.vn/librarycore/config"
    "busmap.vn/librarycore/model"
    "fmt"
    "os"
    "os/signal"
)

func main() {
    fmt.Println("=====================================================")
    fmt.Println("Start API Server......")
    fmt.Println("=====================================================")

    // Fetch configs
    config.FetchEnvironmentVariables()

    // Initialize configs
    model.Initialize()

    go func() {
        r := routers.Initialize()
        port := os.Getenv("PORT")
        if port == "" {
            port = config.Config.ApiPort
        }
        r.Run(fmt.Sprintf(":%s", port))
    }()

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)
    <-c

    fmt.Println("=====================================================")
    fmt.Println("API Server has stopped!")
    defer model.DBInstance.Close()
}
