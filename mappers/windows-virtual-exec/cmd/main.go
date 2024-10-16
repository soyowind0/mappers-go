package main

import (
	"os"
	"os/signal"

	klog "k8s.io/klog/v2"

	"github.com/kubeedge/mappers-go/mappers/windows-virtual-exec/internal/config"
	"github.com/kubeedge/mappers-go/mappers/windows-virtual-exec/internal/core/model"
	"github.com/kubeedge/mappers-go/mappers/windows-virtual-exec/internal/core/mqtt"
	"github.com/kubeedge/mappers-go/mappers/windows-virtual-exec/internal/core/store"
	"github.com/kubeedge/mappers-go/mappers/windows-virtual-exec/internal/missions"
)

func main() {
	var err error
	var c config.Config

	klog.InitFlags(nil)
	defer klog.Flush()

	if err = c.Parse(); err != nil {
		klog.Fatal(err)
	}

	store.InitDB("internal.db")
	if err := store.DB.AutoMigrate(&model.Mission{}); err != nil {
		klog.Errorf("Failed to init db: %v", err)
	}

	if err := mqtt.InitClient(
		c.Mqtt.ServerAddress,
		c.Mqtt.Username,
		c.Mqtt.Password,
		c.Mqtt.Cert,
		c.Mqtt.PrivateKey,
	); err != nil {
		klog.Fatal(err)
	}

	missions.InitCallback(c.NodeName)
	klog.Info("Start to subscribe")
	missions.InitMissions(c.NodeName)

	// waiting kill signal
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	klog.Info("Exit")
}
