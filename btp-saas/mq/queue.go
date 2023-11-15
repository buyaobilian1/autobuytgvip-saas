package mq

import (
	"log"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/internal/config"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/mq/handle"
	"github.com/hibiken/asynq"
)

var srv *asynq.Server
var QueueClient *asynq.Client

func Start(conf config.RedisConf) {
	opt := asynq.RedisClientOpt{
		Addr:     conf.Host,
		Password: conf.Pass,
	}
	initClient(opt)
	go initServer(opt)
}

func Stop() {
	srv.Stop()
}

func initClient(opt asynq.RedisClientOpt) {
	QueueClient = asynq.NewClient(opt)
}

func initServer(opt asynq.RedisClientOpt) {
	srv = asynq.NewServer(
		opt,
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			Logger: CustomLogger{Level: Info},
		},
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc(handle.OrderExpirationPattern, handle.OrderExpirationHandler)
	mux.HandleFunc(handle.RechargeExpirationPattern, handle.RechargeExpirationHandler)
	mux.HandleFunc(handle.GiftTelegramPremiumPattern, handle.GiftTelegramPremiumHandler)
	if err := srv.Run(mux); err != nil {
		log.Fatalf("[queue] could not run server: %v", err)
	}
}

type CustomLoggerLevel = int32

const (
	Debug CustomLoggerLevel = iota
	Info
	Warn
	Error
)

type CustomLogger struct {
	Level CustomLoggerLevel
}

// Debug logs a message at Debug level.
func (l CustomLogger) Debug(args ...interface{}) {
	if l.Level >= Debug {
		log.Println(args...)
	}
}

// Info logs a message at Info level.
func (l CustomLogger) Info(args ...interface{}) {
	if l.Level >= Info {
		log.Println(args...)
	}
}

// Warn logs a message at Warning level.
func (l CustomLogger) Warn(args ...interface{}) {
	if l.Level >= Warn {
		log.Println(args...)
	}
}

// Error logs a message at Error level.
func (l CustomLogger) Error(args ...interface{}) {
	if l.Level >= Error {
		log.Println(args...)
	}
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (l CustomLogger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}
