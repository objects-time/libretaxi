package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"libretaxi/context"
	"libretaxi/objects"
	"log"
)

type Handler interface {
	Handle(user *objects.User, context *context.Context, message *tgbotapi.Message)
}

func isStateChanged(context *context.Context, previousState objects.MenuId, userId int64) (result bool) {
	user := context.Repo.FindUser(userId)

	if user == nil {
		return true
	}

	return user.MenuId != previousState
}

func HandleMessage(context *context.Context, userId int64, message *tgbotapi.Message) {
	log.Printf("Message: '%s'", message.Text)
	previousState := objects.Menu_Ban

	for isStateChanged(context, previousState, userId) == true {
		user := context.Repo.FindUser(userId)

		if user == nil {
			user = &objects.User{
				UserId: userId,
				MenuId: objects.Menu_Init,
			}
		}
		//fmt.Printf("%+v\n", message.Location)

		previousState = user.MenuId
		var handler Handler

		switch user.MenuId {
		case objects.Menu_Init:
			handler = NewInitMenu()
		case objects.Menu_AskLocation:
			handler = NewAskLocationMenu()
		}

		handler.Handle(user, context, message)
	}
}