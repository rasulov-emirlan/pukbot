package telegrambot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/telebot.v3"
)

func (t *Telegrambot) PukCreate() telebot.HandlerFunc {
	type result struct {
		FileID       string `json:"file_id"`
		FileUniqueID string `json:"file_unique_id"`
		FileSize     int64  `json:"file_size"`
		FilePath     string `json:"file_path"`
	}
	type response struct {
		Ok     bool   `json:"ok"`
		Result result `json:"result"`
	}

	return func(c telebot.Context) error {
		log.Println(c.Message().IsReply())
		log.Println(c.Message())
		if c.Message().ReplyTo == nil {
			err := errors.New("invalid message")
			c.Send(err.Error())
			return err
		}
		if c.Message().ReplyTo.Voice == nil {
			err := errors.New("this message should have an audio file in reply")
			c.Send(err.Error())
			return err
		}

		requestLink := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", t.token, c.Message().ReplyTo.Voice.FileID)
		r, err := http.Get(requestLink)
		if err != nil {
			c.Send(err.Error())
			return err
		}

		res := &response{}
		if err := json.NewDecoder(r.Body).Decode(res); err != nil {
			c.Send(err.Error())
			return err
		}
		r.Body.Close()

		requestLink = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", t.token, res.Result.FilePath)

		client := http.Client{
			Timeout: time.Minute * 10,
		}
		resp, err := client.Get(requestLink)
		if err != nil {
			c.Send(err.Error())
			return err
		}
		defer r.Body.Close()

		if _, err := t.puk.Create(context.Background(), c.Chat().ID, c.Message().Sender.ID, resp.Body); err != nil {
			log.Println(err)
			return c.Send(err)
		}
		return c.Reply("uploaded!😃")
	}
}

func (t *Telegrambot) PukList(c telebot.Context) error {
	if len(c.Args()) == 0 {
		c.Send("i need exactly 2 arguments: page number, limit")
	}
	page, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return err
	}
	limit, err := strconv.Atoi(c.Args()[1])
	if err != nil {
		return err
	}
	puks, err := t.puk.List(context.Background(), page, limit)
	if err != nil {
		return err
	}
	for i := 0; i < len(puks); i++ {
		c.Send(puks[i].AudioURL)
	}
	return c.Send("here are your files with farts!😈")
}
