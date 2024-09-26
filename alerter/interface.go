package alerter

import "Peregrine/stru"

type Alert struct {
	Entry  stru.RuleEntry
	Target []stru.Target
}

type alerter map[string]chan Alert

var Alerters alerter

func init() {
	// 初始化 Alerters
	Alerters = make(alerter)
}
