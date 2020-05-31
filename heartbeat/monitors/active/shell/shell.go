package shell

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/elastic/beats/v7/heartbeat/eventext"
	"github.com/elastic/beats/v7/heartbeat/monitors"
	"github.com/elastic/beats/v7/heartbeat/monitors/jobs"
	"github.com/elastic/beats/v7/heartbeat/monitors/wrappers"
	"github.com/elastic/beats/v7/heartbeat/reason"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
)

func init() {
	monitors.RegisterActive("shell", create)
}

var debugf = logp.MakeDebug("shell")

func create(
	name string,
	cfg *common.Config,
) (js []jobs.Job, endpoints int, err error) {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, 0, err
	}

	validator := makeValidator(&config)

	jobs := make([]jobs.Job, len(config.Commands))

	for i, command := range config.Commands {
		jobs[i], err = newShellMonitorJob(command, &config, validator)
		if err != nil {
			return nil, 0, err
		}
	}

	return jobs, len(config.Commands), nil
}

func newShellMonitorJob(
	command string,
	config *Config,
	validator OutputCheck,
) (jobs.Job, error) {

	okstr := ""
	for _, ok := range config.Check.Response.Ok {
		okstr = okstr + ok.String() + ","
	}

	criticalStr := ""
	for _, critical := range config.Check.Response.Critical {
		criticalStr = criticalStr + critical.String() + " ,"
	}

	eventFields := common.MapStr{
		"monitor": common.MapStr{
			"scheme":  "shell",
			"command": command,
		},
		"check": common.MapStr{
			"ok":       okstr,
			"critical": criticalStr,
		},
	}

	customs := config.CustomeFields
	if len(customs) != 0 {
		customMap := common.MapStr{}
		for _, v := range customs {
			splitPos := strings.Index(v, ":")
			if splitPos > 0 && splitPos != len(v)-1 {
				customMap[string(v[0:splitPos])] = string(v[splitPos+1:])
			}
		}
		if len(customMap) > 0 {
			eventFields["custom"] = customMap
		}
	}

	return wrappers.WithFields(eventFields,
		jobs.MakeSimpleJob(func(event *beat.Event) error {
			err := execute(command, event, config.Timeout, validator)
			return err
		})), nil
}

func execute(command string, event *beat.Event, duration time.Duration, validate func(string) error) (errReason reason.Reason) {

	execution := ExecutionRequest{
		Command: command,
		Env:     nil,
		Timeout: duration,
	}

	cmdExec, err := execution.Execute(context.Background(), execution)
	outputStr := cmdExec.Output
	exitStatus := cmdExec.Status

	// ctx, cancel := context.WithTimeout(context.Background(), duration)
	// defer cancel()

	// cmd := exec.CommandContext(ctx, command)
	// outputByte, err := cmd.Output()
	// outputStr := string(outputByte)

	eventext.MergeEventFields(event, common.MapStr{"shell": common.MapStr{
		"response": common.MapStr{
			"output": outputStr,
		},
	}})

	if err != nil {
		return reason.IOFailed(err)
	}

	if exitStatus == 2 {
		var ErrTimeOut = errors.New("Execution timed out")
		return reason.IOFailed(ErrTimeOut)
	}

	if exitStatus == 3 {
		var ErrFallback = errors.New("Execution fallback")
		return reason.IOFailed(ErrFallback)
	}

	err = validate(outputStr)

	return reason.ValidateFailed(err)
}
