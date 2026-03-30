package main

import (
	"changeme/CommAnd"
	"github.com/qtgolang/SunnyNet/src/JsCall"
	"github.com/qtgolang/SunnyNet/src/protobuf/JSON"
	"strings"
)

// 获取进程名列表
func getProcessNames() []string {
	configLock.Lock()
	defer configLock.Unlock()
	names := make([]string, len(GlobalConfig.ProcessNames))
	copy(names, GlobalConfig.ProcessNames)
	return names
}

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

	// 将保存的进程名添加到 SunnyNet（配置中存的是UTF-8，需要转GBK给SunnyNet）
	for _, name := range names {
		if strings.TrimSpace(name) != "" {
			app.App.ProcessAddName(JsCall.ToGBK(name))
		}
	}
}

func processEvent(command string, args *JSON.SyJson) any {
	switch command {
	case "获取进程名列表":
		return getProcessNames()
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
		utf8Name := args.GetData("Name")
		if utf8Name == "{OpenALL}" {
			app.App.ProcessALLName(gx)
			return true
		}
		gbkName := JsCall.ToGBK(utf8Name)
		if gx {
			app.App.ProcessAddName(gbkName)
			// 添加到配置列表（保存UTF-8原始名称，避免JSON序列化时GBK字节被损坏导致乱码）
			configLock.Lock()
			exists := false
			for _, n := range GlobalConfig.ProcessNames {
				if n == utf8Name {
					exists = true
					break
				}
			}
			if !exists {
				GlobalConfig.ProcessNames = append(GlobalConfig.ProcessNames, utf8Name)
			}
			configLock.Unlock()
			saveProcessNamesToConfig()
		} else {
			app.App.ProcessDelName(gbkName)
			// 从配置列表中移除（配置中存的是UTF-8）
			configLock.Lock()
			for i, n := range GlobalConfig.ProcessNames {
				if n == utf8Name {
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
