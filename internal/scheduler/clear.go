package scheduler

import (
	"time"

	"im-service/internal/logic"
	"im-service/pkg/util"
)

// 清理离线用户
func Clear(m map[string][]*logic.Client) {
	util.RecoverPanic()
	var (
		ticker = time.NewTicker(time.Minute * 5)
		// mu     = new(sync.RWMutex)
	)
	for range ticker.C {
		for i := range m {
			var clients []*logic.Client
			// logrus.Info(len(m[i]))
			for _, v := range m[i] {
				if (time.Now().Unix() - v.Updated.Unix()) < 12 {
					clients = append(clients, v)
				}
			}
			m[i] = clients
			if len(m[i]) == 0 {
				m[i] = nil
				continue
			}
		}
	}
}
