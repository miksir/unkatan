package deploy

import (
	"encoding/json"
	"github.com/miksir/unkatan/pkg/helpers"
	zlog "github.com/miksir/unkatan/pkg/log"
	"github.com/miksir/unkatan/pkg/schedule"
	"sync"
	"time"
)

type deployStatus struct {
	command  DeployCommandI
	schedule struct {
		lastStatus bool
		lastReason schedule.TriggerType
	}
	history     *deployHistory
	mux         sync.RWMutex
	logger      zlog.Logger
	logCallback func(from, to DeployCommandI)
}

type saveStruct struct {
	Command     json.RawMessage
	CommandName string
	Schedule    struct {
		LastStatus bool
		LastReason schedule.TriggerType
	}
	History json.RawMessage
}

func InitDeployStatus() *deployStatus {
	st := deployStatus{}
	st.mux = sync.RWMutex{}
	command := DeployOffCommand{}
	command.User = ActionUser{Name: "unknown"}
	st.command = &command
	st.history = initDeployHistory(20)
	return &st
}

func (s *deployStatus) SetLogCallback(callback func(from, to DeployCommandI)) {
	s.logCallback = callback
}

func (s *deployStatus) RestoreState(data []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	saveStruct := saveStruct{}
	err := json.Unmarshal(data, &saveStruct)
	if err != nil {
		return err
	}

	var dc DeployCommandI
	if saveStruct.CommandName == DeployActionOff {
		dc = &DeployOffCommand{}
	} else if saveStruct.CommandName == DeployActionOn {
		dc = &DeployOnCommand{}
	}
	err = json.Unmarshal(saveStruct.Command, dc)
	if err != nil {
		return err
	}

	s.command = dc
	s.schedule.lastReason = saveStruct.Schedule.LastReason
	s.schedule.lastStatus = saveStruct.Schedule.LastStatus

	_ = s.history.restoreState(saveStruct.History)
	return nil
}

func (s *deployStatus) DumpState() ([]byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	commandJson, err := json.Marshal(s.command)
	if err != nil {
		return nil, err
	}

	var saveStruct = saveStruct{
		Command:     commandJson,
		CommandName: CommandName(s.command),
		Schedule: struct {
			LastStatus bool
			LastReason schedule.TriggerType
		}{LastStatus: s.schedule.lastStatus, LastReason: s.schedule.lastReason},
	}
	saveStruct.History, err = s.history.saveState()
	if err != nil {
		return nil, err
	}

	resp, err := json.Marshal(saveStruct)
	return resp, err
}

func (s *deployStatus) GetHistory() []DeployCommandI {
	return s.history.GetList()
}

func (s *deployStatus) SetCommand(command DeployCommandI) {
	s.mux.Lock()

	timeNow := helpers.Now()
	command.SetTime(timeNow)

	s.history.pushCommand(command)

	oldCmd := s.command
	s.command = command

	s.mux.Unlock()

	if s.logCallback != nil {
		s.logCallback(oldCmd, command)
	}
}

func (s *deployStatus) GetCommand() DeployCommandI {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.command
}

func (s *deployStatus) IsOpen(requestUsers ...string) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.command.IsDeployOn(requestUsers...)
}

func (s *deployStatus) SetSchedule(status bool, reason schedule.TriggerType) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.schedule.lastStatus = status
	s.schedule.lastReason = reason
}

func (s *deployStatus) GetSchedule() (status bool, reason schedule.TriggerType) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	status = s.schedule.lastStatus
	reason = s.schedule.lastReason
	return
}

func (s *deployStatus) IsPermanent() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	switch s.command.(type) {
	case *DeployOffCommand:
		return s.command.(*DeployOffCommand).Permanent
	default:
		return false
	}
}

func (s *deployStatus) getUser() ActionUser {
	s.mux.RLock()
	defer s.mux.RUnlock()
	user := s.command.GetUser()
	return user
}

func (s *deployStatus) NonEmptyReason() string {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if s.command.GetReason() == "" {
		return "не указан"
	}
	return s.command.GetReason()
}

func (s *deployStatus) GetTime() time.Time {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.command.GetTime()
}
