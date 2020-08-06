package core

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"time"

	"github.com/linger1216/jelly-schedule/utils"
	"github.com/linger1216/jelly-schedule/utils/snowflake"
)

var (
	WorkflowTableName      = `workflow`
	CreateWorkflowTableDDL = `create table if not exists ` + WorkflowTableName + `
  (
      id                    varchar primary key,
      name                  varchar,
      description           varchar,
			job_ids               varchar,
			cron                  varchar,
			state                 varchar,
      create_time           bigint default extract(epoch from now())::bigint,
      update_time           bigint default extract(epoch from now())::bigint
  );`
	//WorkflowTableSelectColumn  = `id,name,description,array_to_string(, ',', ',') as job_ids,cron,create_time,update_time`
	WorkflowTableSelectColumn  = `*`
	WorkflowTableColumn        = `id,name,description,job_ids,cron,state,create_time,update_time`
	WorkflowTableColumnSize    = len(strings.Split(WorkflowTableColumn, ","))
	WorkflowTableOnConflictDDL = fmt.Sprintf(`
  on conflict (id) 
  do update set
  name = excluded.name,
  description = excluded.description, 
	job_ids = excluded.job_ids,
	cron = excluded.cron,
  state = excluded.state,
  update_time = GREATEST(%s.update_time, excluded.update_time);`, WorkflowTableName)
)

func createWorkflowTableSql() string {
	return CreateWorkflowTableDDL
}

func upsertWorkflowSql(workflows []*WorkFlow) (string, []interface{}, error) {
	size := len(workflows)
	if size == 0 {
		return "", nil, nil
	}

	values := make([]string, 0, size)
	args := make([]interface{}, 0, size*WorkflowTableColumnSize)

	var createTime, updateTime int64
	for i, v := range workflows {
		if v.CreateTime == 0 {
			createTime = time.Now().Unix()
		} else {
			createTime = v.CreateTime
		}

		if v.UpdateTime == 0 {
			updateTime = time.Now().Unix()
		} else {
			updateTime = v.UpdateTime
		}

		if len(v.Id) == 0 {
			v.Id = snowflake.Generate()
		}

		values = append(values, utils.ValueInject(i, WorkflowTableColumnSize))
		jsonBuf, err := jsoniter.ConfigFastest.Marshal(v.JobIds)
		if err != nil {
			return "", nil, err
		}
		args = append(args, v.Id, v.Name, v.Description, string(jsonBuf), v.Cron, v.State, createTime, updateTime)
	}
	query := fmt.Sprintf(`insert into %s (%s) values %s %s`, WorkflowTableName, WorkflowTableColumn,
		strings.Join(values, ","), WorkflowTableOnConflictDDL)
	return query, args, nil
}

func getWorkflowSql(ids []string) (string, []interface{}) {
	query := fmt.Sprintf("select %s from %s where id in (%s);", WorkflowTableSelectColumn, WorkflowTableName, utils.ArrayToSqlIn(ids...))
	return query, nil
}

func listAssetsSql(in *ListWorkflowRequest) (string, []interface{}) {
	firstCond := true
	var buffer bytes.Buffer

	if in.Header > 0 {
		buffer.WriteString(fmt.Sprintf("select count(1) from %s", WorkflowTableName))
	} else {
		buffer.WriteString(fmt.Sprintf("select %s from %s", WorkflowTableSelectColumn, WorkflowTableName))
	}

	if len(in.Names) > 0 {
		query := fmt.Sprintf("%s name in (%s)", utils.CondSql(firstCond), utils.ArrayToSqlIn(in.Names...))
		buffer.WriteString(query)
		firstCond = false
	}

	if len(in.Descriptions) > 0 {
		query := fmt.Sprintf("%s description in (%s)", utils.CondSql(firstCond), utils.ArrayToSqlIn(in.Descriptions...))
		buffer.WriteString(query)
		firstCond = false
	}

	if in.EndTime > 0 {
		query := fmt.Sprintf("%s (update_time >= '%d' and update_time < '%d') ", utils.CondSql(firstCond), in.StartTime, in.EndTime)
		buffer.WriteString(query)
		firstCond = false
	}

	if in.Header == 0 {
		query := fmt.Sprintf(" offset %d limit %d", in.CurrentPage*in.PageSize, in.PageSize)
		buffer.WriteString(query)
	}
	return buffer.String(), nil
}

func deleteWorkflowSql(ids []string) (string, []interface{}) {
	query := fmt.Sprintf("delete from %s where id in (%s);", WorkflowTableName, utils.ArrayToSqlIn(ids...))
	return query, nil
}
