package queue

import (
	"github.com/AskUbuntu/tbot/config"
	"github.com/AskUbuntu/tbot/scraper"

	"path"
	"time"
)

// Queue manages the list of items to be tweeted.
type Queue struct {
	data           *data
	settings       *settings
	messageInChan  <-chan *scraper.Message
	messageOutChan chan<- *scraper.Message
	closeChan      chan bool
}

func (q *Queue) sendIfAvailable() bool {
	q.data.Lock()
	defer q.data.Unlock()
	if len(q.data.QueuedMessages) > 0 {
		q.messageOutChan <- q.data.QueuedMessages[0]
		q.data.LastMessage = time.Now()
		q.data.QueuedMessages = append(q.data.QueuedMessages[1:])
		// TODO: error handling
		q.data.save()
		return true
	}
	return false
}

func (q *Queue) run() {
	for {
		q.data.Lock()
		lastMessage := q.data.LastMessage
		q.data.Unlock()
		q.settings.Lock()
		queueFrequency := q.settings.QueueFrequency
		q.settings.Unlock()
		var (
			now      = time.Now()
			duration = time.Duration(queueFrequency)
			diff     = lastMessage.Add(duration * time.Minute).Sub(now)
		)
		if diff <= 0 && q.sendIfAvailable() {
			continue
		}
		var (
			timer     *time.Timer
			timerChan <-chan time.Time
			quit      = false
		)
		if diff > 0 {
			timer = time.NewTimer(diff)
			timerChan = timer.C
		}
		select {
		case t := <-q.messageInChan:
			q.data.Lock()
			q.data.QueuedMessages = append(q.data.QueuedMessages, t)
			// TODO: error handling
			q.data.save()
			q.data.Unlock()
		case <-timerChan:
		case <-q.closeChan:
			quit = true
		}
		if timer != nil && !timer.Stop() {
			<-timer.C
		}
		if quit {
			break
		}
	}
	close(q.closeChan)
}

// New creates a new queue, loading existing data from disk if available. The
// queue also launches a goroutine to manage tweets.
func New(config *config.Config, inChan <-chan *scraper.Message, outChan chan<- *scraper.Message) (*Queue, error) {
	q := &Queue{
		data:           &data{name: path.Join(config.DataPath, "queue_data.json")},
		settings:       &settings{name: path.Join(config.DataPath, "queue_settings.json")},
		messageInChan:  inChan,
		messageOutChan: outChan,
		closeChan:      make(chan bool),
	}
	if err := q.data.load(); err != nil {
		return nil, err
	}
	if err := q.settings.load(); err != nil {
		return nil, err
	}
	go q.run()
	return q, nil
}

// Messages retrieves the current list of queued messages.
func (q *Queue) Messages() []*scraper.Message {
	q.data.Lock()
	defer q.data.Unlock()
	return q.data.QueuedMessages
}

// Close shuts down the queue and waits for the goroutine to exit.
func (q *Queue) Close() {
	q.closeChan <- true
	<-q.closeChan
}