package main

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"wordcount/FileReader"
)

var(
	res map[string]int
	wg sync.WaitGroup
)

type Ctor struct {
	mux sync.Mutex
}

func (c *Ctor)callWorker(port string,str string){
	defer wg.Done()
	ctorcall,err:=rpc.Dial("tcp","localhost:"+port)
	if err!=nil{
		log.Fatal("dialing:",err)
	}

	reply:=make(map[string]int)

	err=ctorcall.Call("Worker.DoTask",str,&reply)
	fmt.Printf("worker(%s) is counting...\n",port)
	if err!=nil{
		fmt.Println(err)
	}

	c.mux.Lock()

	for k,v:=range(reply){
		fmt.Printf("reply key:%s value:%d\n",k,v)
		if value,ok:=res[k];ok{
			res[k]=value+v
		}else{
			res[k]=v
		}
	}

	c.mux.Unlock()
}

func main(){
	ctor:=new(Ctor)
	res=make(map[string]int)
	ports:=[]string{"3001","3002","3003"}
	split_str,err:=FileReader.FileSplit("WordCount.txt",3)
	if err!=nil{
		fmt.Println(err)
		return
	}
	for i:=0;i<3;i++{
		wg.Add(1)
		fmt.Println(split_str[i]+"end")
		go ctor.callWorker(ports[i],split_str[i])
	}
	wg.Wait()

	file,err:=os.Create("result.txt")
	defer file.Close()

	for k,v :=range(res){
		str:="word:"+k +" value:"+ strconv.Itoa(v) +"\n"
		io.WriteString(file,str)
		fmt.Printf("key:%s value:%d\n",k,v)
	}
}
