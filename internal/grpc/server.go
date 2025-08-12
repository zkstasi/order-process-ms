package grpc

import (
	"context"
	"google.golang.org/grpc"

	"order-ms/internal/model"
	"order-ms/internal/repository"
	pb "order-ms/pkg/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// NewGrpcServer создаёт gRPC сервер и регистрирует на нём User и Order сервисы
func NewGrpcServer() *grpc.Server {
	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, NewUserServer())
	pb.RegisterOrderServiceServer(s, NewOrderServer())
	return s
}

// Helpers: функции, которые переводят внутренние модели пользователя и заказа (model.User/Order) в protobuf-сообщение (pb.User/Order), которое отправляется по gRPC
func toProtoUser(u *model.User) *pb.User {
	if u == nil {
		return nil
	}
	return &pb.User{
		Id:   u.Id,
		Name: u.Name,
	}
}

func toProtoOrder(o *model.Order) *pb.Order {
	if o == nil {
		return nil
	}
	return &pb.Order{
		Id:     o.Id,
		UserId: o.UserID,
		Status: pb.OrderStatus(int32(o.Status)),
	}
}

// User service:

// структура, которая реализует интерфейс gRPC-сервиса UserService
// встраивается UnimplementedUserServiceServer, чтобы не реализовывать все методы сразу
type UserServer struct {
	pb.UnimplementedUserServiceServer
}

// Конструктор, возвращающий новый сервер для UserService
func NewUserServer() pb.UserServiceServer {
	return &UserServer{}
}

// Методы UserServer

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req == nil || req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	u := model.NewUser(req.GetName())
	repository.SaveStorable(u)
	return &pb.CreateUserResponse{User: toProtoUser(u)}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	u := repository.GetUserByID(req.GetId())
	if u == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return toProtoUser(u), nil
}

func (s *UserServer) ListUsers(ctx context.Context, _ *emptypb.Empty) (*pb.ListUsersResponse, error) {
	users := repository.GetUsers()
	out := &pb.ListUsersResponse{}
	for _, u := range users {
		out.Users = append(out.Users, toProtoUser(u))
	}
	return out, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	if req == nil || req.GetId() == "" || req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "id and name required")
	}
	ok := repository.UpdateUserName(req.GetId(), req.GetName())
	if !ok {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	updated := repository.GetUserByID(req.GetId())
	return toProtoUser(updated), nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	ok := repository.DeleteUser(req.GetId())
	if !ok {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &emptypb.Empty{}, nil
}

// Order service

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
}

// Конструктор, возвращающий новый сервер для OrderService
func NewOrderServer() pb.OrderServiceServer {
	return &OrderServer{}
}

// Методы OrderServer

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	if req == nil || req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	o := model.NewOrder(req.GetUserId())
	repository.SaveStorable(o)
	return &pb.CreateOrderResponse{Order: toProtoOrder(o)}, nil
}

func (s *OrderServer) ListOrders(ctx context.Context, _ *emptypb.Empty) (*pb.ListOrdersResponse, error) {
	ords := repository.GetOrders()
	out := &pb.ListOrdersResponse{}
	for _, o := range ords {
		out.Orders = append(out.Orders, toProtoOrder(o))
	}
	return out, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	o := repository.GetOrderByID(req.GetId())
	if o == nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}
	return toProtoOrder(o), nil
}

func (s *OrderServer) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	ok := repository.DeleteOrder(req.GetId())
	if !ok {
		return nil, status.Error(codes.NotFound, "order not found")
	}
	return &emptypb.Empty{}, nil
}

func (s *OrderServer) ConfirmOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	o := repository.GetOrderByID(req.GetId())
	if o == nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}
	ok := repository.ConfirmOrder(req.GetId())
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "cannot confirm order")
	}
	updated := repository.GetOrderByID(req.GetId())
	return toProtoOrder(updated), nil
}

func (s *OrderServer) DeliverOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	o := repository.GetOrderByID(req.GetId())
	if o == nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}
	ok := repository.DeliveredOrder(req.GetId())
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "cannot deliver order")
	}
	updated := repository.GetOrderByID(req.GetId())
	return toProtoOrder(updated), nil
}

func (s *OrderServer) CancelOrder(ctx context.Context, req *pb.GetOrderRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	ok := repository.CancelOrder(req.GetId())
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "cannot cancel order")
	}
	return &emptypb.Empty{}, nil
}
