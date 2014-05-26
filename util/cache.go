package util

import (
	"time"

	"github.com/pmylund/go-cache"
)

var (
	Cache = cache.New(12*time.Hour, 5*time.Minute)
)
