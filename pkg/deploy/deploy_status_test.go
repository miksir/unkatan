package deploy

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)
import "github.com/stretchr/testify/assert"

func TestDeployStatusInit(t *testing.T) {
	a := assert.New(t)

	status := InitDeployStatus()

	a.IsType(&DeployOffCommand{}, status.GetCommand())
	a.False(status.IsOpen())
	a.Equal(0, len(status.GetHistory()))
}

func TestDeployStatusOn(t *testing.T) {
	a := assert.New(t)

	status := InitDeployStatus()
	status.SetCommand(&DeployOnCommand{
		OnlyFor: nil,
		DeployBaseCommand: DeployBaseCommand{
			Reason: "reason for test",
			User: ActionUser{
				Name: "Test User",
			},
		},
	})

	a.IsType(&DeployOnCommand{}, status.GetCommand())
	a.True(status.IsOpen())
	a.Equal("reason for test", status.NonEmptyReason())
	a.Equal("Test User", status.getUser().Name)
	a.False(status.IsPermanent())
	a.True(time.Now().After(status.GetTime()))
	a.Equal(1, len(status.GetHistory()))
}

func TestDeployStatusOff(t *testing.T) {
	a := assert.New(t)

	status := InitDeployStatus()
	status.SetCommand(&DeployOffCommand{
		DeployBaseCommand: DeployBaseCommand{
			Reason: "reason for test",
			User: ActionUser{
				Name: "Test User",
			},
		},
	})

	a.IsType(&DeployOffCommand{}, status.GetCommand())
	a.False(status.IsOpen())
	a.Equal("reason for test", status.NonEmptyReason())
	a.Equal("Test User", status.getUser().Name)
	a.False(status.IsPermanent())
	a.True(time.Now().After(status.GetTime()))
	a.Equal(1, len(status.GetHistory()))
}

func TestDeployStatusOffPermanent(t *testing.T) {
	a := assert.New(t)

	status := InitDeployStatus()
	status.SetCommand(&DeployOffCommand{
		Permanent: true,
		DeployBaseCommand: DeployBaseCommand{
			Reason: "reason for test",
			User: ActionUser{
				Name: "Test User",
			},
		},
	})

	a.IsType(&DeployOffCommand{}, status.GetCommand())
	a.False(status.IsOpen())
	a.Equal("reason for test", status.NonEmptyReason())
	a.Equal("Test User", status.getUser().Name)
	a.True(status.IsPermanent())
	a.True(time.Now().After(status.GetTime()))
	a.Equal(1, len(status.GetHistory()))
}

func TestDeployStatusOnOffSaveState(t *testing.T) {
	a := assert.New(t)

	command := &DeployOnCommand{
		DeployBaseCommand: DeployBaseCommand{
			Reason: "reason for test",
			User: ActionUser{
				SlackId: "ABCDEFG",
				Name:    "Test User",
				Email:   "test@test.ru",
			},
		},
	}

	command2 := &DeployOffCommand{
		Permanent: true,
		DeployBaseCommand: DeployBaseCommand{
			Reason: "reason for test",
			User: ActionUser{
				SlackId: "ABCDEFG",
				Name:    "Test User",
				Email:   "test@test.ru",
			},
		},
	}

	status := InitDeployStatus()
	status.SetCommand(command)
	status.SetCommand(command2)

	state, err := status.DumpState()
	a.Nil(err)

	status2 := InitDeployStatus()
	err = status2.RestoreState(state)
	a.Nil(err)
	state2, err := status.DumpState()
	a.Nil(err)

	a.True(cmp.Equal(state, state2))
}

func TestDeployStatusOnOnlyFor(t *testing.T) {
	a := assert.New(t)

	status := InitDeployStatus()
	status.SetCommand(&DeployOnCommand{
		OnlyFor:           []string{"vasya", "petya"},
		DeployBaseCommand: DeployBaseCommand{},
	})

	a.IsType(&DeployOnCommand{}, status.GetCommand())
	a.False(status.IsOpen())
	a.False(status.IsOpen("ivan"))
	a.True(status.IsOpen("vasya"))
	a.True(status.IsOpen("petya"))
}

func TestDeployStatusTestCallback(t *testing.T) {
	a := assert.New(t)

	c := make(chan struct {
		from DeployCommandI
		to   DeployCommandI
	}, 2)
	defer close(c)

	status := InitDeployStatus()
	status.SetLogCallback(func(from, to DeployCommandI) {
		c <- struct {
			from DeployCommandI
			to   DeployCommandI
		}{from: from, to: to}
	})

	status.SetCommand(&DeployOnCommand{
		OnlyFor: nil,
		DeployBaseCommand: DeployBaseCommand{
			Reason: "A",
		},
	})
	status.SetCommand(&DeployOffCommand{
		DeployBaseCommand: DeployBaseCommand{
			Reason: "B",
		},
	})

	select {
	case t1 := <-c:
		a.IsType(&DeployOffCommand{}, t1.from)
		a.IsType(&DeployOnCommand{}, t1.to)
		a.Equal("A", t1.to.GetReason())
	case <-time.After(500 * time.Millisecond):
		a.FailNow("callback wait timeout")
	}

	select {
	case t2 := <-c:
		a.IsType(&DeployOnCommand{}, t2.from)
		a.Equal("A", t2.from.GetReason())
		a.IsType(&DeployOffCommand{}, t2.to)
		a.Equal("B", t2.to.GetReason())
	case <-time.After(500 * time.Millisecond):
		a.FailNow("callback wait timeout")
	}
}
