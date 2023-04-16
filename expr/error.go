package expr

import "golox/tok"

type RuntimeError struct {
	Token   *tok.Token
	Message string
}
