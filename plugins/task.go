package example

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/FloatTech/zbputils/control"
	"github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var tasksFile = "tasks.txt"

func init() {
	engine := control.Register("task_plugin", &control.Options{
		Brief: "任务管理插件",
		Help:  "用法：添加任务[任务名称]、移除任务[任务名称]、展示任务",
	})

	engine.OnCommand("添加任务").Handle(addTask)
	engine.OnCommand("移除任务").Handle(removeTask)
	engine.OnCommand("展示任务").Handle(showTasks)

	// 读取已保存的任务列表
	loadTasks()
}

func addTask(ctx *ZeroBot.Ctx) {
	taskName := extractTaskName(ctx.RawMessage)
	if taskName != "" {
		tasks = append(tasks, taskName)
		saveTasks() // 添加任务后保存到文件
		ctx.Send(fmt.Sprintf("任务'%s'添加成功！", taskName))
	} else {
		ctx.Send("任务名称不能为空！")
	}
}

func removeTask(ctx *ZeroBot.Ctx) {
	taskName := extractTaskName(ctx.RawMessage)
	if taskName != "" {
		for i, task := range tasks {
			if task == taskName {
				tasks = append(tasks[:i], tasks[i+1:]...)
				saveTasks() // 移除任务后保存到文件
				ctx.Send(fmt.Sprintf("任务'%s'移除成功！", taskName))
				return
			}
		}
		ctx.Send(fmt.Sprintf("未找到名称为'%s'的任务！", taskName))
	} else {
		ctx.Send("任务名称不能为空！")
	}
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
		fmt.Println("无法保存任务列表：", err)
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
		fmt.Println("无法加载任务列表：", err)
	}
}
