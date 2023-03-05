package deploy

import "testing"
import "github.com/stretchr/testify/assert"

func TestDeployBaseCommand_GetName(t *testing.T) {
	a := assert.New(t)
	onCommand := &DeployOnCommand{}
	a.Equal(DeployActionOn, CommandName(onCommand))
	offCommand := &DeployOffCommand{}
	a.Equal(DeployActionOff, CommandName(offCommand))
}
