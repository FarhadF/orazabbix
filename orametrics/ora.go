package orametrics

import (
	"database/sql"
	//	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"github.com/golang/glog"
)

type tsBytes struct {
	Ts    string `json:"TS"`
	Bytes string `json:"bytes"`
}

func Init(connectionString string, zabbixHost string, zabbixPort int, hostName string) {
	defer glog.Flush()
	db, err := sql.Open("oci8", connectionString)
	if err != nil {
		glog.Fatal("Connection Failed!",err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		glog.Fatal("Error connecting to the database: %s\n", err)
		return
	}
	zabbixData := make(map[string]string)
	for k, v := range queries {
		//	zabbixData[k] = runQuery(v, db)
		rows, err := db.Query(v)
		if err != nil {
			glog.Error("Error fetching addition", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var res string
			rows.Scan(&res)
			zabbixData[k] = res
		}
	}
	fmt.Println(zabbixData)
	//discoveryData := make(map[string][]string)
	discoveryData := make(map[string]string)
	for k, v := range discoveryQueries {
		if k == "tablespaces" {
			result := runDiscoveryQuery(v, db)
			var fix string = "{\"data\":["
			count := 1
			len := len(result)
			for _, va := range result {
				if count < len {
					fix = fix + "{\"{#TS}\":\"" + va + "\"},"
				} else {
					fix = fix + "{\"{#TS}\":\"" + va + "\"}"
				}
				count++
			}
			fix = fix + "]}"
			discoveryData[k] = fix
		}
	}
	j := discoveryData["tablespaces"]
	send(zabbixData, zabbixHost, zabbixPort, hostName)
	sendD(j, zabbixHost, zabbixPort, hostName)
	ts_usage_bytes := runTsBytesDiscoveryQuery(ts_usage_bytes, db)
	ts_usage_pct := runTsBytesDiscoveryQuery(ts_usage_pct, db)
	discoveryMetrics := make(map[string]string)
	for _, v := range ts_usage_bytes {
		discoveryMetrics[`ts_usage_bytes[`+v.Ts+`]`] = v.Bytes
	}
	for _, v := range ts_usage_pct {
		discoveryMetrics[`ts_usage_pct[`+v.Ts+`]`] = v.Bytes
	}
	fmt.Println(discoveryMetrics)
	send(discoveryMetrics, zabbixHost, zabbixPort, hostName)
	//	tes, err := json.Marshal(discoveryMetrics)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(string(tes))
}
func runDiscoveryQuery(query string, db *sql.DB) []string {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error fetching addition")
		fmt.Println(err)
		var er []string
		er = append(er, err.Error())
		return er
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var res string
		rows.Scan(&res)
		result = append(result, res)
	}
	return result
}

func runQuery(query string, db *sql.DB) string {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error fetching addition")
		fmt.Println(err)
		return err.Error()
	}
	defer rows.Close()
	var res string
	for rows.Next() {
		rows.Scan(&res)
	}
	return res
}

func runTsBytesDiscoveryQuery(query string, db *sql.DB) []tsBytes {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error fetching addition")
		fmt.Println(err)
		var er []string
		er = append(er, err.Error())
		//return er
	}
	defer rows.Close()
	var result []tsBytes
	for rows.Next() {
		var res tsBytes
		rows.Scan(&res.Ts, &res.Bytes)

		result = append(result, res)
	}
	return result
}
