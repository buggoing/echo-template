package handler

import (
	"time"

	"github.com/PPIO/pi-cloud-monitor-backend/logger"
)

const databaseOperationTimeout = 10 * time.Second

var log = logger.New("handler")
