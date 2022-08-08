package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	NoError = iota
	ValidatorError
	CustomError
)

type RuleError struct {
	Error     error
	ErrorType int
}

type ValidationError struct {
	Field string
	Err   error
}

type Validators struct{}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	result := make([]string, 0)

	for _, validationError := range v {
		ve := fmt.Sprintf("error in the field: %s - %v", validationError.Field, validationError.Err.Error())
		result = append(result, ve)
	}

	return strings.Join(result, "\n")
}

func (v Validators) Len(expectedLenStr string, actualValue interface{}) RuleError {
	expectedLenInt, err := strconv.Atoi(expectedLenStr)
	if err != nil {
		return RuleError{err, CustomError}
	}

	typeOf := reflect.TypeOf(actualValue).String()
	if typeOf == "string" {
		if expectedLenInt != len(fmt.Sprint(actualValue)) {
			return RuleError{errors.New("the length of the line must be exactly " + expectedLenStr), ValidatorError}
		}
	} else if typeOf == "[]string" {
		values, ok := actualValue.([]string)
		if !ok {
			return RuleError{errors.New("it is not possible to convert the value"), CustomError}
		}

		for _, value := range values {
			if expectedLenInt != len(value) {
				return RuleError{
					errors.New("the length of the line must be exactly " + expectedLenStr),
					ValidatorError,
				}
			}
		}
	}

	return RuleError{Error: nil, ErrorType: NoError}
}

func (v Validators) Min(minExpectedLenString string, actualValue interface{}) RuleError {
	minExpectedLenInt, err := strconv.Atoi(minExpectedLenString)
	if err != nil {
		return RuleError{err, CustomError}
	}

	value, ok := actualValue.(int)
	if !ok {
		return RuleError{errors.New("it is not possible to convert the value"), CustomError}
	}

	if value < minExpectedLenInt {
		return RuleError{errors.New("the number cannot be less than " + minExpectedLenString), ValidatorError}
	}

	return RuleError{Error: nil, ErrorType: NoError}
}

func (v Validators) Max(maxExpectedLenString string, actualValue interface{}) RuleError {
	maxExpectedLenInt, err := strconv.Atoi(maxExpectedLenString)
	if err != nil {
		return RuleError{err, CustomError}
	}

	value, ok := actualValue.(int)
	if !ok {
		return RuleError{errors.New("it is not possible to convert the value"), CustomError}
	}

	if value > maxExpectedLenInt {
		return RuleError{errors.New("the number cannot be greater than " + maxExpectedLenString), ValidatorError}
	}

	return RuleError{Error: nil, ErrorType: NoError}
}

func (v Validators) Regexp(expr string, actualValue interface{}) RuleError {
	val, ok := actualValue.(string)
	if !ok {
		return RuleError{errors.New("it is not possible to convert the value"), CustomError}
	}

	re, err := regexp.Compile(expr)
	if err != nil {
		return RuleError{err, CustomError}
	}

	if !re.MatchString(val) {
		return RuleError{errors.New("the string must match the regular expression"), ValidatorError}
	}

	return RuleError{Error: nil, ErrorType: NoError}
}

func inSet(expectedValue string, set []string) bool {
	for _, in := range set {
		if in == expectedValue {
			return true
		}
	}

	return false
}

func (v Validators) In(expected string, actualValue interface{}) RuleError {
	expIns := strings.Split(expected, ",")

	typeOf := reflect.TypeOf(actualValue)

	kind := typeOf.Kind()
	if kind == reflect.Int || kind == reflect.String {
		value := fmt.Sprint(actualValue)
		if !inSet(value, expIns) {
			return RuleError{errors.New(value + " must be part of the set " + expected), ValidatorError}
		}
	}

	if kind == reflect.Slice && typeOf.String() == "[]int" {
		values, ok := actualValue.([]int)
		if !ok {
			return RuleError{errors.New("it is not possible to convert the value"), CustomError}
		}
		for _, value := range values {
			strValue := strconv.Itoa(value)
			if !inSet(strValue, expIns) {
				return RuleError{errors.New(strValue + " must be part of the set " + expected), ValidatorError}
			}
		}
	}

	if kind == reflect.Slice && typeOf.String() == "[]string" {
		values, ok := actualValue.([]string)
		if !ok {
			return RuleError{errors.New("it is not possible to convert the value"), CustomError}
		}

		for _, value := range values {
			if !inSet(value, expIns) {
				return RuleError{errors.New(value + " must be part of the set " + expected), ValidatorError}
			}
		}
	}

	return RuleError{Error: nil, ErrorType: NoError}
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)

	if reflect.Struct != value.Kind() {
		return errors.New("value is not a struct")
	}

	valueTypeOf := reflect.TypeOf(v)
	validationErrors := make(ValidationErrors, 0)

	validators := Validators{}
	my := reflect.ValueOf(validators)
	myT := reflect.TypeOf(validators)

	for i := 0; i < valueTypeOf.NumField(); i++ {
		field := valueTypeOf.Field(i)

		validate, ok := field.Tag.Lookup("validate")
		if !ok || validate == "" {
			continue
		}

		for _, rule := range strings.Split(validate, "|") {
			ruleKeyValue := strings.Split(rule, ":")

			ruleKey := strings.Title(ruleKeyValue[0])
			ruleValue := ruleKeyValue[1]

			_, methodExist := myT.MethodByName(ruleKey)
			if methodExist {
				params := []reflect.Value{reflect.ValueOf(ruleValue), value.Field(i)}

				result := my.MethodByName(ruleKey).Call(params)[0].Interface().(RuleError)

				if result.ErrorType == ValidatorError {
					validationErrors = append(validationErrors, ValidationError{
						Field: field.Name,
						Err:   result.Error,
					})
				} else if result.ErrorType == CustomError {
					return result.Error
				}
			} else {
				return fmt.Errorf("rule '%v' not exist", ruleKey)
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
