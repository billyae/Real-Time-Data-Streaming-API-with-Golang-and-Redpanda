package utils

import (
    "fmt"
    "time"
    "github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
    Logger.SetFormatter(&logrus.JSONFormatter{})
}

func GenerateStreamID() string {
    // Simple random string generator for stream ID
    return fmt.Sprintf("%d", time.Now().UnixNano())
}
