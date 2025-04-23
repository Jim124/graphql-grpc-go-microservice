package account

import (
	"context"
	"fmt"
	"net"

	"github.com/jim124/graphql-grpc-go-microservice/account/protobuf/github.com/graphql-grpc-go-microservice/account/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	protobuf.UnimplementedAccountServiceServer
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	service := grpc.NewServer()
	protobuf.RegisterAccountServiceServer(service, &grpcServer{s, protobuf.UnimplementedAccountServiceServer{}})
	reflection.Register(service)
	return service.Serve(lis)
}

func (s *grpcServer) GetAccount(ctx context.Context, r *protobuf.GetAccountRequest) (*protobuf.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &protobuf.GetAccountResponse{
		Account: &protobuf.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *protobuf.GetAccountsRequest) (*protobuf.GetAccountsResponse, error) {
	res, err := s.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	accounts := []*protobuf.Account{}
	for _, a := range res {
		accounts = append(accounts, &protobuf.Account{
			Id:   a.ID,
			Name: a.Name,
		})
	}
	return &protobuf.GetAccountsResponse{
		Accounts: accounts,
	}, nil
}

func (s *grpcServer) PostAccount(ctx context.Context, r *protobuf.PostAccountRequest) (*protobuf.PostAccountResponse, error) {
	a, err := s.service.PostAccount(ctx, r.Name)
	if err != nil {
		return nil, err
	}
	return &protobuf.PostAccountResponse{
		Account: &protobuf.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}
