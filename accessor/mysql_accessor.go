package accessor

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"news-feed/dto"
	"sync"
)

type MySQLAccessor struct {
	DB *sql.DB
}

func NewMySQLAccessor(user, password, host, dbname string) (*MySQLAccessor, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &MySQLAccessor{DB: db}, nil
}

// Squirrel을 이용한 SELECT 예제 함수
func (ma *MySQLAccessor) ExampleSelectUserByID(id int) (string, error) {
	query, args, err := sq.Select("name").From("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return "", err
	}
	var name string
	err = ma.DB.QueryRow(query, args...).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

// user 테이블에 데이터 추가 함수
func (ma *MySQLAccessor) InsertUser(username, password string) (int64, error) {
	query, args, err := sq.Insert("user").Columns("username", "password").Values(username, password).ToSql()
	if err != nil {
		return 0, err
	}
	result, err := ma.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// post 테이블에 데이터 추가 함수
func (ma *MySQLAccessor) InsertPost(content string, userID int64) (int64, error) {
	query, args, err := sq.Insert("post").Columns("content", "user_id", "created_at").Values(content, userID, sq.Expr("NOW()")).ToSql()
	if err != nil {
		return 0, err
	}
	result, err := ma.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// username, password로 user id를 조회하는 함수
func (ma *MySQLAccessor) GetUserIDByUsernamePassword(username, password string) (int64, error) {
	query, args, err := sq.Select("id").From("user").Where(sq.Eq{"username": username, "password": password}).ToSql()
	if err != nil {
		return 0, err
	}
	var id int64
	err = ma.DB.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// 여러 id 목록을 받아 post 테이블에서 해당 id들의 데이터를 조회하는 함수
func (ma *MySQLAccessor) GetPostsByIDs(ids []int64) ([]dto.Post, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	placeholders := sq.Placeholders(len(ids))
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	query := fmt.Sprintf("SELECT id, content, user_id, created_at FROM post WHERE id IN (%s) ORDER BY created_at DESC", placeholders)
	rows, err := ma.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []dto.Post
	for rows.Next() {
		var p dto.Post
		err := rows.Scan(&p.ID, &p.Content, &p.UserID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// follow 테이블에 데이터 추가 함수
func (ma *MySQLAccessor) InsertFollow(followerID, followeeID int64) (int64, error) {
	query, args, err := sq.Insert("follow").Columns("follower_id", "followee_id").Values(followerID, followeeID).ToSql()
	if err != nil {
		return 0, err
	}
	result, err := ma.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// username이 이미 존재하는지 체크하는 함수
func (ma *MySQLAccessor) IsUserExists(username string) (bool, error) {
	query, args, err := sq.Select("COUNT(*)").From("user").Where(sq.Eq{"username": username}).ToSql()
	if err != nil {
		return false, err
	}
	var count int
	err = ma.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// follower와 followee가 이미 팔로우 중인지 체크하는 함수
func (ma *MySQLAccessor) IsFollowExists(followerID, followeeID int64) (bool, error) {
	query, args, err := sq.Select("COUNT(*)").From("follow").Where(sq.Eq{"follower_id": followerID, "followee_id": followeeID}).ToSql()
	if err != nil {
		return false, err
	}
	var count int
	err = ma.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

var (
	onceMySQL     sync.Once
	mysqlInstance *MySQLAccessor
)

func GetMySQLAccessor(user, password, host, dbname string) (*MySQLAccessor, error) {
	var err error
	onceMySQL.Do(func() {
		mysqlInstance, err = NewMySQLAccessor(user, password, host, dbname)
	})
	return mysqlInstance, err
}
