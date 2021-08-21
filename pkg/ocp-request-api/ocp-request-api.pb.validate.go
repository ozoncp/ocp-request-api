// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/ocp-request-api/ocp-request-api.proto

package ocp_request_api

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
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
)

// Validate checks the field values on ListRequestsV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ListRequestsV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if val := m.GetLimit(); val <= 0 || val > 10000 {
		return ListRequestsV1RequestValidationError{
			field:  "Limit",
			reason: "value must be inside range (0, 10000]",
		}
	}

	if m.GetOffset() < 0 {
		return ListRequestsV1RequestValidationError{
			field:  "Offset",
			reason: "value must be greater than or equal to 0",
		}
	}

	return nil
}

// ListRequestsV1RequestValidationError is the validation error returned by
// ListRequestsV1Request.Validate if the designated constraints aren't met.
type ListRequestsV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListRequestsV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListRequestsV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListRequestsV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListRequestsV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListRequestsV1RequestValidationError) ErrorName() string {
	return "ListRequestsV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListRequestsV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListRequestsV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListRequestsV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListRequestsV1RequestValidationError{}

// Validate checks the field values on ListRequestsV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ListRequestsV1Response) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetRequests() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListRequestsV1ResponseValidationError{
					field:  fmt.Sprintf("Requests[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// ListRequestsV1ResponseValidationError is the validation error returned by
// ListRequestsV1Response.Validate if the designated constraints aren't met.
type ListRequestsV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListRequestsV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListRequestsV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListRequestsV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListRequestsV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListRequestsV1ResponseValidationError) ErrorName() string {
	return "ListRequestsV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListRequestsV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListRequestsV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListRequestsV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListRequestsV1ResponseValidationError{}

// Validate checks the field values on MultiCreateRequestV1Request with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *MultiCreateRequestV1Request) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetRequests() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MultiCreateRequestV1RequestValidationError{
					field:  fmt.Sprintf("Requests[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// MultiCreateRequestV1RequestValidationError is the validation error returned
// by MultiCreateRequestV1Request.Validate if the designated constraints
// aren't met.
type MultiCreateRequestV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MultiCreateRequestV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MultiCreateRequestV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MultiCreateRequestV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MultiCreateRequestV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MultiCreateRequestV1RequestValidationError) ErrorName() string {
	return "MultiCreateRequestV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e MultiCreateRequestV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMultiCreateRequestV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MultiCreateRequestV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MultiCreateRequestV1RequestValidationError{}

// Validate checks the field values on MultiCreateRequestV1Response with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *MultiCreateRequestV1Response) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// MultiCreateRequestV1ResponseValidationError is the validation error returned
// by MultiCreateRequestV1Response.Validate if the designated constraints
// aren't met.
type MultiCreateRequestV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MultiCreateRequestV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MultiCreateRequestV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MultiCreateRequestV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MultiCreateRequestV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MultiCreateRequestV1ResponseValidationError) ErrorName() string {
	return "MultiCreateRequestV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e MultiCreateRequestV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMultiCreateRequestV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MultiCreateRequestV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MultiCreateRequestV1ResponseValidationError{}

// Validate checks the field values on UpdateRequestV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *UpdateRequestV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetRequestId() <= 0 {
		return UpdateRequestV1RequestValidationError{
			field:  "RequestId",
			reason: "value must be greater than 0",
		}
	}

	if m.GetUserId() <= 0 {
		return UpdateRequestV1RequestValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
	}

	// no validation rules for Type

	// no validation rules for Text

	return nil
}

// UpdateRequestV1RequestValidationError is the validation error returned by
// UpdateRequestV1Request.Validate if the designated constraints aren't met.
type UpdateRequestV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateRequestV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateRequestV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateRequestV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateRequestV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateRequestV1RequestValidationError) ErrorName() string {
	return "UpdateRequestV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateRequestV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateRequestV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateRequestV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateRequestV1RequestValidationError{}

// Validate checks the field values on UpdateRequestV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *UpdateRequestV1Response) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// UpdateRequestV1ResponseValidationError is the validation error returned by
// UpdateRequestV1Response.Validate if the designated constraints aren't met.
type UpdateRequestV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateRequestV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateRequestV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateRequestV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateRequestV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateRequestV1ResponseValidationError) ErrorName() string {
	return "UpdateRequestV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateRequestV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateRequestV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateRequestV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateRequestV1ResponseValidationError{}

// Validate checks the field values on CreateRequestV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreateRequestV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetUserId() <= 0 {
		return CreateRequestV1RequestValidationError{
			field:  "UserId",
			reason: "value must be greater than 0",
		}
	}

	// no validation rules for Type

	// no validation rules for Text

	return nil
}

// CreateRequestV1RequestValidationError is the validation error returned by
// CreateRequestV1Request.Validate if the designated constraints aren't met.
type CreateRequestV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateRequestV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateRequestV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateRequestV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateRequestV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateRequestV1RequestValidationError) ErrorName() string {
	return "CreateRequestV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateRequestV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateRequestV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateRequestV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateRequestV1RequestValidationError{}

// Validate checks the field values on CreateRequestV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreateRequestV1Response) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for RequestId

	return nil
}

// CreateRequestV1ResponseValidationError is the validation error returned by
// CreateRequestV1Response.Validate if the designated constraints aren't met.
type CreateRequestV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateRequestV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateRequestV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateRequestV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateRequestV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateRequestV1ResponseValidationError) ErrorName() string {
	return "CreateRequestV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateRequestV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateRequestV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateRequestV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateRequestV1ResponseValidationError{}

// Validate checks the field values on RemoveRequestV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *RemoveRequestV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetRequestId() <= 0 {
		return RemoveRequestV1RequestValidationError{
			field:  "RequestId",
			reason: "value must be greater than 0",
		}
	}

	return nil
}

// RemoveRequestV1RequestValidationError is the validation error returned by
// RemoveRequestV1Request.Validate if the designated constraints aren't met.
type RemoveRequestV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemoveRequestV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemoveRequestV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemoveRequestV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemoveRequestV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemoveRequestV1RequestValidationError) ErrorName() string {
	return "RemoveRequestV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e RemoveRequestV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemoveRequestV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemoveRequestV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemoveRequestV1RequestValidationError{}

// Validate checks the field values on RemoveRequestV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *RemoveRequestV1Response) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// RemoveRequestV1ResponseValidationError is the validation error returned by
// RemoveRequestV1Response.Validate if the designated constraints aren't met.
type RemoveRequestV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemoveRequestV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemoveRequestV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemoveRequestV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemoveRequestV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemoveRequestV1ResponseValidationError) ErrorName() string {
	return "RemoveRequestV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e RemoveRequestV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemoveRequestV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemoveRequestV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemoveRequestV1ResponseValidationError{}

// Validate checks the field values on DescribeRequestV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *DescribeRequestV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetRequestId() <= 0 {
		return DescribeRequestV1RequestValidationError{
			field:  "RequestId",
			reason: "value must be greater than 0",
		}
	}

	return nil
}

// DescribeRequestV1RequestValidationError is the validation error returned by
// DescribeRequestV1Request.Validate if the designated constraints aren't met.
type DescribeRequestV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DescribeRequestV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DescribeRequestV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DescribeRequestV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DescribeRequestV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DescribeRequestV1RequestValidationError) ErrorName() string {
	return "DescribeRequestV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e DescribeRequestV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDescribeRequestV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DescribeRequestV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DescribeRequestV1RequestValidationError{}

// Validate checks the field values on DescribeRequestV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *DescribeRequestV1Response) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetRequest()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DescribeRequestV1ResponseValidationError{
				field:  "Request",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// DescribeRequestV1ResponseValidationError is the validation error returned by
// DescribeRequestV1Response.Validate if the designated constraints aren't met.
type DescribeRequestV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DescribeRequestV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DescribeRequestV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DescribeRequestV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DescribeRequestV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DescribeRequestV1ResponseValidationError) ErrorName() string {
	return "DescribeRequestV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DescribeRequestV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDescribeRequestV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DescribeRequestV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DescribeRequestV1ResponseValidationError{}

// Validate checks the field values on Request with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Request) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	// no validation rules for UserId

	// no validation rules for Type

	// no validation rules for Text

	return nil
}

// RequestValidationError is the validation error returned by Request.Validate
// if the designated constraints aren't met.
type RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RequestValidationError) ErrorName() string { return "RequestValidationError" }

// Error satisfies the builtin error interface
func (e RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RequestValidationError{}
