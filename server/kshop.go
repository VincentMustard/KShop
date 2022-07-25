package main

type KShop struct {
	Consumer  *Consumer
	ApiServer *ApiServer
}

func NewKShop(
	consumer *Consumer,
	apiServer *ApiServer,
) *KShop {
	return &KShop{
		Consumer:  consumer,
		ApiServer: apiServer,
	}
}

func (k *KShop) KStart(SubChannel string) {
	go k.Consumer.Consume(SubChannel)
	go k.ApiServer.Publish()
}
