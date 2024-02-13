package transport

import (
	"fmt"
)

type appConstant fmt.Stringer

//go:generate go-enum -ptr -values

// ENUM(Authorization)
type ContextKeyRequest int

// ENUM(token)
type ContextKeyCookie int

// ENUM(UID)
type contextAuthenticated int
