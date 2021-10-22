package FileReader

import (
	"bufio"
	"io"
	"os"
)

func FileSplit(filepath string,n int) ([]string,error){
	file_split:=make([]string,n)
	file,err:=os.Open(filepath)
	defer file.Close()
	if err!=nil{
		return nil,err
	}
	reader:=bufio.NewReader(file)
	idx:=0
	for{
		str,err:=reader.ReadString(' ')
		if err==nil {
				file_split[idx]+=str
				idx=(idx+1)%n
		}else if err==io.EOF{
			file_split[idx]+=str
			break
		}else{
			return nil,err
		}
	}
	return file_split,nil
}
