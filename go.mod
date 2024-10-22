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

require google.golang.org/grpc v1.67.1

require (
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
