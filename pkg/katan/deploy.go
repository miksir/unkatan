package katan

import (
	"context"
	"fmt"
	"github.com/miksir/unkatan/pkg/announce"
	"github.com/miksir/unkatan/pkg/deploy"
	"github.com/miksir/unkatan/pkg/helpers"
	zlog "github.com/miksir/unkatan/pkg/log"
	"github.com/miksir/unkatan/pkg/schedule"
	"github.com/miksir/unkatan/pkg/storage"
	"go.uber.org/zap"
	"time"
)

var status = deploy.InitDeployStatus()
var deployRegistry storage.Registry
var deploySchedule schedule.Schedule
var logger zlog.Logger
var deployAnnounce announce.Announce

func DeployStatusInit(registry storage.Registry, announce announce.Announce, schedule schedule.Schedule, log zlog.Logger) {
	deployRegistry = registry
	deployAnnounce = announce
	deploySchedule = schedule
	logger = log
	status.SetLogCallback(func(from, to deploy.DeployCommandI) {
		logger.Warn(
			context.Background(),
			fmt.Sprintf(
				"деплой переключен с %s на %s",
				deploy.DeployCommandFullRussianName(from),
				deploy.DeployCommandFullRussianName(to),
			),
			zap.Any("from", from),
			zap.Any("to", to),
		)

		if to.IsDeployOn() != from.IsDeployOn() {
			var nextDateP *time.Time
			if to.GetUser().IsScheduler == true {
				nextDate, _, _, err := deploySchedule.FindNextEvent(to.IsDeployOn(), to.GetTime())
				if err == nil {
					nextDateP = &nextDate
				}
			}
			deployAnnounce.DeployAnnounce(to.IsDeployOn(), to.GetReason(), to.GetUser(), nextDateP)
		}

		data, err := status.DumpState()
		if err == nil {
			_ = deployRegistry.SaveKatanState(data)
		}
	})

	data, err := deployRegistry.RestoreKatanState()
	if err == nil {
		_ = status.RestoreState(data)
	}
}

func DeployStatusCommand(command deploy.DeployCommandI) {
	status.SetCommand(command)

}

func DeployHistory() []deploy.DeployCommandI {
	return status.GetHistory()
}

func DeployStatus() deploy.DeployCommandI {
	command := status.GetCommand()

	if deploy.CommandName(command) == deploy.DeployActionOn {
		// тут проверка по имени команды, а не IsDeployOn так как может быть открыто, но для определенного юзера
		now := helpers.Now()
		scheduleStatus, reason := deploySchedule.CheckSchedule(now)
		if scheduleStatus == false && now.Sub(command.GetTime()).Hours() > 1 {
			// сейчас нерабочее время но деплой открыт, закроем его через час
			command = &deploy.DeployOffCommand{
				DeployBaseCommand: deploy.DeployBaseCommand{
					Reason: schedule.RussianReasonName(reason),
					User:   deploy.ActionUser{Name: "schedule", IsScheduler: true},
				},
			}
			status.SetCommand(command)
		}
	}

	return status.GetCommand()
}

func CheckSchedule() {
	now := helpers.Now()
	scheduleStatus, reason := deploySchedule.CheckSchedule(now)
	wasStatus, _ := status.GetSchedule()

	if scheduleStatus == false {
		command := status.GetCommand()

		if deploy.CommandName(command) == deploy.DeployActionOn {
			// тут проверка по имени команды, а не IsDeployOn так как может быть открыто, но для определенного юзера
			if now.Sub(command.GetTime()).Hours() > 1 {
				// сейчас нерабочее время но деплой открыт, закроем его через час
				wasStatus = true
			}
		}
	}

	if wasStatus != scheduleStatus {

		status.SetSchedule(scheduleStatus, reason)

		if status.IsPermanent() {
			return
		}

		var command deploy.DeployCommandI
		if scheduleStatus {
			command = &deploy.DeployOnCommand{
				DeployBaseCommand: deploy.DeployBaseCommand{
					Reason: schedule.RussianReasonName(reason),
					User:   deploy.ActionUser{Name: "schedule", IsScheduler: true},
				},
			}
		} else {
			command = &deploy.DeployOffCommand{
				DeployBaseCommand: deploy.DeployBaseCommand{
					Reason: schedule.RussianReasonName(reason),
					User:   deploy.ActionUser{Name: "schedule", IsScheduler: true},
				},
			}
		}
		status.SetCommand(command)
	}
}
