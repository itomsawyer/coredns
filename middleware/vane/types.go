package vane

import (
	"errors"
	"time"
)

var (
	errUnreachable     = errors.New("unreachable backend")
	errUnexpectedLogic = errors.New("incorrect or incomplete data logic")
	errFormatError     = errors.New("format error")

	defaultMaxAttampts = 10
	defaultTimeout     = 3 * time.Second
)

const (
	warnNoGoodAnswer = iota
	warnNoneReplies
	warnPolicyMaxAttempts
)
