package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("[ERROR]: Cannot load environment variables - ", err)
	}

	var SlackBotToken = os.Getenv("SLACK_BOT_TOKEN")
	var SlackAppToken = os.Getenv("SLACK_APP_TOKEN")

	bot := slacker.NewClient(SlackBotToken, SlackAppToken)

	go printCommandEvents(bot.CommandEvents())

	bot.Command("work hours left until <finishHour>", &slacker.CommandDefinition{
		Description: "Work hours left calculator",
		Examples:    []string{"work hours left until 18:00"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			finishHour := request.Param("finishHour")
			getHourMin := strings.SplitN(finishHour, ":", 2)

			now := time.Now().Local()

			hour, _ := strconv.Atoi(getHourMin[0])
			min, _ := strconv.Atoi(getHourMin[1])

			if hour < now.Hour() {
				response.Reply("There is something wrong!")
				return
			}

			finishDate := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, time.Local)
			diff := finishDate.Sub(now)

			res := "Stay strong!"
			if diff.Hours() > 1 {
				res += fmt.Sprintf(" %d hours and ", int(diff.Hours()))
			}

			res += fmt.Sprintf(" %d mins left", int(diff.Minutes())%60)
			response.Reply(res)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal("[ERROR]: Cannot listen bot - ", err)
	}
}
