package postgres

import (
	"log"
	"user-service/config"
	"user-service/pkg/db"
	pb "user-service/protos/template-service"
	"user-service/storage/repo"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserReposisitoryTestSuite struct {
	suite.Suite
	CleanUpFunc func()
	Repository  repo.UserStorageI
}

func (s *UserReposisitoryTestSuite) SetupSuite() {
	pgPoll, err, cleanUp := db.ConnectToDB(config.Load())
	if err != nil {
		log.Fatal("Error while connecting database with suite test")
		return
	}
	s.CleanUpFunc = cleanUp
	s.Repository = NewUserRepo(pgPoll)
}

// test func
func (s *UserReposisitoryTestSuite) TestUserCRUD() {
	// struct for create user
	user := &pb.User{
		Name:     "test name",
		LastName: "test last name",
	}

	user.Id = uuid.New().String()

	// CreateUser test
	repouser, err := s.Repository.Create(user)
	s.Suite.NotNil(repouser)
	s.Suite.NoError(err)
	s.Suite.Equal(user, repouser)
	// ----------------------------------------

	// GetUser test
	repouser, err = s.Repository.GetUser(user.Id)
	s.Suite.NotNil(repouser)
	s.Suite.NoError(err)
	s.Suite.Equal(user, repouser)
	// -----------------------------------------
	
	// // Update test
	// repouser, err = s.Repository.UpdateUser(user.)
	// s.Suite.NotNil(repouser)
	// s.Suite.NoError(err)
	// s.Suite.NotEqual(user, repouser)
	// // --------------------------------------------
	
	
	
	// Getall test
	
	allrequest := pb.GetAllRequest{
		Page: 1,
		Limit: 5,
	}
	
	repousers, err := s.Repository.GetAllUsers(&allrequest)
	
	s.Suite.NotNil(repousers)
	s.Suite.NoError(err)
	
	// Delete test
	repouser, err = s.Repository.DeleteUser(user.Id)
	s.Suite.NotNil(repouser)
	s.Suite.NoError(err)
	// --------------------------------------------
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(UserReposisitoryTestSuite))
}
