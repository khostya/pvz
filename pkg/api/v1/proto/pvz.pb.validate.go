// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pvz.proto

package pvz_v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on PVZ with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *PVZ) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PVZ with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in PVZMultiError, or nil if none found.
func (m *PVZ) ValidateAll() error {
	return m.validate(true)
}

func (m *PVZ) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if all {
		switch v := interface{}(m.GetRegistrationDate()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PVZValidationError{
					field:  "RegistrationDate",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PVZValidationError{
					field:  "RegistrationDate",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRegistrationDate()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PVZValidationError{
				field:  "RegistrationDate",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for City

	if len(errors) > 0 {
		return PVZMultiError(errors)
	}

	return nil
}

// PVZMultiError is an error wrapping multiple validation errors returned by
// PVZ.ValidateAll() if the designated constraints aren't met.
type PVZMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PVZMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PVZMultiError) AllErrors() []error { return m }

// PVZValidationError is the validation error returned by PVZ.Validate if the
// designated constraints aren't met.
type PVZValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PVZValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PVZValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PVZValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PVZValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PVZValidationError) ErrorName() string { return "PVZValidationError" }

// Error satisfies the builtin error interface
func (e PVZValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPVZ.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PVZValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PVZValidationError{}

// Validate checks the field values on GetPVZListRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetPVZListRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetPVZListRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetPVZListRequestMultiError, or nil if none found.
func (m *GetPVZListRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetPVZListRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return GetPVZListRequestMultiError(errors)
	}

	return nil
}

// GetPVZListRequestMultiError is an error wrapping multiple validation errors
// returned by GetPVZListRequest.ValidateAll() if the designated constraints
// aren't met.
type GetPVZListRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetPVZListRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetPVZListRequestMultiError) AllErrors() []error { return m }

// GetPVZListRequestValidationError is the validation error returned by
// GetPVZListRequest.Validate if the designated constraints aren't met.
type GetPVZListRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetPVZListRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetPVZListRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetPVZListRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetPVZListRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetPVZListRequestValidationError) ErrorName() string {
	return "GetPVZListRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetPVZListRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetPVZListRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetPVZListRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetPVZListRequestValidationError{}

// Validate checks the field values on GetPVZListResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetPVZListResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetPVZListResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetPVZListResponseMultiError, or nil if none found.
func (m *GetPVZListResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetPVZListResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetPvzs() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetPVZListResponseValidationError{
						field:  fmt.Sprintf("Pvzs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetPVZListResponseValidationError{
						field:  fmt.Sprintf("Pvzs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetPVZListResponseValidationError{
					field:  fmt.Sprintf("Pvzs[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetPVZListResponseMultiError(errors)
	}

	return nil
}

// GetPVZListResponseMultiError is an error wrapping multiple validation errors
// returned by GetPVZListResponse.ValidateAll() if the designated constraints
// aren't met.
type GetPVZListResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetPVZListResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetPVZListResponseMultiError) AllErrors() []error { return m }

// GetPVZListResponseValidationError is the validation error returned by
// GetPVZListResponse.Validate if the designated constraints aren't met.
type GetPVZListResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetPVZListResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetPVZListResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetPVZListResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetPVZListResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetPVZListResponseValidationError) ErrorName() string {
	return "GetPVZListResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetPVZListResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetPVZListResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetPVZListResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetPVZListResponseValidationError{}
