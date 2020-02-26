package handler

import (
	"time"

	"github.com/buggoing/echo-template/logger"
)

const databaseOperationTimeout = 10 * time.Second

var log = logger.New("handler")
