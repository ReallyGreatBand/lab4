package engine

import (
	"sync"
)

type commandsQueue struct {
	waitingForReceive bool
	receiveSignal     chan bool
	mutex             sync.Mutex
	tasks             []Command
}

func (q *commandsQueue) push(cmd Command) {
	q.mutex.Lock()
	q.tasks = append(q.tasks, cmd)
	// Because q.tasks is sharing resource
	length := len(q.tasks)
	q.mutex.Unlock()
	// q.waitingForReceive is used because we must not send signal if peek is not
	// waiting for command. Otherwise next signal allows peek to read from empty queue
	if length == 1 && q.waitingForReceive {
		q.receiveSignal <- true
	}
}

func (q *commandsQueue) pull() {
	q.mutex.Lock()
	q.tasks[0] = nil
	q.tasks = q.tasks[1:]
	q.mutex.Unlock()
}

// Peek returns first command from queue but doesn't delete it from queue
func (q *commandsQueue) peek() Command {
	q.mutex.Lock()
	if len(q.tasks) == 0 {
		// Informs push that peek is waiting for a channel signal
		q.waitingForReceive = true
		q.mutex.Unlock()
		// Waits for queue to fill
		<-q.receiveSignal
		q.waitingForReceive = false
	} else {
		q.mutex.Unlock()
	}
	return q.tasks[0]
}

type EventLoop struct {
	stopFlag    bool
	stopChannel chan bool
	queue       commandsQueue
}

type IHandler interface {
	Post(cmd Command)
}

type Handler struct {
	eventLoop *EventLoop
}

func (h *Handler) Post(cmd Command) {
	h.eventLoop.Post(cmd)
}

func (l *EventLoop) Start() {
	handler := Handler{eventLoop: l}
	// This flag is set in AwaitFinish to break for loop when all commands are finished
	l.stopFlag = false
	// This channel blocks main routine waiting for other commands to finish
	l.stopChannel = make(chan bool)
	// This flag is needed to inform push that peek is waiting for push to append command to queue
	l.queue.waitingForReceive = false
	// This channel is needed for waiting while empty queue refills with commands
	l.queue.receiveSignal = make(chan bool)
	l.queue.tasks = make([]Command, 0, 1)

	go func() {
		for {
			if l.stopFlag && len(l.queue.tasks) == 0 {
				break
			}
			l.queue.peek().Execute(&handler)
			l.queue.pull()
		}
		l.stopChannel <- true
	}()
}

func (l *EventLoop) AwaitFinish() {
	if len(l.queue.tasks) != 0 {
		l.stopFlag = true
		<-l.stopChannel
	}
}

func (l *EventLoop) Post(cmd Command) {
	l.queue.push(cmd)
}
