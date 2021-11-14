package main

import (
	"geekbang/week04/internal/biz"
	"geekbang/week04/internal/data"
	"geekbang/week04/internal/server"
	"geekbang/week04/internal/service"
	"github.com/google/wire"
)

func initApp() {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet))
}
