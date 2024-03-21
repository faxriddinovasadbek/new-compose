package postgres

import (
	"fmt"
	pb "user-service/protos/template-service"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *pb.User) (*pb.User, error) {
	query := `
	INSERT INTO users 
		(id, name, last_name, email, password, user_name, refresh_token) 
	VALUES 
		($1, $2, $3, $4, $5, $6, $7) 
	RETURNING 
		id,
		name, 
		last_name, 
		email, 
		password,
		user_name,
		refresh_token`

	var respouser pb.User

	err := r.db.QueryRow(query, user.Id, user.Name, user.LastName, user.Email, user.Password, user.UserName, user.RefreshToken).Scan(
		&respouser.Id,
		&respouser.Name,
		&respouser.LastName,
		&respouser.Email,
		&respouser.Password,
		&respouser.UserName,
		&respouser.RefreshToken,
	)

	if err != nil {
		return nil, err
	}

	return &respouser, nil
}

func (r *userRepo) GetUser(id string) (*pb.User, error) {
	query := `SELECT id, name, last_name, email, password, user_name FROM users WHERE id = $1`

	var respouser pb.User

	err := r.db.QueryRow(query, id).Scan(
		&respouser.Id,
		&respouser.Name,
		&respouser.LastName,
		&respouser.Email,
		&respouser.Password,
		&respouser.UserName,
	)

	if err != nil {
		return nil, err
	}

	return &respouser, nil
}

func (r *userRepo) DeleteUser(id string) (*pb.User, error) {
	query := `DELETE FROM users WHERE id = $1 RETURNING id, name, last_name, email, password, user_name`

	var respouser pb.User

	err := r.db.QueryRow(query, id).Scan(
		&respouser.Id,
		&respouser.Name,
		&respouser.LastName,
		&respouser.Email,
		&respouser.Password,
		&respouser.UserName,
	)

	if err != nil {
		return nil, err
	}

	return &respouser, nil
}

func (r *userRepo) UpdateUser(req *pb.User) (*pb.User, error) {
	query := `
    UPDATE 
        users
    SET 
		id = $1,
		name = $2, 
		last_name = $3, 
		email = $4, 
		password = $5,
		user_name = $6,
		refresh_token $7,
    WHERE 
        id = $8
    RETURNING 
        id , name, last_name, email, password, user_name, refresh_token`

	var respouser pb.User

	err := r.db.QueryRow(query,
		req.Id,
		req.Name,
		req.LastName,
		req.Email,
		req.Password,
		req.UserName,
		req.RefreshToken,
		req.Id).Scan(
		&respouser.Id,
		&respouser.Name,
		&respouser.LastName,
		&respouser.Email,
		&respouser.Password,
		&respouser.UserName,
		&respouser.RefreshToken,
	)

	if err != nil {
		return nil, err
	}

	return &respouser, nil
}

func (r *userRepo) GetAllUsers(req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	offet := req.Limit * (req.Page - 1)
	query := `SELECT id, name, last_name, email, password, user_name FROM users LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, req.Limit, offet)
	if err != nil {
		return nil, err
	}

	var allusers pb.GetAllResponse
	for rows.Next() {
		var user pb.User
		if err := rows.Scan(&user.Id, &user.Name, &user.LastName, &user.Email, &user.Password, &user.UserName); err != nil {
			return nil, err
		}
		allusers.AllUsers = append(allusers.AllUsers, &user)
	}
	return &allusers, err
}

func (r *userRepo) CheckUniques(req *pb.CheckUniquesRequest) (bool, error) {
	var exists int
	err := r.db.QueryRow(fmt.Sprintf("SELECT COUNT(1) from users where %s = $1", req.Field), req.Value).Scan(&exists)

	if err != nil {
		return true, err
	}

	if exists != 0 {
		return true, nil
	}
	return false, nil
}

func (r *userRepo) GetUserByEmail(req *pb.EmailRequest) (*pb.User, error) {
	query := `
	SELECT 
		id, 
		name, 
		last_name, 
		email, 
		password, 
		user_name,
		refresh_token
	FROM 
		users 
	WHERE 
		email = $1`
	

	var respouser pb.User

	err := r.db.QueryRow(query, req.Email).Scan(
		&respouser.Id,
		&respouser.Name,
		&respouser.LastName,
		&respouser.Email,
		&respouser.Password,
		&respouser.UserName,
		&respouser.RefreshToken,
	
	)

	if err != nil {
		return nil, err
	}

	return &respouser, nil
}

func (r *userRepo) GetUserByRefreshToken(req *pb.RefreshToken) (*pb.User, error){
	query := `
	SELECT 
		id, 
		name, 
		last_name, 
		email, 
		password, 
		user_name,
		refresh_token
	FROM 
		users 
	WHERE 
		refresh_token = $1`
	

	var respouser pb.User

	err := r.db.QueryRow(query, req.RefreshToken).Scan(
		&respouser.Id,
		&respouser.Name,
		&respouser.LastName,
		&respouser.Email,
		&respouser.Password,
		&respouser.UserName,
		&respouser.RefreshToken,
	
	)

	if err != nil {
		return nil, err
	}

	return &respouser, nil
}
