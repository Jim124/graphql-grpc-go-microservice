package account

import (
	"context"
	"fmt"
	"net"

	pb "github.com/jim124/graphql-grpc-go-microservice/account/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	pb.UnimplementedAccountServiceServer
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	service := grpc.NewServer()
	pb.RegisterAccountServiceServer(service, &grpcServer{s, pb.UnimplementedAccountServiceServer{}})
	reflection.Register(service)
	return service.Serve(lis)
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	res, err := s.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	accounts := []*pb.Account{}
	for _, a := range res {
		accounts = append(accounts, &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		})
	}
	return &pb.GetAccountsResponse{
		Accounts: accounts,
	}, nil
}

func (s *grpcServer) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	a, err := s.service.PostAccount(ctx, r.Name)
	if err != nil {
		return nil, err
	}
	return &pb.PostAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}
