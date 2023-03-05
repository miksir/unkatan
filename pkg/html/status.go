package html

import (
	"github.com/miksir/unkatan/pkg/genie"
	"github.com/miksir/unkatan/pkg/helpers"
	"github.com/miksir/unkatan/pkg/katan"
	"log"
	"net/http"
	"time"
)
import "html/template"

type Html struct {
	statusT *template.Template
}

type statusTemplateData struct {
	DeployStatusName   string
	DeployChangedBy    string
	DeployChangeReason string
	DeployHistory      []historyStatusTemplateData
	SreDuty            genie.GenieUser
	SupportDuty        genie.GenieUser
}

type historyStatusTemplateData struct {
	When         time.Time
	StatusName   string
	ChangedBy    string
	ChangeReason string
}

func NewHtml() *Html {
	var err error
	html := Html{}
	html.statusT, err = template.ParseFiles("html/status.html")
	if err != nil {
		log.Fatal(err)
	}
	return &html
}

func (html *Html) HandleMain(w http.ResponseWriter, _ *http.Request) {
	deployCmd := katan.DeployStatus()
	history := katan.DeployHistory()
	historyTemplate := make([]historyStatusTemplateData, 0, len(history))

	for _, item := range history {
		historyTemplate = append(historyTemplate, historyStatusTemplateData{
			When:         item.GetTime(),
			StatusName:   helpers.DeployStatusRussianName(item.IsDeployOn()),
			ChangedBy:    item.GetUser().PlainName(),
			ChangeReason: item.GetReason(),
		})
	}

	data := statusTemplateData{
		DeployStatusName:   helpers.DeployStatusRussianName(deployCmd.IsDeployOn()),
		DeployChangedBy:    deployCmd.GetUser().PlainName(),
		DeployChangeReason: deployCmd.GetReason(),
		DeployHistory:      historyTemplate,
		SreDuty:            genie.GetGenieUser(genie.GenieSRETeam),
		SupportDuty:        genie.GetGenieUser(genie.GenieSupportTeam),
	}

	err := html.statusT.Execute(w, data)
	if err != nil {
		return
	}
}
