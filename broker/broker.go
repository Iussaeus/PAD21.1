package broker

type Broker interface {
	Init(port string)	
	Serve()
}
