package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/url"
)

const (
	commandStart = "start"
	replyStartTemplate="Привет! Чтобы пользоваться данным ботом, тебе необходимо дать доступ к своему аккаунту по этой ссылочке: \n%s"
	replyAlreadyAuth="Ты уже авторизирован! Жду твои ссылочки; )"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error{

	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error{
	msg := tgbotapi.NewMessage(message.Chat.ID, "Ссылка успешно сохранена!")
	_,err:=url.ParseRequestURI(message.Text)
	if err!=nil{
		msg.Text="Ссылка невалидна"
		_,err=b.bot.Send(msg)
		return err
	}

	accessToken,err:=b.getAccessToken(message.Chat.ID)
	if err!=nil{
		msg.Text="Для сохранения ссылок, нужно авторизоваться (команда /start)"
		_,err=b.bot.Send(msg)
		return err
	}

	if err:=b.pocketClient.Add(context.Background(),pocket.AddInput{
		AccessToken: accessToken,
		URL: message.Text,
	});err!=nil{
		msg.Text="Не удалось сохранить ссылку, попробуй еще раз попозже"
		_,err=b.bot.Send(msg)
		return err
	}
	_,err=b.bot.Send(msg)
	return err
}

func (b *Bot)handleStartCommand(message *tgbotapi.Message)error{
	_,err:=b.getAccessToken(message.Chat.ID)
	if err!=nil{
		return b.initAuthProcess(message)
	}
	msg:=tgbotapi.NewMessage(message.Chat.ID,replyAlreadyAuth)
	_,err=b.bot.Send(msg)
	return err
}

func (b *Bot)handleUnknownCommand(message *tgbotapi.Message)error{
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды")

	_,err:=b.bot.Send(msg)
	return err
}