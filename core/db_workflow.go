package core

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
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
      job_ids               varchar[],
      create_time           bigint default extract(epoch from now())::bigint,
      update_time           bigint default extract(epoch from now())::bigint
  );`
	WorkflowTableColumn        = `id,name,description,job_ids,create_time,update_time`
	WorkflowTableColumnSize    = len(strings.Split(WorkflowTableColumn, ","))
	WorkflowTableOnConflictDDL = `
  on conflict (id) 
  do update set 
  name = excluded.name, 
  description = excluded.description, 
  job_ids = excluded.job_ids, 
  update_time = GREATEST(asset.update_time, excluded.update_time);`
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
		args = append(args, v.Id, v.Name, v.description, pq.Array(v.JobIds), createTime, updateTime)
	}
	query := fmt.Sprintf(`insert into %s (%s) values %s %s`, WorkflowTableName, WorkflowTableColumn,
		strings.Join(values, ","), WorkflowTableOnConflictDDL)
	return query, args, nil
}

func getWorkflowSql(ids []string) (string, []interface{}) {
	query := fmt.Sprintf("select * from %s where id in (%s);", WorkflowTableName, utils.ArrayToSqlIn(ids...))
	return query, nil
}

func listAssetsSql(in *ListWorkflowRequest) (string, []interface{}) {
	firstCond := true
	var buffer bytes.Buffer

	if in.Header > 0 {
		buffer.WriteString(fmt.Sprintf("select count(1) from %s", WorkflowTableName))
	} else {
		buffer.WriteString(fmt.Sprintf("select * from %s", WorkflowTableName))
	}

	if len(in.Names) > 0 {
		query := fmt.Sprintf("%s place_id in (%s)", utils.CondSql(firstCond), utils.ArrayToSqlIn(in.Names...))
		buffer.WriteString(query)
		firstCond = false
	}

	if len(in.Descriptions) > 0 {
		query := fmt.Sprintf("%s user_id in (%s)", utils.CondSql(firstCond), utils.ArrayToSqlIn(in.Descriptions...))
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
