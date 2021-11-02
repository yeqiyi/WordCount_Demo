package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"wordcount/FileReader"
)

var(
	res map[string]int
	wg sync.WaitGroup
	conns chan worker
	task_str chan string
)

type Ctor struct {
	mux sync.Mutex
	lisener net.Listener
}

type worker struct{
	port string
	conn net.Conn
}

//不断轮询能够工作的Worker
func (c *Ctor)RegistWorkers(){
	for{
		conn,err:=c.lisener.Accept()
		buf:=make([]byte,512)
		len,err:=conn.Read(buf)
		if err!=nil{
			fmt.Println("error reading",err)
		}
		fmt.Println("worker port:",string(buf[:len]))
		conns<-worker{port: string(buf[:len]),conn: conn}
		//conn.Write([]byte("reply from ctor"))
	}
}

func (c *Ctor)callWorker(addr string){
	defer wg.Done()
	defer func(){
		if err:=recover();err!=nil{
			fmt.Println(err)
		}
	}()
	ctorcall,err:=rpc.Dial("tcp",":"+addr)
	if err!=nil{
		log.Println(err)
		//fmt.Printf("Worker %s shut down\n",addr)
		//panic(err)
	}

	reply:=make(map[string]int)
	fmt.Printf("worker(%s) is counting...\n",addr)
	str:=<-task_str
	err=ctorcall.Call("Worker.DoTask",str,&reply)
	if err!=nil{
		task_str<-str
		//panic(err)
		//fmt.Println(err)
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

	conns<-worker{port: addr}
}

func main(){
	ctor:=new(Ctor)
	res=make(map[string]int)
	conns=make(chan worker,3)
	task_str=make(chan string,3)
	defer close(conns)
	//ports:=[]string{"3001","3002","3003"}

	l,err:=net.Listen("tcp",":3000")
	if err!=nil{
		fmt.Println("listen error:",err)
		return
	}
	fmt.Println("Contractor is listening at 3000")
	ctor.lisener=l
	go ctor.RegistWorkers()
	split_str,err:=FileReader.FileSplit("WordCount.txt",3)
	for i:=0;i<len(split_str);i++{
		task_str<-split_str[i]
	}
	if err!=nil{
		fmt.Println(err)
		return
	}

	for w:=range(conns){
		//任务完成 退出
		fmt.Println("task_str len:",len(task_str))
		if len(task_str)==0{
			break
		}
		//fmt.Println("port = ",port)
		wg.Add(1)
		fmt.Println(w.port)
		go ctor.callWorker(w.port)
	}
	/*
	for i:=0;i<3;i++{
		wg.Add(1)
		fmt.Println(split_str[i]+"end")
		go ctor.callWorker(ports[i],split_str[i])
	}
	 */
	wg.Wait()

	file,err:=os.Create("result.txt")
	defer file.Close()
	for k,v :=range(res){
		str:="word:"+k +" value:"+ strconv.Itoa(v) +"\n"
		io.WriteString(file,str)
		fmt.Printf("key:%s value:%d\n",k,v)
	}
}
