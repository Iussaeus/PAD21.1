package broker

type Broker interface {
	Open(port string)	
	Serve()
}
