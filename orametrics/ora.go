package orametrics

import (
	"database/sql"
	//	"encoding/json"
	_ "github.com/mattn/go-oci8"
	"github.com/golang/glog"
	"time"
	"strconv"
)

type tsBytes struct {
	Ts    string `json:"TS"`
	Bytes string `json:"bytes"`
}

type diskgroups struct {
	Dg           string `json:"DG"`
	UsableFileMB string `json:"USABLE_FILE_MB"`
	OfflineDisks string `json:"OFFLINE_DISKS"`
}

type instance struct {
	INST_ID          string `json:"INST_ID"`
	INSTANCE_NUMBER  string `json:"INSTANCE_NUMBER"`
	INSTANCE_NAME	 string `json:"INSTANCE_NAME"`
	HOST_NAME        string `json:"HOST_NAME"`
	VERSION          string `json:"VERSION"`
	STARTUP_TIME     string `json:"STARTUP_TIME"`
	STATUS           string `json:"STATUS"`
	PARALLEL         string `json:"PARALLEL"`
	THREAD_NO        string `json:"THREAD_NO"`
	ARCHIVER         string `json:"ARCHIVER"`
	LOG_SWITCH_WAIT  sql.NullString `json:"LOG_SWITCH_WAIT"`
	LOGINS           string `json:"LOGINS"`
	SHUTDOWN_PENDING string `json:"SHUTDOWN_PENDING"`
	DATABASE_STATUS  string `json:"DATABASE_STATUS"`
	INSTANCE_ROLE    string `json:"INSTANCE_ROLE"`
	ACTIVE_STATE     string `json:"ACTIVE_STATE"`
	BLOCKED          string `json:"BLOCKED"`
	CON_ID           string `json:"CON_ID"`
	INSTANCE_MODE    string `json:"INSTANCE_MODE"`
	EDITION          string `json:"EDITION"`
	FAMILY           sql.NullString `json:"FAMILY"`
	DATABASE_TYPE    string `json:"DATABASE_TYPE"`
}

func Init(connectionString string, zabbixHost string, zabbixPort int, hostName string) {
	start := time.Now()
	defer glog.Flush()
	db, err := sql.Open("oci8", connectionString)
	if err != nil {
		glog.Fatal("Connection Failed!", err)
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
			glog.Error("Error fetching addition. ", err, k, v)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var res string
			rows.Scan(&res)
			zabbixData[k] = res
		}
	}
	if zabbixData["pool_dict_cache"] == "" {
		zabbixData["pool_dict_cache"] = "0"
	}
	if zabbixData["pool_lib_cache"] == "" {
		zabbixData["pool_lib_cache"] = "0"
	}
	if zabbixData["pool_sql_area"] == "" {
		zabbixData["pool_sql_area"] = "0"
	}
	//glog.Info("zabbixData:", zabbixData)
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
			resultd := runDiscoveryQuery(v, db)
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
		if k == "instances" {
			resultI := runDiscoveryQuery(v, db)
			var fixd string = "{\"data\":["
			countI := 1
			lend := len(resultI)
			for _, vi := range resultI {
				if countI < lend {
					fixd = fixd + "{\"{#INS}\":\"" + vi + "\"},"
				} else {
					fixd = fixd + "{\"{#INS}\":\"" + vi + "\"}"
				}
				countI++
			}
			fixd = fixd + "]}"
			discoveryData[k] = fixd
		}
	}
	//j := discoveryData["tablespaces"]
	//d := discoveryData["diskgroups"]

	ts_usage_bytes := runTsBytesDiscoveryQuery(ts_usage_bytes, db)
	ts_maxsize_bytes := runTsBytesDiscoveryQuery(ts_maxsize_bytes, db)
	ts_usage_pct := runTsBytesDiscoveryQuery(ts_usage_pct, db)
	diskGroupsMetrics := runDiskGroupsMetrics(diskgroup_metrics, db)
	instanceMetrics := runInstanceMetrics(instance_metrics, db)
	discoveryMetrics := make(map[string]string)
	for _, v := range ts_usage_bytes {
		discoveryMetrics[`ts_usage_bytes[`+v.Ts+`]`] = v.Bytes
	}
	for _, v := range ts_maxsize_bytes {
		discoveryMetrics[`ts_maxsize_bytes[`+v.Ts+`]`] = v.Bytes
	}
	for _, v := range ts_usage_pct {
		discoveryMetrics[`ts_usage_pct[`+v.Ts+`]`] = v.Bytes
	}
	for _, v := range diskGroupsMetrics {
		bytes, _ := strconv.Atoi(v.UsableFileMB)
		bytes = bytes * 1048576
		bytesS := strconv.Itoa(bytes)
		discoveryMetrics[`usable_file_mb[`+v.Dg+`]`] = bytesS
		discoveryMetrics[`offline_disks[`+v.Dg+`]`] = v.OfflineDisks
	}
	for _, v := range instanceMetrics {
		discoveryMetrics[`INST_ID[`+v.INSTANCE_NAME+`]`] = v.INST_ID
		discoveryMetrics[`INSTANCE_NUMBER[`+v.INSTANCE_NAME+`]`] = v.INSTANCE_NUMBER
		discoveryMetrics[`HOST_NAME[`+v.INSTANCE_NAME+`]`] = v.HOST_NAME
		discoveryMetrics[`VERSION[`+v.INSTANCE_NAME+`]`] = v.VERSION
		discoveryMetrics[`STARTUP_TIME[`+v.INSTANCE_NAME+`]`] = v.STARTUP_TIME
		discoveryMetrics[`STATUS[`+v.INSTANCE_NAME+`]`] = v.STATUS
		discoveryMetrics[`PARALLEL[`+v.INSTANCE_NAME+`]`] = v.PARALLEL
		discoveryMetrics[`THREAD_NO[`+v.INSTANCE_NAME+`]`] = v.THREAD_NO
		discoveryMetrics[`ARCHIVER[`+v.INSTANCE_NAME+`]`] = v.ARCHIVER
		if v.LOG_SWITCH_WAIT.Valid == true {
			discoveryMetrics[`LOG_SWITCH_WAIT[`+v.INSTANCE_NAME+`]`] = v.LOG_SWITCH_WAIT.String
		} else {
			discoveryMetrics[`LOG_SWITCH_WAIT[`+v.INSTANCE_NAME+`]`] = "0"
		}
		discoveryMetrics[`LOGINS[`+v.INSTANCE_NAME+`]`] = v.LOGINS
		discoveryMetrics[`SHUTDOWN_PENDING[`+v.INSTANCE_NAME+`]`] = v.SHUTDOWN_PENDING
		discoveryMetrics[`DATABASE_STATUS[`+v.INSTANCE_NAME+`]`] = v.DATABASE_STATUS
		discoveryMetrics[`INSTANCE_ROLE[`+v.INSTANCE_NAME+`]`] = v.INSTANCE_ROLE
		discoveryMetrics[`ACTIVE_STATE[`+v.INSTANCE_NAME+`]`] = v.ACTIVE_STATE
		discoveryMetrics[`BLOCKED[`+v.INSTANCE_NAME+`]`] = v.BLOCKED
		discoveryMetrics[`CON_ID[`+v.INSTANCE_NAME+`]`] = v.CON_ID
		discoveryMetrics[`INSTANCE_MODE[`+v.INSTANCE_NAME+`]`] = v.INSTANCE_MODE
		discoveryMetrics[`EDITION[`+v.INSTANCE_NAME+`]`] = v.EDITION
		discoveryMetrics[`FAMILY[`+v.INSTANCE_NAME+`]`] = v.FAMILY.String
		discoveryMetrics[`DATABASE_TYPE[`+v.INSTANCE_NAME+`]`] = v.DATABASE_TYPE
	}
	//glog.Info("discoveryMetrics: ", discoveryMetrics)
	//glog.Info("discoveryData: ", discoveryData)
	for k, v := range discoveryMetrics {
		zabbixData[k] = v
	}
	for k, v := range discoveryData {
		zabbixData[k] = v
	}
	//send(discoveryMetrics, zabbixHost, zabbixPort, hostName)
	glog.Info("zabbixData Combined: ", zabbixData)
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
		glog.Error("Error fetching addition", err, query)
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
		glog.Error("Error fetching addition", err)
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
		glog.Error("Error fetching addition", err, query)
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

func runInstanceMetrics(query string, db *sql.DB) []instance {
	var result []instance

	rows, err := db.Query(query)
	if err != nil {
		glog.Error("Error fetching addition", err, query)
		var er []string
		er = append(er, err.Error())
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var res instance
		err := rows.Scan(&res.INST_ID, &res.INSTANCE_NUMBER, &res.INSTANCE_NAME, &res.HOST_NAME, &res.VERSION, &res.STARTUP_TIME,
			&res.STATUS, &res.PARALLEL, &res.THREAD_NO, &res.ARCHIVER, &res.LOG_SWITCH_WAIT, &res.LOGINS,
			&res.SHUTDOWN_PENDING, &res.DATABASE_STATUS, &res.INSTANCE_ROLE, &res.ACTIVE_STATE, &res.BLOCKED, &res.CON_ID,
			&res.INSTANCE_MODE, &res.EDITION, &res.FAMILY, &res.DATABASE_TYPE)
		result = append(result, res)
		if err != nil {
			glog.Error(err)
		}

	}
	return result
}
