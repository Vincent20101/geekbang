# Week13 毕业项目

利用go-micro框架设计的电商微服务系统



user： 用户服务

product : 商品服务

order: 订单服务

cart : 购物车服务

cartapi: 购物车网关服务

payment:  支付服务

paymentapi: 支付网关服务



由于时间问题，目前完善的微服务可以查看： 购物车服务跟支付服务，其他微服务设计都跟这两个一样，就不重复造轮子了。

Ps.  完整的一个微服务请可参考： cart 与 cartapi网关



~~~go
func main() {
	//注册中心
	consul := consul2.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	//链路追踪
	t,io,err := common.NewTracer("go.micro.api.cartApi","localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	//启动端口
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0","9096"),hystrixStreamHandler)
		if err !=nil {
			log.Error(err)
		}
	}()


	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.cartApi"),
		micro.Version("latest"),
		micro.Address("0.0.0.0:8086"),
		//添加 consul 注册中心
		micro.Registry(consul),
		//添加链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		//添加熔断
		micro.WrapClient(NewClientHystrixWrapper()),
		//添加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	// Initialise service
	service.Init()

	cartService:=go_micro_service_cart.NewCartService("go.micro.service.cart",service.Client())

	cartService.AddCart(context.TODO(),&go_micro_service_cart.CartInfo{

		UserId:    3,
		ProductId: 4,
		SizeId:    5,
		Num:       5,
	})

	// Register Handler
	if err := cartApi.RegisterCartApiHandler(service.Server(), &handler.CartApi{CartService:cartService});err !=nil {
		log.Error(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
~~~

