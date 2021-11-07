package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
)

//http start
func StartHttpServer(serv *http.Server) error {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "Hello world")
	})
	err := serv.ListenAndServe()
	if err == nil {
		fmt.Println("http server starting")
	}
	fmt.Println("http server starting")
	//return errors.New("测试一下")
	return err
}

func main() {

	//root
	ctx := context.Background()
	//新的cancel，用于控制下文，获取系统的signal后，连带停止errgroup
	ctx, cancel := context.WithCancel(ctx)

	//WithContext 就是使用 WithCancel 创建一个可以取消的 context 将 cancel 赋值给 Group 保存起来，然后再将 context 返回回去
	//注意这里有一个坑，在后面的代码中不要把这个 ctx 当做父 context 又传给下游，因为 errgroup 取消了，这个 context 就没用了，会导致下游复用的时候出错
	g, errCtx := errgroup.WithContext(ctx)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "Hello world")
	})

	//模拟URL请求停止服务，有点多余，写着玩
	servStop := make(chan struct{})
	mux.HandleFunc("/stop", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("请求停止服务")
		fmt.Fprintf(writer, "stop")
		servStop <- struct{}{}
	})

	serv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	//启动服务器
	g.Go(func() error {
		return serv.ListenAndServe()
		//return StartHttpServer(serv)
	})

	//监听errCtx.Done 当Context 被 canceled 或是 times out 的时候，Done 返回一个被 closed 的channel
	g.Go(func() error {
		//这里监听errCtx和ctx
		select {
		case <-errCtx.Done():
			fmt.Println("errgroup exit")
		case <-servStop:
			fmt.Println("url request exit")
		}
		return serv.Shutdown(ctx)
	})

	//监听linux信号
	sign := make(chan os.Signal)
	signal.Notify(sign)
	g.Go(func() error {
		select {
		case <-errCtx.Done():
			return errCtx.Err()
		case s := <-sign:
			fmt.Printf("os sign : %v\n", s)
			cancel()
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Println("group error: ", err)
		return
	}
	fmt.Println("all group done!")
}
