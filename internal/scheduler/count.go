package scheduler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"im-service/internal/logic"
	"im-service/pkg/db/redis"
	"im-service/pkg/mq"
	"im-service/pkg/pb"
	"im-service/pkg/util"
)

var (
	online  int
	baggage []byte
)

func Count(m map[string]int)  {
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		// i 房间号
		for i := range m {
			// 赛事正在进行中：基数 20 + 在线人数 * 3
			if startEnd(i) {
				online = 20 + m[i]*3
			} else {
				online = m[i]
			}
			pbOnline := pb.Encode(&pb.MessageWrapper{
				Topic: logic.TopicOnline,
				MsgBody: &pb.MessageBody{
					MessageType:    pb.MessageType_MT_TEXT,
					MessageContent: fmt.Sprintf("%v", online),
				},
			})
			ds := pb.Encode(&pb.Downstream{Type: pb.PackageType_PT_MESSAGE, Data: pbOnline})
			baggage = util.BytesCombine([]byte(i+logic.Separator), ds)
			mq.Pub(logic.TopicOnline, baggage)
		}
	}
}

/*func Count(m map[string]map[string]bool) {
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		for i := range m {
			// 赛事正在进行中：基数 20 + 在线人数 * 3
			if startEnd(i) {
				online = 20 + len(m[i])*3
			} else {
				online = len(m[i])
			}
			pbOnline := pb.Encode(&pb.MessageWrapper{
				Topic: logic.TopicOnline,
				MsgBody: &pb.MessageBody{
					MessageType:    pb.MessageType_MT_TEXT,
					MessageContent: fmt.Sprintf("%v", online),
				},
			})
			ds := pb.Encode(&pb.Downstream{Type: pb.PackageType_PT_MESSAGE, Data: pbOnline})
			baggage = util.BytesCombine([]byte(i+logic.Separator), ds)
			mq.Pub(logic.TopicOnline, baggage)
		}
	}
}*/

func startEnd(roomId string) (ok bool) {
	var (
		array  []int
		series = redis.HGet("series", "seriesGo")
	)
	// roomId string to int and json string to array
	i, _ := strconv.Atoi(roomId)
	_ = json.Unmarshal([]byte(series), &array)
	for _, v := range array {
		if v == i {
			ok = true
			break
		} else {
			ok = false
		}
	}
	return ok
}
