package service

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "order-ms/pkg/proto"
)

type GrpcClient struct {
	conn        *grpc.ClientConn
	userClient  pb.UserServiceClient
	orderClient pb.OrderServiceClient
}

func NewGrpcClient(addr string) (*GrpcClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcClient{
		conn:        conn,
		userClient:  pb.NewUserServiceClient(conn),
		orderClient: pb.NewOrderServiceClient(conn),
	}, nil
}

func (c *GrpcClient) Close() error {
	return c.conn.Close()
}

// CreateUser
func (c *GrpcClient) CreateUserExample() (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.userClient.CreateUser(ctx, &pb.CreateUserRequest{Name: "Alice"})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

// GetUser
func (c *GrpcClient) GetUserExample(id string) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.userClient.GetUser(ctx, &pb.GetUserRequest{Id: id})
}

// ListUsers
func (c *GrpcClient) ListUsersExample() ([]*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.userClient.ListUsers(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}

// UpdateUser
func (c *GrpcClient) UpdateUserExample(id, name string) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.userClient.UpdateUser(ctx, &pb.UpdateUserRequest{Id: id, Name: name})
}

// DeleteUser
func (c *GrpcClient) DeleteUserExample(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.userClient.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
	return err
}

// CreateOrder
func (c *GrpcClient) CreateOrderExample(userID string) (*pb.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.orderClient.CreateOrder(ctx, &pb.CreateOrderRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return resp.Order, nil
}

// GetOrder
func (c *GrpcClient) GetOrderExample(id string) (*pb.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.orderClient.GetOrder(ctx, &pb.GetOrderRequest{Id: id})
}

// ListOrders
func (c *GrpcClient) ListOrdersExample() ([]*pb.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.orderClient.ListOrders(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.Orders, nil
}

// DeleteOrder
func (c *GrpcClient) DeleteOrderExample(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.orderClient.DeleteOrder(ctx, &pb.DeleteOrderRequest{Id: id})
	return err
}

// ConfirmOrder
func (c *GrpcClient) ConfirmOrderExample(id string) (*pb.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.orderClient.ConfirmOrder(ctx, &pb.GetOrderRequest{Id: id})
}

// DeliverOrder
func (c *GrpcClient) DeliverOrderExample(id string) (*pb.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.orderClient.DeliverOrder(ctx, &pb.GetOrderRequest{Id: id})
}

// CancelOrder
func (c *GrpcClient) CancelOrderExample(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.orderClient.CancelOrder(ctx, &pb.GetOrderRequest{Id: id})
	return err
}
