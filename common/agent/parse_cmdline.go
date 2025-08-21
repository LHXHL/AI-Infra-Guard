package agent

import (
	"encoding/json"
	"fmt"
<<<<<<< HEAD
	"github.com/google/uuid"
	"os"

	"github.com/Tencent/AI-Infra-Guard/internal/gologger"
=======
	"github.com/Tencent/AI-Infra-Guard/internal/gologger"
	"github.com/google/uuid"
	"path"
>>>>>>> opensource
)

type CmdNewPlanStep struct {
	Title  string `json:"title"`
	StepId string `json:"stepId"`
}

type CmdStatusUpdate struct {
	Brief       string `json:"brief"`
	Description string `json:"description"`
	StepId      string `json:"stepId"`
<<<<<<< HEAD
=======
	Status      string `json:"status"`
>>>>>>> opensource
}

type CmdToolUsed struct {
	ToolId   string `json:"tool_id"`
	ToolName string `json:"tool_name"`
	Brief    string `json:"brief"`
	Status   string `json:"status"`
	StepId   string `json:"stepId"`
}

type CmdActionLog struct {
	ToolId   string `json:"tool_id"`
	ToolName string `json:"tool_name"`
	Log      string `json:"log"`
	StepId   string `json:"stepId"`
}

type CmdContent struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

<<<<<<< HEAD
var statusId string

func ParseStdoutLine(tasks []SubTask, line string, callbacks TaskCallbacks) {
=======
type CmdConfig struct {
	StatusId string
	Status   string
}

type PromptContent struct {
	Results    []PromptResults `json:"results"`
	Total      int             `json:"total"`
	Score      float32         `json:"score"`
	Attachment string          `json:"attachment"`
	Jailbreak  int             `json:"jailbreak"`
}
type PromptResults struct {
	Status        string `json:"status"`
	ModelName     string `json:"modelName"`
	Vulnerability string `json:"vulnerability"`
	AttackMethod  string `json:"attackMethod"`
	Input         string `json:"input"`
	Output        string `json:"output"`
	Reason        string `json:"reason"`
}

func ParseStdoutLine(server, rootDir string, tasks []SubTask, line string, callbacks TaskCallbacks, config *CmdConfig) {
>>>>>>> opensource
	var cmd CmdContent
	if err := json.Unmarshal([]byte(line), &cmd); err != nil {
		fmt.Println(line)
		return
	}
	switch cmd.Type {
	case AgentMsgTypeNewPlanStep:
		var content CmdNewPlanStep
		if err := json.Unmarshal(cmd.Content, &content); err != nil {
			gologger.WithError(err).Errorln("Failed to AgentMsgTypeNewPlanStep unmarshal command", cmd.Content)
			return
		}
		callbacks.NewPlanStepCallback(content.StepId, content.Title)
		// 更新任务状态
		for i, _ := range tasks {
			if tasks[i].StepId < content.StepId {
				tasks[i].Status = SubTaskStatusDone
			} else if tasks[i].StepId == content.StepId {
				tasks[i].Status = SubTaskStatusDoing
			}
		}
		callbacks.PlanUpdateCallback(tasks)
	case AgentMsgTypeStatusUpdate:
		var content CmdStatusUpdate
		if err := json.Unmarshal(cmd.Content, &content); err != nil {
			gologger.WithError(err).Errorln("Failed to AgentMsgTypeStatusUpdate unmarshal command", cmd.Content)
			return
		}
<<<<<<< HEAD
		statusId = uuid.NewString()
		callbacks.StepStatusUpdateCallback(content.StepId, statusId, AgentStatusCompleted, content.Brief, content.Description)
=======
		if content.Status == AgentStatusRunning {
			config.StatusId = uuid.NewString()
			config.Status = "running"
		} else if content.Status == AgentStatusCompleted {
			if config.Status == "completed" {
				config.StatusId = uuid.NewString()
			}
			config.Status = "completed"
		}
		callbacks.StepStatusUpdateCallback(content.StepId, config.StatusId, content.Status, content.Brief, content.Description)
>>>>>>> opensource
	case AgentMsgTypeToolUsed:
		var content CmdToolUsed
		if err := json.Unmarshal(cmd.Content, &content); err != nil {
			gologger.WithError(err).Errorln("Failed to AgentMsgTypeToolUsed unmarshal command", cmd.Content)
			return
		}
		tool := CreateTool(content.ToolId, content.ToolId, statusString(content.Status), content.Brief, content.Brief, "", "")
<<<<<<< HEAD
		callbacks.ToolUsedCallback(content.StepId, statusId, content.Brief, []Tool{tool})
=======
		callbacks.ToolUsedCallback(content.StepId, config.StatusId, content.Brief, []Tool{tool})
>>>>>>> opensource
	case AgentMsgTypeActionLog:
		var content CmdActionLog
		if err := json.Unmarshal(cmd.Content, &content); err != nil {
			gologger.WithError(err).Errorln("Failed to AgentMsgTypeActionLog unmarshal command", cmd.Content)
			return
		}
		callbacks.ToolUseLogCallback(content.ToolId, content.ToolName, content.StepId, content.Log)
	case AgentMsgTypeResultUpdate:
		var content map[string]interface{}
		if err := json.Unmarshal(cmd.Content, &content); err != nil {
			gologger.WithError(err).Errorln("Failed to AgentMsgTypeResultUpdate unmarshal command", cmd.Content)
			return
		}
		for i, _ := range tasks {
			tasks[i].Status = SubTaskStatusDone
		}
		callbacks.PlanUpdateCallback(tasks)
<<<<<<< HEAD

		content["msgType"] = "markdown"
		filename := content["content"].(string)
		data, err := os.ReadFile(filename)
		if err != nil {
			gologger.WithError(err).Errorln("Failed to read result file", filename)
			return
		}
		content["content"] = string(data)
		callbacks.ResultCallback(content)
=======
		var ret []PromptContent
		dd, err := json.Marshal(content["content"])
		if err != nil {
			gologger.WithError(err).Errorln("Failed to parse result file json")
			return
		}
		if err := json.Unmarshal(dd, &ret); err != nil {
			gologger.WithError(err).Errorln("Failed to parse result file")
			return
		}
		gologger.Infoln("开始上传文件")
		for i, v := range ret {
			if v.Attachment == "" {
				continue
			}
			info, err := UploadFile(server, path.Join(rootDir, v.Attachment))
			if err != nil {
				gologger.WithError(err).Errorln("Failed to upload file")
				return
			}
			gologger.Infoln("上传文件成功")
			v.Attachment = info.Data.FileUrl
			ret[i] = v
		}
		dd, err = json.Marshal(ret)
		if err != nil {
			gologger.WithError(err).Errorln("Failed to parse result file json")
			return
		}
		var content2 []map[string]interface{}
		if err := json.Unmarshal(dd, &content2); err != nil {
			gologger.WithError(err).Errorln("Failed to parse result file json")
			return
		}
		content["content"] = content2
		callbacks.ResultCallback(content)
	case AgentMsgTypeError:
		content := string(cmd.Content)
		callbacks.ErrorCallback(content)
>>>>>>> opensource
	}
}
