package missions

import (
	"encoding/json"
	"fmt"
	"path"

	mq "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"

	"github.com/kubeedge/mappers-go/mappers/windows-virtual-exec/internal/core/mqtt"
	"github.com/kubeedge/mappers-go/mappers/windows-virtual-exec/internal/dto"
)

func InitCallback(nodeName string) {
	_err := mqtt.Client.Subscribe(fmt.Sprintf(mqtt.TopicRecNodeDeviceUpdate, nodeName), onMembershipUpdateMessage)
	if _err != nil {
		klog.Error("Subscribe error: ", _err)
	} else {
		klog.Info("Subscribe topic: ", fmt.Sprintf(mqtt.TopicRecNodeDeviceUpdate, nodeName))
	}
	_err = mqtt.Client.Subscribe(fmt.Sprintf(mqtt.TopicRecModeDeviceListResponse, nodeName), onMembershipListMessage)
	if _err != nil {
		klog.Error("Subscribe error: ", _err)
	} else {
		klog.Info("Subscribe topic: ", fmt.Sprintf(mqtt.TopicRecModeDeviceListResponse, nodeName))
	}
	_err = mqtt.Client.Subscribe(fmt.Sprintf(mqtt.TopicRevTwinUpdateDelta, "+"), onTwinDelta)
	if _err != nil {
		klog.Error("Subscribe error: ", _err)
	} else {
		klog.Info("Subscribe topic: ", fmt.Sprintf(mqtt.TopicRevTwinUpdateDelta, "+"))
	}
	_err = mqtt.Client.Subscribe(fmt.Sprintf(mqtt.TopicRecTwinInfoResponse, "+"), onTwinInfo)
	if _err != nil {
		klog.Error("Subscribe error: ", _err)
	} else {
		klog.Info("Subscribe topic: ", fmt.Sprintf(mqtt.TopicRecTwinInfoResponse, "+"))
	}
}

func onMembershipUpdateMessage(_ mq.Client, message mq.Message) {
	klog.V(2).Info("Receive message from topic: ", message.Topic())
	nodeId := mqtt.GetNodeID(message.Topic())
	if nodeId == "" {
		klog.Error("Wrong topic")
		return
	}
	klog.V(2).Info("Node id: ", nodeId)
	var req dto.DeviceListUpdate
	if _err := json.Unmarshal(message.Payload(), &req); _err != nil {
		klog.Error("Unmarshal error: ", _err)
		return
	}

	klog.Info("Receive device list update: ", "nodeId: ", nodeId, " update: ", len(req.AddedDevices), " delete: ", len(req.RemovedDevices))

	for _, device := range req.RemovedDevices {
		RemoveMission(device.Id)
	}

	for _, device := range req.RemovedDevices {
		if _, ok := cache.Load(device.Id); ok {
			klog.Info("Device already exists: ", device.Id)
			continue
		}
		klog.Info("Waiting twin update to create device: ", device.Id)
	}
}

func onMembershipListMessage(_ mq.Client, message mq.Message) {
	klog.V(2).Info("Receive message from topic: ", message.Topic())
	nodeId := mqtt.GetNodeID(message.Topic())
	if nodeId == "" {
		klog.Error("Wrong topic")
		return
	}
	klog.V(2).Info("Node id: ", nodeId)
	var req dto.DeviceList
	if _err := json.Unmarshal(message.Payload(), &req); _err != nil {
		klog.Error("Unmarshal error: ", _err)
		return
	}

	klog.Info("Receive device list: ", "nodeId: ", nodeId, " count: ", len(req.Devices))
	for _, device := range req.Devices {
		_err := mqtt.Client.Publish(fmt.Sprintf(mqtt.TopicPubTwinInfoRequest, device.Id), mqtt.CreateEmptyMessage())
		if _err != nil {
			klog.Error("Publish error: ", _err)
			return
		}
	}
}

func onTwinDelta(_ mq.Client, message mq.Message) {
	klog.V(2).Info("Receive message from topic: ", message.Topic())
	id := mqtt.GetDeviceID(message.Topic())
	if id == "" {
		klog.Error("Wrong topic")
		return
	}
	klog.V(2).Info("Mission id: ", id)
	var req dto.MissionDelta
	if _err := json.Unmarshal(message.Payload(), &req); _err != nil {
		klog.Error("Unmarshal error: ", _err)
		return
	}

	// check params
	if req.Twin.ExecCommand == nil || req.Twin.ExecFileName == nil || req.Twin.ExecFileContent == nil {
		klog.Error("Twin format error")
		return
	}

	if req.Twin.ExecCommand.Expected == nil || req.Twin.ExecCommand.Expected.Value == nil {
		klog.Error("Twin ExecCommand format error")
		return
	}

	if req.Twin.ExecFileName.Expected == nil || req.Twin.ExecFileName.Expected.Value == nil {
		klog.Error("Twin ExecFileName format error")
		return
	}

	if req.Twin.ExecFileContent.Expected == nil || req.Twin.ExecFileContent.Expected.Value == nil {
		klog.Error("Twin ExecFileContent format error")
		return
	}

	_err := mqtt.Client.Publish(fmt.Sprintf(mqtt.TopicPubTwinInfoRequest, id), mqtt.CreateEmptyMessage())
	if _err != nil {
		klog.Error("Publish error: ", _err)
		return
	}
}

// onTwinInfo indeed trigger a new mission
func onTwinInfo(_ mq.Client, message mq.Message) {
	klog.V(2).Info("Receive message from topic: ", message.Topic())
	id := mqtt.GetDeviceID(message.Topic())
	if id == "" {
		klog.Error("Wrong topic")
		return
	}
	klog.V(2).Info("Mission id: ", id)
	var req dto.MissionDelta
	if _err := json.Unmarshal(message.Payload(), &req); _err != nil {
		klog.Error("Unmarshal error: ", _err)
		return
	}
	if req.Twin.ExecCommand == nil || req.Twin.ExecFileName == nil || req.Twin.ExecFileContent == nil || req.Twin.Output == nil || req.Twin.Status == nil {
		klog.Error("Twin format error")
		return
	}

	if req.Twin.ExecCommand.Expected == nil || req.Twin.ExecCommand.Expected.Value == nil {
		klog.Error("Twin ExecCommand format error")
		return
	}

	if req.Twin.ExecFileName.Expected == nil || req.Twin.ExecFileName.Expected.Value == nil {
		klog.Error("Twin ExecFileName format error")
		return
	}

	if req.Twin.ExecFileContent.Expected == nil || req.Twin.ExecFileContent.Expected.Value == nil {
		klog.Error("Twin ExecFileContent format error")
		return
	}

	_, _err := NewMission(MissionConfig{
		UniqueName:       id,
		Command:          *req.Twin.ExecCommand.Expected.Value,
		FileContent:      *req.Twin.ExecFileContent.Expected.Value,
		FileName:         *req.Twin.ExecFileName.Expected.Value,
		WorkingDirectory: path.Join("tmp", id),
	})
	if _err != nil {
		klog.Error("NewMission error: ", _err)
		return
	}
}
