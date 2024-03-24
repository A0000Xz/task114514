package example

import (
	"fmt"
	"io/ioutil"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var tasks []string
var tasksFile = "tasks.txt"

func init() {
	engine := control.Register("task_plugin", &control.Options{
		Brief: "任务管理插件",
		Help:  "用法：添加任务[任务名称]、移除任务[任务名称]、展示任务",
	})

	engine.OnCommand("添加任务").Handle(addTask)
	engine.OnCommand("移除任务").Handle(removeTask)
	engine.OnCommand("展示任务").Handle(showTasks)

	loadTasks()
}

func addTask(ctx *ZeroBot.Ctx) {
	taskName := extractTaskName(ctx.RawMessage)
	if taskName == "" {
		ctx.Send("任务名称不能为空！")
		return
	}

	tasks = append(tasks, taskName)
	saveTasks()
	ctx.Send(fmt.Sprintf("任务'%s'添加成功！", taskName))
}

func removeTask(ctx *ZeroBot.Ctx) {
	taskName := extractTaskName(ctx.RawMessage)
	if taskName == "" {
		ctx.Send("任务名称不能为空！")
		return
	}

	for i, task := range tasks {
		if task == taskName {
			tasks = append(tasks[:i], tasks[i+1:]...)
			saveTasks()
			ctx.Send(fmt.Sprintf("任务'%s'移除成功！", taskName))
			return
		}
	}
	ctx.Send(fmt.Sprintf("未找到名称为'%s'的任务！", taskName))
}

func showTasks(ctx *ZeroBot.Ctx) {
	if len(tasks) > 0 {
		taskList := strings.Join(tasks, "\n")
		ctx.Send("当前任务列表：\n" + taskList)
	} else {
		ctx.Send("当前无任务！")
	}
}

func extractTaskName(message string) string {
	startIndex := strings.Index(message, "[")
	endIndex := strings.Index(message, "]")
	if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
		return message[startIndex+1 : endIndex]
	}
	return ""
}

func saveTasks() {
	taskData := strings.Join(tasks, "\n")
	err := ioutil.WriteFile(tasksFile, []byte(taskData), 0644)
	if err != nil {
		ctx.Send(fmt.Sprintf("无法保存任务列表：%v", err))
	}
}

func loadTasks() {
	taskData, err := ioutil.ReadFile(tasksFile)
	if err == nil {
		taskList := strings.Split(string(taskData), "\n")
		tasks = make([]string, 0, len(taskList))
		for _, task := range taskList {
			if task != "" {
				tasks = append(tasks, task)
			}
		}
	} else {
		ctx.Send(fmt.Sprintf("无法加载任务列表：%v", err))
	}
}
