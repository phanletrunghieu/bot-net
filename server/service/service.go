package service

import (
	"github.com/phanletrunghieu/botnet/server/service/boss"
	"github.com/phanletrunghieu/botnet/server/service/client"
)

// Service main service
type Service struct {
	BossServie   boss.Service
	ClientServie client.Service
}
