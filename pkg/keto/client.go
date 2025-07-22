package keto

import (
	"fmt"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Host      string `env:"HOST"`
	ReadPort  int    `env:"READ_PORT, default=4466"`
	WritePort int    `env:"WRITE_PORT, default=4467"`
}

type Keto struct {
	ReadConn  *grpc.ClientConn
	WriteConn *grpc.ClientConn
	Read      rts.ReadServiceClient
	Check     rts.CheckServiceClient
	Write     rts.WriteServiceClient
}

func NewKeto(cfg Config) (*Keto, error) {
	readConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%v", cfg.Host, cfg.ReadPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("new grpc client: %w", err)
	}

	writeConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%v", cfg.Host, cfg.WritePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
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

func (k *Keto) Close() {
	k.ReadConn.Close()
	k.WriteConn.Close()
}
