package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"t-kt/internal/configs"
	"time"

	"github.com/icholy/digest"
)

const (
	reqBody          = `{"version": "1.0"}`
	contentType      = "applicaton/json"
	loginURL         = "http://%s/API/Web/Login"
	heartbeatURL     = "http://%s/API/Login/Heartbeat"
	recordConfGetURL = "http://%s/API/RecordConfig/Get"
	recordConfSetURL = "http://%s/API/RecordConfig/Set"
)

type RecordConfig struct {
	Result string `json:"result"`
	Data   Data   `json:"data"`
}

type Data struct {
	ChannelInfo ChannelInfo `json:"channel_info"`
}

type ChannelInfo struct {
	Channel Channel `json:"CH1"`
}

type Channel struct {
	RecordSwitch   bool   `json:"record_switch"`
	StreamMode     string `json:"stream_mode"`
	Prerecord      bool   `json:"prerecord"`
	NetBreakRecord bool   `json:"net_break_record"`
}

func NewArchiveSwitcher(statusChecker *bool, conf configs.IPCConf) (func(ctx context.Context) CmdResult, error) {
	sleepDuration := conf.ArchiveSwitchTime
	username, password, addr, _ := conf.GetDNSConf()
	// if err != nil {
	// 	return nil, err
	// }
	return func(ctx context.Context) CmdResult {
		c, _ := cookiejar.New(nil)
		client := &http.Client{
			Transport: &digest.Transport{
				Username: username,
				Password: password,
			},
			Jar:     c,
			Timeout: 5 * time.Second,
		}

		reader := strings.NewReader(reqBody)

		resp, err := client.Post(fmt.Sprintf(loginURL, addr), "", reader)
		if err != nil {
			return CmdResult{err: err}
		}
		cookie := resp.Cookies()
		csrfToken := resp.Header.Get("X-csrftoken")

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					req, err := newRequest(client, "POST", fmt.Sprintf(heartbeatURL, addr), cookie, csrfToken, reader)
					if err != nil {
						return
					}
					time.Sleep(5 * time.Second)
					client.Do(req)
				}
			}
		}()

		recordStatus := true
		data := RecordConfig{}
		req, err := newRequest(client, "POST", fmt.Sprintf(recordConfGetURL, addr), cookie, csrfToken, reader)
		if err != nil {
			return CmdResult{err: err}
		}
		resp, _ = client.Do(req)
		if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return CmdResult{err: err}
		}
		data.Data.ChannelInfo.Channel.Prerecord = false

		for {
			select {
			case <-ctx.Done():
				return CmdResult{info: "disable"}
			default:
				data.Data.ChannelInfo.Channel.RecordSwitch = recordStatus
				recordStatus = !recordStatus

				buf, err := json.Marshal(data)
				if err != nil {
					return CmdResult{err: err}
				}
				reader := bytes.NewReader(buf)
				req, err := newRequest(client, "POST", fmt.Sprintf(recordConfSetURL, addr), cookie, csrfToken, reader)
				if err != nil {
					return CmdResult{err: err}
				}
				_, err = client.Do(req)
				if err != nil {
					return CmdResult{err: err}
				}
				if sleepDuration == 0 {
					time.Sleep(time.Duration(rand.IntN(30)) * time.Second)
				} else {
					time.Sleep(time.Duration(sleepDuration) * time.Second)
				}
			}
		}
	}, nil
}

func newRequest(client *http.Client, method string, url string, cookie []*http.Cookie, CSRFToken string, reader io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-csrftoken", CSRFToken)
	req.Header.Set("Content-Type", contentType)
	client.Jar.SetCookies(req.URL, cookie)
	return req, nil
}
