ackage task

import (
	"fmt"
	"ioutil"
	"strings"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var tasks []string
var tasksFile = "tasks.txt"

// 自定义限制函数, 括号内填入(时间,触发次数)
var examplelimit = ctxext.NewLimiterManager(time.Second*5, 1)

// 这里就是插件主体了
func init() {
	// 既然是zbp, 那就从接入control开始, 在这里注册你的插件以及设置是否默认开启和填写帮助和数据存放路径
	engine := control.Register("task", &ctrl.Options[*zero.Ctx]{
		// 控制插件是否默认启用 true为默认不启用 false反之
		DisableOnDefault: true,
		// 插件的简介
		Brief: "任务管理插件",
		// 插件的帮助 管理员发送 /用法 example 可见
		Help:  "用法：添加任务[任务名称]、移除任务[任务名称]、展示任务",
		// 插件的背景图, 支持http和本地路径
		// Banner: "",
		// 插件的数据存放路径, 分为公共和私有, 都会在/data下创建目录, 公有需要首字母大写, 私有需要首字母小写
		PublicDataFolder: "task",
		// PrivateDataFolder: "example",		// 避免问题所以注释了
		// 自定义插件开启时的回复
		OnEnable: func(ctx *zero.Ctx) {
			ctx.Send("任务系统已启用")
		},
		// 自定义插件关闭时的回复
		OnDisable: func(ctx *zero.Ctx) {
			ctx.Send("任务系统已禁用")
		},
	})

	engine.OnCommand("添加任务").Handle(addTask)
	engine.OnCommand("移除任务").Handle(removeTask)
	engine.OnCommand("展示任务").Handle(showTasks)

	loadTasks()

func addTask(ctx *zero.Ctx) {
	taskName := extractTaskName(ctx.RawMessage)
	if taskName == "" {
		ctx.SendChain(message.Text("任务名称不能为空！"))
		return
	}

	tasks = append(tasks, taskName)
	saveTasks()
	ctx.SendChain(message.Text("任务'%s'添加成功！", taskName))
}

func removeTask(ctx *zero.Ctx) {
	taskName := extractTaskName(ctx.RawMessage)
	if taskName == "" {
		ctx.SendChain(message.Text("任务名称不能为空！"))
		return
	}

	for i, task := range tasks {
		if task == taskName {
			tasks = append(tasks[:i], tasks[i+1:]...)
			saveTasks()
			ctx.SendChain(message.Text("任务'%s'移除成功！", taskName))
			return
		}
	}
	ctx.SendChain(message.Text("未找到名称为'%s'的任务！", taskName))
}

func showTasks(ctx *zero.Ctx) {
	if len(tasks) > 0 {
		taskList := strings.Join(tasks, "\n")
		ctx.SendChain(message.Text("当前任务列表：\n" + taskList))
	} else {
		ctx.SendChain(message.Text("当前无任务！"))
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
}
