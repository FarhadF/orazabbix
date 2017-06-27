package orametrics

var queries = map[string]string{
	"archivelog_switch":    "select count(*) from gv$log_history where first_time >= (sysdate - 1 / 24)",
	"uptime":               "select to_char ( (sysdate - startup_time) *24 , 'FM99999999999999990') from gv$instance",
	"dbblockgets":          "select sum (value) from gv$sysstat where name = 'db block gets'",
	"dblockchanges":        "select sum (value) from gv$sysstat where name = 'db block changes'",
	"dbconsistentgets":     "select sum (value) from gv$sysstat where name = 'consistent gets'",
	"dbphysicalreads":      "select sum (value) from gv$sysstat where name = 'physical reads'",
	"dbhitratio":           "select ( sum (case name when 'consistent gets' then value else 0 end) + sum (case name when 'db block gets' then value else 0 end) - sum (case name when 'physical reads' then value else 0 end)) / ( sum (case name when 'consistent gets' then value else 0 end) + sum (case name when 'db block gets' then value else 0 end)) * 100 from gv$sysstat",
	"hitratio_body":        "select gethitratio * 100 from gv$librarycache where namespace = 'BODY'",
	"hitratio_sqlarea":     "select gethitratio * 100 from gv$librarycache where namespace = 'SQL AREA'",
	"hitratio_trigger":     "select gethitratio * 100 from gv$librarycache where namespace = 'TRIGGER'",
	"hitratio_table_proc":  "select gethitratio * 100 from gv$librarycache where namespace = 'TABLE/PROCEDURE'",
	"miss_latch":           "select sum (misses) from gv$latch",
	"pga_aggregate_target": "select value from gv$pgastat where name = 'aggregate PGA target parameter'",
	"pga": "select value from gv$pgastat where name = 'total PGA inuse'",
	"phio_datafile_reads":    "select sum (value) from gv$sysstat where name = 'physical reads direct'",
	"phio_datafile_writes":   "select sum (value) from gv$sysstat where name = 'physical writes direct'",
	"phio_redo_writes":       "select sum (value) from gv$sysstat where name = 'redo writes'",
	"pinhitratio_body":       "select pins / (pins + reloads) * 100 from gv$librarycache where namespace = 'BODY'",
	"pinhitratio_sqlarea":    "select pins / (pins + reloads) * 100  from gv$librarycache  where namespace = 'SQL AREA'",
	"pinhitratio_trigger":    "select pins / (pins + reloads) * 100  from gv$librarycache  where namespace = 'TRIGGER'",
	"pinhitratio_table_proc": "select pins / (pins + reloads) * 100  from gv$librarycache  where namespace = 'TABLE/PROCEDURE'",
	//pool_dict_cache maybe empty, check and insert zero instead
	"pool_dict_cache": "select bytes  from gv$sgastat  where pool = 'shared pool' and name = 'dictionary cache'",
}
