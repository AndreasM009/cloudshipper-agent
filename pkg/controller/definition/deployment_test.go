package definition

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var yamlDefinition = `
jobs:
  testing:
    displayname: MySteps
    working-directory: ./test
    steps:
      - command: AzPowerShellCore
        working-directory: ./test
        displayname: MyStep
        with:
          arguments: -ResourceGroup myRG
          scriptToRun: ./myscript.ps1
parameters: {}
`

func TestDefinitionFromYaml(t *testing.T) {
	definition, err := NewFromYaml([]byte(yamlDefinition))
	assert.Nil(t, err)
	assert.NotNil(t, definition)

	job := definition.Jobs["testing"]
	assert.NotNil(t, job)

	assert.Equal(t, 1, len(job.Steps))

	step := job.Steps[0]
	assert.Equal(t, strings.ToLower("AzPowerShellCore"), strings.ToLower(step.Command))
	assert.Equal(t, strings.ToLower("MyStep"), strings.ToLower(step.Displayname))

	assert.NotNil(t, step.With)
	assert.Equal(t, "-ResourceGroup myRG", step.With["arguments"])
	assert.Equal(t, "./myscript.ps1", step.With["scriptToRun"])
}
