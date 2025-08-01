// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: tzero/v1/payment_intent/recipient/recipient.proto

package recipient

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

	common "github.com/t-0-network/provider-sdk-go/api/gen/proto/tzero/v1/common"
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

	_ = common.PaymentMethodType(0)
)

// Validate checks the field values on CreatePaymentIntentRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreatePaymentIntentRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreatePaymentIntentRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreatePaymentIntentRequestMultiError, or nil if none found.
func (m *CreatePaymentIntentRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreatePaymentIntentRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PaymentReference

	// no validation rules for PayInCurrency

	if all {
		switch v := interface{}(m.GetPayInAmount()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreatePaymentIntentRequestValidationError{
					field:  "PayInAmount",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreatePaymentIntentRequestValidationError{
					field:  "PayInAmount",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPayInAmount()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreatePaymentIntentRequestValidationError{
				field:  "PayInAmount",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for PayOutCurrency

	if all {
		switch v := interface{}(m.GetPayOutMethod()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreatePaymentIntentRequestValidationError{
					field:  "PayOutMethod",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreatePaymentIntentRequestValidationError{
					field:  "PayOutMethod",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPayOutMethod()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreatePaymentIntentRequestValidationError{
				field:  "PayOutMethod",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CreatePaymentIntentRequestMultiError(errors)
	}

	return nil
}

// CreatePaymentIntentRequestMultiError is an error wrapping multiple
// validation errors returned by CreatePaymentIntentRequest.ValidateAll() if
// the designated constraints aren't met.
type CreatePaymentIntentRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreatePaymentIntentRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreatePaymentIntentRequestMultiError) AllErrors() []error { return m }

// CreatePaymentIntentRequestValidationError is the validation error returned
// by CreatePaymentIntentRequest.Validate if the designated constraints aren't met.
type CreatePaymentIntentRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePaymentIntentRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePaymentIntentRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePaymentIntentRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePaymentIntentRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePaymentIntentRequestValidationError) ErrorName() string {
	return "CreatePaymentIntentRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePaymentIntentRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePaymentIntentRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePaymentIntentRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePaymentIntentRequestValidationError{}

// Validate checks the field values on CreatePaymentIntentResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreatePaymentIntentResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreatePaymentIntentResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreatePaymentIntentResponseMultiError, or nil if none found.
func (m *CreatePaymentIntentResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreatePaymentIntentResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PaymentIntentId

	for idx, item := range m.GetPayInPaymentMethods() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CreatePaymentIntentResponseValidationError{
						field:  fmt.Sprintf("PayInPaymentMethods[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CreatePaymentIntentResponseValidationError{
						field:  fmt.Sprintf("PayInPaymentMethods[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CreatePaymentIntentResponseValidationError{
					field:  fmt.Sprintf("PayInPaymentMethods[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return CreatePaymentIntentResponseMultiError(errors)
	}

	return nil
}

// CreatePaymentIntentResponseMultiError is an error wrapping multiple
// validation errors returned by CreatePaymentIntentResponse.ValidateAll() if
// the designated constraints aren't met.
type CreatePaymentIntentResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreatePaymentIntentResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreatePaymentIntentResponseMultiError) AllErrors() []error { return m }

// CreatePaymentIntentResponseValidationError is the validation error returned
// by CreatePaymentIntentResponse.Validate if the designated constraints
// aren't met.
type CreatePaymentIntentResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePaymentIntentResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePaymentIntentResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePaymentIntentResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePaymentIntentResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePaymentIntentResponseValidationError) ErrorName() string {
	return "CreatePaymentIntentResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePaymentIntentResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePaymentIntentResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePaymentIntentResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePaymentIntentResponseValidationError{}

// Validate checks the field values on ConfirmPaymentRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ConfirmPaymentRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ConfirmPaymentRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ConfirmPaymentRequestMultiError, or nil if none found.
func (m *ConfirmPaymentRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ConfirmPaymentRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PaymentIntentId

	// no validation rules for PaymentReference

	// no validation rules for PaymentMethod

	if len(errors) > 0 {
		return ConfirmPaymentRequestMultiError(errors)
	}

	return nil
}

// ConfirmPaymentRequestMultiError is an error wrapping multiple validation
// errors returned by ConfirmPaymentRequest.ValidateAll() if the designated
// constraints aren't met.
type ConfirmPaymentRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ConfirmPaymentRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ConfirmPaymentRequestMultiError) AllErrors() []error { return m }

// ConfirmPaymentRequestValidationError is the validation error returned by
// ConfirmPaymentRequest.Validate if the designated constraints aren't met.
type ConfirmPaymentRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ConfirmPaymentRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ConfirmPaymentRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ConfirmPaymentRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ConfirmPaymentRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ConfirmPaymentRequestValidationError) ErrorName() string {
	return "ConfirmPaymentRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ConfirmPaymentRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConfirmPaymentRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ConfirmPaymentRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ConfirmPaymentRequestValidationError{}

// Validate checks the field values on ConfirmPaymentIntentResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ConfirmPaymentIntentResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ConfirmPaymentIntentResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ConfirmPaymentIntentResponseMultiError, or nil if none found.
func (m *ConfirmPaymentIntentResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ConfirmPaymentIntentResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return ConfirmPaymentIntentResponseMultiError(errors)
	}

	return nil
}

// ConfirmPaymentIntentResponseMultiError is an error wrapping multiple
// validation errors returned by ConfirmPaymentIntentResponse.ValidateAll() if
// the designated constraints aren't met.
type ConfirmPaymentIntentResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ConfirmPaymentIntentResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ConfirmPaymentIntentResponseMultiError) AllErrors() []error { return m }

// ConfirmPaymentIntentResponseValidationError is the validation error returned
// by ConfirmPaymentIntentResponse.Validate if the designated constraints
// aren't met.
type ConfirmPaymentIntentResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ConfirmPaymentIntentResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ConfirmPaymentIntentResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ConfirmPaymentIntentResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ConfirmPaymentIntentResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ConfirmPaymentIntentResponseValidationError) ErrorName() string {
	return "ConfirmPaymentIntentResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ConfirmPaymentIntentResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConfirmPaymentIntentResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ConfirmPaymentIntentResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ConfirmPaymentIntentResponseValidationError{}

// Validate checks the field values on RejectPaymentIntentRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RejectPaymentIntentRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RejectPaymentIntentRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RejectPaymentIntentRequestMultiError, or nil if none found.
func (m *RejectPaymentIntentRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *RejectPaymentIntentRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PaymentIntentId

	// no validation rules for PaymentReference

	// no validation rules for Reason

	if len(errors) > 0 {
		return RejectPaymentIntentRequestMultiError(errors)
	}

	return nil
}

// RejectPaymentIntentRequestMultiError is an error wrapping multiple
// validation errors returned by RejectPaymentIntentRequest.ValidateAll() if
// the designated constraints aren't met.
type RejectPaymentIntentRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RejectPaymentIntentRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RejectPaymentIntentRequestMultiError) AllErrors() []error { return m }

// RejectPaymentIntentRequestValidationError is the validation error returned
// by RejectPaymentIntentRequest.Validate if the designated constraints aren't met.
type RejectPaymentIntentRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RejectPaymentIntentRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RejectPaymentIntentRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RejectPaymentIntentRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RejectPaymentIntentRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RejectPaymentIntentRequestValidationError) ErrorName() string {
	return "RejectPaymentIntentRequestValidationError"
}

// Error satisfies the builtin error interface
func (e RejectPaymentIntentRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRejectPaymentIntentRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RejectPaymentIntentRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RejectPaymentIntentRequestValidationError{}

// Validate checks the field values on RejectPaymentIntentResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RejectPaymentIntentResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RejectPaymentIntentResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RejectPaymentIntentResponseMultiError, or nil if none found.
func (m *RejectPaymentIntentResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *RejectPaymentIntentResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return RejectPaymentIntentResponseMultiError(errors)
	}

	return nil
}

// RejectPaymentIntentResponseMultiError is an error wrapping multiple
// validation errors returned by RejectPaymentIntentResponse.ValidateAll() if
// the designated constraints aren't met.
type RejectPaymentIntentResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RejectPaymentIntentResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RejectPaymentIntentResponseMultiError) AllErrors() []error { return m }

// RejectPaymentIntentResponseValidationError is the validation error returned
// by RejectPaymentIntentResponse.Validate if the designated constraints
// aren't met.
type RejectPaymentIntentResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RejectPaymentIntentResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RejectPaymentIntentResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RejectPaymentIntentResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RejectPaymentIntentResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RejectPaymentIntentResponseValidationError) ErrorName() string {
	return "RejectPaymentIntentResponseValidationError"
}

// Error satisfies the builtin error interface
func (e RejectPaymentIntentResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRejectPaymentIntentResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RejectPaymentIntentResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RejectPaymentIntentResponseValidationError{}

// Validate checks the field values on
// CreatePaymentIntentResponse_PaymentMethod with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CreatePaymentIntentResponse_PaymentMethod) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on
// CreatePaymentIntentResponse_PaymentMethod with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in
// CreatePaymentIntentResponse_PaymentMethodMultiError, or nil if none found.
func (m *CreatePaymentIntentResponse_PaymentMethod) ValidateAll() error {
	return m.validate(true)
}

func (m *CreatePaymentIntentResponse_PaymentMethod) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PaymentUrl

	// no validation rules for ProviderId

	// no validation rules for PaymentMethod

	if len(errors) > 0 {
		return CreatePaymentIntentResponse_PaymentMethodMultiError(errors)
	}

	return nil
}

// CreatePaymentIntentResponse_PaymentMethodMultiError is an error wrapping
// multiple validation errors returned by
// CreatePaymentIntentResponse_PaymentMethod.ValidateAll() if the designated
// constraints aren't met.
type CreatePaymentIntentResponse_PaymentMethodMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreatePaymentIntentResponse_PaymentMethodMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreatePaymentIntentResponse_PaymentMethodMultiError) AllErrors() []error { return m }

// CreatePaymentIntentResponse_PaymentMethodValidationError is the validation
// error returned by CreatePaymentIntentResponse_PaymentMethod.Validate if the
// designated constraints aren't met.
type CreatePaymentIntentResponse_PaymentMethodValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreatePaymentIntentResponse_PaymentMethodValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreatePaymentIntentResponse_PaymentMethodValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreatePaymentIntentResponse_PaymentMethodValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreatePaymentIntentResponse_PaymentMethodValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreatePaymentIntentResponse_PaymentMethodValidationError) ErrorName() string {
	return "CreatePaymentIntentResponse_PaymentMethodValidationError"
}

// Error satisfies the builtin error interface
func (e CreatePaymentIntentResponse_PaymentMethodValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreatePaymentIntentResponse_PaymentMethod.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreatePaymentIntentResponse_PaymentMethodValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreatePaymentIntentResponse_PaymentMethodValidationError{}
