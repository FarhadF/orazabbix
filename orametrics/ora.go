package orametrics

import (
	"database/sql"
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
		rows, err := db.Query(v)
		if err != nil {
			fmt.Println("Error fetching addition")
			fmt.Println(err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var sum string
			rows.Scan(&sum)
			fmt.Println(sum)
			zabbixData[k] = sum

		}
	}
	fmt.Println(zabbixData)
	send(zabbixData, zabbixHost, zabbixPort)

}
