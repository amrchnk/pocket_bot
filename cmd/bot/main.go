package main

import (
	"github.com/amrchnk/pocket_bot/pkg/repository"
	"github.com/amrchnk/pocket_bot/pkg/repository/boltDB"
	"github.com/amrchnk/pocket_bot/pkg/server"
	"github.com/amrchnk/pocket_bot/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("2061629377:AAHm7XED8V-VZoB-TeGNtVM2VLUAbHQxxkk")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	pocketClient,err:=pocket.NewClient("99455-5321bfb0b0228b3ec1d7f658")
	if err!=nil{
		log.Fatal(err)
	}

	db,err:=initDB()
	if err!=nil{
		log.Fatal(err)
	}
	tokenRepository:=boltDB.NewTokenRepository(db)

	tgBot := telegram.NewBot(bot,pocketClient,tokenRepository,"http://localhost:8080/")

	authorizationServer:=server.NewAuthorizationServer(pocketClient,tokenRepository,"https://t.me/AnekPocketBot")

	go func(){
		if err := tgBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}

}

func initDB()(*bolt.DB,error){
	db,err:=bolt.Open("boltDB.db",0600,nil)
	if err!=nil{
		return nil,err
	}

	if err:=db.Update(func(tx *bolt.Tx) error {
		_,err:=tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err!=nil{
			return err
		}

		_,err=tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err!=nil{
			return err
		}

		return nil
	});err!=nil{
		return nil,err
	}
	return db,nil
}
