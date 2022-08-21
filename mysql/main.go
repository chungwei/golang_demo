package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func main() {
	InitMysql(loadMysqlConf())
	defer ReleaseMysqlConn()

	bizHandler()
}

func bizHandler() {
	ctx := context.Background()

	// insert
	user := &UserEntity{
		UserId:     0,
		Username:   "张三",
		Mobile:     fmt.Sprint(time.Now().Unix()),
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	uid, err := insertUser(ctx, []*UserEntity{user})
	fmt.Println("insertUser res: id=", uid, ",err=", err)

	// update
	where := map[string]interface{}{"id": uid, "pre_username": "张三"}
	update := map[string]interface{}{"username": "李四"}
	row, err := updateUser(ctx, update, where)
	fmt.Println("updateUser res: row=", row, ",err=", err)

	// update in trx
	userLog := &UserLogEntity{
		Id:         0,
		UserId:     uid,
		Username:   "李四",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	where = map[string]interface{}{"id": uid, "pre_username": "李四"}
	update = map[string]interface{}{"username": "王五"}
	row, err = updateUserInTrx(ctx, update, where, userLog)
	fmt.Println("updateUserInTrx res: row=", row, ",err=", err)

	// select one/list
	where = map[string]interface{}{"id": uid}
	users, err := queryUserList(ctx, where) // list
	user = users[0]                         // one
	fmt.Println("queryUserList:", user)

	// select in
	users, err = queryUserListInIds(ctx, []int64{1000, 2, 1003, 4}) // list
	fmt.Println("queryUserListInIds", users)

}

type UserEntity struct {
	UserId     int64  `db:"id" json:"userId"`              //ID,自增主键
	Username   string `db:"username" json:"username"`      //用户名
	Mobile     string `db:"mobile" json:"mobile"`          //手机号
	CreateTime string `db:"create_time" json:"createTime"` //注册时间
}

func insertUser(ctx context.Context, ens []*UserEntity) (int64, error) {
	conn := GetMysqlConn()
	if conn == nil {
		return int64(0), errors.New(`InsertUser->GetMysqlConn=nil`)
	}
	qs := "INSERT INTO b_user (id, username, mobile, create_time) " +
		" VALUES (:id, :username, :mobile, :create_time)"
	ret, err := conn.NamedExec(qs, ens)
	if err != nil {
		return 0, err
	}

	return ret.LastInsertId()
}

func queryUserListInIds(ctx context.Context, ids []int64) (users []*UserEntity, err error) {
	conn := GetMysqlConn()
	if conn == nil {
		return nil, errors.New(`queryUserListInIds->GetMysqlConn=nil`)
	}

	qs, args, err := sqlx.In("SELECT id, username, mobile, create_time FROM b_user WHERE id IN (?) LIMIT 10", ids)
	if err != nil {
		return nil, err
	}
	qs = conn.Rebind(qs)
	err = conn.Select(&users, qs, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
	}
	return users, err
}

func queryUserList(ctx context.Context, where map[string]interface{}) (users []*UserEntity, err error) {
	conn := GetMysqlConn()
	if conn == nil {
		return nil, errors.New(`queryUserList->GetMysqlConn=nil`)
	}

	whereStr := ""
	if len(where) == 0 {
		whereStr = "id > 0 AND "
	}
	for k, _ := range where {
		whereStr += fmt.Sprintf("%v = :%v AND ", k, k)
	}

	qs := "SELECT id, username, mobile, create_time FROM b_user" +
		" WHERE " + strings.TrimRight(whereStr, " AND ") +
		" ORDER BY id DESC" +
		" LIMIT 1000"
	nstmt, err := conn.PrepareNamed(qs)
	defer ReleaseStmt(nstmt)
	if err != nil {
		return
	}

	err = nstmt.Select(&users, where)
	if err != nil {
		if sql.ErrNoRows == err {
			err = nil
		}
	}
	return
}

func updateUser(ctx context.Context, update map[string]interface{}, where map[string]interface{}) (row int64, err error) {
	conn := GetMysqlConn()
	if conn == nil {
		err = errors.New(`updateUser->GetMysqlConn=nil`)
		return
	}

	argsMap := make(map[string]interface{})
	setStr := ""
	for k, v := range update {
		setStr += fmt.Sprintf("%v = :%v, ", strings.Replace(k, "pre_", "", 1), k)
		argsMap[k] = v
	}
	whereStr := ""
	for k, v := range where {
		whereStr += fmt.Sprintf("%v = :%v AND ", strings.Replace(k, "pre_", "", 1), k)
		argsMap[k] = v
	}

	qs := "UPDATE b_user " +
		" SET " + strings.TrimRight(setStr, ", ") +
		" WHERE " + strings.TrimRight(whereStr, " AND ")
	res, err := conn.NamedExec(qs, argsMap)
	if err != nil {
		return
	}

	return res.RowsAffected()
}

func updateUserInTrx(ctx context.Context, update map[string]interface{}, where map[string]interface{}, logEntity *UserLogEntity) (row int64, err error) {
	conn := GetMysqlConn()
	if conn == nil {
		err = errors.New(`updateUser->GetMysqlConn=nil`)
		return
	}

	argsMap := make(map[string]interface{})
	setStr := ""
	for k, v := range update {
		setStr += fmt.Sprintf("%v = :%v, ", strings.Replace(k, "pre_", "", 1), k)
		argsMap[k] = v
	}
	whereStr := ""
	for k, v := range where {
		whereStr += fmt.Sprintf("%v = :%v AND ", strings.Replace(k, "pre_", "", 1), k)
		argsMap[k] = v
	}
	tx := conn.MustBegin()
	qs := "UPDATE b_user " +
		" SET " + strings.TrimRight(setStr, ", ") +
		" WHERE " + strings.TrimRight(whereStr, " AND ")
	res, err := conn.NamedExec(qs, argsMap)
	if err != nil {
		return
	}
	if logEntity != nil {
		_, err = insertUserLog(ctx, tx, logEntity)
	}
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}

	return res.RowsAffected()
}

type UserLogEntity struct {
	Id         int64  `db:"id" json:"id"`                  //ID,自增主键
	UserId     int64  `db:"user_id" json:"userId"`         //用户ID
	Username   string `db:"username" json:"username"`      //用户名
	Mobile     string `db:"mobile" json:"mobile"`          //手机号
	CreateTime string `db:"create_time" json:"createTime"` //注册时间
}

func insertUserLog(ctx context.Context, tx *sqlx.Tx, en *UserLogEntity) (int64, error) {
	if tx == nil {
		return int64(0), errors.New(`insertUserLog->GetMysqlConn=nil`)
	}
	qs := "INSERT INTO b_user_log (id, user_id, username, mobile, create_time) " +
		" VALUES (:id, :user_id, :username, :mobile, :create_time)"
	ret, err := tx.NamedExec(qs, &en)
	if err != nil {
		return 0, err
	}

	return ret.LastInsertId()
}
