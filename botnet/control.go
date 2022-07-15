package botnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
)

type Control struct {
	Bots     []*Bot
	BotCount int
}

type Task func(bot *Bot)

func NewControl(bots []*Bot) *Control {
	return &Control{
		Bots:     bots,
		BotCount: len(bots),
	}
}

func (control *Control) GetCommunityInfo(link string) (*Community, error) {
	response, errResp := GetUnauthorized(fmt.Sprintf("g/s/community/link-identify?q=%s", link))
	if errResp != nil {
		return nil, errors.New("Не удалось получить информацию о сообществе, возникла ошибка при отправке запроса")
	}
	rawBody, errRead := io.ReadAll(response.Body)
	if errRead != nil {
		return nil, errors.New("Не удалось получить информацию о сообществе, ответ сервера полностью нечитаемый")
	}
	var body CommunityResponse
	errUnmarshal := json.Unmarshal(rawBody, &body)
	if errUnmarshal != nil {
		return nil, errors.New(fmt.Sprintf("Не удалось получить информацию о сообществе, ответ сервера не в JSON формате (%d)", response.StatusCode))
	} else if response.StatusCode == 400 {
		return nil, errors.New("Этого сообщества не существует")
	}
	return &body.Community, nil
}

func (control *Control) GetObjectInfo(link string) (*LinkInfoV2, error) {
	response, errResp := GetUnauthorized(fmt.Sprintf("g/s/link-resolution?q=%s", link))
	if errResp != nil {
		return nil, errors.New("Не удалось получить информацию об объекте, возникла ошибка при отправке запроса")
	}
	rawBody, errRead := io.ReadAll(response.Body)
	if errRead != nil {
		return nil, errors.New("Не удалось получить информацию об объекте, ответ сервера полностью нечитаемый")
	}
	var body LinkInfoResponse
	errUnmarshal := json.Unmarshal(rawBody, &body)
	if errUnmarshal != nil {
		return nil, errors.New(fmt.Sprintf("Не удалось получить информацию об объекте, ответ сервера не в JSON формате (%d)", response.StatusCode))
	} else if response.StatusCode == 400 {
		return nil, errors.New("Этого объекта не существует")
	}
	return &body.LinkInfoV2, nil
}

func (control *Control) Execute(task Task) {
	for _, bot := range control.Bots {
		task(bot)
	}
}

func (control *Control) ExecuteAsync(task Task) {
	var wait sync.WaitGroup
	for _, bot := range control.Bots {
		wait.Add(1)
		go func(a *Bot) {
			// time.Sleep(time.Duration(RandInt(1, 10)) * time.Second)
			task(a)
			wait.Done()
		}(bot)
	}
	wait.Wait()
}

func (control *Control) Authorize() {
	control.Execute(func(bot *Bot) {
		if !bot.Authorize() {
			control.BotCount--
		}
	})
}

func (control *Control) JoinCommunity(ndcId int) {
	control.ExecuteAsync(func(bot *Bot) {
		bot.JoinCommunity(ndcId)
	})
}

func (control *Control) LeaveCommunity(ndcId int) {
	control.ExecuteAsync(func(bot *Bot) {
		// недоступно
	})
}

func (control *Control) ChangeNickname(nickname string) {
	control.ExecuteAsync(func(bot *Bot) {
		// недоступно
	})
}

func (control *Control) ChangeAvatar(icon string) {
	control.ExecuteAsync(func(bot *Bot) {
		// недоступно
	})
}

func (control *Control) JoinChat(ndcId int, threadId string) {
	control.ExecuteAsync(func(bot *Bot) {
		bot.JoinChat(ndcId, threadId)
	})
}

func (control *Control) LeaveChat(ndcId int, threadId string) {
	control.ExecuteAsync(func(bot *Bot) {
		// недоступно
	})
}

func (control *Control) SendMessage(ndcId int, threadId string, messageType int, content string) {
	control.ExecuteAsync(func(bot *Bot) {
		bot.SendMessage(ndcId, threadId, messageType, content)
	})
}
