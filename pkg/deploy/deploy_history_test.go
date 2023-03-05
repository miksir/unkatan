package deploy

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	time1 = time.Now()
	time2 = time1.Add(1 * time.Hour)
	time3 = time2.Add(1 * time.Hour)
	time4 = time3.Add(1 * time.Hour)
)

func Test_deployHistory_GetList(t *testing.T) {
	a := assert.New(t)
	dh := preparedHistory()
	dh.pushCommand(&DeployOnCommand{
		DeployBaseCommand: DeployBaseCommand{
			Time: time4,
		},
	})
	listH := dh.GetList()
	a.Equal(3, len(listH))
	a.Equal(time2, listH[0].GetTime())
	a.Equal(time3, listH[1].GetTime())
	a.Equal(time4, listH[2].GetTime())
}

func Test_deployHistory_saveState(t *testing.T) {
	a := assert.New(t)
	dh := preparedHistory()
	data, err := dh.saveState()
	a.Nil(err)

	dh2 := initDeployHistory(3)
	err = dh2.restoreState(data)
	a.Nil(err)

	a.True(cmp.Equal(dh.GetList(), dh2.GetList()))
}

func preparedHistory() *deployHistory {
	dh := initDeployHistory(3)
	dh.pushCommand(&DeployOffCommand{
		DeployBaseCommand: DeployBaseCommand{
			Time: time1,
		},
	})
	dh.pushCommand(&DeployOnCommand{
		DeployBaseCommand: DeployBaseCommand{
			Time: time2,
		},
	})
	dh.pushCommand(&DeployOffCommand{
		DeployBaseCommand: DeployBaseCommand{
			Time: time3,
		},
	})
	return dh
}
