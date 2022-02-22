package apollo

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4/storage"
)

type listenerResponse struct {
	Content []byte
	Error   error
}

type listener struct {
	eventCh              chan *listenerResponse
	onConfigChange       func(*storage.ChangeEvent)
	onNewestConfigChange func(*storage.FullChangeEvent)
}

// OnChange 增加变更监控
func (l listener) OnChange(event *storage.ChangeEvent) {
	changes := make(map[string]string, len(event.Changes))
	for k, v := range event.Changes {
		if k == "content" || v.ChangeType == storage.DELETED {
			continue
		}
		changes[k] = fmt.Sprint(v.NewValue)
	}

	setting, err := marshalToConfigType(changes, instance.configType)
	if err != nil {
		l.eventCh <- &listenerResponse{
			Content: nil,
			Error:   err,
		}
	} else {
		l.eventCh <- &listenerResponse{
			Content: setting,
			Error:   nil,
		}
	}

	if l.onConfigChange != nil {
		l.onConfigChange(event)
	}
}

// OnNewestChange 监控最新变更
func (l listener) OnNewestChange(event *storage.FullChangeEvent) {
	changes := make(map[string]string, len(event.Changes))
	for k, v := range event.Changes {
		changes[k] = fmt.Sprint(v)
	}

	setting, err := marshalToConfigType(changes, instance.configType)
	if err != nil {
		l.eventCh <- &listenerResponse{
			Content: nil,
			Error:   err,
		}
	} else {
		l.eventCh <- &listenerResponse{
			Content: setting,
			Error:   nil,
		}
	}

	if l.onNewestConfigChange != nil {
		l.onNewestConfigChange(event)
	}
}
