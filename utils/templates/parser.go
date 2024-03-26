package templates

import (
	"bytes"
	"encoding/json"
	"strconv"
	"text/template"
	"time"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
)

var tmpl *template.Template

// ParserTemplate 处理告警推送的消息模版
func ParserTemplate(defineName string, alert models.AlertCurEvent, templateStr string) string {

	firstTriggerTime := time.Unix(alert.FirstTriggerTime, 0).Format(globals.Layout)
	recoverTime := time.Unix(alert.RecoverTime, 0).Format(globals.Layout)
	alert.FirstTriggerTimeFormat = firstTriggerTime
	alert.RecoverTimeFormat = recoverTime

	tmpl = template.Must(template.New("tmpl").Parse(templateStr))

	var (
		buf bytes.Buffer
		err error
	)

	if defineName == "Card" {
		err = tmpl.Execute(&buf, alert)
		// 当前告警的 json 反序列化成 map 对象, 用于解析报警事件详情中的 ${xx} 变量
		data := parserEvent(alert)
		return cmd.ParserVariables(buf.String(), data)
	}

	err = tmpl.ExecuteTemplate(&buf, defineName, alert)
	if err != nil {
		globals.Logger.Sugar().Error("告警模版执行失败 ->", err.Error())
		return ""
	}

	// 前面只会渲染出模版框架, 下面来渲染告警数据内容
	if defineName == "Event" {
		data := parserEvent(alert)
		return cmd.ParserVariables(buf.String(), data)
	}

	return buf.String()

}

func parserEvent(alert models.AlertCurEvent) map[string]interface{} {

	data := make(map[string]interface{})

	eventJson := cmd.JsonMarshal(alert)
	err := json.Unmarshal([]byte(eventJson), &data)
	if err != nil {
		globals.Logger.Sugar().Error("parserEvent Unmarshal failed: ", err)
	}

	if alert.DatasourceType == "AliCloudSLS" {
		// 需要转义, 日志中可能会出现特殊符号
		alarmInfo := strconv.Quote(data["annotations"].(string))
		data["annotations"] = alarmInfo[1 : len(alarmInfo)-1]
	}

	return data

}
