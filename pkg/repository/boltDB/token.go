package boltDB

import (
	"errors"
	"github.com/amrchnk/pocket_bot/pkg/repository"
	"github.com/boltdb/bolt"
	"strconv"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository)Save(chatID int64, token string,bucket repository.Bucket) error{
	return r.db.Update(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(bucket))
		return b.Put(strconv.AppendInt(nil, chatID, 10),[]byte(token))
	})
}
func (r *TokenRepository)Get(chatID int64,bucket repository.Bucket) (string, error){
	var token string
	err:=r.db.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(bucket))
		data:=b.Get(strconv.AppendInt(nil, chatID, 10))
		token=string(data)
		return nil
	})
	if err!=nil{
		return "",err
	}

	if token==""{
		return "",errors.New("token is not found")
	}

	return token,nil
}

