package random_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/zerogo-hub/zero-helper/random"
	"github.com/zerogo-hub/zero-helper/time"
)

func TestSnowflake(t *testing.T) {
	// 多个协程并发生成不重复的 uuid

	// 节点数量
	workerNum := 5
	// 每一个节点生成的数量
	uuidNum := 10000

	var wg sync.WaitGroup
	wg.Add(workerNum)

	ch := make(chan uint64, workerNum*uuidNum+1)

	for i := 0; i < workerNum; i++ {
		go func(workderID int) {
			defer wg.Done()

			snowflake, _ := random.NewSnowflake(workderID)

			for j := 0; j < uuidNum; j++ {
				uuid, _ := snowflake.UUID()
				ch <- uuid
			}

		}(i)
	}

	wg.Wait()

	// 检查是否有重复
	m := make(map[uint64]int, workerNum*uuidNum)

	for i := 0; i < workerNum*uuidNum; i++ {
		uuid := <-ch
		if _, exist := m[uuid]; exist {
			t.Fatalf("%d repeated", uuid)
		}

		m[uuid] = 0
	}
}

func TestSnowflakeBy(t *testing.T) {
	originTime := time.MS()
	snowflake, _ := random.NewSnowflakeBy(1, originTime, 10, 2, nil, nil)
	uuid, _ := snowflake.UnsafeUUID()
	t.Log(uuid)
}

// TestSnowflakeTimeback 测试时间回拨
func TestSnowflakeTimeback(t *testing.T) {
	workerIDs := []int{1, 2, 3, 4, 5}

	nextWorkerIDFunc := func() (int, error) {
		if len(workerIDs) == 0 {
			return 0, errors.New("empty workerIDs")
		}

		newWorkerID := workerIDs[0]
		workerIDs = workerIDs[1:]
		return newWorkerID, nil
	}

	backWorkerIDFunc := func(workerID int) error {
		workerIDs = append(workerIDs, workerID)
		return nil
	}

	workerID, _ := nextWorkerIDFunc()
	snowflake, _ := random.NewSnowflake(workerID)
	snowflake.SetWorkIDFunc(nextWorkerIDFunc, backWorkerIDFunc)

	// 正常情况
	v1, _ := snowflake.UUID()
	t.Logf("v1: %d", v1)

	// 时间回退
	random.SetSnowflakeTestTimebackward(true)
	v2, _ := snowflake.UUID()
	t.Logf("v2: %d", v2)
}
