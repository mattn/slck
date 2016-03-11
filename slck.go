package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pastjean/slackapi/api"
	"github.com/pastjean/slackapi/rtm"
)

func main() {
	token := os.Getenv("SLACK_TOKEN")

	resp, err := api.RtmStart(token)
	if err != nil {
		log.Fatal(err)
	}

	var channelListResponse *api.ChannelListResponse
	var userListResponse *api.UserListResponse

	slackrtm := rtm.NewSlackRTM(resp)
	slackrtm.OnMessageEvents(func(evt rtm.MessageEvent) {
		d, _ := strconv.ParseFloat(evt.TS, 64)
		t := time.Unix(int64(d), 0)
		channel := evt.Channel
		user := evt.User

		if channelListResponse == nil {
			channelListResponse, err = api.GetChannelList(token)
			if err != nil {
				log.Print(err)
			}
		}
		if channelListResponse != nil {
			for _, c := range channelListResponse.Channels {
				if c.ID == evt.Channel {
					channel = c.Name
					break
				}
			}
		}
		if userListResponse == nil {
			userListResponse, err = api.GetUserList(token)
			if err != nil {
				log.Print(err)
			}
		}
		if userListResponse != nil {
			for _, m := range userListResponse.Users {
				if m.ID == evt.User {
					user = m.Name
					break
				}
			}
		}
		fmt.Printf("%s: #%s: %s: %s\n", t,
			channel,
			user,
			evt.Text)
	})
	log.Fatal(slackrtm.Start())
}