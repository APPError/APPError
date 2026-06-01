package apperror

import "reflect"

type GeneralError interface {
    error
    LoadData() any
    LoadDescription() ErrorDescription
    LoadDetail() any
    LoadErrors() []error
    LoadInformation() any
    LoadMessage() string
    LoadProfile() any
    LoadStatus() any
    WithData(any) GeneralError
    WithDescription(ErrorDescription) GeneralError
    WithDetail(any) GeneralError
    WithErrors([]error) GeneralError
    WithInformation(any) GeneralError
    WithMessage(string) GeneralError
    WithProfile(any) GeneralError
    WithStatus(any) GeneralError
}

type ProjectError interface {
    GeneralError
}

type RuntimeError interface {
    ProjectError
}

type ServiceError interface {
    RuntimeError
}

type UniformError interface {
    ServiceError
}

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
func Join(errors ...error) UniformError {
    if errors == nil || len(errors) == 0 {
        return nil
    }
    n := 0
    for _, err := range errors {
        if err != nil {
            n++
        }
    }
    if n == 0 {
        return nil
    }
    validErrors := make([]error, 0, n)
    for _, err := range errors {
        if err != nil {
            validErrors = append(validErrors, err)
        }
    }
    return NewUniformExceptionBuilder().Errors(validErrors).BuildUniformError()
}

// Is reports whether any error in err's tree matches target.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly
// calling its Errors(), LoadErrors(), Unwrap() []error, or Unwrap() error method.
// When err wraps multiple errors, Is examines err followed by a depth-first
// traversal of its children.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. An Is method should only
// shallowly compare err and the target and not call Unwrap on either.
func Is(err, target error) bool {
    for {
        if err == nil || target == nil {
            return err == target
        }
        if reflect.TypeOf(err).Comparable() && reflect.TypeOf(target).Comparable() && err == target {
            return true
        }
        if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
            return true
        }
        if x, ok := err.(interface{ Errors() []error }); ok {
            for _, child := range x.Errors() {
                if child != nil && Is(child, target) {
                    return true
                }
            }
        }
        if x, ok := err.(interface{ LoadErrors() []error }); ok {
            for _, child := range x.LoadErrors() {
                if child != nil && Is(child, target) {
                    return true
                }
            }
        }
        if x, ok := err.(interface{ Unwrap() []error }); ok {
            for _, child := range x.Unwrap() {
                if child != nil && Is(child, target) {
                    return true
                }
            }
        }
        if x, ok := err.(interface{ Unwrap() error }); ok {
            child := x.Unwrap()
            if child != nil {
                err = child
                continue
            }
        }
        return false
    }
}

// As finds the first error in err's tree that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly
// calling its Errors(), LoadErrors(), Unwrap() []error, or Unwrap() error method.
// When err wraps multiple errors, As examines err followed by a depth-first traversal
// of its children.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(any) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
func As(err error, target any) bool {
    if err == nil || target == nil {
        return false
    }
    val := reflect.ValueOf(target)
    typ := val.Type()
    if typ.Kind() != reflect.Ptr || val.IsNil() {
        return false
    }
    targetType := typ.Elem()
    if targetType.Kind() != reflect.Interface &&
        !targetType.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
        return false
    }
    return as(err, target, val, targetType)
}

func as(err error, target any, targetVal reflect.Value, targetType reflect.Type) bool {
    for {
        if reflect.TypeOf(err).AssignableTo(targetType) {
            targetVal.Elem().Set(reflect.ValueOf(err))
            return true
        }
        if x, ok := err.(interface{ As(any) bool }); ok && x.As(target) {
            return true
        }
        if x, ok := err.(interface{ Errors() []error }); ok {
            for _, child := range x.Errors() {
                if child != nil && as(child, target, targetVal, targetType) {
                    return true
                }
            }
        }
        if x, ok := err.(interface{ LoadErrors() []error }); ok {
            for _, child := range x.LoadErrors() {
                if child != nil && as(child, target, targetVal, targetType) {
                    return true
                }
            }
        }
        if x, ok := err.(interface{ Unwrap() []error }); ok {
            for _, child := range x.Unwrap() {
                if child != nil && as(child, target, targetVal, targetType) {
                    return true
                }
            }
        }
        if x, ok := err.(interface{ Unwrap() error }); ok {
            child := x.Unwrap()
            if child != nil {
                err = child
                continue
            }
        }
        return false
    }
}