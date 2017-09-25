package orametrics

import (
	"database/sql"
	//	"encoding/json"
	_ "github.com/mattn/go-oci8"
	"github.com/golang/glog"
	"time"
)

type tsBytes struct {
	Ts    string `json:"TS"`
	Bytes string `json:"bytes"`
}

type diskgroups struct {
	Dg	string`json:"DG"`
	UsableFileMB string `json:"USABLE_FILE_MB"`
	OfflineDisks string `json:"OFFLINE_DISKS"`
}

func Init(connectionString string, zabbixHost string, zabbixPort int, hostName string) {
	start := time.Now()
	defer glog.Flush()
	db, err := sql.Open("oci8", connectionString)
	if err != nil {
		glog.Fatal("Connection Failed!",err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		glog.Fatal("Error connecting to the database:", err)
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
	glog.Info("zabbixData:",zabbixData)
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
		if k == "diskgroups" {
			resultd := runDiscoveryQuery(v,db)
			var fixd string = "{\"data\":["
			countd := 1
			lend := len(resultd)
			for _, vd := range resultd {
				if countd < lend {
					fixd = fixd + "{\"{#DG}\":\"" + vd + "\"},"
				} else {
					fixd = fixd + "{\"{#DG}\":\"" + vd + "\"}"
				}
				countd++
			}
			fixd = fixd + "]}"
			discoveryData[k] = fixd
		}
	}
	//j := discoveryData["tablespaces"]
	//d := discoveryData["diskgroups"]


	ts_usage_bytes := runTsBytesDiscoveryQuery(ts_usage_bytes, db)
	ts_usage_pct := runTsBytesDiscoveryQuery(ts_usage_pct, db)
	diskGroupsMetrics := runDiskGroupsMetrics(diskgroup_metrics,db)
	discoveryMetrics := make(map[string]string)
	for _, v := range ts_usage_bytes {
		discoveryMetrics[`ts_usage_bytes[`+v.Ts+`]`] = v.Bytes
	}
	for _, v := range ts_usage_pct {
		discoveryMetrics[`ts_usage_pct[`+v.Ts+`]`] = v.Bytes
	}
	for _, v := range diskGroupsMetrics {
		discoveryMetrics[`usable_file_mb[`+v.Dg + `]`] = v.UsableFileMB
		discoveryMetrics[`offline_disks[`+v.Dg + `]`] = v.OfflineDisks
	}
	glog.Info("discoveryMetrics: ",discoveryMetrics)
	glog.Info("discoveryData: ", discoveryData)
	for k, v := range discoveryMetrics {
		zabbixData[k] = v
	}
	for k, v := range discoveryData {
		zabbixData[k] = v
	}
	//send(discoveryMetrics, zabbixHost, zabbixPort, hostName)
	glog.Info("zabbixData Combined: ",zabbixData)
	send(zabbixData, zabbixHost, zabbixPort, hostName)
	//sendD(j,"tablespaces", zabbixHost, zabbixPort, hostName)
	//sendD(d,"diskgroups",zabbixHost,zabbixPort,hostName)
	glog.Info(time.Since(start))
	//	tes, err := json.Marshal(discoveryMetrics)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(string(tes))
}
func runDiscoveryQuery(query string, db *sql.DB) []string {
	rows, err := db.Query(query)
	if err != nil {
		glog.Error("Error fetching addition",err)
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
		glog.Error("Error fetching addition", err)
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
		glog.Error("Error fetching addition",err)
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

func runDiskGroupsMetrics(query string, db *sql.DB) []diskgroups {
	rows, err := db.Query(query)
	if err != nil {
		glog.Error("Error fetching addition",err)
		var er []string
		er = append(er, err.Error())
		//return er
	}
	defer rows.Close()
	var result []diskgroups
	for rows.Next() {
		var res diskgroups
		rows.Scan(&res.Dg, &res.UsableFileMB, &res.OfflineDisks)
        result = append(result, res)
	}
	return result
}