package syncmanager

func NewSyncManager(concurrent int) *SyncManager {
	syncManager := new(SyncManager)

	syncManager.concurrent = concurrent
	syncManager.addQueue = make(chan struct{})
	syncManager.doneQueue = make(chan struct{})
	syncManager.exitCh = make(chan struct{})

	go syncManager.daemon()

	return syncManager
}

type SyncManager struct {
	QueueList []*Queue

	addQueue  chan struct{}
	doneQueue chan struct{}
	exitCh    chan struct{}

	concurrent int
	isWaiting  bool
}

func (s *SyncManager) Add(f func(...any), args ...any) {
	s.QueueList = append(s.QueueList, &Queue{
		Func: f,
		Args: args,
	})
	s.addQueue <- struct{}{}
}

func (s *SyncManager) Wait() {
	s.isWaiting = true
	<-s.exitCh
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
				s.doneQueue <- struct{}{}
			}()

			running++
		}

		if s.isWaiting && running == 0 {
			close(s.exitCh)
			break
		}
	}
}
