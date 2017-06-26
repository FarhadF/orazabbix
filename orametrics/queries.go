package orametrics

var queries = map[string]string{
	"archivelog_switch": "select count(*) from gv$log_history where first_time >= (sysdate - 1 / 24)",
	"uptime":            "select to_char ( (sysdate - startup_time) * 86400, 'FM99999999999999990') from gv$instance",
	"dbblockgets":       "select sum (value) from gv$sysstat where name = 'db block gets'",
}
