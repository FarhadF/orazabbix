package orametrics

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-oci8"
)

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
	discoveryData := make(map[string]map[string][]string)
	for k, v := range discoveryQueries {
		middle := make(map[string][]string)
		result := runDiscoveryQuery(v, db)
		length := len(result)
		fix := make([]string, length)
		for _, value := range result {

			metric := runQuery(ts_usage_pct, db)
			f := value + ":" + metric
			fix = append(fix, f)

		}
		middle["data"] = result
		discoveryData[k] = middle
	}
	j, _ := json.Marshal(discoveryData)
	fmt.Println(string(j))
	send(zabbixData, zabbixHost, zabbixPort, hostName)
	sendD(discoveryData, zabbixHost, zabbixPort, hostName)
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
