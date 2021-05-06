package channel

import (
	"fmt"
	"math/rand"
	"time"
)

// Channel的本质是一个队列；
// Channel是线程安全的，也就是自带锁定功能；
// Channel和切片/字典一样，必须创建后才能使用，否则会报错；
// Channel是引用类型，地址传递；
// Channel必须关闭后才能遍历；
// Channel再go协程中可以发生管道阻塞；
// 双向管道：var myChan chan int
// 单向只读管道：var myChan <-chan int
// 单向只写管道：var myChan chan<- int
// 双向管道可以转化为单向管道，单向管道不可以转化为双向管道

// Channel 管道
func Channel() {
	/*
	   1.什么是管道:
	   管道就是一个队列, 具备先进先出的原则
	   是线程安全的, 也就是自带锁定功能

	   2.管道作用:
	   在Go语言的协程中, 一般都使用管道来保证多个协程的同步, 或者多个协程之间的通讯

	   3.如何声明一个管道, 和如何创建一个管道
	   管道在Go语言中和切片/字典一样也是一种数据类型
	   管道和切片/字典非常相似, 都可以用来存储数据, 都需要make之后才能使用
	   3.1管道声明格式:
	   var 变量名称 chan 数据类型
	   var myCh chan int
	   如上代码的含义: 声明一个名称叫做myCh的管道变量, 管道中可以存储int类型的数据
	   3.2管道的创建:
	   make(chan 数据类型, 容量)
	   myCh = make(chan int, 3);
	   路上代码的含义: 创建一个容量为3, 并且可以保存int类型数据的管道

	   4.管道的使用
	   4.1如何往管道中存储(写入)数据?
	   myCh<-被写入的数据
	   4.2如何从管道中获取(读取)数据?
	   <-myCh
	   对管道的操作是IO操作
	   例如: 过去的往文件中写入或者读取数据, 也是IO操作
	   例如: 过去的往屏幕上输出内容, 或者从屏幕获取内容, 也是IO操作
	   stdin / stdout / stderr

	   注意点:
	   和切片不同, 在切片中make函数的第二个参数表示的切片的长度(已经存储了多少个数据),
	   而第三个参数才是指定切片的容量
	   但是在管道中, make函数的第二个参数就是指定管道的容量, 默认长度就是0
	*/

	//1.定义一个管道
	//var myChan chan int
	//2.使用make创建管道
	myChan := make(chan int, 3)
	//3.往管道中存储数据
	myChan <- 1
	myChan <- 2
	myChan <- 3
	//从管道中取出数据
	fmt.Println(<-myChan)
	fmt.Println(<-myChan)
	fmt.Println(<-myChan)

	//定义一个管道
	// var myChan chan int
	//直接使用管道
	//注意点: 会报错,管道定义完成后不创建是无法直接使用的
	//myChan<-666
	//fmt.Println(<-myChan)

	//创建管道
	// myChan = make(chan int, 3)
	//只要往管道中写入了数据, 那么len就会增加
	myChan <- 2
	fmt.Println("len = ", len(myChan), "cap = ", cap(myChan))
	myChan <- 4
	fmt.Println("len = ", len(myChan), "cap = ", cap(myChan))
	myChan <- 6
	fmt.Println("len = ", len(myChan), "cap = ", cap(myChan))
	//注意点: 如果len等于cap, 那么就不能往管道中再写入数据了, 否则会报错
	// myChan <- 8

	//管道未写入数据,使用管道去取数据会报错
	//从管道中取数据,len会减少
	//<-myChan
	fmt.Println(<-myChan)
	fmt.Println("len=", len(myChan), "cap = ", cap(myChan))
	fmt.Println(<-myChan)
	fmt.Println("len=", len(myChan), "cap = ", cap(myChan))
	fmt.Println(<-myChan)
	fmt.Println("len=", len(myChan), "cap = ", cap(myChan))
	//注意点: 取数据个数也不可以超出写入的数据个数,否则会报错
	//fmt.Println(<-myChan)
}

// Loop 循环遍历管道的两种方式，for or for...range
func Loop() {
	myChan := make(chan int, 3)
	myChan <- 1
	myChan <- 2
	myChan <- 3
	close(myChan)
	for v := range myChan {
		fmt.Println(v)
	}
	myChan = make(chan int, 3)
	myChan <- 4
	myChan <- 5
	myChan <- 6
	close(myChan)
	for {
		if v, ok := <-myChan; ok {
			fmt.Println(v)
			fmt.Println(ok)
		} else {
			fmt.Println(ok)
			break
		}
	}
}

// Block 管道阻塞
func Block() {
	myChan := make(chan int, 2)
	myChan <- 1
	myChan <- 2
	fmt.Println("管道满之前的代码")
	//这里不会报错, 会阻塞, 等到将管道中的数据读出去之后, 有的新的空间再往管道中写
	myChan <- 3
	fmt.Println("管道满之后的代码")
	//这里不会报错,会阻塞,等到有人往管道中写入数据之后,有新的数据之后才会读取
	fmt.Println("读取之前的代码")
	<-myChan
	fmt.Println("读取之后的代码")
}

// ConcurrentSerial 并发串行
func ConcurrentSerial() {
	myChan := make(chan bool, 1)
	go func() {
		fmt.Printf("Hello ")
		time.Sleep(5 * time.Second)
		myChan <- true
	}()
	go func() {
		<-myChan
		fmt.Println("World!")
	}()
}

// ConcurrentSerial2 并发串行
func ConcurrentSerial2() {
	out := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			out <- rand.Intn(5)
		}
		close(out)
	}()
	go func() {
		for i := range out {
			fmt.Println(i)
		}
	}()
}

// ProducerAndConsumer 生产者/消费者模型
func ProducerAndConsumer() {
	myChan := make(chan int, 10)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("Producer 1: %d\n", i)
			myChan <- i
			time.Sleep(time.Second)
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("Producer 2: %d\n", i)
			myChan <- i
			time.Sleep(time.Second)
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("Consumer 1: %d\n", <-myChan)
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("Consumer 2: %d\n", <-myChan)
		}
	}()
}
