package repo

import (
    pb "user-service/protos/template-service"
)

//UserStorageI ...
type UserStorageI interface {
    Create(*pb.User) (*pb.User, error)
    GetUser(id string) (*pb.User,error)
    DeleteUser(id string) (*pb.User, error)
    UpdateUser(*pb.User) (*pb.User, error)
    GetAllUsers(*pb.GetAllRequest) (*pb.GetAllResponse, error)
    CheckUniques(req *pb.CheckUniquesRequest) (bool, error)
    GetUserByEmail(*pb.EmailRequest) (*pb.User, error)
    GetUserByRefreshToken(req *pb.RefreshToken) (*pb.User, error)

}
