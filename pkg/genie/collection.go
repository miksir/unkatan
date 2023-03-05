package genie

import (
	"github.com/miksir/unkatan/pkg/lconfig"
	zlog "github.com/miksir/unkatan/pkg/log"
)

const (
	GenieSRETeam     = "CityMobil Main SRE"
	GenieSupportTeam = "CityMobil Support Team"
)

var collection map[string]*genie

func InitCollection(cfg lconfig.Reader, logger zlog.Logger) {
	collection = make(map[string]*genie)
	collection[GenieSRETeam] = NewGenie(GenieSRETeam, cfg, logger)
	collection[GenieSupportTeam] = NewGenie(GenieSupportTeam, cfg, logger)
}

func GetGenieUser(name string) GenieUser {
	genie, found := collection[name]
	if !found {
		return GenieUser{}
	}
	return genie.GetUser()
}

func UpdateCollection() {
	for _, genie := range collection {
		_ = genie.Update()
	}
}
