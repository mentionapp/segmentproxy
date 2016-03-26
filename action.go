package segmentproxy

import (
	"encoding/json"

	"github.com/segmentio/analytics-go"
)

// Action represents common operations for all actions
type Action interface {
	// Unmarshal parses a paylaod into the Action instance
	Unmarshal(data []byte) error
	// Send calls the appropriate analytics.Client method for the Action
	Send(client *analytics.Client) error
	// GetEmail returns the email found in the payload, if any
	GetEmail() string
	// SetUserID allows to override the user ID associated with the action
	SetUserID(id string)
}

type Email struct {
	Email string `json:"email"`
}

type Identify struct {
	analytics.Identify
	Email
}

type Group struct {
	analytics.Group
	Email
}

type Track struct {
	analytics.Track
	Email
}

func (e Email) GetEmail() string {
	return e.Email
}

func (i *Identify) Unmarshal(data []byte) error {
	return json.Unmarshal(data, i)
}
func (i *Identify) Send(client *analytics.Client) error {
	return client.Identify(&i.Identify)
}
func (i *Identify) SetUserID(id string) {
	i.UserId = id
}

func (i *Group) Unmarshal(data []byte) error {
	return json.Unmarshal(data, i)
}
func (i *Group) Send(client *analytics.Client) error {
	return client.Group(&i.Group)
}
func (i *Group) SetUserID(id string) {
	i.UserId = id
}

func (i *Track) Unmarshal(data []byte) error {
	return json.Unmarshal(data, i)
}
func (i *Track) Send(client *analytics.Client) error {
	return client.Track(&i.Track)
}
func (i *Track) SetUserID(id string) {
	i.UserId = id
}
