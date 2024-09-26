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

var alarmChan chan stru.AlarmTrigger

func Run(ch chan stru.AlarmTrigger, rule []stru.Rule) {
	if len(rule) == 0 {
		log.Fatalln("没有指定规则。。")
	}
	alarmChan = ch
	for _, r := range rule {
		go entryProcessor(asset.GetAsset(r.AssetName), r.Entry, r.TriggerCount, r.ProbeInterval)
	}
}

func entryProcessor(asset string, entrys []stru.RuleEntry, count, interval int) {
	if len(entrys) == 0 {
		log.Fatalln("没有设置条目规则。。")
	}
	for _, e := range entrys {
		go detectionTask(asset, e, count, interval)
	}
}
func detectionTask(asset string, rule stru.RuleEntry, count, interval int) {
	// 内部计数器，用于判断是否触发报警
	var errorCount int = 0
	uri := "/api/v1/query"
	encodedQuery := url.QueryEscape(rule.Expr) // URL 编码
	urlWithQuery := fmt.Sprintf("%s?query=%s", asset+uri, encodedQuery)

	for {
		fmt.Println("=============")
		fmt.Println(errorCount)
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
		if errorCount >= count {
			fmt.Println("触发报警")
			fmt.Println(respData.Data.Result[0].Metric.Instance)
			fmt.Println(respData.Data.Result[0].Value[1])
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
