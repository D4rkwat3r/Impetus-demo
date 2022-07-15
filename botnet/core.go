package botnet

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Callback func(response *http.Response)

const (
	base            = "https://service.narvii.com/api/v1/"
	baseWeb         = "https://aminoapps.com/api/"
	lang            = "ru"
	userAgent       = "Dalvik/2.1.0 (Linux; U; Android 7.1.2; SM-G977N Build/beyond1qlteue-user 7; com.narvii.amino.master/3.4.33597)"
	contentTypeJSON = "application/json; charset=utf-8"
)

var client = &http.Client{}

func normalize(sid *string, contentType *string, body []byte, request *http.Request) {
	request.Header.Set("NDCDEVICEID", NdcDevice())
	request.Header.Set("NDCLANG", lang)
	request.Header.Set("User-Agent", userAgent)
	if sid != nil {
		request.Header.Set("NDCAUTH", *sid)
	}
	if contentType != nil {
		request.Header.Set("Content-Type", *contentType)
		request.Header.Set("NDC-MSG-SIG", NdcSignature(body))
	}
}

func Send(bot *Bot, httpMethod string, endpoint string, contentType *string, body []byte) (*http.Response, error) {
	request, reqErr := http.NewRequest(httpMethod, base+endpoint, bytes.NewBuffer(body))
	if reqErr != nil {
		return nil, reqErr
	}
	if bot != nil {
		normalize(bot.SID, contentType, body, request)
	} else {
		normalize(nil, contentType, body, request)
	}
	response, respErr := client.Do(request)
	if respErr != nil {
		return nil, respErr
	}
	return response, nil
}

func GetUnauthorized(endpoint string) (*http.Response, error) {
	return Send(nil, "GET", endpoint, nil, make([]byte, 0))
}

func Get(bot *Bot, endpoint string) (*http.Response, error) {
	return Send(bot, "GET", endpoint, nil, make([]byte, 0))
}

func Post(bot *Bot, endpoint string, contentType string, body []byte) (*http.Response, error) {
	return Send(bot, "POST", endpoint, &contentType, body)
}

func PostJson(bot *Bot, endpoint string, body interface{}) (*http.Response, error) {
	serialized, _ := json.Marshal(body)
	return Post(bot, endpoint, contentTypeJSON, serialized)
}

func Delete(bot *Bot, endpoint string) (*http.Response, error) {
	return Send(bot, "DELETE", endpoint, nil, nil)
}

func WebPost(bot *Bot, endpoint string, body interface{}) (*http.Response, error) {
	serialized, _ := json.Marshal(body)
	request, reqErr := http.NewRequest("POST", baseWeb+endpoint, bytes.NewBuffer(serialized))
	if reqErr != nil {
		return nil, reqErr
	}
	if bot != nil {
		request.Header.Set("Cookie", *bot.SID)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Requested-With", "XMLHttpRequest")
	}
	response, respErr := client.Do(request)
	if respErr != nil {
		return nil, respErr
	}
	return response, nil
}
