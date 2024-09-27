package monitor

import (
	"Peregrine/asset"
	"Peregrine/stru"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
		log.Fatalln("没有指定规则。。")
	}
	rule = r

	for i, r := range rule {
		go entryProcessor(asset.GetAsset(r.AssetName), r.AlerterTarget, r.Entry, r.TriggerCount, r.ProbeInterval, i)
	}
	go receiver(ch)
}

func receiver(ch chan stru.AlarmTrigger) {

	for {
		trigger := <-alarmChan
		fmt.Println("告警告警2")
		fmt.Println(trigger.Time.Sub(alertMap[trigger.ID]))
		fmt.Println(time.Duration(rule[trigger.ID].For) * time.Minute)
		if trigger.Time.Sub(alertMap[trigger.ID]) > time.Duration(rule[trigger.ID].For)*time.Minute {
			fmt.Println("告警告警aaa")
			// 更新告警时间
			alertMap[trigger.ID] = trigger.Time
			fmt.Println(trigger)
			ch <- trigger
		}
	}
}

func entryProcessor(asset string, target string, entrys []stru.RuleEntry, count, interval, ruleID int) {
	if len(entrys) == 0 {
		log.Fatalln("没有设置条目规则。。")
	}
	for _, e := range entrys {
		alertMap[ruleID] = time.Now().Add(-5 * time.Minute)
		go detectionTask(asset, target, e, count, interval, ruleID)
	}
}
func detectionTask(asset string, target string, rule stru.RuleEntry, count, interval, ruleID int) {
	// 内部计数器，用于判断是否触发报警
	var errorCount int = 0
	uri := "/api/v1/query"
	encodedQuery := url.QueryEscape(rule.Expr) // URL 编码
	urlWithQuery := fmt.Sprintf("%s?query=%s", asset+uri, encodedQuery)

	for {
		var lock bool = false
		time.Sleep(time.Duration(interval) * time.Second)

		code, body, err := runDetection(urlWithQuery)
		if err != nil {
			log.Println("请求失败", err)
			errorCount++
			lock = true
		}
		if !lock && code != http.StatusOK {
			errorCount++
			log.Println("请求状态码错误", code)
		}
		var respData stru.PrometheusResp
		err = json.Unmarshal([]byte(body), &respData)
		if !lock && err != nil {
			log.Println("返回解析失败", err)
			errorCount++
		}
		if !lock && respData.Data.Result != nil {
			errorCount++
		}
		// 告警返回内容
		if errorCount >= count {
			fmt.Println("告警告警1")
			var trigger = stru.AlarmTrigger{
				AlerterTarget: target,
				Entry:         rule,
				ID:            ruleID,
				Time:          time.Now(),
			}
			if len(respData.Data.Result) != 0 {
				for _, v := range respData.Data.Result {
					trigger.Instance = append(trigger.Instance, v.Metric.Instance)
					trigger.Value = append(trigger.Value, fmt.Sprintf("%v", v.Value[0]))
					trigger.Value = append(trigger.Value, fmt.Sprintf("%v", v.Value[1]))
				}
			}
			fmt.Println("告警触发条件满足")
			alarmChan <- trigger
			fmt.Println("告警发送成功")
			errorCount = 0
			continue
		}
	}
}

func runDetection(url string) (statusCode int, bodyData string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("请求失败。。。")
		return 0, "", nil

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("资源请求失败。。")
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)

	return resp.StatusCode, buf.String(), nil
}
