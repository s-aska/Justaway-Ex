package routes

import (
	"database/sql"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func (r *Router) ApiDeviceTokenRegister(c echo.Context) error {
	m := r.NewModel()

	db, err := m.Open()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	userIdStr := auth(c, db)

	deviceName := c.QueryParam("deviceName")
	deviceType := c.QueryParam("deviceType")
	deviceToken := c.QueryParam("deviceToken")

	if deviceType != "APNS" && deviceType != "APNS_SANDBOX" && deviceType != "GCM" {
		return c.String(400, "invalid deviceType:"+deviceType)
	}

	if deviceToken == "" {
		return c.String(400, "missing deviceToken")
	}

	if deviceName == "" {
		return c.String(400, "missing deviceName")
	}

	_, err = db.Exec(`
		INSERT INTO notification_device(user_id, name, token, platform, created_at)
		VALUES(?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name)
	`, userIdStr, deviceName, deviceToken, deviceType, time.Now().Unix())

	return c.JSON(200, map[string]bool{"Success": true})
}

func (r *Router) ApiActivityList(c echo.Context) error {
	m := r.NewModel()

	db, err := m.Open()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	userIdStr := auth(c, db)

	if userIdStr == "" {
		return c.String(401, "invalid X-Justaway-API-Token header")
	}

	activities := m.LoadActivities(userIdStr, c.QueryParam("max_id"), c.QueryParam("since_id"))

	return c.JSON(200, activities)
}

func (r *Router) Revoke(c echo.Context) error {
	m := r.NewModel()

	db, err := m.Open()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	userIdStr := auth(c, db)

	if userIdStr == "" {
		return c.String(401, "invalid X-Justaway-API-Token header")
	}

	var crawlerUrl string
	err = db.QueryRow(`
		SELECT crawler.url
		FROM account
		LEFT OUTER JOIN crawler ON crawler.id = account.crawler_id
		WHERE user_id = ? LIMIT 1
	`, userIdStr).Scan(&crawlerUrl)
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("UPDATE account SET status = 'REVOKE' WHERE user_id = ?", userIdStr)
	if err != nil {
		panic(err.Error())
	}

	req, _ := http.NewRequest("GET", crawlerUrl+userIdStr+"/stop", nil)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	return c.JSON(200, map[string]bool{"Success": true})
}

func auth(c echo.Context, db *sql.DB) string {
	apiToken := c.Request().Header().Get("X-Justaway-API-Token")
	if apiToken == "" {
		return ""
	}

	var userIdStr string
	err := db.QueryRow(`
		SELECT user_id FROM api_token WHERE api_token = ? LIMIT 1
	`, apiToken).Scan(&userIdStr)
	if err != nil {
		return ""
	}

	_, err = db.Exec(`UPDATE api_token SET authenticated_at = ? WHERE api_token = ?`, time.Now().Unix(), apiToken)

	return userIdStr
}
