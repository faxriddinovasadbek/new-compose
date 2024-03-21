package storage

import (
    "github.com/jmoiron/sqlx"
    "user-service/storage/postgres"
    "user-service/storage/repo"
)

//IStorage ...
type IStorage interface {
    User() repo.UserStorageI
}

type storagePg struct {
    db         *sqlx.DB
    userRepo   repo.UserStorageI
}

//NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *storagePg {
    return &storagePg{
        db:         db,
        userRepo:   postgres.NewUserRepo(db),
    }
}

func (s storagePg) User() repo.UserStorageI {
    return s.userRepo
}
