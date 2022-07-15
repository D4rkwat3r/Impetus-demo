package botnet

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

type Account struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	NativeDeviceID string `json:"nativeDeviceID"`
}

const (
	logo    = "   ____               __         \n  /  _/_ _  ___  ___ / /___ _____\n _/ //  ' \\/ _ \\/ -_) __/ // (_-<\n/___/_/_/_/ .__/\\__/\\__/\\_,_/___/\n         /_/                     "
	version = "demo"
)

var (
	scanner   = bufio.NewScanner(os.Stdin)
	menuItems = []string{
		"Войти в сообщество",
		"Покинуть сообщество",
		"Изменить никнейм",
		"Войти в чат",
		"Покинуть чат",
		"Отправить сообщение",
		"Спамить входом-выходом",
		"Спамить сообщениями",
	}
)

func InputValue(hint string) string {
	fmt.Println(hint)
	fmt.Print("> ")
	scanner.Scan()
	return scanner.Text()
}

func PressEnter() {
	fmt.Println(color.CyanString("Нажмите на Enter для продолжения..."))
	scanner.Scan()
	scanner.Text()
}

func ClearWindows() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func ClearUnix() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func ClearConsole() {
	if runtime.GOOS == "windows" {
		ClearWindows()
	} else {
		ClearUnix()
	}
}

func RenderLogo() {
	ClearConsole()
	color.Yellow("%s --- %s\n", logo, version)
}

func RenderStats(control *Control) {
	color.Yellow("Число аккаунтов: %d\n\n", control.BotCount)
}

func RenderMenu() {
	for i, item := range menuItems {
		if i == 1 || i == 2 || i == 4 || i == 6 || i == 7 {
			color.Red("%d. %s (недоступно)\n", i+1, item)
		} else {
			color.Cyan("%d. %s\n", i+1, item)
		}
	}
}

func GetFromFile(fileName string) ([]*Bot, error) {
	var accounts []Account
	var bots []*Bot
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Не удаётся открыть файл %s, возможно, его не существует в этой директории\n", fileName)
		return nil, err
	}
	err = json.Unmarshal(data, &accounts)
	if err != nil {
		fmt.Printf("Файл %s содержит некорректный JSON, чтение невозможно\n", fileName)
		return nil, err
	}
	for _, account := range accounts {
		bots = append(bots, NewBot(account.Email, account.Password, account.NativeDeviceID, nil, nil, nil))
	}
	return bots, nil
}

func LoadAccounts() *Control {
	text := InputValue("Введите название файла с JSON-представленными аккаунтами")
	accounts, err := GetFromFile(text)
	if err != nil {
		LoadAccounts()
	}
	control := NewControl(accounts)
	control.Authorize()
	return control
}

func GetControl() (*Control, error) {
	fromDB, errRead := Read()
	var control *Control
	if errRead != nil {
		fmt.Println("Не удалось прочитать данные из локальной базы данных")
		return nil, errRead
	}
	if len(fromDB) == 0 {
		for {
			text := InputValue("Ботнет не может работатать без аккаунтов. Загрузить их? (y/n)")
			if text == "y" {
				control = LoadAccounts()
				break
			} else if text == "n" {
				os.Exit(0)
			}
		}
	} else {
		control = NewControl(fromDB)
	}
	return control, nil
}

func RequestCommunity(control *Control) (*Community, error) {
	community, err := control.GetCommunityInfo(InputValue("Введите ссылку на сообщество"))
	if err != nil {
		color.Red("%s\n", err.Error())
		return nil, err
	}
	return community, nil
}

func RequestChat(control *Control) (*Community, *LinkInfoV2, error) {
	community, err := control.GetCommunityInfo(InputValue("Введите ссылку на сообщество с чатом"))
	if err != nil {
		color.Red("%s\n", err.Error())
		return nil, nil, err
	}
	chat, err := control.GetObjectInfo(InputValue("Введите ссылку на чат"))
	if err != nil {
		color.Red("%s\n", err.Error())
		return nil, nil, err
	}
	return community, chat, nil
}

func RequestMessageType() int {
	for {
		msgType, err := strconv.Atoi(InputValue("Введите тип сообщения"))
		if err != nil {
			color.Red("Неверный ввод")
			continue
		}
		return msgType
	}
}

func CommunityMembership(control *Control, newStatus bool) {
	community, err := RequestCommunity(control)
	if err == nil {
		if newStatus {
			control.JoinCommunity(community.NdcId)
		} else {
			// недоступно
		}
	}
}

func ChatMembership(control *Control, newStatus bool) {
	community, chat, err := RequestChat(control)
	if err == nil {
		if newStatus {
			control.JoinChat(community.NdcId, chat.Extensions.LinkInfo.ObjectId)
		} else {
			// недоступно
		}
	}
}

func Message(control *Control) {
	community, chat, err := RequestChat(control)
	if err == nil {
		msgType := RequestMessageType()
		control.SendMessage(community.NdcId,
			chat.Extensions.LinkInfo.ObjectId,
			msgType,
			InputValue("Введите текст сообщения"))
	}
}

func JLSpam(control *Control) {
	// недоступно
}

func MessageSpam(control *Control) {
	// недоступно
}

func MainLifecycle(control *Control) {
	RenderLogo()
	RenderStats(control)
	RenderMenu()
	userSelect := InputValue("\nЧто вы хотите сделать?")
	switch userSelect {
	case "1":
		CommunityMembership(control, true)
	case "2":
		// недоступно
	case "3":
		// недоступно
	case "4":
		ChatMembership(control, true)
	case "5":
		// недоступно
	case "6":
		Message(control)
	case "7":
		// недоступно
	case "8":
		// недоступно
	default:
		MainLifecycle(control)
	}
	PressEnter()
	MainLifecycle(control)
}

func ConsoleUIStart() {
	errInit := InitDB()
	if errInit != nil {
		fmt.Println("Не удалось получить доступ к локальной базе данных")
		return
	}
	control, err := GetControl()
	if err != nil {
		return
	}
	MainLifecycle(control)
}
