package main

import (
	"changeme/CommAnd"
	"github.com/qtgolang/SunnyNet/src/JsCall"
	"github.com/qtgolang/SunnyNet/src/protobuf/JSON"
	"strings"
)

// 保存进程名列表到配置文件
func saveProcessNamesToConfig() {
	configLock.Lock()
	defer configLock.Unlock()
	// 过滤掉空字符串
	var validNames []string
	for _, name := range GlobalConfig.ProcessNames {
		if strings.TrimSpace(name) != "" {
			validNames = append(validNames, name)
		}
	}
	GlobalConfig.ProcessNames = validNames
	_ = GlobalConfig.saveToFile()
}

// 加载进程名列表到 SunnyNet
func loadProcessNamesFromConfig() {
	configLock.Lock()
	names := make([]string, len(GlobalConfig.ProcessNames))
	copy(names, GlobalConfig.ProcessNames)
	configLock.Unlock()

	// 将保存的进程名添加到 SunnyNet
	for _, name := range names {
		if strings.TrimSpace(name) != "" {
			app.App.ProcessAddName(name)
		}
	}
}

func processEvent(command string, args *JSON.SyJson) any {
	switch command {
	case "加载驱动":
		result := app.App.StartProcess()
		// 加载保存的进程名列表
		loadProcessNamesFromConfig()
		return result
	case "枚举进程":
		p := CommAnd.EnumerateProcesses()
		return p
	case "进程驱动添加PID":
		gx := args.GetData("isSelected") == "true"
		app.App.ProcessCancelAll()
		if gx {
			app.App.ProcessAddPid(getInt(args.GetData("PID")))
		} else {
			app.App.ProcessDelPid(getInt(args.GetData("PID")))
		}
		return true
	case "进程驱动添加进程名":
		gx := args.GetData("isSet") == "true"
		Name := JsCall.ToGBK(args.GetData("Name"))
		if Name == "{OpenALL}" {
			app.App.ProcessALLName(gx)
			return true
		}
		if gx {
			app.App.ProcessAddName(Name)
			// 添加到配置列表
			configLock.Lock()
			// 检查是否已存在
			exists := false
			for _, n := range GlobalConfig.ProcessNames {
				if n == Name {
					exists = true
					break
				}
			}
			if !exists {
				GlobalConfig.ProcessNames = append(GlobalConfig.ProcessNames, Name)
			}
			configLock.Unlock()
			saveProcessNamesToConfig()
		} else {
			app.App.ProcessDelName(Name)
			// 从配置列表中移除
			configLock.Lock()
			for i, n := range GlobalConfig.ProcessNames {
				if n == Name {
					GlobalConfig.ProcessNames = append(GlobalConfig.ProcessNames[:i], GlobalConfig.ProcessNames[i+1:]...)
					break
				}
			}
			configLock.Unlock()
			saveProcessNamesToConfig()
		}
		return true
	default:
		return ReplaceRulesEvent(command, args)
	}
}
