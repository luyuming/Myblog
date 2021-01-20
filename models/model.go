package models

import (
	"myblog/dao"
	"time"
)

type User struct {
	ID       uint
	Name     string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Tel      string `form:"phone" json:"phone"`
	Email    string `form:"email" json:"email"`
	Records  []Record
}

type Record struct {
	ID         uint
	Title      string `form:"title" json:"title"`
	Content    string `form:"content" json:"content"`
	UserID     uint
	Created_at time.Time
	Comments   []Comment
}

type Comment struct {
	ID         uint
	Content    string
	Name       string
	RecordID   uint
	Created_at time.Time
}

type RecordUser struct {
	Record   Record
	UserName string
	Star     int
}

func GetUser(name string) (user *User, err error) {
	user = new(User)
	err = dao.DB.Where("name = ?", name).Find(user).Error
	if err != nil {
		return nil, err
	}
	return
}

func AddUser(user *User) (err error) {
	err = dao.DB.Create(user).Error
	return
}

func GetRecord(id string) (record *Record, err error) {
	record = new(Record)
	err = dao.DB.Where("ID = ?", id).Find(record).Error
	if err != nil {
		return nil, err
	}
	return
}

func AddRecord(record *Record) (err error) {
	err = dao.DB.Create(record).Error
	return
}

func UserAssocRecord(user *User, record *Record) (err error) {
	err = dao.DB.Model(user).Association("Records").Append(record).Error
	return
}

func RecordAssocComment(record *Record, comment *Comment) (err error) {
	err = dao.DB.Model(record).Association("Comments").Append(comment).Error
	return
}

func RecordRelatedUser(record *Record, user *User) (err error) {
	err = dao.DB.Model(record).Related(user).Error
	return
}

func GetUserRecord(id string) (user *User, err error) {
	user = new(User)
	err = dao.DB.Where("ID = ?", id).Preload("Records").Find(user).Error
	return
}

func GetRecordComment(id string) (record *Record, err error) {
	record = new(Record)
	err = dao.DB.Where("ID = ?", id).Preload("Comments").Find(record).Error
	return
}
func IncrKey(key string) (val int64, err error) {
	val, err = dao.RDB.Incr(key).Result()
	return
}

func DecrKey(key string) (val int64, err error) {
	val, err = dao.RDB.Decr(key).Result()
	return
}

func CheckKeyExist(key string) (val int64, err error) {
	val, err = dao.RDB.Exists(key).Result()
	return
}

func GetKV(key string) (val string, err error) {
	val, err = dao.RDB.Get(key).Result()
	return
}

func SetKV(key string, val interface{}, ex time.Duration) (err error) {
	err = dao.RDB.Set(key, val, ex).Err()
	return
}

func GetCount(tablename string, count *int) (err error) {
	err = dao.DB.Table(tablename).Count(count).Error
	return
}

func GetLimitRecordsOffset(limit int, offset int, showrecords *[]Record) (err error) {
	err = dao.DB.Offset(offset).Limit(limit).Find(showrecords).Error
	return
}
