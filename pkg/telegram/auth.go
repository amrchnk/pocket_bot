package telegram

import (
	"context"
	"fmt"
	"github.com/amrchnk/pocket_bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func(b *Bot)initAuthProcess(message *tgbotapi.Message)error{
	authLink,err:=b.generateAuthorizationLink(message.Chat.ID)
	if err!=nil{
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID,fmt.Sprintf(replyStartTemplate,authLink))

	_,err=b.bot.Send(msg)
	return err
}

func(b *Bot)getAccessToken(chat_id int64)(string,error){
	return b.tokenRepository.Get(chat_id,repository.AccessTokens)
}

func (b *Bot)generateAuthorizationLink(chatID int64)(string,error){
	redirectURL:=b.generateRedirectURL(chatID)
	requestToken,err:=b.pocketClient.GetRequestToken(context.Background(),redirectURL)
	if err!=nil{
		return "",err
	}

	if err:=b.tokenRepository.Save(chatID,requestToken,repository.RequestTokens);err!=nil{
		return "",err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken,redirectURL)
}

func (b *Bot)generateRedirectURL(chatId int64)string{
	return fmt.Sprintf("%s?chat_id=%d",b.redirectURL,chatId)
}

