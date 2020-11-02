package definition

import (
	"github.com/andreasM009/cloudshipper-agent/pkg/commands"
)

// RuntimeCommand runtime command
type RuntimeCommand struct {
	Command     interface{}
	CommandType commands.CommandType
	Index       int
	Name        string
	Displayname string
}

// RuntimeJob job to run
type RuntimeJob struct {
	Name             string
	Displayname      string
	WorkingDirectory string
	Commands         []*RuntimeCommand
}

// RuntimeDeployment runtime deployment
type RuntimeDeployment struct {
	DeploymentName string
	ID             string
	TenantID       string
	Parameters     map[string]string
	Jobs           []*RuntimeJob
	DefinitionID   string
}

// NewFromDefinition new from definition
func NewFromDefinition(deployment *Deployment, tenantID, definitionID, name, id string, parameters map[string]string) (*RuntimeDeployment, error) {
	rd := &RuntimeDeployment{
		DeploymentName: name,
		ID:             id,
		TenantID:       tenantID,
		Parameters:     parameters,
		Jobs:           make([]*RuntimeJob, len(deployment.Jobs)),
		DefinitionID:   definitionID,
	}

	jobIdx := 0

	for key, elem := range deployment.Jobs {
		rj := &RuntimeJob{
			Name:             key,
			Displayname:      elem.Displayname,
			WorkingDirectory: elem.WorkingDirectory,
		}
		rd.Jobs[jobIdx] = rj
		jobIdx++

		rj.Commands = make([]*RuntimeCommand, len(elem.Steps))

		for idx, step := range elem.Steps {
			cmd := &RuntimeCommand{
				Index:       idx,
				Displayname: step.Displayname,
			}

			// WorkingDirectory from step or job?
			wrkDir := step.WorkingDirectory
			if wrkDir == "" {
				wrkDir = rj.WorkingDirectory
			}

			inherided := make(map[string]string)
			inherided["workingdirectory"] = wrkDir

			cmdType := commands.Parse(step.Command)
			cmd.CommandType = cmdType
			cmd.Name = cmdType.String()

			acmd, err := commands.GetRegistry().CreateFromProps(cmd.CommandType, step.With, inherided)
			if err != nil {
				return nil, err
			}

			cmd.Command = acmd
			rj.Commands[idx] = cmd
		}
	}

	return rd, nil
}
