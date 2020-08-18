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

/*
	// 执行几次结束
	SuccessLimit int64 `json:"executeLimit" yaml:"executeLimit" `
	// 碰到错误的方式
	ErrorPolicy string `json:"errorPolicy" yaml:"errorPolicy"`
	// 可以指定由哪个执行器执行
	BelongExecutor string `json:"belongExecutor" yaml:"belongExecutor" `
*/
var (
	WorkflowTableName      = `workflow`
	CreateWorkflowTableDDL = `create table if not exists ` + WorkflowTableName + `
  (
      id                    varchar primary key,
      name                  varchar,
      description           varchar,
			job_ids               varchar,
			cron                  varchar,
      para                  varchar,
      success_limit         int,
      failed_limit           int,
      belong_executor       varchar,
			state                 varchar,
      create_time           bigint default extract(epoch from now())::bigint,
      update_time           bigint default extract(epoch from now())::bigint
  );`
	WorkflowTableSelectColumn  = `*`
	WorkflowTableColumn        = `id,name,description,job_ids,cron,para,success_limit,failed_limit,belong_executor,state,create_time,update_time`
	WorkflowTableColumnSize    = len(strings.Split(WorkflowTableColumn, ","))
	WorkflowTableOnConflictDDL = fmt.Sprintf(`
  on conflict (id) 
  do update set
  name = excluded.name,
  description = excluded.description, 
	job_ids = excluded.job_ids,
	cron = excluded.cron,
  para = excluded.para,
  success_limit = excluded.success_limit,
  failed_limit = excluded.failed_limit,
  belong_executor = excluded.belong_executor,
  state = excluded.state,
  update_time = GREATEST(%s.update_time, excluded.update_time);`, WorkflowTableName)
)

func createWorkflowTableSql() string {
	return CreateWorkflowTableDDL
}

func upsertWorkflowSql(workflows ...*WorkFlow) (string, []interface{}, error) {
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
		args = append(args, v.Id, v.Name, v.Description, string(jsonBuf), v.Cron, v.Para, v.SuccessLimit, v.FailedLimit,
			v.BelongExecutor, v.State, createTime, updateTime)
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

	if len(in.States) > 0 {
		query := fmt.Sprintf("%s state in (%s)", utils.CondSql(firstCond), utils.ArrayToSqlIn(in.States...))
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

// 行级锁
// 用于executor获得任务
func getWorkFLowForUpdate(state string, n int) string {
	return fmt.Sprintf(`select * from %s where state = '%s' limit %d for update;`, WorkflowTableName, state, n)
}

func getWorkFLowByExecutorBelongForUpdate(belong, state string, n int) string {
	return fmt.Sprintf(`select * from %s where state = '%s and executor_belong = %s' limit %d for update;`,
		WorkflowTableName, state, belong, n)
}

// helper
