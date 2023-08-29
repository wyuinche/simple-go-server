package model

type Order struct {
	OID  int64 `json:"oid"`
	UID  int64 `json:"uid"`
	Date int64 `json:"date"`
}

type OrderProduct struct {
	OID int64 `json:"oid"`
	PID int64 `json:"pid"`
}
