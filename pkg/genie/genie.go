package genie

import (
	"encoding/json"
	"fmt"
	"github.com/miksir/unkatan/pkg/lconfig"
	zlog "github.com/miksir/unkatan/pkg/log"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strings"
	"sync"
)

var scheduleUrl = "https://api.opsgenie.com/v2/schedules/%s_schedule/on-calls?scheduleIdentifierType=name&flat=true"
var userInfoUrl = "https://api.opsgenie.com/v2/users/%s"
var userContactsUrl = "https://api.opsgenie.com/v2/users/%s/contacts"

type onCallGenieRespose struct {
	Data struct {
		OnCallRecipients []string `json:"onCallRecipients"`
	} `json:"data"`
}

type genieUserInfo struct {
	Data struct {
		FullName string `json:"fullName"`
	} `json:"data"`
}

type genieUserContacts struct {
	Data []struct {
		Method string `json:"method"`
		To     string `json:"to"`
	} `json:"data"`
}

type genie struct {
	cfg        lconfig.Reader
	log        zlog.Logger
	groupName  string
	mux        *sync.RWMutex
	email      string
	fullName   string
	voicePhone string
}

type GenieUser struct {
	FullName string
	Phone    string
}

func NewGenie(name string, cfg lconfig.Reader, logger zlog.Logger) *genie {
	genie := genie{
		cfg:       cfg,
		log:       logger,
		groupName: name,
		mux:       &sync.RWMutex{},
	}
	return &genie
}

func (g *genie) GetUser() GenieUser {
	g.mux.RLock()
	defer g.mux.RUnlock()
	return GenieUser{
		FullName: g.fullName,
		Phone:    g.voicePhone,
	}
}

func (g *genie) Update() error {
	var err error
	g.mux.Lock()
	defer g.mux.Unlock()

	g.log.Info(nil, "[GENIE] Updating duty...",
		zap.String("group", g.groupName),
		zap.String("oldEmail", g.email),
	)

	currentEmail := g.email
	err = g.updateSchedule()
	if err != nil {
		return err
	}
	if currentEmail == g.email {
		return nil
	}

	_ = g.updateInfo()
	_ = g.updateContact()

	g.log.Info(nil, "[GENIE] New duty found",
		zap.String("group", g.groupName),
		zap.String("email", g.email),
		zap.String("fullName", g.fullName),
		zap.String("phone", g.voicePhone),
	)

	return nil
}

func (g *genie) updateSchedule() error {
	url := fmt.Sprintf(scheduleUrl, url2.PathEscape(g.groupName))
	response := &onCallGenieRespose{}
	err := g.doRequest(url, &response)
	if err != nil {
		return err
	}
	if len(response.Data.OnCallRecipients) > 0 {
		g.email = response.Data.OnCallRecipients[0]
	}
	return nil
}

func (g *genie) updateInfo() error {
	g.fullName = ""
	if g.email == "" {
		return nil
	}

	url := fmt.Sprintf(userInfoUrl, url2.PathEscape(g.email))
	response := &genieUserInfo{}

	err := g.doRequest(url, &response)
	if err != nil {
		return err
	}

	g.fullName = response.Data.FullName
	return nil
}

func (g *genie) updateContact() error {
	g.voicePhone = ""
	if g.email == "" {
		return nil
	}

	url := fmt.Sprintf(userContactsUrl, url2.PathEscape(g.email))
	response := &genieUserContacts{}

	err := g.doRequest(url, &response)
	if err != nil {
		return err
	}

	for _, item := range response.Data {
		if item.Method == "voice" {
			g.voicePhone = strings.Replace(item.To, "7-", "+7", 1)
			break
		}
	}
	return nil
}

func (g *genie) doRequest(url string, data interface{}) error {
	var err error

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		g.log.Error(nil, "http.NewRequest", zap.Error(err), zap.String("url", url))
		return nil
	}
	req.Header.Add("Authorization", "GenieKey "+g.cfg.GetString("key"))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		g.log.Error(nil, "http.DefaultClient.Dot", zap.Error(err), zap.String("url", url))
		return nil
	}

	defer func() { _ = res.Body.Close() }()

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		g.log.Error(nil, "ioutil.ReadAll", zap.Error(err), zap.String("url", url), zap.Any("res", res))
		return err
	}

	err = json.Unmarshal(result, data)
	if err != nil {
		g.log.Error(nil, "json.Unmarshal", zap.Error(err), zap.String("url", url), zap.String("result", string(result)))
		return err
	}

	return nil
}
