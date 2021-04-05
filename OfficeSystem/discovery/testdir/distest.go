package main

import (
	"fmt"
	"sync"
	"time"
)

import "OfficeSystem/discovery"

func main()  {
	endpoints :=[]string{"localhost:2379"}
	master,err := discovery.NewMaster(endpoints,"/")
	if err !=nil {
		fmt.Println("err--------------------",err)
	}
	time.Sleep(3*time.Second)
	var wg sync.WaitGroup
	wg.Add(2)
	wg.Wait()
	fmt.Println("suc------------------",len(master.Nodes))
}

