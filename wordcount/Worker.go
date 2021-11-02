package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strings"
	"time"
)

type Worker struct {
	addr string
	port string
	conn net.Conn
}

var(
	wport = flag.String("port","","setup port")
)

var mp map[string]int

func (w *Worker)DoTask(request string,reply *map[string]int) error{
	mp=make(map[string]int)
	fmt.Println(request)
	str_slice:=strings.Fields(request)
	for i:=0;i<len(str_slice);i++{
		if v,ok:=mp[str_slice[i]];ok{
			time.Sleep(500*time.Millisecond)
			mp[str_slice[i]]=v+1
		}else{
			mp[str_slice[i]]=1
		}
	}
	*reply=mp

	return nil
}

func main(){
	flag.Parse()
	conn,err:=net.Dial("tcp","localhost:3000")
	conn.Write([]byte(*wport))
	fmt.Println(conn.LocalAddr())
	if err!=nil{
		fmt.Println("Error Dialing:",err)
		return
	}
	server:=rpc.NewServer()
	err=server.RegisterName("Worker",&Worker{port: *wport,conn: conn})
	if err!=nil{
		log.Fatal("Format of Worker is not correct:",err)
	}
	//rpc.RegisterName("Worker",Worker{port: *wport})
	fmt.Printf("Worker is ready and listening at %s\n",*wport)
	listener,err:=net.Listen("tcp",":"+*wport)
	fmt.Println(listener.Addr())
	if err!=nil{
		log.Fatal("ListeningTCP error:",err)
	}

	server.Accept(listener)
	if err!=nil{
		log.Fatal("Accept error:",err)
	}


	//rpc.ServeConn(conn)
}
