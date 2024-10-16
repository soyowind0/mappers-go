package missions

import (
	"fmt"

	"k8s.io/klog"

	"github.com/kubeedge/mappers-go/mappers/windows-virtual-exec/internal/core/mqtt"
)

func InitMissions(nodeName string) {
	if err := mqtt.GetClient().Publish(fmt.Sprintf(mqtt.TopicPubNodeDeviceListRequest, nodeName), mqtt.CreateEmptyMessage()); err != nil {
		klog.Errorf("Failed to init missions on %s", nodeName)
	}
}
