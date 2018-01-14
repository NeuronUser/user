package user_db

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/sql/wrap"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"strings"
	"time"
)

var _ = sql.ErrNoRows
var _ = mysql.ErrOldProtocol

type BaseQuery struct {
	forUpdate     bool
	forShare      bool
	where         string
	limit         string
	order         string
	groupByFields []string
}

func (q *BaseQuery) buildQueryString() string {
	buf := bytes.NewBufferString("")

	if q.where != "" {
		buf.WriteString(" WHERE ")
		buf.WriteString(q.where)
	}

	if q.groupByFields != nil && len(q.groupByFields) > 0 {
		buf.WriteString(" GROUP BY ")
		buf.WriteString(strings.Join(q.groupByFields, ","))
	}

	if q.order != "" {
		buf.WriteString(" order by ")
		buf.WriteString(q.order)
	}

	if q.limit != "" {
		buf.WriteString(q.limit)
	}

	if q.forUpdate {
		buf.WriteString(" FOR UPDATE ")
	}

	if q.forShare {
		buf.WriteString(" LOCK IN SHARE MODE ")
	}

	return buf.String()
}

const OAUTH_STATE_TABLE_NAME = "oauth_state"

type OAUTH_STATE_FIELD string

const OAUTH_STATE_FIELD_ID = OAUTH_STATE_FIELD("id")
const OAUTH_STATE_FIELD_OAUTH_STATE = OAUTH_STATE_FIELD("oauth_state")
const OAUTH_STATE_FIELD_STATE_USED = OAUTH_STATE_FIELD("state_used")
const OAUTH_STATE_FIELD_CREATE_TIME = OAUTH_STATE_FIELD("create_time")
const OAUTH_STATE_FIELD_UPDATE_TIME = OAUTH_STATE_FIELD("update_time")
const OAUTH_STATE_FIELD_QUERY_STRING = OAUTH_STATE_FIELD("query_string")

const OAUTH_STATE_ALL_FIELDS_STRING = "id,oauth_state,state_used,create_time,update_time,query_string"

var OAUTH_STATE_ALL_FIELDS = []string{
	"id",
	"oauth_state",
	"state_used",
	"create_time",
	"update_time",
	"query_string",
}

type OauthState struct {
	Id          uint64 //size=20
	OauthState  string //size=128
	StateUsed   int32  //size=1
	CreateTime  time.Time
	UpdateTime  time.Time
	QueryString string //size=256
}

type OauthStateQuery struct {
	BaseQuery
	dao *OauthStateDao
}

func NewOauthStateQuery(dao *OauthStateDao) *OauthStateQuery {
	q := &OauthStateQuery{}
	q.dao = dao

	return q
}

func (q *OauthStateQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*OauthState, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *OauthStateQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*OauthState, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *OauthStateQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *OauthStateQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *OauthStateQuery) ForUpdate() *OauthStateQuery {
	q.forUpdate = true
	return q
}

func (q *OauthStateQuery) ForShare() *OauthStateQuery {
	q.forShare = true
	return q
}

func (q *OauthStateQuery) GroupBy(fields ...OAUTH_STATE_FIELD) *OauthStateQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *OauthStateQuery) Limit(startIncluded int64, count int64) *OauthStateQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *OauthStateQuery) OrderBy(fieldName OAUTH_STATE_FIELD, asc bool) *OauthStateQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *OauthStateQuery) OrderByGroupCount(asc bool) *OauthStateQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *OauthStateQuery) w(format string, a ...interface{}) *OauthStateQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *OauthStateQuery) Left() *OauthStateQuery  { return q.w(" ( ") }
func (q *OauthStateQuery) Right() *OauthStateQuery { return q.w(" ) ") }
func (q *OauthStateQuery) And() *OauthStateQuery   { return q.w(" AND ") }
func (q *OauthStateQuery) Or() *OauthStateQuery    { return q.w(" OR ") }
func (q *OauthStateQuery) Not() *OauthStateQuery   { return q.w(" NOT ") }

func (q *OauthStateQuery) Id_Equal(v uint64) *OauthStateQuery { return q.w("id='" + fmt.Sprint(v) + "'") }
func (q *OauthStateQuery) Id_NotEqual(v uint64) *OauthStateQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) Id_Less(v uint64) *OauthStateQuery { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *OauthStateQuery) Id_LessEqual(v uint64) *OauthStateQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) Id_Greater(v uint64) *OauthStateQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) Id_GreaterEqual(v uint64) *OauthStateQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_Equal(v string) *OauthStateQuery {
	return q.w("oauth_state='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_NotEqual(v string) *OauthStateQuery {
	return q.w("oauth_state<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_Less(v string) *OauthStateQuery {
	return q.w("oauth_state<'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_LessEqual(v string) *OauthStateQuery {
	return q.w("oauth_state<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_Greater(v string) *OauthStateQuery {
	return q.w("oauth_state>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) OauthState_GreaterEqual(v string) *OauthStateQuery {
	return q.w("oauth_state>='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) StateUsed_Equal(v int32) *OauthStateQuery {
	return q.w("state_used='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) StateUsed_NotEqual(v int32) *OauthStateQuery {
	return q.w("state_used<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) StateUsed_Less(v int32) *OauthStateQuery {
	return q.w("state_used<'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) StateUsed_LessEqual(v int32) *OauthStateQuery {
	return q.w("state_used<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) StateUsed_Greater(v int32) *OauthStateQuery {
	return q.w("state_used>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) StateUsed_GreaterEqual(v int32) *OauthStateQuery {
	return q.w("state_used>='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_Equal(v time.Time) *OauthStateQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_NotEqual(v time.Time) *OauthStateQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_Less(v time.Time) *OauthStateQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_LessEqual(v time.Time) *OauthStateQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_Greater(v time.Time) *OauthStateQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) CreateTime_GreaterEqual(v time.Time) *OauthStateQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_Equal(v time.Time) *OauthStateQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_NotEqual(v time.Time) *OauthStateQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_Less(v time.Time) *OauthStateQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_LessEqual(v time.Time) *OauthStateQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_Greater(v time.Time) *OauthStateQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) UpdateTime_GreaterEqual(v time.Time) *OauthStateQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) QueryString_Equal(v string) *OauthStateQuery {
	return q.w("query_string='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) QueryString_NotEqual(v string) *OauthStateQuery {
	return q.w("query_string<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) QueryString_Less(v string) *OauthStateQuery {
	return q.w("query_string<'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) QueryString_LessEqual(v string) *OauthStateQuery {
	return q.w("query_string<='" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) QueryString_Greater(v string) *OauthStateQuery {
	return q.w("query_string>'" + fmt.Sprint(v) + "'")
}
func (q *OauthStateQuery) QueryString_GreaterEqual(v string) *OauthStateQuery {
	return q.w("query_string>='" + fmt.Sprint(v) + "'")
}

type OauthStateDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewOauthStateDao(db *DB) (t *OauthStateDao, err error) {
	t = &OauthStateDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *OauthStateDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}
func (dao *OauthStateDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO oauth_state (oauth_state,state_used,query_string) VALUES (?,?,?)")
	return err
}

func (dao *OauthStateDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE oauth_state SET oauth_state=?,state_used=?,query_string=? WHERE id=?")
	return err
}

func (dao *OauthStateDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM oauth_state WHERE id=?")
	return err
}

func (dao *OauthStateDao) Insert(ctx context.Context, tx *wrap.Tx, e *OauthState) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.OauthState, e.StateUsed, e.QueryString)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *OauthStateDao) Update(ctx context.Context, tx *wrap.Tx, e *OauthState) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.OauthState, e.StateUsed, e.QueryString, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthStateDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthStateDao) scanRow(row *wrap.Row) (*OauthState, error) {
	e := &OauthState{}
	err := row.Scan(&e.Id, &e.OauthState, &e.StateUsed, &e.CreateTime, &e.UpdateTime, &e.QueryString)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *OauthStateDao) scanRows(rows *wrap.Rows) (list []*OauthState, err error) {
	list = make([]*OauthState, 0)
	for rows.Next() {
		e := OauthState{}
		err = rows.Scan(&e.Id, &e.OauthState, &e.StateUsed, &e.CreateTime, &e.UpdateTime, &e.QueryString)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *OauthStateDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*OauthState, error) {
	querySql := "SELECT " + OAUTH_STATE_ALL_FIELDS_STRING + " FROM oauth_state " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *OauthStateDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*OauthState, err error) {
	querySql := "SELECT " + OAUTH_STATE_ALL_FIELDS_STRING + " FROM oauth_state " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *OauthStateDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM oauth_state " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *OauthStateDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM oauth_state " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *OauthStateDao) GetQuery() *OauthStateQuery {
	return NewOauthStateQuery(dao)
}

const OAUTH_TOKENS_TABLE_NAME = "oauth_tokens"

type OAUTH_TOKENS_FIELD string

const OAUTH_TOKENS_FIELD_ID = OAUTH_TOKENS_FIELD("id")
const OAUTH_TOKENS_FIELD_ACCOUNT_ID = OAUTH_TOKENS_FIELD("account_id")
const OAUTH_TOKENS_FIELD_AUTHORIZATION_CODE = OAUTH_TOKENS_FIELD("authorization_code")
const OAUTH_TOKENS_FIELD_ACCESS_TOKEN = OAUTH_TOKENS_FIELD("access_token")
const OAUTH_TOKENS_FIELD_REFRESH_TOKEN = OAUTH_TOKENS_FIELD("refresh_token")
const OAUTH_TOKENS_FIELD_CREATE_TIME = OAUTH_TOKENS_FIELD("create_time")
const OAUTH_TOKENS_FIELD_UPDATE_TIME = OAUTH_TOKENS_FIELD("update_time")

const OAUTH_TOKENS_ALL_FIELDS_STRING = "id,account_id,authorization_code,access_token,refresh_token,create_time,update_time"

var OAUTH_TOKENS_ALL_FIELDS = []string{
	"id",
	"account_id",
	"authorization_code",
	"access_token",
	"refresh_token",
	"create_time",
	"update_time",
}

type OauthTokens struct {
	Id                uint64 //size=20
	AccountId         string //size=128
	AuthorizationCode string //size=128
	AccessToken       string //size=1024
	RefreshToken      string //size=128
	CreateTime        time.Time
	UpdateTime        time.Time
}

type OauthTokensQuery struct {
	BaseQuery
	dao *OauthTokensDao
}

func NewOauthTokensQuery(dao *OauthTokensDao) *OauthTokensQuery {
	q := &OauthTokensQuery{}
	q.dao = dao

	return q
}

func (q *OauthTokensQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*OauthTokens, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *OauthTokensQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*OauthTokens, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *OauthTokensQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *OauthTokensQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *OauthTokensQuery) ForUpdate() *OauthTokensQuery {
	q.forUpdate = true
	return q
}

func (q *OauthTokensQuery) ForShare() *OauthTokensQuery {
	q.forShare = true
	return q
}

func (q *OauthTokensQuery) GroupBy(fields ...OAUTH_TOKENS_FIELD) *OauthTokensQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *OauthTokensQuery) Limit(startIncluded int64, count int64) *OauthTokensQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *OauthTokensQuery) OrderBy(fieldName OAUTH_TOKENS_FIELD, asc bool) *OauthTokensQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *OauthTokensQuery) OrderByGroupCount(asc bool) *OauthTokensQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *OauthTokensQuery) w(format string, a ...interface{}) *OauthTokensQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *OauthTokensQuery) Left() *OauthTokensQuery  { return q.w(" ( ") }
func (q *OauthTokensQuery) Right() *OauthTokensQuery { return q.w(" ) ") }
func (q *OauthTokensQuery) And() *OauthTokensQuery   { return q.w(" AND ") }
func (q *OauthTokensQuery) Or() *OauthTokensQuery    { return q.w(" OR ") }
func (q *OauthTokensQuery) Not() *OauthTokensQuery   { return q.w(" NOT ") }

func (q *OauthTokensQuery) Id_Equal(v uint64) *OauthTokensQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) Id_NotEqual(v uint64) *OauthTokensQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) Id_Less(v uint64) *OauthTokensQuery { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *OauthTokensQuery) Id_LessEqual(v uint64) *OauthTokensQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) Id_Greater(v uint64) *OauthTokensQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) Id_GreaterEqual(v uint64) *OauthTokensQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccountId_Equal(v string) *OauthTokensQuery {
	return q.w("account_id='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccountId_NotEqual(v string) *OauthTokensQuery {
	return q.w("account_id<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccountId_Less(v string) *OauthTokensQuery {
	return q.w("account_id<'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccountId_LessEqual(v string) *OauthTokensQuery {
	return q.w("account_id<='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccountId_Greater(v string) *OauthTokensQuery {
	return q.w("account_id>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccountId_GreaterEqual(v string) *OauthTokensQuery {
	return q.w("account_id>='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AuthorizationCode_Equal(v string) *OauthTokensQuery {
	return q.w("authorization_code='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AuthorizationCode_NotEqual(v string) *OauthTokensQuery {
	return q.w("authorization_code<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AuthorizationCode_Less(v string) *OauthTokensQuery {
	return q.w("authorization_code<'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AuthorizationCode_LessEqual(v string) *OauthTokensQuery {
	return q.w("authorization_code<='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AuthorizationCode_Greater(v string) *OauthTokensQuery {
	return q.w("authorization_code>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AuthorizationCode_GreaterEqual(v string) *OauthTokensQuery {
	return q.w("authorization_code>='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccessToken_Equal(v string) *OauthTokensQuery {
	return q.w("access_token='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccessToken_NotEqual(v string) *OauthTokensQuery {
	return q.w("access_token<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccessToken_Less(v string) *OauthTokensQuery {
	return q.w("access_token<'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccessToken_LessEqual(v string) *OauthTokensQuery {
	return q.w("access_token<='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccessToken_Greater(v string) *OauthTokensQuery {
	return q.w("access_token>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) AccessToken_GreaterEqual(v string) *OauthTokensQuery {
	return q.w("access_token>='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) RefreshToken_Equal(v string) *OauthTokensQuery {
	return q.w("refresh_token='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) RefreshToken_NotEqual(v string) *OauthTokensQuery {
	return q.w("refresh_token<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) RefreshToken_Less(v string) *OauthTokensQuery {
	return q.w("refresh_token<'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) RefreshToken_LessEqual(v string) *OauthTokensQuery {
	return q.w("refresh_token<='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) RefreshToken_Greater(v string) *OauthTokensQuery {
	return q.w("refresh_token>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) RefreshToken_GreaterEqual(v string) *OauthTokensQuery {
	return q.w("refresh_token>='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) CreateTime_Equal(v time.Time) *OauthTokensQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) CreateTime_NotEqual(v time.Time) *OauthTokensQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) CreateTime_Less(v time.Time) *OauthTokensQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) CreateTime_LessEqual(v time.Time) *OauthTokensQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) CreateTime_Greater(v time.Time) *OauthTokensQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) CreateTime_GreaterEqual(v time.Time) *OauthTokensQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) UpdateTime_Equal(v time.Time) *OauthTokensQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) UpdateTime_NotEqual(v time.Time) *OauthTokensQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) UpdateTime_Less(v time.Time) *OauthTokensQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) UpdateTime_LessEqual(v time.Time) *OauthTokensQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) UpdateTime_Greater(v time.Time) *OauthTokensQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *OauthTokensQuery) UpdateTime_GreaterEqual(v time.Time) *OauthTokensQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type OauthTokensDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewOauthTokensDao(db *DB) (t *OauthTokensDao, err error) {
	t = &OauthTokensDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *OauthTokensDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}
func (dao *OauthTokensDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO oauth_tokens (account_id,authorization_code,access_token,refresh_token) VALUES (?,?,?,?)")
	return err
}

func (dao *OauthTokensDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE oauth_tokens SET account_id=?,authorization_code=?,access_token=?,refresh_token=? WHERE id=?")
	return err
}

func (dao *OauthTokensDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM oauth_tokens WHERE id=?")
	return err
}

func (dao *OauthTokensDao) Insert(ctx context.Context, tx *wrap.Tx, e *OauthTokens) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.AccountId, e.AuthorizationCode, e.AccessToken, e.RefreshToken)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *OauthTokensDao) Update(ctx context.Context, tx *wrap.Tx, e *OauthTokens) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.AccountId, e.AuthorizationCode, e.AccessToken, e.RefreshToken, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthTokensDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *OauthTokensDao) scanRow(row *wrap.Row) (*OauthTokens, error) {
	e := &OauthTokens{}
	err := row.Scan(&e.Id, &e.AccountId, &e.AuthorizationCode, &e.AccessToken, &e.RefreshToken, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *OauthTokensDao) scanRows(rows *wrap.Rows) (list []*OauthTokens, err error) {
	list = make([]*OauthTokens, 0)
	for rows.Next() {
		e := OauthTokens{}
		err = rows.Scan(&e.Id, &e.AccountId, &e.AuthorizationCode, &e.AccessToken, &e.RefreshToken, &e.CreateTime, &e.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *OauthTokensDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*OauthTokens, error) {
	querySql := "SELECT " + OAUTH_TOKENS_ALL_FIELDS_STRING + " FROM oauth_tokens " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *OauthTokensDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*OauthTokens, err error) {
	querySql := "SELECT " + OAUTH_TOKENS_ALL_FIELDS_STRING + " FROM oauth_tokens " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *OauthTokensDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM oauth_tokens " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *OauthTokensDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM oauth_tokens " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *OauthTokensDao) GetQuery() *OauthTokensQuery {
	return NewOauthTokensQuery(dao)
}

const REFRESH_TOKEN_TABLE_NAME = "refresh_token"

type REFRESH_TOKEN_FIELD string

const REFRESH_TOKEN_FIELD_ID = REFRESH_TOKEN_FIELD("id")
const REFRESH_TOKEN_FIELD_ACCOUNT_ID = REFRESH_TOKEN_FIELD("account_id")
const REFRESH_TOKEN_FIELD_REFRESH_TOKEN = REFRESH_TOKEN_FIELD("refresh_token")
const REFRESH_TOKEN_FIELD_CREATE_TIME = REFRESH_TOKEN_FIELD("create_time")
const REFRESH_TOKEN_FIELD_UPDATE_TIME = REFRESH_TOKEN_FIELD("update_time")

const REFRESH_TOKEN_ALL_FIELDS_STRING = "id,account_id,refresh_token,create_time,update_time"

var REFRESH_TOKEN_ALL_FIELDS = []string{
	"id",
	"account_id",
	"refresh_token",
	"create_time",
	"update_time",
}

type RefreshToken struct {
	Id           uint64 //size=20
	AccountId    string //size=128
	RefreshToken string //size=128
	CreateTime   time.Time
	UpdateTime   time.Time
}

type RefreshTokenQuery struct {
	BaseQuery
	dao *RefreshTokenDao
}

func NewRefreshTokenQuery(dao *RefreshTokenDao) *RefreshTokenQuery {
	q := &RefreshTokenQuery{}
	q.dao = dao

	return q
}

func (q *RefreshTokenQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*RefreshToken, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *RefreshTokenQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*RefreshToken, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *RefreshTokenQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *RefreshTokenQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *RefreshTokenQuery) ForUpdate() *RefreshTokenQuery {
	q.forUpdate = true
	return q
}

func (q *RefreshTokenQuery) ForShare() *RefreshTokenQuery {
	q.forShare = true
	return q
}

func (q *RefreshTokenQuery) GroupBy(fields ...REFRESH_TOKEN_FIELD) *RefreshTokenQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *RefreshTokenQuery) Limit(startIncluded int64, count int64) *RefreshTokenQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *RefreshTokenQuery) OrderBy(fieldName REFRESH_TOKEN_FIELD, asc bool) *RefreshTokenQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *RefreshTokenQuery) OrderByGroupCount(asc bool) *RefreshTokenQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *RefreshTokenQuery) w(format string, a ...interface{}) *RefreshTokenQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *RefreshTokenQuery) Left() *RefreshTokenQuery  { return q.w(" ( ") }
func (q *RefreshTokenQuery) Right() *RefreshTokenQuery { return q.w(" ) ") }
func (q *RefreshTokenQuery) And() *RefreshTokenQuery   { return q.w(" AND ") }
func (q *RefreshTokenQuery) Or() *RefreshTokenQuery    { return q.w(" OR ") }
func (q *RefreshTokenQuery) Not() *RefreshTokenQuery   { return q.w(" NOT ") }

func (q *RefreshTokenQuery) Id_Equal(v uint64) *RefreshTokenQuery {
	return q.w("id='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_NotEqual(v uint64) *RefreshTokenQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_Less(v uint64) *RefreshTokenQuery {
	return q.w("id<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_LessEqual(v uint64) *RefreshTokenQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_Greater(v uint64) *RefreshTokenQuery {
	return q.w("id>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) Id_GreaterEqual(v uint64) *RefreshTokenQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_Equal(v string) *RefreshTokenQuery {
	return q.w("account_id='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_NotEqual(v string) *RefreshTokenQuery {
	return q.w("account_id<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_Less(v string) *RefreshTokenQuery {
	return q.w("account_id<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_LessEqual(v string) *RefreshTokenQuery {
	return q.w("account_id<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_Greater(v string) *RefreshTokenQuery {
	return q.w("account_id>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) AccountId_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("account_id>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_Equal(v string) *RefreshTokenQuery {
	return q.w("refresh_token='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_NotEqual(v string) *RefreshTokenQuery {
	return q.w("refresh_token<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_Less(v string) *RefreshTokenQuery {
	return q.w("refresh_token<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_LessEqual(v string) *RefreshTokenQuery {
	return q.w("refresh_token<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_Greater(v string) *RefreshTokenQuery {
	return q.w("refresh_token>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) RefreshToken_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("refresh_token>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_Equal(v time.Time) *RefreshTokenQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_NotEqual(v time.Time) *RefreshTokenQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_Less(v time.Time) *RefreshTokenQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_LessEqual(v time.Time) *RefreshTokenQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_Greater(v time.Time) *RefreshTokenQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) CreateTime_GreaterEqual(v time.Time) *RefreshTokenQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_Equal(v time.Time) *RefreshTokenQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_NotEqual(v time.Time) *RefreshTokenQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_Less(v time.Time) *RefreshTokenQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_LessEqual(v time.Time) *RefreshTokenQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_Greater(v time.Time) *RefreshTokenQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UpdateTime_GreaterEqual(v time.Time) *RefreshTokenQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type RefreshTokenDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewRefreshTokenDao(db *DB) (t *RefreshTokenDao, err error) {
	t = &RefreshTokenDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *RefreshTokenDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}
func (dao *RefreshTokenDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO refresh_token (account_id,refresh_token) VALUES (?,?)")
	return err
}

func (dao *RefreshTokenDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE refresh_token SET account_id=?,refresh_token=? WHERE id=?")
	return err
}

func (dao *RefreshTokenDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM refresh_token WHERE id=?")
	return err
}

func (dao *RefreshTokenDao) Insert(ctx context.Context, tx *wrap.Tx, e *RefreshToken) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.AccountId, e.RefreshToken)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *RefreshTokenDao) Update(ctx context.Context, tx *wrap.Tx, e *RefreshToken) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.AccountId, e.RefreshToken, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *RefreshTokenDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *RefreshTokenDao) scanRow(row *wrap.Row) (*RefreshToken, error) {
	e := &RefreshToken{}
	err := row.Scan(&e.Id, &e.AccountId, &e.RefreshToken, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *RefreshTokenDao) scanRows(rows *wrap.Rows) (list []*RefreshToken, err error) {
	list = make([]*RefreshToken, 0)
	for rows.Next() {
		e := RefreshToken{}
		err = rows.Scan(&e.Id, &e.AccountId, &e.RefreshToken, &e.CreateTime, &e.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *RefreshTokenDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*RefreshToken, error) {
	querySql := "SELECT " + REFRESH_TOKEN_ALL_FIELDS_STRING + " FROM refresh_token " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *RefreshTokenDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*RefreshToken, err error) {
	querySql := "SELECT " + REFRESH_TOKEN_ALL_FIELDS_STRING + " FROM refresh_token " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *RefreshTokenDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM refresh_token " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *RefreshTokenDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM refresh_token " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *RefreshTokenDao) GetQuery() *RefreshTokenQuery {
	return NewRefreshTokenQuery(dao)
}

const USER_TOKEN_TABLE_NAME = "user_token"

type USER_TOKEN_FIELD string

const USER_TOKEN_FIELD_ID = USER_TOKEN_FIELD("id")
const USER_TOKEN_FIELD_ACCOUNT_ID = USER_TOKEN_FIELD("account_id")
const USER_TOKEN_FIELD_EXPIRES_TIME = USER_TOKEN_FIELD("expires_time")
const USER_TOKEN_FIELD_USER_TOKEN = USER_TOKEN_FIELD("user_token")
const USER_TOKEN_FIELD_CREATE_TIME = USER_TOKEN_FIELD("create_time")
const USER_TOKEN_FIELD_UPDATE_TIME = USER_TOKEN_FIELD("update_time")

const USER_TOKEN_ALL_FIELDS_STRING = "id,account_id,expires_time,user_token,create_time,update_time"

var USER_TOKEN_ALL_FIELDS = []string{
	"id",
	"account_id",
	"expires_time",
	"user_token",
	"create_time",
	"update_time",
}

type UserToken struct {
	Id          uint64 //size=20
	AccountId   string //size=128
	ExpiresTime time.Time
	UserToken   string //size=1024
	CreateTime  time.Time
	UpdateTime  time.Time
}

type UserTokenQuery struct {
	BaseQuery
	dao *UserTokenDao
}

func NewUserTokenQuery(dao *UserTokenDao) *UserTokenQuery {
	q := &UserTokenQuery{}
	q.dao = dao

	return q
}

func (q *UserTokenQuery) QueryOne(ctx context.Context, tx *wrap.Tx) (*UserToken, error) {
	return q.dao.QueryOne(ctx, tx, q.buildQueryString())
}

func (q *UserTokenQuery) QueryList(ctx context.Context, tx *wrap.Tx) (list []*UserToken, err error) {
	return q.dao.QueryList(ctx, tx, q.buildQueryString())
}

func (q *UserTokenQuery) QueryCount(ctx context.Context, tx *wrap.Tx) (count int64, err error) {
	return q.dao.QueryCount(ctx, tx, q.buildQueryString())
}

func (q *UserTokenQuery) QueryGroupBy(ctx context.Context, tx *wrap.Tx) (rows *wrap.Rows, err error) {
	return q.dao.QueryGroupBy(ctx, tx, q.groupByFields, q.buildQueryString())
}

func (q *UserTokenQuery) ForUpdate() *UserTokenQuery {
	q.forUpdate = true
	return q
}

func (q *UserTokenQuery) ForShare() *UserTokenQuery {
	q.forShare = true
	return q
}

func (q *UserTokenQuery) GroupBy(fields ...USER_TOKEN_FIELD) *UserTokenQuery {
	q.groupByFields = make([]string, len(fields))
	for i, v := range fields {
		q.groupByFields[i] = string(v)
	}
	return q
}

func (q *UserTokenQuery) Limit(startIncluded int64, count int64) *UserTokenQuery {
	q.limit = fmt.Sprintf(" limit %d,%d", startIncluded, count)
	return q
}

func (q *UserTokenQuery) OrderBy(fieldName USER_TOKEN_FIELD, asc bool) *UserTokenQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += string(fieldName) + " "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *UserTokenQuery) OrderByGroupCount(asc bool) *UserTokenQuery {
	if q.order != "" {
		q.order += ","
	}
	q.order += "count(1) "
	if asc {
		q.order += "asc"
	} else {
		q.order += "desc"
	}

	return q
}

func (q *UserTokenQuery) w(format string, a ...interface{}) *UserTokenQuery {
	q.where += fmt.Sprintf(format, a...)
	return q
}

func (q *UserTokenQuery) Left() *UserTokenQuery  { return q.w(" ( ") }
func (q *UserTokenQuery) Right() *UserTokenQuery { return q.w(" ) ") }
func (q *UserTokenQuery) And() *UserTokenQuery   { return q.w(" AND ") }
func (q *UserTokenQuery) Or() *UserTokenQuery    { return q.w(" OR ") }
func (q *UserTokenQuery) Not() *UserTokenQuery   { return q.w(" NOT ") }

func (q *UserTokenQuery) Id_Equal(v uint64) *UserTokenQuery { return q.w("id='" + fmt.Sprint(v) + "'") }
func (q *UserTokenQuery) Id_NotEqual(v uint64) *UserTokenQuery {
	return q.w("id<>'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) Id_Less(v uint64) *UserTokenQuery { return q.w("id<'" + fmt.Sprint(v) + "'") }
func (q *UserTokenQuery) Id_LessEqual(v uint64) *UserTokenQuery {
	return q.w("id<='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) Id_Greater(v uint64) *UserTokenQuery { return q.w("id>'" + fmt.Sprint(v) + "'") }
func (q *UserTokenQuery) Id_GreaterEqual(v uint64) *UserTokenQuery {
	return q.w("id>='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) AccountId_Equal(v string) *UserTokenQuery {
	return q.w("account_id='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) AccountId_NotEqual(v string) *UserTokenQuery {
	return q.w("account_id<>'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) AccountId_Less(v string) *UserTokenQuery {
	return q.w("account_id<'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) AccountId_LessEqual(v string) *UserTokenQuery {
	return q.w("account_id<='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) AccountId_Greater(v string) *UserTokenQuery {
	return q.w("account_id>'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) AccountId_GreaterEqual(v string) *UserTokenQuery {
	return q.w("account_id>='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) ExpiresTime_Equal(v time.Time) *UserTokenQuery {
	return q.w("expires_time='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) ExpiresTime_NotEqual(v time.Time) *UserTokenQuery {
	return q.w("expires_time<>'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) ExpiresTime_Less(v time.Time) *UserTokenQuery {
	return q.w("expires_time<'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) ExpiresTime_LessEqual(v time.Time) *UserTokenQuery {
	return q.w("expires_time<='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) ExpiresTime_Greater(v time.Time) *UserTokenQuery {
	return q.w("expires_time>'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) ExpiresTime_GreaterEqual(v time.Time) *UserTokenQuery {
	return q.w("expires_time>='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UserToken_Equal(v string) *UserTokenQuery {
	return q.w("user_token='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UserToken_NotEqual(v string) *UserTokenQuery {
	return q.w("user_token<>'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UserToken_Less(v string) *UserTokenQuery {
	return q.w("user_token<'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UserToken_LessEqual(v string) *UserTokenQuery {
	return q.w("user_token<='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UserToken_Greater(v string) *UserTokenQuery {
	return q.w("user_token>'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UserToken_GreaterEqual(v string) *UserTokenQuery {
	return q.w("user_token>='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) CreateTime_Equal(v time.Time) *UserTokenQuery {
	return q.w("create_time='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) CreateTime_NotEqual(v time.Time) *UserTokenQuery {
	return q.w("create_time<>'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) CreateTime_Less(v time.Time) *UserTokenQuery {
	return q.w("create_time<'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) CreateTime_LessEqual(v time.Time) *UserTokenQuery {
	return q.w("create_time<='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) CreateTime_Greater(v time.Time) *UserTokenQuery {
	return q.w("create_time>'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) CreateTime_GreaterEqual(v time.Time) *UserTokenQuery {
	return q.w("create_time>='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UpdateTime_Equal(v time.Time) *UserTokenQuery {
	return q.w("update_time='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UpdateTime_NotEqual(v time.Time) *UserTokenQuery {
	return q.w("update_time<>'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UpdateTime_Less(v time.Time) *UserTokenQuery {
	return q.w("update_time<'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UpdateTime_LessEqual(v time.Time) *UserTokenQuery {
	return q.w("update_time<='" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UpdateTime_Greater(v time.Time) *UserTokenQuery {
	return q.w("update_time>'" + fmt.Sprint(v) + "'")
}
func (q *UserTokenQuery) UpdateTime_GreaterEqual(v time.Time) *UserTokenQuery {
	return q.w("update_time>='" + fmt.Sprint(v) + "'")
}

type UserTokenDao struct {
	logger     *zap.Logger
	db         *DB
	insertStmt *wrap.Stmt
	updateStmt *wrap.Stmt
	deleteStmt *wrap.Stmt
}

func NewUserTokenDao(db *DB) (t *UserTokenDao, err error) {
	t = &UserTokenDao{}
	t.logger = log.TypedLogger(t)
	t.db = db
	err = t.init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (dao *UserTokenDao) init() (err error) {
	err = dao.prepareInsertStmt()
	if err != nil {
		return err
	}

	err = dao.prepareUpdateStmt()
	if err != nil {
		return err
	}

	err = dao.prepareDeleteStmt()
	if err != nil {
		return err
	}

	return nil
}
func (dao *UserTokenDao) prepareInsertStmt() (err error) {
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO user_token (account_id,expires_time,user_token) VALUES (?,?,?)")
	return err
}

func (dao *UserTokenDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE user_token SET account_id=?,expires_time=?,user_token=? WHERE id=?")
	return err
}

func (dao *UserTokenDao) prepareDeleteStmt() (err error) {
	dao.deleteStmt, err = dao.db.Prepare(context.Background(), "DELETE FROM user_token WHERE id=?")
	return err
}

func (dao *UserTokenDao) Insert(ctx context.Context, tx *wrap.Tx, e *UserToken) (id int64, err error) {
	stmt := dao.insertStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	result, err := stmt.Exec(ctx, e.AccountId, e.ExpiresTime, e.UserToken)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *UserTokenDao) Update(ctx context.Context, tx *wrap.Tx, e *UserToken) (err error) {
	stmt := dao.updateStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, e.AccountId, e.ExpiresTime, e.UserToken, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *UserTokenDao) Delete(ctx context.Context, tx *wrap.Tx, id uint64) (err error) {
	stmt := dao.deleteStmt
	if tx != nil {
		stmt = tx.Stmt(ctx, stmt)
	}

	_, err = stmt.Exec(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (dao *UserTokenDao) scanRow(row *wrap.Row) (*UserToken, error) {
	e := &UserToken{}
	err := row.Scan(&e.Id, &e.AccountId, &e.ExpiresTime, &e.UserToken, &e.CreateTime, &e.UpdateTime)
	if err != nil {
		if err == wrap.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return e, nil
}

func (dao *UserTokenDao) scanRows(rows *wrap.Rows) (list []*UserToken, err error) {
	list = make([]*UserToken, 0)
	for rows.Next() {
		e := UserToken{}
		err = rows.Scan(&e.Id, &e.AccountId, &e.ExpiresTime, &e.UserToken, &e.CreateTime, &e.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return nil, err
	}

	return list, nil
}

func (dao *UserTokenDao) QueryOne(ctx context.Context, tx *wrap.Tx, query string) (*UserToken, error) {
	querySql := "SELECT " + USER_TOKEN_ALL_FIELDS_STRING + " FROM user_token " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	return dao.scanRow(row)
}

func (dao *UserTokenDao) QueryList(ctx context.Context, tx *wrap.Tx, query string) (list []*UserToken, err error) {
	querySql := "SELECT " + USER_TOKEN_ALL_FIELDS_STRING + " FROM user_token " + query
	var rows *wrap.Rows
	if tx == nil {
		rows, err = dao.db.Query(ctx, querySql)
	} else {
		rows, err = tx.Query(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return nil, err
	}

	return dao.scanRows(rows)
}

func (dao *UserTokenDao) QueryCount(ctx context.Context, tx *wrap.Tx, query string) (count int64, err error) {
	querySql := "SELECT COUNT(1) FROM user_token " + query
	var row *wrap.Row
	if tx == nil {
		row = dao.db.QueryRow(ctx, querySql)
	} else {
		row = tx.QueryRow(ctx, querySql)
	}
	if err != nil {
		dao.logger.Error("sqlDriver", zap.Error(err))
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *UserTokenDao) QueryGroupBy(ctx context.Context, tx *wrap.Tx, groupByFields []string, query string) (rows *wrap.Rows, err error) {
	querySql := "SELECT " + strings.Join(groupByFields, ",") + ",count(1) FROM user_token " + query
	if tx == nil {
		return dao.db.Query(ctx, querySql)
	} else {
		return tx.Query(ctx, querySql)
	}
}

func (dao *UserTokenDao) GetQuery() *UserTokenQuery {
	return NewUserTokenQuery(dao)
}

type DB struct {
	wrap.DB
	OauthState   *OauthStateDao
	OauthTokens  *OauthTokensDao
	RefreshToken *RefreshTokenDao
	UserToken    *UserTokenDao
}

func NewDB(connectionString string) (d *DB, err error) {
	if connectionString == "" {
		return nil, fmt.Errorf("connectionString nil")
	}

	d = &DB{}

	db, err := wrap.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	d.DB = *db

	err = d.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	d.OauthState, err = NewOauthStateDao(d)
	if err != nil {
		return nil, err
	}

	d.OauthTokens, err = NewOauthTokensDao(d)
	if err != nil {
		return nil, err
	}

	d.RefreshToken, err = NewRefreshTokenDao(d)
	if err != nil {
		return nil, err
	}

	d.UserToken, err = NewUserTokenDao(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
