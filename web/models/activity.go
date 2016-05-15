package models

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"strings"
)

type Activity struct {
	Id                uint64         `json:"id"`
	Event             string         `json:"event"`
	TargetId          uint64         `json:"target_id"`
	SourceId          uint64         `json:"source_id"`
	TargetObjectId    uint64         `json:"target_object_id"`
	RetweetedStatusId JsonNullUInt64 `json:"retweeted_status_id"`
	CreatedOn         int            `json:"created_at"`
}

func (m *Model) LoadActivities(userIdStr string, maxIdStr string, sinceIdStr string) []*Activity {
	db, err := m.Open()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	toId := func(id string) string {
		if id == "" {
			return ""
		} else if strings.Contains(id, ":") {
			fields := strings.Split(id, ":")
			stmtOut, err := db.Prepare("SELECT id FROM activity WHERE target_object_id = ? AND event = ? AND source_id = ? LIMIT 1")
			if err != nil {
				fmt.Println(err.Error())
				return ""
			}
			defer stmtOut.Close()
			var dbId string
			fmt.Printf("target_object_id:%s event:%s source_id:%s\n", fields[0], fields[1], fields[2])
			err = stmtOut.QueryRow(fields[0], fields[1], fields[2]).Scan(&dbId)
			if err != nil {
				fmt.Println(err.Error())
				return ""
			}
			return dbId
		} else {
			return id
		}
	}

	maxId := toId(maxIdStr)
	sinceId := toId(sinceIdStr)

	stmt := sq.
		Select("id, event, target_id, source_id, target_object_id, retweeted_status_id, created_at").
		From("activity").
		Where(sq.Eq{"target_id": userIdStr}).
		OrderBy("id DESC").
		Limit(100)

	if maxId != "" {
		stmt = stmt.Where(sq.LtOrEq{"id": maxId})
	}

	if sinceId != "" {
		stmt = stmt.Where(sq.GtOrEq{"id": sinceId})
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("[LoadActivities] maxId:%s sinceId:%s sql:%s\n", maxId, sinceId, sql)

	rows, err := db.Query(sql, args...)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	activities := []*Activity{}

	for rows.Next() {
		var id uint64
		var event string
		var targetId uint64
		var sourceId uint64
		var targetObjectId uint64
		var retweetedStatusId JsonNullUInt64
		var createdOn int
		err = rows.Scan(&id, &event, &targetId, &sourceId, &targetObjectId, &retweetedStatusId, &createdOn)
		if err != nil {
			panic(err.Error())
		}
		activity := &Activity{
			Id:                id,
			Event:             event,
			TargetId:          targetId,
			SourceId:          sourceId,
			TargetObjectId:    targetObjectId,
			RetweetedStatusId: retweetedStatusId,
			CreatedOn:         createdOn,
		}
		activities = append(activities, activity)
	}

	return activities
}

func (m *Model) LoadFavoriterIds(userIdStr string, idStr string) []string {
	db, err := m.Open()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmt := sq.
		Select("source_id").
		From("activity").
		Where(sq.Eq{"target_id": userIdStr}).
		Where(sq.Eq{"target_object_id": idStr}).
		Where(sq.Eq{"event": []string{"favorite", "favorited_retweet"}}).
		OrderBy("id DESC").
		Limit(100)

	sql, args, err := stmt.ToSql()
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("[LoadFavoriterIds] id:%s sql:%s\n", idStr, sql)

	rows, err := db.Query(sql, args...)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	ids := []string{}

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		ids = append(ids, id)
	}

	return ids
}
