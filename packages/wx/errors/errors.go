/*
Package errors provides custom error types used throughout the application.
This file defines detailed error structs to distinguish and handle errors
related to URI parameter parsing, value conversion, and service initialization.

Author: [Your Name or Team]
Created: 2025-08-16
License: MIT License (or your project's license)

Examples of errors in this package:
- UriParamParseError: error when a URI parameter is not found.
- UriParamConvertError: error when converting parameter value types.
- BadRequestError: invalid request error.
- ParamMissMatchError: parameter mismatch error.
- ServiceInitError: error during service initialization.
*/

package errors

import (
	"fmt"
	"reflect"
)

type UriParamParseError struct {
	ParamName    string
	TypeOfStruct reflect.Type
}

func NewUriParamParseError(paramName string, typeOfStruct reflect.Type) error {
	return &UriParamParseError{
		ParamName:    paramName,
		TypeOfStruct: typeOfStruct,
	}
}
func (e *UriParamParseError) Error() string {
	return fmt.Sprintf("%s was not found in %s", e.ParamName, e.TypeOfStruct.String())
}

type UriParamConvertError struct {
	ParamName        string
	ValueSetType     reflect.Type
	fielValueSetType reflect.Type
}

func NewUriParamConvertError(paramName string, valueSetType reflect.Type, fielValueSetType reflect.Type) error {
	return &UriParamConvertError{
		ParamName:        paramName,
		ValueSetType:     valueSetType,
		fielValueSetType: fielValueSetType,
	}
}
func (e *UriParamConvertError) Error() string {
	return fmt.Sprintf("error converting from %s to %s", e.ValueSetType.String(), e.fielValueSetType.String())
}

type BadRequestError struct {
	Message string
}

func NewBadRequestError(message string) error {
	return &BadRequestError{
		Message: message,
	}
}

func (e *BadRequestError) Error() string {
	return e.Message

}

type ParamMissMatchError struct {
	Message string
}

func NewParamMissMatchError(message string) error {
	return &ParamMissMatchError{
		Message: message,
	}
}

func (e *ParamMissMatchError) Error() string {
	return e.Message

}

type ServiceInitError struct {
	Message string
}

func NewServiceInitError(message string) error {
	return &ServiceInitError{
		Message: message,
	}
}

func (e *ServiceInitError) Error() string {
	return e.Message

}
