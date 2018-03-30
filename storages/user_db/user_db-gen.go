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
	"os"
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

const REFRESH_TOKEN_TABLE_NAME = "refresh_token"

type REFRESH_TOKEN_FIELD string

const REFRESH_TOKEN_FIELD_ID = REFRESH_TOKEN_FIELD("id")
const REFRESH_TOKEN_FIELD_OAUTH_STATE = REFRESH_TOKEN_FIELD("oauth_state")
const REFRESH_TOKEN_FIELD_QUERY_STRING = REFRESH_TOKEN_FIELD("query_string")
const REFRESH_TOKEN_FIELD_USER_AGENT = REFRESH_TOKEN_FIELD("user_agent")
const REFRESH_TOKEN_FIELD_ACCOUNT_ID = REFRESH_TOKEN_FIELD("account_id")
const REFRESH_TOKEN_FIELD_REFRESH_TOKEN = REFRESH_TOKEN_FIELD("refresh_token")
const REFRESH_TOKEN_FIELD_IS_LOGOUT = REFRESH_TOKEN_FIELD("is_logout")
const REFRESH_TOKEN_FIELD_GMT_LOGOUT = REFRESH_TOKEN_FIELD("gmt_logout")
const REFRESH_TOKEN_FIELD_OAUTH_AUTHORIZATION_CODE = REFRESH_TOKEN_FIELD("oauth_authorization_code")
const REFRESH_TOKEN_FIELD_OAUTH_REFRESH_TOKEN = REFRESH_TOKEN_FIELD("oauth_refresh_token")
const REFRESH_TOKEN_FIELD_CREATE_TIME = REFRESH_TOKEN_FIELD("create_time")
const REFRESH_TOKEN_FIELD_UPDATE_TIME = REFRESH_TOKEN_FIELD("update_time")

const REFRESH_TOKEN_ALL_FIELDS_STRING = "id,oauth_state,query_string,user_agent,account_id,refresh_token,is_logout,gmt_logout,oauth_authorization_code,oauth_refresh_token,create_time,update_time"

var REFRESH_TOKEN_ALL_FIELDS = []string{
	"id",
	"oauth_state",
	"query_string",
	"user_agent",
	"account_id",
	"refresh_token",
	"is_logout",
	"gmt_logout",
	"oauth_authorization_code",
	"oauth_refresh_token",
	"create_time",
	"update_time",
}

type RefreshToken struct {
	Id                     uint64 //size=20
	OauthState             string //size=128
	QueryString            string //size=256
	UserAgent              string //size=256
	AccountId              string //size=128
	RefreshToken           string //size=128
	IsLogout               int32  //size=1
	GmtLogout              time.Time
	OauthAuthorizationCode string //size=128
	OauthRefreshToken      string //size=128
	CreateTime             time.Time
	UpdateTime             time.Time
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
func (q *RefreshTokenQuery) OauthState_Equal(v string) *RefreshTokenQuery {
	return q.w("oauth_state='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthState_NotEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_state<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthState_Less(v string) *RefreshTokenQuery {
	return q.w("oauth_state<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthState_LessEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_state<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthState_Greater(v string) *RefreshTokenQuery {
	return q.w("oauth_state>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthState_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_state>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) QueryString_Equal(v string) *RefreshTokenQuery {
	return q.w("query_string='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) QueryString_NotEqual(v string) *RefreshTokenQuery {
	return q.w("query_string<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) QueryString_Less(v string) *RefreshTokenQuery {
	return q.w("query_string<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) QueryString_LessEqual(v string) *RefreshTokenQuery {
	return q.w("query_string<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) QueryString_Greater(v string) *RefreshTokenQuery {
	return q.w("query_string>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) QueryString_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("query_string>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserAgent_Equal(v string) *RefreshTokenQuery {
	return q.w("user_agent='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserAgent_NotEqual(v string) *RefreshTokenQuery {
	return q.w("user_agent<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserAgent_Less(v string) *RefreshTokenQuery {
	return q.w("user_agent<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserAgent_LessEqual(v string) *RefreshTokenQuery {
	return q.w("user_agent<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserAgent_Greater(v string) *RefreshTokenQuery {
	return q.w("user_agent>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) UserAgent_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("user_agent>='" + fmt.Sprint(v) + "'")
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
func (q *RefreshTokenQuery) IsLogout_Equal(v int32) *RefreshTokenQuery {
	return q.w("is_logout='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) IsLogout_NotEqual(v int32) *RefreshTokenQuery {
	return q.w("is_logout<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) IsLogout_Less(v int32) *RefreshTokenQuery {
	return q.w("is_logout<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) IsLogout_LessEqual(v int32) *RefreshTokenQuery {
	return q.w("is_logout<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) IsLogout_Greater(v int32) *RefreshTokenQuery {
	return q.w("is_logout>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) IsLogout_GreaterEqual(v int32) *RefreshTokenQuery {
	return q.w("is_logout>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) GmtLogout_Equal(v time.Time) *RefreshTokenQuery {
	return q.w("gmt_logout='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) GmtLogout_NotEqual(v time.Time) *RefreshTokenQuery {
	return q.w("gmt_logout<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) GmtLogout_Less(v time.Time) *RefreshTokenQuery {
	return q.w("gmt_logout<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) GmtLogout_LessEqual(v time.Time) *RefreshTokenQuery {
	return q.w("gmt_logout<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) GmtLogout_Greater(v time.Time) *RefreshTokenQuery {
	return q.w("gmt_logout>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) GmtLogout_GreaterEqual(v time.Time) *RefreshTokenQuery {
	return q.w("gmt_logout>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthAuthorizationCode_Equal(v string) *RefreshTokenQuery {
	return q.w("oauth_authorization_code='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthAuthorizationCode_NotEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_authorization_code<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthAuthorizationCode_Less(v string) *RefreshTokenQuery {
	return q.w("oauth_authorization_code<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthAuthorizationCode_LessEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_authorization_code<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthAuthorizationCode_Greater(v string) *RefreshTokenQuery {
	return q.w("oauth_authorization_code>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthAuthorizationCode_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_authorization_code>='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthRefreshToken_Equal(v string) *RefreshTokenQuery {
	return q.w("oauth_refresh_token='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthRefreshToken_NotEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_refresh_token<>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthRefreshToken_Less(v string) *RefreshTokenQuery {
	return q.w("oauth_refresh_token<'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthRefreshToken_LessEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_refresh_token<='" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthRefreshToken_Greater(v string) *RefreshTokenQuery {
	return q.w("oauth_refresh_token>'" + fmt.Sprint(v) + "'")
}
func (q *RefreshTokenQuery) OauthRefreshToken_GreaterEqual(v string) *RefreshTokenQuery {
	return q.w("oauth_refresh_token>='" + fmt.Sprint(v) + "'")
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
	dao.insertStmt, err = dao.db.Prepare(context.Background(), "INSERT INTO refresh_token (oauth_state,query_string,user_agent,account_id,refresh_token,is_logout,gmt_logout,oauth_authorization_code,oauth_refresh_token) VALUES (?,?,?,?,?,?,?,?,?)")
	return err
}

func (dao *RefreshTokenDao) prepareUpdateStmt() (err error) {
	dao.updateStmt, err = dao.db.Prepare(context.Background(), "UPDATE refresh_token SET oauth_state=?,query_string=?,user_agent=?,account_id=?,refresh_token=?,is_logout=?,gmt_logout=?,oauth_authorization_code=?,oauth_refresh_token=? WHERE id=?")
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

	result, err := stmt.Exec(ctx, e.OauthState, e.QueryString, e.UserAgent, e.AccountId, e.RefreshToken, e.IsLogout, e.GmtLogout, e.OauthAuthorizationCode, e.OauthRefreshToken)
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

	_, err = stmt.Exec(ctx, e.OauthState, e.QueryString, e.UserAgent, e.AccountId, e.RefreshToken, e.IsLogout, e.GmtLogout, e.OauthAuthorizationCode, e.OauthRefreshToken, e.Id)
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
	err := row.Scan(&e.Id, &e.OauthState, &e.QueryString, &e.UserAgent, &e.AccountId, &e.RefreshToken, &e.IsLogout, &e.GmtLogout, &e.OauthAuthorizationCode, &e.OauthRefreshToken, &e.CreateTime, &e.UpdateTime)
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
		err = rows.Scan(&e.Id, &e.OauthState, &e.QueryString, &e.UserAgent, &e.AccountId, &e.RefreshToken, &e.IsLogout, &e.GmtLogout, &e.OauthAuthorizationCode, &e.OauthRefreshToken, &e.CreateTime, &e.UpdateTime)
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
	RefreshToken *RefreshTokenDao
	UserToken    *UserTokenDao
}

func NewDB() (d *DB, err error) {
	d = &DB{}

	connectionString := os.Getenv("DB")
	if connectionString == "" {
		return nil, fmt.Errorf("DB env nil")
	}
	connectionString += "/user?parseTime=true"
	db, err := wrap.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	d.DB = *db

	err = d.Ping(context.Background())
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
