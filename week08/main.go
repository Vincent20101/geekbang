package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	redisDb *redis.Client
	infomemory string
	loopNum int = 10000
	wg sync.WaitGroup
)

//正则匹配下used_memory的数值
func GetUsedMemory(redisInfo string) int  {
	reg := regexp.MustCompile(`used_memory:\d+`)
	if reg == nil { //解释失败，返回nil
		fmt.Println("regexp err")
		return 0
	}

	result := reg.Find([]byte(redisInfo))
	//fmt.Println(result)
	res , _ := strconv.Atoi(string(result[12:]))
	return res
}

func GetRedisMemory() int {
	redisInfo := redisDb.Info("memory").String()
	return GetUsedMemory(redisInfo)
}

//指定长度的字符转
func ZdyString(n int) string  {
	str := "a"
	return strings.Repeat(str,n)
}

func BatchInserData(bytCount int) {
	val := ZdyString(bytCount)
	i := 0
	for i < loopNum {
		key := "test_" + strconv.Itoa(i)
		redisDb.Set(key,val,0)
		i++
	}
}

func show(n int) {
	redisDb.FlushDB()
	startMemory := GetRedisMemory()
	BatchInserData(n)
	endMemory := GetRedisMemory()
	fmt.Println("内存差：",endMemory - startMemory)
	fmt.Printf("%.10f",(endMemory - startMemory) / loopNum)
	fmt.Println()
}

func main() {
	redisDb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redisDb.Close()
	//fmt.Printf("%.2f",123.123)

	pong, err := redisDb.Ping().Result()
	log.Println(pong, err) // Output: PONG <nil>

	fmt.Println("10字节")
	show(10)

	fmt.Println("20字节")
	show(20)

	fmt.Println("50字节")
	show(50)

	fmt.Println("100字节")
	show(100)

	fmt.Println("200字节")
	show(200)

	fmt.Println("1k字节")
	show(1024)

	fmt.Println("5k字节")
	show(1024 * 5)
}
