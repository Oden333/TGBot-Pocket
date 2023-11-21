package main

import (
	"log"

	"github.com/Oden333/TG_bot/pkg/config"
	"github.com/Oden333/TG_bot/pkg/repository"
	"github.com/Oden333/TG_bot/pkg/repository/boltdb"
	"github.com/Oden333/TG_bot/pkg/server"
	"github.com/Oden333/TG_bot/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg)
	bot, err := tgbotapi.NewBotAPI("6538630363:AAE-_BalSqvxcSEmOFQ98Mn0V576nNqU3oo")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	pocketClient, err := pocket.NewClient("109198-270c17a8c5baae383e6ed67")
	if err != nil {
		log.Fatal(err)
	}
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost/")
	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, "https://t.me/PocketManager_Bot")

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}

}

func initDB() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}
