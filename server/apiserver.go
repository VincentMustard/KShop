package main

type ApiServer struct {
}

func NewApiServer() *ApiServer {
	return &ApiServer{}
}

func (a *ApiServer) Publish() {}
