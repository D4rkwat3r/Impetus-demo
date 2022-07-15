package botnet

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Bot struct {
	Login          string
	Password       string
	NativeDeviceID string
	UID            *string
	SID            *string
	Nickname       *string
}

func NewBot(login string, password string, nativeDeviceID string, uid *string, sid *string, nickname *string) *Bot {
	return &Bot{
		Login:          login,
		Password:       password,
		NativeDeviceID: nativeDeviceID,
		UID:            uid,
		SID:            sid,
		Nickname:       nickname,
	}
}

func (bot *Bot) InfoMessage(text string, successMessage bool, failureMessage bool) {
	format := "[%s]: %s\n"
	if successMessage {
		color.Green(fmt.Sprintf(format, bot.Login, text))
	} else if failureMessage {
		color.Red(fmt.Sprintf(format, bot.Login, text))
	} else {
		fmt.Printf(format, bot.Login, text)
	}
}

func (bot *Bot) Success(text string) {
	bot.InfoMessage(text, true, false)
}

func (bot *Bot) Failure(text string) {
	bot.InfoMessage(text, false, true)
}

func (bot *Bot) NetworkError(action string) {
	color.Red(fmt.Sprintf("[%s] %s не удалось из-за сетевой ошибки\n", bot.Login, action))
}

func (bot *Bot) HttpError(action string, statusCode int) {
	color.Red(fmt.Sprintf("[%s] %s не удалось, ответ сервера не содержит JSON (%d)\n", bot.Login, action, statusCode))
}

func (bot *Bot) ApiError(action string, message string) {
	color.Red(fmt.Sprintf("[%s] %s не удалось, сообщение: \"%s\"\n", bot.Login, action, message))
}

func (bot *Bot) CheckResponseV1Api(response *http.Response, err error, action string) {
	if err == nil {
		rawBody, _ := io.ReadAll(response.Body)
		var body BaseResponse
		err = json.Unmarshal(rawBody, &body)
		if err != nil {
			bot.HttpError(action, response.StatusCode)
		} else if body.ApiStatusCode != 0 {
			bot.ApiError(action, body.ApiMessage)
		} else {
			bot.Success(action + " удалось")
		}
	} else {
		bot.NetworkError(action)
	}
}

func (bot *Bot) CheckResponseWeb(response *http.Response, err error, action string) {
	if err == nil {
		rawBody, _ := io.ReadAll(response.Body)
		var body BaseWebResponse
		err = json.Unmarshal(rawBody, &body)
		if err != nil {
			bot.HttpError(action, response.StatusCode)
		} else if body.Code != 200 {
			bot.ApiError(action, body.ApiMessage)
		} else {
			bot.Success(action + " удалось")
		}
	} else {
		bot.NetworkError(action)
	}
}

func (bot *Bot) CheckResponse(response *http.Response, err error, action string) {
	responseServer := response.Header.Get("server")
	if strings.Contains(responseServer, "openresty") {
		bot.CheckResponseWeb(response, err, action)
	} else {
		bot.CheckResponseV1Api(response, err, action)
	}
}

func (bot *Bot) ChangeCommunityMembershipType(ndcId int, mType string, mTypeText string) {
	data := CommunityMembershipChange{NdcId: strconv.Itoa(ndcId)}
	response, err := WebPost(bot, mType, data)
	bot.CheckResponse(response, err, mTypeText)
}

func (bot *Bot) JoinCommunity(ndcId int) {
	bot.ChangeCommunityMembershipType(ndcId, "join", "Войти в сообщество")
}

func (bot *Bot) LeaveCommunity(ndcId int) {
	// недоступно
}

func (bot *Bot) Authorize() bool {
	data := Auth{
		Email:      bot.Login,
		V:          2,
		Secret:     "0 " + bot.Password,
		DeviceID:   bot.NativeDeviceID,
		ClientType: 100,
		Action:     "normal",
		Timestamp:  TimeInMillis(),
	}
	response, err := PostJson(bot, "g/s/auth/login", data)
	if err == nil {
		rawBody, _ := io.ReadAll(response.Body)
		var body AccountResponse
		err = json.Unmarshal(rawBody, &body)
		if err != nil {
			bot.HttpError("Авторизоваться", response.StatusCode)
			return false
		} else if body.ApiStatusCode != 0 {
			bot.ApiError("Авторизоваться", body.ApiMessage)
			return false
		}
		sessionString := fmt.Sprintf("sid=%s", body.SID)
		bot.SID = &sessionString
		bot.Nickname = &body.UserProfile.Nickname
		bot.UID = &body.UserProfile.Uid
		err = Write(Bot{
			Login:          bot.Login,
			Password:       bot.Password,
			NativeDeviceID: bot.NativeDeviceID,
			UID:            bot.UID,
			SID:            bot.SID,
			Nickname:       bot.Nickname,
		})
		if err != nil {
			bot.Failure("Авторизоваться не удалось (ошибка при записи в локальную базу данных)")
			return false
		}
		bot.Success(fmt.Sprintf("Успешная авторизация, никнейм: %s\n", *bot.Nickname))
		return true
	}
	bot.NetworkError("Авторизоваться")
	return false
}

func (bot *Bot) EditProfile(data interface{}, action string) {
	// недоступно
}

func (bot *Bot) ChangeNickname(nickname string) {
	// недоступно
}

func (bot *Bot) ChangeAvatar(icon string) {
	// недоступно
}

func (bot *Bot) ChangeChatMembershipType(ndcId int, threadId string, newStatus bool) {
	var action string
	var endpoint string
	data := ChatMembershipChange{
		NdcId:    "x" + strconv.Itoa(ndcId),
		ThreadId: threadId,
	}
	if newStatus {
		action, endpoint = "Войти в чат", "join-thread"
	} else {
		// недоступно
		return
	}
	response, err := WebPost(bot, endpoint, data)
	bot.CheckResponse(response, err, action)
}

func (bot *Bot) JoinChat(ndcId int, threadId string) {
	bot.ChangeChatMembershipType(ndcId, threadId, true)
}

func (bot *Bot) LeaveChat(ndcId int, threadId string) {
	// недоступно
}

func (bot *Bot) SendMessage(ndcId int, threadId string, messageType int, content string) {
	data := MessageSend{
		NdcId:    "x" + strconv.Itoa(ndcId),
		ThreadId: threadId,
		Message: struct {
			Type    int    `json:"type"`
			Content string `json:"content"`
		}{Type: messageType, Content: content},
	}
	response, err := WebPost(bot, "add-chat-message", data)
	bot.CheckResponse(response, err, "Отправить сообщение")
}
