package httpAuth

import (
	"fmt"
	sso_v1_ssov1 "github.com/Senkoker/sso_proto/proto/proto_go/protobufcontract/protobufcontract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

type AuthGrpcConnect struct {
	NewAuthClient sso_v1_ssov1.AuthClient
}

func NewAuthClient(host string, port int) *AuthGrpcConnect {
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	cc, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	newClient := sso_v1_ssov1.NewAuthClient(cc)
	return &AuthGrpcConnect{NewAuthClient: newClient}
}
