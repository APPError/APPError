package apperror

import (
    "encoding/json"
    "errors"
    "fmt"
    "reflect"
)

// copyAny 尝试对 any 类型执行受控深拷贝。
// 优先检查 Clone() any 接口，其次 Copy() any 接口。
// 对于未实现以上接口的 slice/map，执行逐元素浅层拷贝防止外部篡改。
// 对于其他类型直接返回原值（引用类型需自行实现 Clone/Copy 才能完全隔离）。
func copyAny(v any) any {
    if v == nil {
        return nil
    }
    if cloner, ok := v.(interface{ Clone() any }); ok {
        return cloner.Clone()
    }
    if copier, ok := v.(interface{ Copy() any }); ok {
        return copier.Copy()
    }

    // 对 slice / map 做递归深拷贝，防止外部修改共享底层数组
    rv := reflect.ValueOf(v)
    if rv.Kind() == reflect.Slice {
        if rv.IsNil() {
            return v
        }
        dst := reflect.MakeSlice(rv.Type(), rv.Len(), rv.Len())
        for i := 0; i < rv.Len(); i++ {
            dst.Index(i).Set(reflect.ValueOf(copyAny(rv.Index(i).Interface())))
        }
        return dst.Interface()
    }
    if rv.Kind() == reflect.Map {
        if rv.IsNil() {
            return v
        }
        dst := reflect.MakeMap(rv.Type())
        for _, key := range rv.MapKeys() {
            dst.SetMapIndex(key, reflect.ValueOf(copyAny(rv.MapIndex(key).Interface())))
        }
        return dst.Interface()
    }
    if rv.Kind() == reflect.Ptr {
        if rv.IsNil() {
            return v
        }
        elem := copyAny(rv.Elem().Interface())
        newPtr := reflect.New(rv.Type().Elem())
        newPtr.Elem().Set(reflect.ValueOf(elem))
        return newPtr.Interface()
    }
    return v
}

// safeDeepEqual 安全地执行深比较，通过 recover 防护 reflect.DeepEqual 对含函数/通道等
// 不可比较类型的 panic 风险。
func safeDeepEqual(a, b any) (ok bool) {
    defer func() {
        if r := recover(); r != nil {
            ok = false
        }
    }()
    return reflect.DeepEqual(a, b)
}

// needsWrap 判断 err 是否需要包装为 UniformException。
// 当 err 的类型没有导出字段时，json.Marshal 无法序列化有效信息，需要包装。
// jsonSerializable 递归检查 reflect.Type 是否能被 encoding/json 安全序列化。
// 函数、通道、复数等类型无法被 json.Marshal 处理，会被返回 UnsupportedTypeError。
func jsonSerializable(t reflect.Type) bool {
    switch t.Kind() {
    case reflect.Func, reflect.Chan, reflect.Complex64, reflect.Complex128:
        return false
    case reflect.Ptr, reflect.Slice, reflect.Array:
        return jsonSerializable(t.Elem())
    case reflect.Map:
        return jsonSerializable(t.Key()) && jsonSerializable(t.Elem())
    case reflect.Struct:
        for i := 0; i < t.NumField(); i++ {
            f := t.Field(i)
            if !f.IsExported() {
                continue
            }
            if !jsonSerializable(f.Type) {
                return false
            }
        }
        return true
    default:
        // 基本类型（Bool, Int*, Uint*, Float*, String, Interface 等）均可安全序列化
        return true
    }
}

func needsWrap(err error) bool {
    if err == nil {
        return false
    }
    // 如果已实现 json.Marshaler 接口，可自行控制序列化结果，无需包装
    if _, ok := err.(json.Marshaler); ok {
        return false
    }
    t := reflect.TypeOf(err)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    switch t.Kind() {
    case reflect.Struct:
        hasExported := false
        for i := 0; i < t.NumField(); i++ {
            if t.Field(i).IsExported() {
                hasExported = true
                if !jsonSerializable(t.Field(i).Type) {
                    return true
                }
            }
        }
        return !hasExported
    case reflect.Map, reflect.Slice, reflect.Array, reflect.String, reflect.Bool,
        reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
        reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
        reflect.Float32, reflect.Float64:
        return false
    default:
        return true
    }
}

// normalizeErrors 深拷贝 src 切片并过滤 nil，再将其中需要包装的 error 转换为 jsonSafeError，
// 保证 JSON 序列化能正确输出错误信息，且返回的切片不与 src 共享底层数组。
func normalizeErrors(src []error) []error {
    if src == nil {
        return nil
    }
    dst := make([]error, 0, len(src))
    for _, err := range src {
        if err == nil {
            continue
        }
        // 已经是 jsonSafeError 的不再重复包装（幂等保护）
        if _, ok := err.(*jsonSafeError); ok {
            dst = append(dst, err)
            continue
        }
        if needsWrap(err) {
            dst = append(
                dst, &jsonSafeError{
                    msg: err.Error(),
                    err: err, // 保留原始错误引用，供 Is/As 委派
                },
            )
            continue
        }
        dst = append(dst, err)
    }
    if len(dst) == 0 {
        return nil
    }
    return dst
}

// jsonSafeError 包装没有导出字段的 error，保留原始错误引用用于 Is/As 匹配，
// 同时提供安全的 JSON 序列化（输出为字符串而非 {}）。
type jsonSafeError struct {
    msg string // 错误消息，用于 Error() 和 JSON 序列化
    err error  // 原始错误，用于 Is/As 委派
}

func (e *jsonSafeError) Error() string { return e.msg }

func (e *jsonSafeError) Unwrap() error { return e.err }

// Is 委派给原始 error，确保 errors.Is(result, originalErr) 能匹配
func (e *jsonSafeError) Is(target error) bool {
    return errors.Is(e.err, target)
}

// As 委派给原始 error，确保 errors.As(result, &targetType) 能匹配
func (e *jsonSafeError) As(target any) bool {
    return As(e.err, target)
}

// MarshalJSON 序列化为消息字符串，避免 {} 空对象
func (e *jsonSafeError) MarshalJSON() ([]byte, error) {
    return json.Marshal(e.msg)
}

// uniformException 统一异常
type uniformException struct {
    Data             any              `json:"data,omitempty"`        // Data 数据
    Description      ErrorDescription `json:"description,omitempty"` // Description 描述
    Detail           any              `json:"detail,omitempty"`      // Detail 详情
    Errors           []error          `json:"errors,omitempty"`      // Errors 错误集合
    Information      any              `json:"information,omitempty"` // Information 信息
    Message          string           `json:"message,omitempty"`     // Message 消息
    Profile          any              `json:"profile,omitempty"`     // Profile
    Status           any              `json:"status,omitempty"`      // Status 状态
    equationFunction func(error, error) bool
}

// copy 深拷贝 uniformException，确保 Errors 切片不与原对象共享底层数组，
// 并对 Data、Detail、Information、Profile、Status 执行受控深拷贝。
func (e *uniformException) copy() *uniformException {
    ex := new(uniformException)
    if e != nil {
        *ex = *e
        ex.Data = copyAny(e.Data)
        ex.Detail = copyAny(e.Detail)
        ex.Information = copyAny(e.Information)
        ex.Profile = copyAny(e.Profile)
        ex.Status = copyAny(e.Status)
        ex.Errors = e.LoadErrors()
        // 深拷贝 Description 内部的可变引用类型字段（Data、Detail、Information、Status、Type）
        ex.Description.Data = copyAny(e.Description.Data)
        ex.Description.Detail = copyAny(e.Description.Detail)
        ex.Description.Information = copyAny(e.Description.Information)
        ex.Description.Status = copyAny(e.Description.Status)
        ex.Description.Type = copyAny(e.Description.Type)
    }
    return ex
}

func (e *uniformException) Error() string {
    if e == nil {
        return ""
    }
    if bytes, err := json.Marshal(e); err == nil {
        return string(bytes)
    }
    return fmt.Sprintf(
        `uniformException{
        Data:%v
        Description:%v
        Detail:%v
        Errors:%v
        Information:%v
        Message:%s
        Profile:%v
        Status:%v
    }`,
        e.Data,
        e.Description,
        e.Detail,
        e.Errors,
        e.Information,
        e.Message,
        e.Profile,
        e.Status,
    )
}

func (e *uniformException) String() string {
    return e.Error()
}

func (e *uniformException) Is(target error) bool {
    if e == nil {
        return false
    }
    if e.equationFunction != nil {
        return e.equationFunction(e, target)
    }
    if other, ok := target.(GeneralError); ok {
        getters := []func(GeneralError) any{
            func(ge GeneralError) any { return ge.LoadData() },
            func(ge GeneralError) any { return ge.LoadDescription() },
            func(ge GeneralError) any { return ge.LoadDetail() },
            func(ge GeneralError) any { return ge.LoadErrors() },
            func(ge GeneralError) any { return ge.LoadInformation() },
            func(ge GeneralError) any { return ge.LoadMessage() },
            func(ge GeneralError) any { return ge.LoadProfile() },
            func(ge GeneralError) any { return ge.LoadStatus() },
        }
        for _, getter := range getters {
            if !safeDeepEqual(getter(e), getter(other)) {
                return false
            }
        }
        return true
    }
    return e == target
}

func (e *uniformException) Unwrap() []error {
    return e.LoadErrors()
}

func (e *uniformException) LoadData() any {
    if e != nil {
        return copyAny(e.Data)
    }
    return nil
}

func (e *uniformException) LoadDescription() ErrorDescription {
    if e != nil {
        return e.Description
    }
    return ErrorDescription{}
}

func (e *uniformException) LoadDetail() any {
    if e != nil {
        return copyAny(e.Detail)
    }
    return nil
}

func (e *uniformException) LoadErrors() []error {
    if e == nil || e.Errors == nil {
        return nil
    }
    dst := make([]error, len(e.Errors))
    copy(dst, e.Errors)
    return dst
}

func (e *uniformException) LoadInformation() any {
    if e != nil {
        return copyAny(e.Information)
    }
    return nil
}

func (e *uniformException) LoadMessage() string {
    if e != nil {
        return e.Message
    }
    return ""
}

func (e *uniformException) LoadProfile() any {
    if e != nil {
        return copyAny(e.Profile)
    }
    return nil
}

func (e *uniformException) LoadStatus() any {
    if e != nil {
        return copyAny(e.Status)
    }
    return nil
}

func (e *uniformException) LoadEquationFunction() func(error, error) bool {
    if e != nil {
        return e.equationFunction
    }
    return nil
}

func (e *uniformException) WithData(data any) GeneralError {
    ex := e.copy()
    ex.Data = copyAny(data)
    return ex
}

func (e *uniformException) WithDescription(description ErrorDescription) GeneralError {
    ex := e.copy()
    ex.Description = description
    return ex
}

func (e *uniformException) WithDetail(detail any) GeneralError {
    ex := e.copy()
    ex.Detail = copyAny(detail)
    return ex
}

func (e *uniformException) WithErrors(errors []error) GeneralError {
    ex := e.copy()
    ex.Errors = normalizeErrors(errors)
    return ex
}

func (e *uniformException) WithInformation(information any) GeneralError {
    ex := e.copy()
    ex.Information = copyAny(information)
    return ex
}

func (e *uniformException) WithMessage(message string) GeneralError {
    ex := e.copy()
    ex.Message = message
    return ex
}

func (e *uniformException) WithProfile(profile any) GeneralError {
    ex := e.copy()
    ex.Profile = copyAny(profile)
    return ex
}

func (e *uniformException) WithStatus(status any) GeneralError {
    ex := e.copy()
    ex.Status = copyAny(status)
    return ex
}

func (e *uniformException) WithEquationFunction(equation func(error, error) bool) GeneralError {
    ex := e.copy()
    ex.equationFunction = equation
    return ex
}

func (e *uniformException) ConvertToGeneralError() GeneralError {
    return e
}

func (e *uniformException) ConvertToProjectError() ProjectError {
    return e
}

func (e *uniformException) ConvertToRuntimeError() RuntimeError {
    return e
}

func (e *uniformException) ConvertToServiceError() ServiceError {
    return e
}

func (e *uniformException) ConvertToUniformError() UniformError {
    return e
}

func NewUniformExceptionBuilder() UniformExceptionBasicBuilder {
    return &uniformExceptionBuilder{}
}

type UniformExceptionBasicBuilder interface {
    Data(any) UniformExceptionBasicBuilder
    Description(ErrorDescription) UniformExceptionBasicBuilder
    Detail(any) UniformExceptionBasicBuilder
    Errors([]error) UniformExceptionBasicBuilder
    Information(any) UniformExceptionBasicBuilder
    Message(string) UniformExceptionBasicBuilder
    Profile(any) UniformExceptionBasicBuilder
    Status(any) UniformExceptionBasicBuilder
    EquationFunction(func(error, error) bool) UniformExceptionBasicBuilder
    BuildGeneralError() GeneralError
    BuildProjectError() ProjectError
    BuildRuntimeError() RuntimeError
    BuildServiceError() ServiceError
    BuildUniformError() UniformError
}

type uniformExceptionBuilder struct {
    uniformException uniformException
}

func (b *uniformExceptionBuilder) Data(data any) UniformExceptionBasicBuilder {
    b.uniformException.Data = copyAny(data)
    return b
}

func (b *uniformExceptionBuilder) Description(description ErrorDescription) UniformExceptionBasicBuilder {
    b.uniformException.Description = description
    return b
}

func (b *uniformExceptionBuilder) Detail(detail any) UniformExceptionBasicBuilder {
    b.uniformException.Detail = copyAny(detail)
    return b
}

func (b *uniformExceptionBuilder) Errors(errors []error) UniformExceptionBasicBuilder {
    b.uniformException.Errors = normalizeErrors(errors)
    return b
}

func (b *uniformExceptionBuilder) Information(information any) UniformExceptionBasicBuilder {
    b.uniformException.Information = copyAny(information)
    return b
}

func (b *uniformExceptionBuilder) Message(message string) UniformExceptionBasicBuilder {
    b.uniformException.Message = message
    return b
}

func (b *uniformExceptionBuilder) Profile(profile any) UniformExceptionBasicBuilder {
    b.uniformException.Profile = copyAny(profile)
    return b
}

func (b *uniformExceptionBuilder) Status(status any) UniformExceptionBasicBuilder {
    b.uniformException.Status = copyAny(status)
    return b
}

func (b *uniformExceptionBuilder) EquationFunction(equation func(error, error) bool) UniformExceptionBasicBuilder {
    b.uniformException.equationFunction = equation
    return b
}

func (b *uniformExceptionBuilder) BuildGeneralError() GeneralError {
    return b.uniformException.copy()
}

func (b *uniformExceptionBuilder) BuildProjectError() ProjectError {
    return b.uniformException.copy()
}

func (b *uniformExceptionBuilder) BuildRuntimeError() RuntimeError {
    return b.uniformException.copy()
}

func (b *uniformExceptionBuilder) BuildServiceError() ServiceError {
    return b.uniformException.copy()
}

func (b *uniformExceptionBuilder) BuildUniformError() UniformError {
    return b.uniformException.copy()
}
