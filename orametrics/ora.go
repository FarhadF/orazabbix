package orametrics

import (
	"database/sql"
	//	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-oci8"
)

type tsBytes struct {
	Ts    string `json:"TS"`
	Bytes string `json:"bytes"`
}

func Init(connectionString string, zabbixHost string, zabbixPort int, hostName string) {
	db, err := sql.Open("oci8", connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		return
	}
	zabbixData := make(map[string]string)
	for k, v := range queries {
		//	zabbixData[k] = runQuery(v, db)
		rows, err := db.Query(v)
		if err != nil {
			fmt.Println("Error fetching addition")
			fmt.Println(err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var res string
			rows.Scan(&res)
			fmt.Println(res)
			zabbixData[k] = res

		}
	}
	fmt.Println(zabbixData)
	//discoveryData := make(map[string][]string)
	discoveryData := make(map[string]string)
	for k, v := range discoveryQueries {
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
	j := discoveryData["tablespaces"]
	fmt.Println("----")
	fmt.Println(string(j))
	send(zabbixData, zabbixHost, zabbixPort, hostName)
	sendD(j, zabbixHost, zabbixPort, hostName)
	ts_usage_bytes := runTsBytesDiscoveryQuery(ts_usage_bytes, db)
	discoveryMetrics := make(map[string]string)
	for _, v := range ts_usage_bytes {
		discoveryMetrics[`ts_usage_bytes[`+v.Ts+`]`] = v.Bytes
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
		fmt.Println(res)
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
		fmt.Println(res)
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
		fmt.Println(res)
		result = append(result, res)
	}
	return result
}
