// Code generated by go-enum DO NOT EDIT.
// Version: 0.6.0
// Revision: 919e61c0174b91303753ee3898569a01abb32c97
// Build Date: 2023-12-18T15:54:43Z
// Built By: goreleaser

package transport

import (
	"errors"
	"fmt"
)

const (
	// ContextKeyCookieToken is a ContextKeyCookie of type Token.
	ContextKeyCookieToken ContextKeyCookie = iota
)

var ErrInvalidContextKeyCookie = errors.New("not a valid ContextKeyCookie")

const _ContextKeyCookieName = "token"

// ContextKeyCookieValues returns a list of the values for ContextKeyCookie
func ContextKeyCookieValues() []ContextKeyCookie {
	return []ContextKeyCookie{
		ContextKeyCookieToken,
	}
}

var _ContextKeyCookieMap = map[ContextKeyCookie]string{
	ContextKeyCookieToken: _ContextKeyCookieName[0:5],
}

// String implements the Stringer interface.
func (x ContextKeyCookie) String() string {
	if str, ok := _ContextKeyCookieMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ContextKeyCookie(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ContextKeyCookie) IsValid() bool {
	_, ok := _ContextKeyCookieMap[x]
	return ok
}

var _ContextKeyCookieValue = map[string]ContextKeyCookie{
	_ContextKeyCookieName[0:5]: ContextKeyCookieToken,
}

// ParseContextKeyCookie attempts to convert a string to a ContextKeyCookie.
func ParseContextKeyCookie(name string) (ContextKeyCookie, error) {
	if x, ok := _ContextKeyCookieValue[name]; ok {
		return x, nil
	}
	return ContextKeyCookie(0), fmt.Errorf("%s is %w", name, ErrInvalidContextKeyCookie)
}

func (x ContextKeyCookie) Ptr() *ContextKeyCookie {
	return &x
}

const (
	// ContextKeyRequestAuthorization is a ContextKeyRequest of type Authorization.
	ContextKeyRequestAuthorization ContextKeyRequest = iota
)

var ErrInvalidContextKeyRequest = errors.New("not a valid ContextKeyRequest")

const _ContextKeyRequestName = "Authorization"

// ContextKeyRequestValues returns a list of the values for ContextKeyRequest
func ContextKeyRequestValues() []ContextKeyRequest {
	return []ContextKeyRequest{
		ContextKeyRequestAuthorization,
	}
}

var _ContextKeyRequestMap = map[ContextKeyRequest]string{
	ContextKeyRequestAuthorization: _ContextKeyRequestName[0:13],
}

// String implements the Stringer interface.
func (x ContextKeyRequest) String() string {
	if str, ok := _ContextKeyRequestMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ContextKeyRequest(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ContextKeyRequest) IsValid() bool {
	_, ok := _ContextKeyRequestMap[x]
	return ok
}

var _ContextKeyRequestValue = map[string]ContextKeyRequest{
	_ContextKeyRequestName[0:13]: ContextKeyRequestAuthorization,
}

// ParseContextKeyRequest attempts to convert a string to a ContextKeyRequest.
func ParseContextKeyRequest(name string) (ContextKeyRequest, error) {
	if x, ok := _ContextKeyRequestValue[name]; ok {
		return x, nil
	}
	return ContextKeyRequest(0), fmt.Errorf("%s is %w", name, ErrInvalidContextKeyRequest)
}

func (x ContextKeyRequest) Ptr() *ContextKeyRequest {
	return &x
}

const (
	// ContextAuthenticatedUID is a contextAuthenticated of type UID.
	ContextAuthenticatedUID contextAuthenticated = iota
)

var ErrInvalidcontextAuthenticated = errors.New("not a valid contextAuthenticated")

const _contextAuthenticatedName = "UID"

// contextAuthenticatedValues returns a list of the values for contextAuthenticated
func contextAuthenticatedValues() []contextAuthenticated {
	return []contextAuthenticated{
		ContextAuthenticatedUID,
	}
}

var _contextAuthenticatedMap = map[contextAuthenticated]string{
	ContextAuthenticatedUID: _contextAuthenticatedName[0:3],
}

// String implements the Stringer interface.
func (x contextAuthenticated) String() string {
	if str, ok := _contextAuthenticatedMap[x]; ok {
		return str
	}
	return fmt.Sprintf("contextAuthenticated(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x contextAuthenticated) IsValid() bool {
	_, ok := _contextAuthenticatedMap[x]
	return ok
}

var _contextAuthenticatedValue = map[string]contextAuthenticated{
	_contextAuthenticatedName[0:3]: ContextAuthenticatedUID,
}

// ParsecontextAuthenticated attempts to convert a string to a contextAuthenticated.
func ParsecontextAuthenticated(name string) (contextAuthenticated, error) {
	if x, ok := _contextAuthenticatedValue[name]; ok {
		return x, nil
	}
	return contextAuthenticated(0), fmt.Errorf("%s is %w", name, ErrInvalidcontextAuthenticated)
}

func (x contextAuthenticated) Ptr() *contextAuthenticated {
	return &x
}
