package utils

import (
    "github.com/sirupsen/logrus"
    "os"
    "time"
    "fmt"
)

var Logger = logrus.New()

func init() {
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        Logger.Out = os.Stdout
        Logger.Warn("Failed to log to file, using default stderr")
    } else {
        Logger.Out = file
    }

    Logger.SetFormatter(&logrus.JSONFormatter{})
    Logger.SetLevel(logrus.InfoLevel)
}

func GenerateStreamID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}
