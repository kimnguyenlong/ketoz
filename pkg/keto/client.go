package keto

import (
	"fmt"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Keto struct {
	ReadConn  *grpc.ClientConn
	WriteConn *grpc.ClientConn
	Read      rts.ReadServiceClient
	Check     rts.CheckServiceClient
	Write     rts.WriteServiceClient
}

func NewKeto() (*Keto, error) {
	readConn, err := grpc.NewClient("127.0.0.1:4466", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("new grpc client: %w", err)
	}

	writeConn, err := grpc.NewClient("127.0.0.1:4467", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("new grpc client: %w", err)
	}

	return &Keto{
		ReadConn:  readConn,
		WriteConn: writeConn,
		Read:      rts.NewReadServiceClient(readConn),
		Check:     rts.NewCheckServiceClient(readConn),
		Write:     rts.NewWriteServiceClient(writeConn),
	}, nil
}
