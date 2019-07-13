package main

import (
	"goStorm/Logger"
	"goStorm/Register"
	"runtime"
	"time"
)

var logger = Logger.GetLogger("goStorm")

func main() {

	//num := 1
	//// debug.SetMaxThreads(num + 1000) //设置最大线程数
	//// 注册工作池，传入任务
	//// 参数1 worker并发个数
	//p := Concurrence.NewWorkerPool(num)
	//p.Run()
	//datanum := 1000
	//go func() {
	//	for i := 1; i <= datanum; i++ {
	//		nowTime := time.Now()
	//		nowTimestamp := strconv.FormatInt(nowTime.UnixNano(), 10)
	//		req := &blackHoleReq{
	//			ReqCode:      nowTimestamp,
	//			ReqStartTime: nowTime.UnixNano(),
	//		}
	//		p.JobQueue <- req
	//	}
	//}()
	//for {
	//	logger.Info("runtime.NumGoroutine() :", runtime.NumGoroutine())
	//	time.Sleep(1 * time.Second)
	//}

	// master 测试
	master := Register.Master{
		AllSlaves:make(map[string]*Register.SlaveInfo),
	}
	r :=   Register.SetupRouter(&master)
	// Listen and Server in 0.0.0.0:9526
	go func() {
		err := r.Run(":9526")
		if err != nil{
			logger.Fatalf("服务器异常退出! err=%v\n", err)
		}
	}()

	slave := Register.Slave{
		RegistInfo: &Register.RegistInfo{},
		SlaveInfo: &Register.SlaveInfo{},
	}
	go func() {
		// 延迟注册 等待master启动
		time.Sleep(10 * time.Second)
		slave.RegToMaster()
		for  {
			slave.SendHb()
			// 每5s发送一次心跳
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		logger.Info("runtime.NumGoroutine() :", runtime.NumGoroutine())
		time.Sleep(1 * time.Second)
	}


}


