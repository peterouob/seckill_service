package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/peterouob/seckill_service/app/seckill-service/pkg/model"
	"github.com/peterouob/seckill_service/pkg/logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ConsumeHandler interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}

type ConsumerGroup struct {
	batchSize int
	flushTime time.Duration
	db        *gorm.DB
	rdb       *redis.Client
	ready     chan bool
}

var _ ConsumeHandler = (*ConsumerGroup)(nil)

func (l *ConsumerGroup) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("Consumer group session setup for member: %s\n", session.MemberID())
	return nil
}

func (l *ConsumerGroup) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session cleanup for member: ", session.MemberID())
	return nil
}

func (l *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	batch := make([]*sarama.ConsumerMessage, 0, l.batchSize)
	ticker := time.NewTicker(l.flushTime)
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				if len(batch) != 0 {
					l.commit(session, batch)
				}
				return nil
			}
			batch = append(batch, msg)
			if len(batch) >= l.batchSize {
				l.commit(session, batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				l.commit(session, batch)
				batch = batch[:0]
			}
		case <-session.Context().Done():
			if len(batch) > 0 {
				l.commit(session, batch)
			}
			return nil
		}
	}
}

func NewSeckillHandler(batchSize int, flushTime time.Duration, db *gorm.DB) *ConsumerGroup {
	return &ConsumerGroup{
		batchSize: batchSize,
		flushTime: flushTime,
		db:        db,
		ready:     make(chan bool, 1),
	}
}

func (l *ConsumerGroup) commit(session sarama.ConsumerGroupSession, batch []*sarama.ConsumerMessage) {
	var orders []model.Order
	for _, topic := range batch {
		var order model.Order
		if err := json.Unmarshal(topic.Value, &order); err != nil {
			logger.Error("error to unmarshal json order", err)
		}
		orders = append(orders, order)
	}

	if err := l.db.Create(&orders).Error; err != nil {
		logger.Error("error to create order", err)
	}

	for _, msg := range batch {
		session.MarkMessage(msg, "")
	}

	session.Commit()
}
