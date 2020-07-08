package orametrics

var queries = map[string]string{
	"alive":                  "select 1 from dual",
	"archivelog_switch":      "select count(*) from gv$log_history where first_time >= (sysdate - 1 / 24)",
	"uptime":                 "select to_char ( (sysdate - startup_time) *24*60 , 'FM99999999999999990') from gv$instance",
	"dbblockgets":            "select sum (value) from gv$sysstat where name = 'db block gets'",
	"dbblockchanges":         "select sum (value) from gv$sysstat where name = 'db block changes'",
	"dbconsistentgets":       "select sum (value) from gv$sysstat where name = 'consistent gets'",
	"dbphysicalreads":        "select sum (value) from gv$sysstat where name = 'physical reads'",
	"dbhitratio":             "select ( sum (case name when 'consistent gets' then value else 0 end) + sum (case name when 'db block gets' then value else 0 end) - sum (case name when 'physical reads' then value else 0 end)) / ( sum (case name when 'consistent gets' then value else 0 end) + sum (case name when 'db block gets' then value else 0 end)) * 100 from gv$sysstat",
	"hitratio_body":          "select gethitratio * 100 from gv$librarycache where namespace = 'BODY'",
	"hitratio_sqlarea":       "select gethitratio * 100 from gv$librarycache where namespace = 'SQL AREA'",
	"hitratio_trigger":       "select gethitratio * 100 from gv$librarycache where namespace = 'TRIGGER'",
	"hitratio_table_proc":    "select gethitratio * 100 from gv$librarycache where namespace = 'TABLE/PROCEDURE'",
	"miss_latch":             "select sum (misses) from gv$latch",
	"pga_aggregate_target":   "select value from gv$pgastat where name = 'aggregate PGA target parameter'",
	"pga":                    "select value from gv$pgastat where name = 'total PGA inuse'",
	"phio_datafile_reads":    "select sum (value) from gv$sysstat where name = 'physical reads direct'",
	"phio_datafile_writes":   "select sum (value) from gv$sysstat where name = 'physical writes direct'",
	"phio_redo_writes":       "select sum (value) from gv$sysstat where name = 'redo writes'",
	"pinhitratio_body":       "select pins / (pins + reloads) * 100 from gv$librarycache where namespace = 'BODY'",
	"pinhitratio_sqlarea":    "select pins / (pins + reloads) * 100  from gv$librarycache  where namespace = 'SQL AREA'",
	"pinhitratio_trigger":    "select pins / (pins + reloads) * 100  from gv$librarycache  where namespace = 'TRIGGER'",
	"pinhitratio_table_proc": "select pins / (pins + reloads) * 100  from gv$librarycache  where namespace = 'TABLE/PROCEDURE'",
	//pool_dict_cache maybe empty, check and insert zero instead
	"pool_dict_cache": "select bytes  from gv$sgastat  where pool = 'shared pool' and name = 'dictionary cache'",
	"pool_free_mem":   "select bytes from gv$sgastat where pool = 'shared pool' and name = 'free memory'",
	//pool_lib_cache maybe empty, check and insert zero instead
	"pool_lib_cache": "select bytes from gv$sgastat where pool = 'shared pool' and name = 'library cache'",
	//pool_sql_area maybe empty, check and insert zero instead
	"pool_sql_area":          "select bytes from gv$sgastat where pool = 'shared pool' and name = 'sql area'",
	"pool_misc":              "select sum (bytes) from gv$sgastat where pool = 'shared pool' and name not in ('library cache', 'dictionary cache', 'free memory', 'sql area')",
	"maxprocs":               "select value from gv$parameter where name = 'processes'",
	"procnum":                "select count (*) from gv$process",
	"maxsession":             "select value from gv$parameter where name = 'sessions'",
	"session":                "select count (*) from gv$session",
	"session_system":         "select count (*) from gv$session where type = 'BACKGROUND'",
	"session_active":         "select count (*) from gv$session where type != 'BACKGROUND' and status = 'ACTIVE'",
	"session_inactive":       "select count (*) from gv$session where type != 'BACKGROUND' and status = 'INACTIVE'",
	"sga_buffer_cache":       "select sum (bytes) from gv$sgastat where name in ('db_block_buffers', 'buffer_cache')",
	"sga_fixed":              "select sum (bytes) from gv$sgastat where name = 'fixed_sga'",
	"sga_java_pool":          "select sum (bytes) from gv$sgastat where pool = 'java pool'",
	"sga_large_pool":         "select sum (bytes) from gv$sgastat where pool = 'large pool'",
	"sga_shared_pool":        "select sum (bytes) from gv$sgastat where pool = 'shared pool'",
	"sga_log_buffer":         "select sum (bytes) from gv$sgastat where name = 'log_buffer'",
	"waits_directpath_read":  "select total_waits from gv$system_event where event = 'direct path read'",
	"waits_file_io":          "select nvl (sum (total_waits), 0) from gv$system_event where event in ('file identify', 'file open')",
	"waits_controlfileio":    "select sum (total_waits) from gv$system_event where event in ('control file sequential read' , 'control file single write' , 'control file parallel write')",
	"waits_logwrite":         "select sum (total_waits) from gv$system_event where event in ('log file single write', 'log file parallel write')",
	"waits_logsync":          "select sum(total_waits) from gv$system_event where event = 'log file sync'",
	"waits_multiblock_read":  "select sum (total_waits) from gv$system_event where event = 'db file scattered read'",
	"waits_singleblock_read": "select sum (total_waits) from gv$system_event where event = 'db file sequential read'",
	//waits_sqlnet maybe empty, check and insert zero instead
	"waits_sqlnet": "select count(*)  from (select rootid from (select level lvl   , connect_by_root (inst_id || '.' || sid) rootid   , seconds_in_wait from gv$session start with blocking_session is null connect by nocycle prior inst_id = blocking_instance  and prior sid = blocking_session) where lvl > 1 group by rootid having sum(seconds_in_wait) > 300)",
	//blocking_sessions maybe empty, check and insert zero instead
	"blocking_sessions": `select count(*) from (select rootid from (select level lvl, connect_by_root (inst_id || '.' || sid) rootid, seconds_in_wait from gv$session start with blocking_session is null connect by nocycle prior inst_id = blocking_instance and prior sid = blocking_session) where lvl > 1 group by rootid having sum(seconds_in_wait) > 300)`,
	"blocking_sessions_full": `select    lpad(' ', (level - 1) * 4)
		           || 'INST_ID         :  '
		           || inst_id
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'SERVICE_NAME    :  '
		           || service_name
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'SID,SERIAL      :  '
		           || sid
		           || ','
		           || serial#
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'USERNAME        :  '
		           || username
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'OSUSER          :  '
		           || osuser
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'MACHINE         :  '
		           || machine
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'PROGRAM         :  '
		           || program
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'MODULE          :  '
		           || module
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'SQL_ID          :  '
		           || sql_id
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'EVENT           :  '
		           || event
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'SECONDS_IN_WAIT :  '
		           || seconds_in_wait
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'STATE           :  '
		           || state
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || 'STATUS          :  '
		           || status
		           || chr(10)
		           || lpad(' ', (level - 1) * 4)
		           || '========================='
		           || chr(10)
		              blocking_sess_info
		      from (
		                  select inst_id || '.' || sid id
		                       , case
		                            when blocking_instance is not null
		                            then
		                               blocking_instance || '.' || blocking_session
		                         end
		                            parent_id
		                       , inst_id
		                       , service_name
		                       , sid
		                       , serial#
		                       , username
		                       , osuser
		                       , machine
		                       , program
		                       , module
		                       , sql_id
		                       , event
		                       , seconds_in_wait
		                       , state
		                       , status
		                       , level lvl
		                       , connect_by_isleaf isleaf
		                       , connect_by_root (inst_id || '.' || sid) rootid
		                    from gv$session
		              start with blocking_session is null
		              connect by nocycle prior inst_id = blocking_instance
		                             and prior sid = blocking_session
		           )
		     where lvl || isleaf <> '11'
		       and rootid in
		              (
		                   select rootid
		                     from (
		                                 select level lvl
		                                      , connect_by_root (inst_id || '.' || sid) rootid
		                                      , seconds_in_wait
		                                   from gv$session
		                             start with blocking_session is null
		                             connect by nocycle prior inst_id = blocking_instance
		                                            and prior sid = blocking_session
		                          )
		                    where lvl > 1
		                 group by rootid
		                   having sum(seconds_in_wait) > 300
		              )
		connect by nocycle prior id = parent_id
		start with parent_id is null`,
	"dbversion": "select banner from gv$version where banner like '%Oracle Database%'",
}
var discoveryQueries = map[string]string{
	"tablespaces": "select name ts from gv$tablespace",
	"diskgroups":  "select name from v$asm_diskgroup",
	"instances":   "select instance_name from gv$instance",
}

var (
	ts_usage_pct      string = "select tablespace_name ts, round(used_percent, 5) pct from dba_tablespace_usage_metrics"
	ts_usage_bytes    string = "select ta.tablespace_name as ts, ta.used_space * tb.block_size as bytes from dba_tablespace_usage_metrics ta join dba_tablespaces tb on ta.tablespace_name = tb.tablespace_name"
	diskgroup_metrics string = "select name as Dg,USABLE_FILE_MB as UsableFileMB, OFFLINE_DISKS as UsableFileMB from v$asm_diskgroup"
	instance_metrics  string = "select INST_ID, INSTANCE_NUMBER, INSTANCE_NAME, HOST_NAME, VERSION, STARTUP_TIME, STATUS, PARALLEL, THREAD# AS THREAD_NO, ARCHIVER, LOG_SWITCH_WAIT, LOGINS, SHUTDOWN_PENDING, DATABASE_STATUS, INSTANCE_ROLE, ACTIVE_STATE, BLOCKED, CON_ID, INSTANCE_MODE, EDITION, FAMILY, DATABASE_TYPE from gv$instance"
)
