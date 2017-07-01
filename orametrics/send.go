package orametrics

import (
	//"encoding/json"
	//	"fmt"
	. "github.com/blacked/go-zabbix"
	"time"
)

func send(zabbixData map[string]string, zabbixHost string, zabbixPort int, hostName string) {
	var metrics []*Metric
	for k, v := range zabbixData {
		metrics = append(metrics, NewMetric(hostName, k, v, time.Now().Unix()))
		//metrics = append(metrics, NewMetric("server1", "status", "OK"))
	}
	// Create instance of Packet class
	packet := NewPacket(metrics)

	// Send packet to zabbix
	z := NewSender(zabbixHost, zabbixPort)
	z.Send(packet)
}

func sendD(j string, zabbixHost string, zabbixPort int, hostName string) {
	var metrics []*Metric
	//	for k, v := range discoveryData {
	//		j, _ := json.Marshal(v)
	//		fmt.Println(k)
	//		fmt.Println(v)
	metrics = append(metrics, NewMetric(hostName, "tablespaces", string(j), time.Now().Unix()))
	//metrics = append(metrics, NewMetric("server1", "status", "OK"))
	//}
	// Create instance of Packet class
	packet := NewPacket(metrics)

	// Send packet to zabbix
	z := NewSender(zabbixHost, zabbixPort)
	z.Send(packet)
}
