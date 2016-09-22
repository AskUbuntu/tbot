package main

import (
	"path"
	"time"
)

// Queue manages the list of items to be tweeted. The QueueFrequency setting
// determines how often items are removed from the list and tweeted. The struct
// is serializable so that state can be preserved across restarts. When a tweet
// is ready, it will be sent on the Tweet channel.
type Queue struct {
	settings     *Settings
	name         string
	notifyChan   chan bool   // Settings changed
	tweetInChan  chan string // Tweet added to queue
	closeChan    chan bool   // Shutdown the queue
	LastTweet    time.Time   `json:"last_tweet"` // Time of last tweet
	QueuedTweets []string    `json:"tweets"`     // List of tweets in queue
	Tweet        chan string `json:"-"`          // Tweet to be sent
}

// run queues messages waiting to be tweeted and dispatches them when ready.
func (q *Queue) run() {
	for {
		// Determine the difference between the earliest time a tweet could be
		// sent and the current time - if the value is zero or below, a tweet
		// can be sent immediately if one is available - otherwise wait for the
		// timer to expire or a tweet to come in (or a shutdown)
		var (
			n    = time.Now()
			diff = q.LastTweet.Add(
				time.Duration(q.settings.QueueFrequency) * time.Minute,
			).Sub(n)
		)
		if diff <= 0 && len(q.QueuedTweets) > 0 {
			q.Tweet <- q.QueuedTweets[0]
			q.LastTweet = n
			q.QueuedTweets = append(q.QueuedTweets[1:])
			q.Save()
			continue
		}
		var (
			timer = time.NewTimer(diff)
			quit  = false
		)
		select {
		case t := <-q.tweetInChan:
			q.QueuedTweets = append(q.QueuedTweets, t)
			q.Save()
		case <-timer.C:
		case <-q.notifyChan:
		case <-q.closeChan:
			quit = true
		}
		if !timer.Stop() {
			<-timer.C
		}
		if quit {
			break
		}
	}
	q.closeChan <- true
}

// NewQueue creates a new queue, loading existing data from disk if available.
// The queue also launches a goroutine to manage tweets.
func NewQueue(config *Config, settings *Settings) (*Queue, error) {
	q := &Queue{
		settings:    settings,
		name:        path.Join(config.DataPath, "queue.json"),
		notifyChan:  make(chan bool),
		tweetInChan: make(chan string),
		closeChan:   make(chan bool),
		Tweet:       make(chan string),
	}
	_, err := LoadJSON(q.name, q)
	if err != nil {
		return nil, err
	}
	settings.Register(q.notifyChan)
	go q.run()
	return q, nil
}

// Save writes the queue to disk.
func (q *Queue) Save() error {
	return SaveJSON(q.name, q)
}

// Close shuts down the queue and waits for the goroutine to exit.
func (q *Queue) Close() {
	q.closeChan <- true
	<-q.closeChan
}
