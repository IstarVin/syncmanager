package syncmanager

import (
	"math/rand"
	"testing"
	"time"
)

func TestSyncManager(t *testing.T) {
	syncManager := NewSyncManager(2)

	for i := 0; i < 10; i++ {
		randInt := rand.Intn(5)

		syncManager.Add(func(a ...any) {
			i := a[0].(int)
			sec := a[1].(int)

			println(i, " started: ", sec, " seconds")
			time.Sleep(time.Duration(sec) * time.Second)
			println("\t", i, "finished")
		}, i, randInt)
	}

	syncManager.Wait()
}
