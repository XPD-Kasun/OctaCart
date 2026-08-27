package shared

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
