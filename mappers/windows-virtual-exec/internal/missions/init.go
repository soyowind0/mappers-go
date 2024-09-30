package missions

import (
	"fmt"

	"github.com/kubeedge/mappers-go/mappers/windows-virtual-exec/internal/core/mqtt"
)

func InitMissions(nodeName string) {
	mqtt.Client.Publish(fmt.Sprintf(mqtt.TopicPubNodeDeviceListRequest, nodeName), mqtt.CreateEmptyMessage())
}
