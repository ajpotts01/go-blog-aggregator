package api

import "github.com/ajpotts01/go-blog-aggregator/internal/database"

type ApiConfig struct {
	DbConn            *database.Queries
	MaxFeedsProcessed int
}
