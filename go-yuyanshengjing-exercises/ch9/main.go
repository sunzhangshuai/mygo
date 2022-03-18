package ch9

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"

	"exercises/ch8"
)

// Exercises 练习
type Exercises struct {
}

// Task1 给gopl.io/ch9/bank1程序添加一个Withdraw(amount int)取款函数。
// 其返回结果应该要表明事务是成功了还是因为没有足够资金失败了。这条消息会被发送给monitor的goroutine，且消息需要包含取款的额度和一个新的channel，这个新channel会被monitor goroutine来把boolean结果发回给Withdraw。
func (e *Exercises) Task1() {
	wg := sync.WaitGroup{}
	wg.Add(6)

	go func() {
		defer wg.Done()
		Deposit(100)
	}()
	go func() {
		defer wg.Done()
		Deposit(200)
	}()
	go func() {
		defer wg.Done()
		fmt.Println(Balance())
	}()
	go func() {
		defer wg.Done()
		fmt.Println(Withdraw(100))
	}()
	go func() {
		defer wg.Done()
		fmt.Println(Withdraw(300))
	}()
	go func() {
		defer wg.Done()
		fmt.Println(Balance())
	}()

	wg.Wait()
}

// Task2 重写2.6.2节中的PopCount的例子，使用sync.Once，只在第一次需要用到的时候进行初始化。
//（虽然实际上，对PopCount这样很小且高度优化的函数进行同步可能代价没法接受。）
func (e *Exercises) Task2() {
	var pc [256]byte
	once := sync.Once{}
	once.Do(func() {
		for i := range pc {
			pc[i] = pc[i/2] + byte(i&1)
		}
	})
}

// Task3 扩展Func类型和(*Memo).Get方法，支持调用方提供一个可选的done channel，使其具备通过该channel来取消整个操作的能力（§8.9）。
// 一个被取消了的Func的调用结果不应该被缓存。
func (e *Exercises) Task3() {
	fuc := func(key string, done chan struct{}) (interface{}, error) {
		var res int
		var err error

		d := make(chan struct{})

		go func() {
			res, err = strconv.Atoi(key)
			d <- struct{}{}
		}()

		select {
		case _, ok := <-done:
			if !ok {
				return nil, fmt.Errorf("操作被取消了")
			} else {
				return nil, fmt.Errorf("程序有误")
			}
		case <-d:
			break
		}
		return res, err
	}

	memo := NewMemo(fuc)
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 1; i <= 10; i++ {
		wg.Add(4)
		for j := 1; j <= 4; j++ {
			go func(key string) {
				defer wg.Done()
				fmt.Println(memo.Get(key, done))
			}(strconv.Itoa(i))
		}
	}
	close(done)
	wg.Wait()
}

// Task4 创建一个流水线程序，支持用channel连接任意数量的goroutine，在跑爆内存之前，可以创建多少流水线阶段？一个变量通过整个流水线需要用多久？
func (e *Exercises) Task4() {
	type data struct {
		num  int
		time time.Time
	}

	numCount := 100
	processCount := 1000000

	chanList := make([]chan *data, processCount)
	for i := 0; i < processCount; i++ {
		chanList[i] = make(chan *data)
	}

	for i := 1; i <= numCount; i++ {
		go func(i int) {
			dat := &data{
				num:  i,
				time: time.Now(),
			}
			chanList[0] <- dat
		}(i)
		for j := 1; j < processCount; j++ {
			go func(j int) {
				dat := <-chanList[j-1]
				chanList[j] <- dat
			}(j)
		}
	}

	for i := 1; i <= numCount; i++ {
		dat := <-chanList[processCount-1]
		fmt.Println(dat.num, time.Now().Sub(dat.time).String())
	}
}

// Task5 写一个有两个goroutine的程序，两个goroutine会向两个无buffer channel反复地发送ping-pong消息。这样的程序每秒可以支持多少次通信？
func (e *Exercises) Task5() {
	var preNum int
	var num int

	ticker := time.Tick(1 * time.Second)

	ping := make(chan string)
	pong := make(chan string)

	go func() {
		for {
			ping <- "ping"
			<-pong
		}
	}()

	go func() {
		for {
			<-ping
			pong <- "pong"
			num++
		}
	}()

	for i := 0; i < 60; i++ {
		<-ticker
		n := num
		println(n - preNum)
		preNum = n
	}
}

// Task6 测试一下计算密集型的并发程序（练习8.5那样的）会被GOMAXPROCS怎样影响到。在你的电脑上最佳的值是多少？你的电脑CPU有多少个核心？
// go run main.go -ch 9 -task 6 sunchen
func (e *Exercises) Task6() {
	ex := &ch8.Exercises{}
	time1 := time.Now()

	for i := 1; i <= 1000000; i *= 2 {
		runtime.GOMAXPROCS(i)
		ex.Task5()
		time2 := time.Now()
		fmt.Println(i, time2.Sub(time1).String())
		time1 = time2
	}
}
