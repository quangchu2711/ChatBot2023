package main

import (
    "log"
    "fmt"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/ghodss/yaml"
    "io/ioutil"
)

type Mqtt struct {
    Broker string
    User string
    Password string
}

type Telegram struct {
    BotName string
    BotToken string
    DeviceControlGroupID int64
}

type FileConfig struct {
    TelegramConfig Telegram
    MqttConfig Mqtt
}
var cfg FileConfig

var bot *tgbotapi.BotAPI
var updates tgbotapi.UpdatesChannel
var mqttClient mqtt.Client

func mqttBegin(broker string, user string, pw string) mqtt.Client {
    var opts *mqtt.ClientOptions = new(mqtt.ClientOptions)

    opts = mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf(broker))
    opts.SetUsername(user)
    opts.SetPassword(pw)
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    return client
}

func sendNormalMsgToTeleGroup(groupID int64, msg string) {
    message := tgbotapi.NewMessage(groupID, msg)
    bot.Send(message)
}

func telegramBotBegin(bot_token string) (tgbotapi.UpdatesChannel, *tgbotapi.BotAPI) {
    var telegramBot *tgbotapi.BotAPI
    var telgramUpdates tgbotapi.UpdatesChannel

    telegramBot, _ = tgbotapi.NewBotAPI(bot_token)
    log.Printf("Authorized on account %s", telegramBot.Self.UserName)
    u := tgbotapi.NewUpdate(0)
    telgramUpdates = telegramBot.GetUpdatesChan(u)

    return telgramUpdates, telegramBot
}

func yamlFileHandle() {
    yfile, err1 := ioutil.ReadFile("config.yaml")

    if err1 != nil {

      log.Fatal(err1)
    }
    err2 := yaml.Unmarshal(yfile, &cfg)

    if err2 != nil {

      log.Fatal(err2)
    }
}

func main() {
    yamlFileHandle()
    mqttClient = mqttBegin(cfg.MqttConfig.Broker, cfg.MqttConfig.User, cfg.MqttConfig.Password)
    fmt.Println("MQTT Connected")
    updates, bot = telegramBotBegin(cfg.TelegramConfig.BotToken)
    var groupID int64

    for update := range updates {
        if update.Message != nil {
            groupID = update.Message.Chat.ID
            fmt.Println(groupID)
        }
    }
}