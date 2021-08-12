package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var wantToDelete = false
var toRead = ""
var tokenBuy = false
var audioMSG = false
var continueTo = true
var persCode = false
var hostName = "https://bank.goto.msk.ru"
var transactionIDD = 0
var lastFileName = ""
var dayAM = 0

type firstRequest struct {
	Token       string `json:"token"`
	Account_id  int64  `json:"account_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type secondRequest struct {
	Transaction_id int `json:"transaction_id"`
	Code           int `json:"code"`
}

type responseOf struct {
	State          string `json:"state"`
	Transaction_id int    `json:"transaction_id"`
}

type badResponseOf struct {
	State string `json:"state"`
	Error string `json:"error"`
}
var fileNames map[int64] string
var timeStamps map[int64] int64
func BotBuild() {
	fmt.Println("TESTING")
	counter := 0
	data, er := os.ReadFile("botTokenTelegram.txt")
	if er != nil {
		fmt.Println("bad")
		return
	}
	token := string(data)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	data, er = os.ReadFile("tradingToken.txt")
	if er != nil {
		fmt.Println("bad")
		return
	}
	tradingToken := string(data)
	fmt.Println(tradingToken)
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		continueTo = true
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if !update.Message.IsCommand() {
			fmt.Println("NOT COMMAND")
			switch true {
				case audioMSG:
				continueTo = false
				persCode = false
				switch true {
				case update.Message.Voice != nil:
					if update.Message.Voice.Duration > 31 {
						msg.Text = "Слишком долгое сообщение, не больше 15 секунд)"
						continueTo = false
						audioMSG = true
						persCode = false
					} else {
						fmt.Println("OP")
						file, errr := bot.GetFile(tgbotapi.FileConfig{FileID: update.Message.Voice.FileID})
						if errr != nil {
							msg.Text = "Упс( С вашим голосовым есть проблемы,попробуйте снова"
							continue
						}
						url := file.Link(token)
						toUse := time.Now().UnixMilli()
						fileID := fmt.Sprintf("AD%d", toUse)
						lastFileName = fileID
						if err = DownloadFile(fileID+".ogg", url); err != nil {
							msg.Text = "Произошла ошибка, попробуйте снова("
							continue
						} else {
							msg.Text = "Напишите, сколько дней должна играть ваша реклама)"
							audioMSG = false
							fmt.Println(counter)
							wantToDelete = true
						}
					}
				case update.Message.Audio != nil:
					if update.Message.Audio.Duration > 31 {
						fmt.Println("GOT HERE")
						msg.Text = "Слишком долгое сообщение, не больше 15 секунд)"
						audioMSG = true
						persCode = false
					} else {
						fmt.Println("OP")
						continueTo = false
						file, errr := bot.GetFile(tgbotapi.FileConfig{FileID: update.Message.Audio.FileID})
						if errr != nil {
							msg.Text = "Упс( С вашим голосовым есть проблемы,попробуйте снова"
							continue
						}
						url := file.Link(token)
						toUse := time.Now().UnixMilli()
						fileID := fmt.Sprintf("AD%d", toUse)
						lastFileName = fileID
						if err = DownloadFile(fileID+".ogg", url); err != nil {
							msg.Text = "Произошла ошибка, попробуйте снова("
							continue
						} else {
							msg.Text = "Напишите, сколько дней должна играть ваша реклама)"
							audioMSG = false
							fmt.Println(counter)
							wantToDelete = true
						}
					}
				default:
					msg.Text = "Пожалуйста, попробуйте снова"
				}
				case wantToDelete:
					text := update.Message.Text
					switch text {

					case "0":
						wantToDelete = false
						persCode = false
						continueTo = false
						fmt.Println(lastFileName)
						os.Remove(lastFileName + ".ogg")
						msg.Text = "Принято, не отправляйте нам код операции, и повторите попытку через команду /buyad :("
					default:
						numOfDays, err := strconv.Atoi(text)
						if err != nil {
							msg.Text = "Повторите попытку"
							persCode = false
							continueTo = false
							continue
						}

						id := update.Message.From.ID
						fmt.Println(id)
						req := &firstRequest{Token: "2bb2df12-fe10-43b3-8977-ef47edd389ee",
							Account_id:  int64(id),
							Amount:      200 * numOfDays,
							Description: "Payment for POPUGAY advertisement"}
						REQ, err := json.Marshal(req)
						if err != nil {
							msg.Text = "Произошла ошибка, извините("
							if _, err := bot.Send(msg); err != nil {
								fmt.Println(err)
								log.Panic(err)
							}
							continue
						}
						resp, err := http.Post(hostName+"/api/ask", "application/json", bytes.NewBuffer(REQ))
						if err != nil {
							msg.Text = "Произошла ошибка, извините"
							if _, err := bot.Send(msg); err != nil {
								fmt.Println(err)
							}
							continue
						}
						defer resp.Body.Close()
						if resp.StatusCode != 200 {
							data, _ := io.ReadAll(resp.Body)
							fmt.Println(string(data))
							fmt.Println("BAD REQUEST")
							msg.Text = "Произошла ошибка, извините"
							if _, err := bot.Send(msg); err != nil {
								fmt.Println(err)
							}
							continue
						}

						data, _ := io.ReadAll(resp.Body)
						br := responseOf{}

						fmt.Println(string(data))
						_ = json.Unmarshal(data, &br)
						transactionIDD = br.Transaction_id
						fmt.Println(transactionIDD)
						msg.Text = "Спасибо, принято) Теперь отправьте номер транзакции, который вам должен был прийти сообщением от GOTO bank"
						persCode = true
						continueTo = false
						wantToDelete = false
						dayAM = numOfDays
					}
				case persCode:
					continueTo = false
					operCode := update.Message.Text
					if operCode == "" {
						msg.Text = "Пришлите, пожалуйста, правильный код)"
					} else {
						operId, Err := strconv.Atoi(operCode)
						if Err != nil {
							msg.Text = "Попробуйте еще раз и перепроверьте код("
						} else {
							req := secondRequest{Transaction_id: transactionIDD, Code: operId}

							REQ, _ := json.Marshal(req)
							fmt.Println(string(REQ))
							resp, err := http.Post(hostName+"/api/verify", "application/json", bytes.NewBuffer(REQ))
							defer resp.Body.Close()
							if err != nil {
								msg.Text = "Ошибка("
								continue
							} else {
								if resp.StatusCode == 200 {
									msg.Text = "Получено, спасибо за использование)"
									fmt.Println(operId)
									persCode = false
									amount, err := os.ReadFile("numberInMap.txt")
									if err != nil {
										fmt.Println("OP")
										return
									}
									am, err := strconv.Atoi(string (amount))
									am++
									stree := fmt.Sprintf("%d", am)
									_ = os.WriteFile("numberInMap.txt", []byte (stree), 0644)


									m := make(map[int64] string)
									data, err := os.ReadFile("reklMap.txt")
									if err != nil {
										fmt.Println(err)
									}
									err = json.Unmarshal(data, &m)
									m[int64(am)] = lastFileName
									REQ, _ = json.Marshal(m)
									err = os.WriteFile("reklMap.txt", REQ, 0644)

									data, err = os.ReadFile("reklTimeMap.txt")
									ma := make(map[int64] int64)
									err = json.Unmarshal(data, &ma)
									ma[int64(am)] = int64(dayAM)
									REQ, _ = json.Marshal(ma)
									err = os.WriteFile("reklTimeMap.txt", REQ, 0644)

								} else {
									data, _ := io.ReadAll(resp.Body)
									fmt.Println(string(data))

									br := badResponseOf{}
									_ = json.Unmarshal(data, &br)

									msg.Text = "Ошибка("
									fmt.Println("Ошибка(")
								}
							}
						}

					}

				default:
					fmt.Println("BREAKPOINT")
					msg.Text = "Упс( Ознакомьтесь с коммандами, я не понимаю ваш запрос"

				}
			}

			if continueTo {
				switch update.Message.Command() {
				case "help":
					msg.Text = "Привет! С моей помощью ты можешь заказать рекламу на платформе POPUGAY 1.0."
				case "price":
					msg.Text = "200 Готублей"
				case "buyad":
					msg.Text = "Отправьте голосовое сообщение, которое будет рекламой)"
					audioMSG = true
				case "start":
					msg.Text = "Приветствую вас! У этого бота есть 3 команды - /buyad, /help и /price. /help - напоминает о боте. /buyad - позволяет купить рекламу. /price - напоминает цену покупки. Учтите, что длина сообщения не может превышать 15 секунд."
				default:
					msg.Text = "Плохой запрос, попробуйте снова("
				}
			}
			if _, err := bot.Send(msg); err != nil {
				fmt.Println(err)
				log.Panic(err)
			}
		}
	}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
