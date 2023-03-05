package slack_commands

import (
	"github.com/miksir/unkatan/pkg/deploy"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKatanDeployCommand(t *testing.T) {
	tests := []struct {
		commandString   string
		expectedCommand deploy.DeployCommandI
		expectedErr     bool
	}{
		{
			"off test message 1",
			&deploy.DeployOffCommand{
				Permanent: false,
				DeployBaseCommand: deploy.DeployBaseCommand{
					Reason: "test message 1",
				},
			},
			false,
		},
		{
			"off --permanent test message 2",
			&deploy.DeployOffCommand{
				Permanent: true,
				DeployBaseCommand: deploy.DeployBaseCommand{
					Reason: "test message 2",
				},
			},
			false,
		},
		{
			"on test message 3",
			&deploy.DeployOnCommand{
				OnlyFor: []string{},
				DeployBaseCommand: deploy.DeployBaseCommand{
					Reason: "test message 3",
				},
			},
			false,
		},
		{
			"on --only-for=@d.kelmi test message 4",
			&deploy.DeployOnCommand{
				OnlyFor: []string{"d.kelmi"},
				DeployBaseCommand: deploy.DeployBaseCommand{
					Reason: "test message 4",
				},
			},
			false,
		},
		{
			"on -u @d.kelmi,@ivan test message 5",
			&deploy.DeployOnCommand{
				OnlyFor: []string{"d.kelmi", "ivan"},
				DeployBaseCommand: deploy.DeployBaseCommand{
					Reason: "test message 5",
				},
			},
			false,
		},
		{
			"on -u @d.kelmi -u @ivan test message 6",
			&deploy.DeployOnCommand{
				OnlyFor: []string{"d.kelmi", "ivan"},
				DeployBaseCommand: deploy.DeployBaseCommand{
					Reason: "test message 6",
				},
			},
			false,
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.commandString, func(t *testing.T) {
			command, err := processKatanDeployCommand(slack.SlashCommand{
				Text: "deploy " + tt.commandString,
			})
			if (err != nil) != tt.expectedErr {
				t.Errorf("waited for error = %t, got error = %v, ", tt.expectedErr, err)
				return
			}
			if command != nil {
				command.SetTime(tt.expectedCommand.GetTime())
			}
			a.Equal(tt.expectedCommand, command)
		})
	}
}
