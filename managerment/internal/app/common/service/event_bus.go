// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	eventBus "github.com/asaskevich/EventBus"
)

var localEventBus eventBus.Bus

func EventBus() eventBus.Bus {
	if localEventBus == nil {
		panic("implement not found for interface EventBus, forgot register?")
	}
	return localEventBus
}

func RegisterEventBus(i eventBus.Bus) {
	localEventBus = i
}