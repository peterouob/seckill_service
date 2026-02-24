package mq

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/IBM/sarama"
)

type ConsumeHandler interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}

type ConsumerGroup struct {
	ready       chan bool
	message     chan *sarama.ConsumerMessage
	batchSize   int
	flushTime   time.Duration
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	mu          sync.Mutex
	batchBuffer []*sarama.ConsumerMessage
	lastCommit  time.Time
	rdb         *redis.Client
}

var _ ConsumeHandler = (*ConsumerGroup)(nil)

func (l *ConsumerGroup) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("Consumer group session setup for member: %s\n", session.MemberID())
	close(l.ready)
	l.wg.Add(1)
	go l.processBatches(session)
	return nil
}

func (l *ConsumerGroup) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session cleanup for member: ", session.MemberID())
	l.cancel()
	l.wg.Wait()
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.batchBuffer) > 0 {
		log.Println("Processing remaining messages in buffer during cleanup...")
		l.commit(session, l.batchBuffer)
		l.batchBuffer = nil
	}
	return nil
}

func (l *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	batchSize := l.batchSize
	batch := make([]*sarama.ConsumerMessage, 0, batchSize)
	ticker := time.NewTicker(l.flushTime)
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			batch = append(batch, msg)
			if len(batch) >= batchSize {
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

func NewLikeHandler(batchSize int, flushTime time.Duration) *ConsumerGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConsumerGroup{
		ready:      make(chan bool),
		message:    make(chan *sarama.ConsumerMessage, 10000),
		batchSize:  batchSize,
		flushTime:  flushTime,
		ctx:        ctx,
		cancel:     cancel,
		lastCommit: time.Now(),
	}
}

func (l *ConsumerGroup) processBatches(session sarama.ConsumerGroupSession) {
	defer l.wg.Done()
	ticker := time.NewTicker(l.flushTime)
	defer ticker.Stop()

	for {
		select {
		case msg := <-l.message:

			var curBatch []*sarama.ConsumerMessage
			l.mu.Lock()
			l.batchBuffer = append(l.batchBuffer, msg)

			if len(l.batchBuffer) >= l.batchSize {
				log.Printf("Batch size reached (%d), processing batch.\n", len(l.batchBuffer))
				curBatch = l.batchBuffer
				l.batchBuffer = nil
				l.lastCommit = time.Now()
			}
			l.mu.Unlock()

			if curBatch != nil {
				l.commit(session, curBatch)
			}
		case <-ticker.C:
			var curBatch []*sarama.ConsumerMessage
			l.mu.Lock()
			if len(l.batchBuffer) > 0 && time.Since(l.lastCommit) >= l.flushTime {
				log.Printf("Batch interval reached, processing batch (%d messages).\n", len(l.batchBuffer))
				curBatch = l.batchBuffer
				l.batchBuffer = nil
				l.lastCommit = time.Now()
			}
			l.mu.Unlock()
			if curBatch != nil {
				l.commit(session, curBatch)
			}
		case <-l.ctx.Done():
			log.Println("Batch processing goroutine shutting down.")
			return
		case <-session.Context().Done():
			log.Println("Session context cancelled, batch processing goroutine shutting down.")
			return
		}
	}
}

func (l *ConsumerGroup) commit(session sarama.ConsumerGroupSession, batch []*sarama.ConsumerMessage) {
	//counts := make(map[string]int)
	//for _, topic := range batch {
	//	var social model.SocialAction
	//	if err := json.Unmarshal(topic.Value, &social); err != nil {
	//		log.Println("error in json unmarshal")
	//	}
	//	counts[fmt.Sprintf("%s:%s", social.PostID, social.ThumbType)] += social.Num
	//}
	//
	//if errors.Is(model.ErrPipe, model.RunScript(l.ctx, counts)) {
	//	log.Println("error in run script")
	//}
	for _, msg := range batch {
		session.MarkMessage(msg, "")
	}
	session.Commit()
}
