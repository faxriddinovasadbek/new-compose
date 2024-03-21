package service

import (
	"context"
	l "user-service/pkg/logger"
	pb "user-service/protos/template-service"
	"user-service/storage"

	"github.com/jmoiron/sqlx"
)

//UserService ...
type UserService struct {
	storage storage.IStorage
	logger  l.Logger
	pb.UnimplementedUserServiceServer
}

//NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := s.storage.User().Create(req)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s UserService) GetUser(ctx context.Context, req *pb.GetRequest) (*pb.User, error) {
	user, err := s.storage.User().GetUser(req.UserId)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s UserService) DeleteUser(ctx context.Context, req *pb.GetRequest) (*pb.User, error) {
	user, err := s.storage.User().DeleteUser(req.UserId)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s UserService) UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := s.storage.User().UpdateUser(req)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return user, nil

}
func (s UserService) GetAllUsers(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	users, err := s.storage.User().GetAllUsers(req)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return users, nil
}

func (s UserService) CheckUniques(ctx context.Context, req *pb.CheckUniquesRequest) (*pb.CheckUniquesResponse, error){
	tf, err := s.storage.User().CheckUniques(req)
	if err != nil{
		s.logger.Error(err.Error())
	}
	return &pb.CheckUniquesResponse{IsExist: tf},nil
}

func (s UserService) GetUserByEmail(ctx context.Context, req *pb.EmailRequest) (*pb.User, error){
	user, err := s.storage.User().GetUserByEmail(req)
	if err != nil{
		s.logger.Error(err.Error())
	}
	return user, nil
}

func (s UserService) GetUserByRefreshToken(ctx context.Context, req *pb.RefreshToken) (*pb.User, error){
	user, err := s.storage.User().GetUserByRefreshToken(req)
	if err != nil{
		s.logger.Error(err.Error())
	}
	return user, nil
}
