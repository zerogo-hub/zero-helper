package random_test

import (
	"errors"
	"sync"
	"testing"

	zerorandom "github.com/zerogo-hub/zero-helper/random"
	zerotime "github.com/zerogo-hub/zero-helper/time"
)

func TestSnowflake(t *testing.T) {
	// 多个协程并发生成不重复的 uuid

	// 节点数量
	workerNum := 5
	// 每一个节点生成的数量
	uuidNum := 1000

	var wg sync.WaitGroup
	wg.Add(workerNum)

	ch := make(chan uint64, workerNum*uuidNum+1)

	for i := 0; i < workerNum; i++ {
		go func(workderID int) {
			defer wg.Done()

			snowflake, _ := zerorandom.NewSnowflake(workderID)

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
	originTime := zerotime.MS()
	snowflake, _ := zerorandom.NewSnowflakeBy(1, originTime, 10, 2, nil, nil)
	uuid, _ := snowflake.UnsafeUUID()
	t.Log(uuid)

	if _, err := zerorandom.NewSnowflakeBy(-1, originTime, 10, 2, nil, nil); err == nil {
		t.Error("test invalid workID failed")
	}
}

func TestShortSnowflake(t *testing.T) {
	snowflake, _ := zerorandom.NewBit46Snowflake(1)
	uuids := make(map[uint64]struct{}, 100)

	for i := 0; i < 100; i++ {
		id, err := snowflake.UnsafeUUID()
		if err != nil {
			t.Fatal(err.Error())
		}
		if _, exist := uuids[id]; exist {
			t.Fatalf("%d repeated", id)
		}
		uuids[id] = struct{}{}
	}

	if _, err := zerorandom.NewBit46Snowflake(-1); err == nil {
		t.Error("test invalid workID failed")
	}
}

func TestShortSnowflakeRepeat(t *testing.T) {
	snowflake0, _ := zerorandom.NewBit46Snowflake(0)
	snowflake1, _ := zerorandom.NewBit46Snowflake(1)

	// 检查 1000 个中是否有重复生成
	n := 1000
	ids1 := make([]uint64, 0, n*2)
	ids2 := make([]uint64, 0, n)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		for i := 0; i < n; i++ {
			id, err := snowflake0.UUID()
			if err != nil {
				t.Error(err.Error())
				continue
			}
			ids1 = append(ids1, id)
		}
	}()
	go func() {
		defer wg.Done()

		for i := 0; i < n; i++ {
			id, err := snowflake1.UUID()
			if err != nil {
				t.Error(err.Error())
				continue
			}
			ids2 = append(ids2, id)
		}
	}()

	wg.Wait()

	ids1 = append(ids1, ids2...)

	existIDs := make(map[uint64]struct{}, len(ids1))

	for _, id := range ids1 {
		if _, ok := existIDs[id]; ok {
			t.Fatal("Duplicate id")
		}
		existIDs[id] = struct{}{}
	}
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
	snowflake, _ := zerorandom.NewSnowflake(workerID)
	snowflake.SetWorkIDFunc(nextWorkerIDFunc, backWorkerIDFunc)

	// 正常情况
	v1, _ := snowflake.UUID()
	t.Logf("v1: %d", v1)

	// 时间回退
	zerorandom.SetSnowflakeTestTimebackward(true)
	v2, _ := snowflake.UUID()
	t.Logf("v2: %d", v2)
}
