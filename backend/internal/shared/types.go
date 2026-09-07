package shared

import "time"

// The type of a claim is string. i.e. A claim is described as a string
type Claim string

// In any bounded context, a person or client interacting with is called a user
// and has a id of int. i.e. UserId type is int. If need uuid, it should be here.
type UserId int

type Money int64

type KV[Key comparable, Val any] struct {
	Key Key
	Val Val
}

// Represent Id of any user that is capable of configuring admin. Basically, any staff users merchant creates or implicit admin for system tasks. These do not include any customer ids.
type AdminId int

type Pagination struct {
	Page    int
	PerPage int
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PerPage
}

type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
	ActorId() UserId
}
