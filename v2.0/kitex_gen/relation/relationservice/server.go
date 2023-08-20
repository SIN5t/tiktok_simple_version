// Code generated by Kitex v0.7.0. DO NOT EDIT.
package relationservice

import (
	server "github.com/cloudwego/kitex/server"
	relation "tiktok_v2/kitex_gen/relation"
)

// NewServer creates a server.Server with the given service and options.
func NewServer(handler relation.RelationService, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}
