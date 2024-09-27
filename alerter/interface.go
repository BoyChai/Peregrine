package alerter

import "Peregrine/stru"

type alerter map[string]chan stru.AlarmContext

var Alerters alerter

func init() {
	// 初始化 Alerters
	Alerters = make(alerter)
}
