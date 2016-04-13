package routes

import (
	"github.com/labstack/echo"
)

func (r *Router) ApiActivityList(c echo.Context) error {
	apiToken := c.Request().Header().Get("X-Justaway-API-Token")
	if apiToken == "" {
		return c.String(401, "Missing X-Justaway-API-Token header")
	}

	m := r.NewModel()

	db, err := m.Open()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var userIdStr string
	err = db.QueryRow(`
		SELECT user_id FROM api_token WHERE api_token = ? LIMIT 1
	`, apiToken).Scan(&userIdStr)
	if err != nil {
		return c.String(401, "Invalid X-Justaway-API-Token header")
	}

	activities := m.LoadActivities(userIdStr, c.QueryParam("max_id"), c.QueryParam("since_id"))

	return c.JSON(200, activities)
}
