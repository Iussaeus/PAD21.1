module pad

go 1.23.0

replace broker => ./broker/broker

replace client => ./client/client

replace broker_tcp => ./broker_tcp/broker_tcp

replace client_tcp => ./client_tcp/client_tcp

replace message_tcp => ./message_tcp/message_tcp

replace broker_grpc => ./broker_grpc/broker_grpc

replace client_grpc => ./client_grpc/client_grpc

replace helpers => ./helpers/helpers

replace proxy => ./proxy/proxy

replace dataWarehouse => ./data_warehouse/data_warehouse

require google.golang.org/grpc v1.67.1

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

require (
	github.com/lib/pq v1.10.9
	github.com/redis/go-redis/v9 v9.7.0
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
