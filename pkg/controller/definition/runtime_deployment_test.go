package definition

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var yamlDeploymentDefinition = `
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
          subscription: azsubscription
          tenant: tenant.onmicrosoft.com
          serviceprincipal: spname
          secret: secret
parameters: {}
`

func TestRuntimeDefinitionFromDeployment(t *testing.T) {
	definition, err := NewFromYaml([]byte(yamlDeploymentDefinition))
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

	var definitionParams = make(map[string]string)
	definitionParams["param1"] = "value1"
	definitionParams["param2"] = "value2"

	runtimeDefinition, err := NewFromDefinition(definition, "1", "test", "1", definitionParams)

	assert.Nil(t, err)
	assert.NotNil(t, runtimeDefinition)
}
