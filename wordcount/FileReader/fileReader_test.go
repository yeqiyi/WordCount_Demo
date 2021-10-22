package FileReader

import (
	"fmt"
	"testing"
)


func TestReadFile(t *testing.T){
	filesplit:=make([]string,3)
	filesplit,err:=FileSplit("./WordCount.txt",3)
	if err!=nil{
		t.Error("errors!")
	}
	fmt.Println(len(filesplit))
	for _,str:=range(filesplit){
		fmt.Println(str)
	}
}
