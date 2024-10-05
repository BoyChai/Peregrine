package monitor

import (
	"Peregrine/asset"
	"Peregrine/log"
	"Peregrine/stru"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// 告警
var alarmChan = make(chan stru.AlarmTrigger)

var alertMap = make(map[int]time.Time)

var rule []stru.Rule

func Run(ch chan stru.AlarmTrigger, r []stru.Rule) {
	if len(r) == 0 {
		log.Error("没有指定探测规则。。")
	}
	rule = r

	for i, r := range rule {
		go entryProcessor(r.AssetName, r.AlerterTarget, r.AlerterWay, r.Entry, r.TriggerCount, r.ProbeInterval, i)
	}
	go receiver(ch)
}

func receiver(ch chan stru.AlarmTrigger) {

	for {
		trigger := <-alarmChan
		if trigger.Time.Sub(alertMap[trigger.EntryID]) > time.Duration(rule[trigger.RuleID].For)*time.Minute {
			// 更新告警时间
			alertMap[trigger.EntryID] = trigger.Time
			ch <- trigger
		}
	}
}

func entryProcessor(assetName, target, way string, entrys []stru.RuleEntry, count, interval, ruleID int) {
	if len(entrys) == 0 {
		log.Error("没有设置条目规则。。")
	}
	for id, e := range entrys {
		alertMap[id] = time.Now().Add(-5 * time.Minute)
		go detectionTask(assetName, target, way, e, count, interval, ruleID, id)
	}
}
func detectionTask(assetName, target, way string, rule stru.RuleEntry, count, interval, ruleID, entryID int) {
	// 内部计数器，用于判断是否触发报警
	var errorCount int = 0
	uri := "/api/v1/query"
	encodedQuery := url.QueryEscape(rule.Expr) // URL 编码
	urlWithQuery := fmt.Sprintf("%s?query=%s", asset.GetAsset(assetName)+uri, encodedQuery)
	log.Debug(urlWithQuery)
	for {
		var lock bool = false
		time.Sleep(time.Duration(interval) * time.Second)

		code, body, err := runDetection(urlWithQuery)
		if err != nil {
			log.Error("请求失败", err)
			errorCount++
			lock = true
		}
		if !lock && code != http.StatusOK {
			errorCount++
			log.Error("请求状态码错误", code)
		}
		var respData stru.PrometheusResp
		err = json.Unmarshal([]byte(body), &respData)
		if !lock && err != nil {
			log.Error("返回解析失败", err)
			errorCount++
		}
		if !lock && respData.Data.Result != nil {
			errorCount++
		}
		// 告警返回内容
		if errorCount >= count {
			log.Debug("返回告警内容", urlWithQuery)
			var trigger = stru.AlarmTrigger{
				AssetName:     assetName,
				AlerterTarget: target,
				AlerterWay:    way,
				Entry:         rule,
				RuleID:        ruleID,
				EntryID:       entryID,
				Time:          time.Now(),
			}
			if len(respData.Data.Result) != 0 {
				for _, v := range respData.Data.Result {
					trigger.Instance = append(trigger.Instance, v.Metric.Instance)
					trigger.Value = append(trigger.Value, fmt.Sprintf("%v", v.Value[0]))
					trigger.Value = append(trigger.Value, fmt.Sprintf("%v", v.Value[1]))
				}
			}
			alarmChan <- trigger
			errorCount = 0
			continue
		}
	}
}

func runDetection(url string) (statusCode int, bodyData string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error("请求失败。。。")
		return 0, "", nil

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("资源请求失败。。")
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)

	return resp.StatusCode, buf.String(), nil
}
