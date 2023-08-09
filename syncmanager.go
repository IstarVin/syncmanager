package syncmanager

func NewSyncManager(concurrent int) *SyncManager {
	syncManager := new(SyncManager)

	syncManager.concurrent = concurrent
	syncManager.addQueue = make(chan byte)
	syncManager.doneQueue = make(chan byte)
	syncManager.exitCh = make(chan byte)

	go syncManager.daemon()

	return syncManager
}

func NewSyncManagerWait(concurrent int) *SyncManager {
	syncManager := NewSyncManager(concurrent)

	syncManager.shouldWait = true
	syncManager.Wait = func() {
		syncManager.isWaiting = true
		<-syncManager.exitCh
	}

	return syncManager
}

type SyncManager struct {
	QueueList []*Queue
	Wait      func()

	addQueue  chan byte
	doneQueue chan byte
	exitCh    chan byte

	concurrent int
	shouldWait bool
	isWaiting  bool
}

func (s *SyncManager) Add(f func(...any), args ...any) {
	s.QueueList = append(s.QueueList, &Queue{
		Func: f,
		Args: args,
	})
	s.addQueue <- 0
}

func (s *SyncManager) daemon() {
	var running int
	for {
		select {
		case <-s.addQueue:
		case <-s.doneQueue:
			running--
		}

		if running < s.concurrent && len(s.QueueList) != 0 {
			queue := s.QueueList[0]
			s.QueueList = s.QueueList[1:]

			go func() {
				queue.Func(queue.Args...)
				s.doneQueue <- 0
			}()

			running++
		}

		if s.shouldWait && s.isWaiting && running == 0 {
			s.exitCh <- 0
			break
		}
	}
}
