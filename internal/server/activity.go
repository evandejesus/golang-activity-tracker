package server

import (
	"fmt"
	"time"
)

type Activity struct {
	Time        time.Time `json:"time"`
	Description string    `json:"description"`
	ID          uint64    `json:"id"`
}

type Activities struct {
	activities []Activity
}

var ErrIDNotFound = fmt.Errorf("ID not found")

func (c *Activities) Insert(activity Activity) uint64 {
	activity.ID = uint64(len(c.activities))
	c.activities = append(c.activities, activity)
	return activity.ID
}

func (c *Activities) Retrieve(id uint64) (Activity, error) {
	if id >= uint64(len(c.activities)) {
		return Activity{}, ErrIDNotFound
	}
	return c.activities[id], nil
}
