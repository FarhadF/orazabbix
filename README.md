# orazabbix
Golang implementation of oracle database monitoring on zabbix

## Build Steps:
1. Clone repo
2. Setup oracle instant client and environment variables. [Guide](https://gocodecloud.com/blog/2016/08/09/accessing-an-oracle-db-in-go/)
3. `go build main.go`

```
Usage:
  orazabbix [flags]

Flags:
  -c, --connectionstring string   ConnectionString to the Database, Format: username/password@ip:port/sid (default "system/oracle@localhost:1521/xe")
  -h, --help                      help for orazabbix
  -H, --host string               Hostname of the monitored object in zabbix server (default "server1")
  -p, --port int                  Zabbix Server/Proxy Port (default 10051)
  -v, --version                   Prints version information
  -z, --zabbix string             Zabbix Server/Proxy Hostname or IP address (default "localhost")
  ```
  
  ## Installation:
  1. Import template file in zabbix server.
  2. Add cron entry:
  
  ```
  * * * * * oracle source /home/oracle/.bash_profile;/home/oracle/main -c <user>/<password>@<instance ip or hostname>:<port>/<sid> -z <zabbix server ip> -p <zabbix server(trapper) port>-H <hostname in zabbix server> 2>&1 >/dev/null
  ```
  3. Restart cron service
  4. Latest data in zabbix frontend should start populating after a minute.
  
