package models

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type ConnectionWebSocket struct {
	Type     string    `json:"type"`
	UserId   uuid.UUID `json:"userId"`
	Username string    `json:"username"`
}

func (c *ConnectionWebSocket) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string `json:"type"`
		UserId   string `json:"userId"`
		Username string `json:"username"`
	}{
		Type:     c.Type,
		UserId:   c.UserId.String(),
		Username: c.Username,
	})
}

func (c *ConnectionWebSocket) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type     string `json:"type"`
		UserId   string `json:"userId"`
		Username string `json:"username"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.Type = aux.Type
	c.UserId, _ = uuid.FromString(aux.UserId)
	c.Username = aux.Username
	return nil
}

type DisconnectionWebSocket struct {
	Type   string    `json:"type"`
	UserId uuid.UUID `json:"userId"`
}

func (d *DisconnectionWebSocket) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string `json:"type"`
		UserId string `json:"userId"`
	}{
		Type:   d.Type,
		UserId: d.UserId.String(),
	})
}

func (d *DisconnectionWebSocket) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type   string `json:"type"`
		UserId string `json:"userId"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	d.Type = aux.Type
	d.UserId, _ = uuid.FromString(aux.UserId)
	return nil
}

type NewMessageWebSocket struct {
	Type       string    `json:"type"`
	SenderId   string    `json:"senderId"`
	ReceiverId string    `json:"receiverId"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sentAt"`
}

func (n *NewMessageWebSocket) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string    `json:"type"`
		SenderId   string    `json:"senderId"`
		ReceiverId string    `json:"receiverId"`
		Content    string    `json:"content"`
		SentAt     time.Time `json:"sentAt"`
	}{
		Type:       n.Type,
		SenderId:   n.SenderId,
		ReceiverId: n.ReceiverId,
		Content:    n.Content,
		SentAt:     n.SentAt,
	})
}

func (n *NewMessageWebSocket) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type       string    `json:"type"`
		SenderId   string    `json:"senderId"`
		ReceiverId string    `json:"receiverId"`
		Content    string    `json:"content"`
		SentAt     time.Time `json:"sentAt"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	n.Type = aux.Type
	n.SenderId = aux.SenderId
	n.ReceiverId = aux.ReceiverId
	n.Content = aux.Content
	n.SentAt = aux.SentAt
	return nil
}

type UserTypingWebSocket struct {
	Type       string `json:"type"`
	SenderId   string `json:"senderId"`
	ReceiverId string `json:"receiverId"`
	Typing     bool   `json:"typing"`
}

func (u *UserTypingWebSocket) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string `json:"type"`
		SenderId   string `json:"senderId"`
		ReceiverId string `json:"receiverId"`
		Typing     bool   `json:"typing"`
	}{
		Type:       u.Type,
		SenderId:   u.SenderId,
		ReceiverId: u.ReceiverId,
		Typing:     u.Typing,
	})
}

func (u *UserTypingWebSocket) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type       string `json:"type"`
		SenderId   string `json:"senderId"`
		ReceiverId string `json:"receiverId"`
		Typing     bool   `json:"typing"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	u.Type = aux.Type
	u.SenderId = aux.SenderId
	u.ReceiverId = aux.ReceiverId
	u.Typing = aux.Typing
	return nil
}

type OnlineStatusWebSocket struct {
	Type   string    `json:"type"`
	UserId uuid.UUID `json:"userId"`
	Status bool      `json:"status"`
}

func (o *OnlineStatusWebSocket) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string `json:"type"`
		UserId string `json:"userId"`
		Status bool   `json:"status"`
	}{
		Type:   o.Type,
		UserId: o.UserId.String(),
		Status: o.Status,
	})
}

func (o *OnlineStatusWebSocket) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type   string `json:"type"`
		UserId string `json:"userId"`
		Status bool   `json:"status"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	o.Type = aux.Type
	o.UserId, _ = uuid.FromString(aux.UserId)
	o.Status = aux.Status
	return nil
}

type OldMessagesWebSocket struct {
	Type       string `json:"type"`
	SenderId   string `json:"senderId"`
	ReceiverId string `json:"receiverId"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
}

func (o *OldMessagesWebSocket) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string `json:"type"`
		SenderId   string `json:"senderId"`
		ReceiverId string `json:"receiverId"`
		Offset     int    `json:"offset"`
		Limit      int    `json:"limit"`
	}{
		Type:       o.Type,
		SenderId:   o.SenderId,
		ReceiverId: o.ReceiverId,
		Offset:     o.Offset,
		Limit:      o.Limit,
	})
}

func (o *OldMessagesWebSocket) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type       string `json:"type"`
		SenderId   string `json:"senderId"`
		ReceiverId string `json:"receiverId"`
		Offset     int    `json:"offset"`
		Limit      int    `json:"limit"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	o.Type = aux.Type
	o.SenderId = aux.SenderId
	o.ReceiverId = aux.ReceiverId
	o.Offset = aux.Offset
	o.Limit = aux.Limit
	return nil
}

type ErrorWebSocket struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e *ErrorWebSocket) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}{
		Type:    e.Type,
		Message: e.Message,
	})
}

func (e *ErrorWebSocket) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.Type = aux.Type
	e.Message = aux.Message
	return nil
}
