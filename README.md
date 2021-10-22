# WordCount_Demo

An simple WordCount Demo written in Go
this simple system consists of one Contractor and three Workers，Contractor use RPC to call the WordCount service offered by Worker.

[中文说明](https://github.com/yeqiyi/WordCount_Demo/wiki/WordCount-Demo)

# How to start

## Start Worker
```go run Worker.go -port 3001```
```go run Worker.go -port 3002```
```go run Worker.go -port 3003```

## Start Contractor
```go run Contractor.go```

The result will be saved in ```result.txt```

