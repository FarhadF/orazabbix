package orametrics

import (
	. "github.com/blacked/go-zabbix"
	"time"
)

func send(zabbixData map[string]string, zabbixHost string, zabbixPort int) {
	var metrics []*Metric
	for k, v := range zabbixData {
		metrics = append(metrics, NewMetric("server1", k, v, time.Now().Unix()))
		//metrics = append(metrics, NewMetric("server1", "status", "OK"))
	}
	// Create instance of Packet class
	packet := NewPacket(metrics)

	// Send packet to zabbix
	z := NewSender(zabbixHost, zabbixPort)
	z.Send(packet)
}
