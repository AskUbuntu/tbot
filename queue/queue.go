package queue

import (
	"github.com/AskUbuntu/tbot/config"
	"github.com/AskUbuntu/tbot/scraper"

	"errors"
	"path"
	"reflect"
	"time"
)

// Queue manages the list of items to be tweeted.
type Queue struct {
	data     *data
	settings *settings
	messages chan<- *scraper.Message
	trigger  chan bool
}

func (q *Queue) sendIfAvailable() bool {
	q.data.Lock()
	defer q.data.Unlock()
	if len(q.data.QueuedMessages) > 0 {
		q.messages <- q.data.QueuedMessages[0]
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
		case <-timerChan:
		case quit = <-q.trigger:
		}
		if timer != nil {
			timer.Stop()
		}
		if quit {
			break
		}
	}
	close(q.trigger)
}

// New creates a new queue, loading existing data from disk if available. The
// queue also launches a goroutine to manage tweets.
func New(config *config.Config, messages chan<- *scraper.Message) (*Queue, error) {
	q := &Queue{
		data:     &data{name: path.Join(config.DataPath, "queue_data.json")},
		settings: &settings{name: path.Join(config.DataPath, "queue_settings.json")},
		messages: messages,
		trigger:  make(chan bool),
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

// Add inserts a message into the queue.
func (q *Queue) Add(m *scraper.Message) error {
	q.data.Lock()
	defer q.data.Unlock()
	q.data.QueuedMessages = append(q.data.QueuedMessages, m)
	if err := q.data.save(); err != nil {
		return err
	}
	q.trigger <- false
	return nil
}

// Delete removes the specified message from the queue.
func (q *Queue) Delete(id int) error {
	q.data.Lock()
	defer q.data.Unlock()
	for i, m := range q.data.QueuedMessages {
		if m.ID == id {
			q.data.QueuedMessages = append(
				q.data.QueuedMessages[:i],
				q.data.QueuedMessages[i+1:]...,
			)
			return q.data.save()
		}
	}
	return errors.New("invalid message index")
}

// Settings retrieves the current settings for the queue.
func (q *Queue) Settings() Settings {
	q.settings.Lock()
	defer q.settings.Unlock()
	return q.settings.Settings
}

// SetSettings stores the current settings for the queue.
func (q *Queue) SetSettings(settings Settings) error {
	q.settings.Lock()
	defer q.settings.Unlock()
	if !reflect.DeepEqual(settings, q.settings.Settings) {
		q.settings.Settings = settings
		q.trigger <- false
		return q.settings.save()
	}
	return nil
}

// LastTweet retrieves the time of the last tweet being sent.
func (q *Queue) LastTweet() time.Time {
	q.data.Lock()
	defer q.data.Unlock()
	return q.data.LastMessage
}

// Close shuts down the queue and waits for the goroutine to exit.
func (q *Queue) Close() {
	q.trigger <- true
	<-q.trigger
}
