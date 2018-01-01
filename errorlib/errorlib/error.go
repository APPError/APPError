package errorlib

import (
    "encoding/json"
    "fmt"
    "strings"
)

// GeneralError 通用错误接口
type GeneralError interface {
    error
    LoadID() any
    LoadCode() any
    LoadData() any
    LoadDescription() ErrorDescription
    LoadDetail() any
    LoadErrors() []error
    LoadInformation() any
    LoadMessage() string
    LoadProfile() any
    LoadReason() any
    LoadStatus() any
    LoadType() any
    LoadEquationFunction() func(error, error) bool
    WithID(any) GeneralError
    WithCode(any) GeneralError
    WithData(any) GeneralError
    WithDescription(ErrorDescription) GeneralError
    WithDetail(any) GeneralError
    WithErrors([]error) GeneralError
    WithInformation(any) GeneralError
    WithMessage(string) GeneralError
    WithProfile(any) GeneralError
    WithReason(any) GeneralError
    WithStatus(any) GeneralError
    WithType(any) GeneralError
    WithEquationFunction(func(error, error) bool) GeneralError
}

// ProjectError 项目错误接口
type ProjectError interface {
    GeneralError
}

// RuntimeError 运行时错误接口
type RuntimeError interface {
    ProjectError
}

// ServiceError 服务错误接口
type ServiceError interface {
    RuntimeError
}

// UniformError 统一错误接口
type UniformError interface {
    ServiceError
}

// GeneralException 通用异常接口
type GeneralException interface {
    UniformError
}

// ProjectException 项目异常接口
type ProjectException interface {
    GeneralException
}

// RuntimeException 运行时异常接口
type RuntimeException interface {
    ProjectException
}

// ServiceException 服务异常接口
type ServiceException interface {
    RuntimeException
}

// UniformException 统一异常接口
type UniformException interface {
    ServiceException
}

// uniformException 统一异常
type uniformException struct {
    ID               any              `json:"id,omitempty"`          // ID 标识
    Code             any              `json:"code,omitempty"`        // Code 编码
    Data             any              `json:"data,omitempty"`        // Data 数据
    Description      ErrorDescription `json:"description,omitempty"` // Description 描述
    Detail           any              `json:"detail,omitempty"`      // Detail 详情
    Errors           []error          `json:"errors,omitempty"`      // Errors 错误集合
    Information      any              `json:"information,omitempty"` // Information 信息
    Message          string           `json:"message,omitempty"`     // Message 消息
    Profile          any              `json:"profile,omitempty"`     // Profile
    Reason           any              `json:"reason,omitempty"`      // Reason 原因
    Status           any              `json:"status,omitempty"`      // Status 状态
    Type             any              `json:"type,omitempty"`        // Type 类型
    equationFunction func(error, error) bool
}

// copy 深拷贝 uniformException，确保 Errors 切片不与原对象共享底层数组，
// 并对 Data、Detail、Information、Profile、Status 执行受控深拷贝。
func (e *uniformException) copy() *uniformException {
    if e == nil {
        return nil
    }
    ex := new(uniformException)
    *ex = *e
    ex.ID = copyAny(e.ID)
    ex.Code = copyAny(e.Code)
    ex.Data = copyAny(e.Data)
    ex.Description = e.Description
    if e.Description != nil {
        ex.Description = (&descriptionBuilder{}).loadFromErrorDescription(e.Description).BuildErrorDescription()
    }
    ex.Detail = copyAny(e.Detail)
    ex.Errors = normalizeErrors(e.LoadErrors())
    ex.Information = copyAny(e.Information)
    ex.Profile = copyAny(e.Profile)
    ex.Reason = copyAny(e.Reason)
    ex.Status = copyAny(e.Status)
    ex.Type = copyAny(e.Type)
    return ex
}

func (e *uniformException) Error() string {
    if e == nil {
        return ""
    }
    // 使用 copy() 深拷贝并规范化 Errors，防止 JSON 序列化时隐式循环
    if bytes, err := json.Marshal(e.copy()); err == nil {
        return string(bytes)
    }
    //  json.Marshal 失败时，逐字段尝试序列化以产生合法 JSON
    jsonVal := func(v any) string {
        if v == nil {
            return "null"
        }
        if b, err := json.Marshal(v); err == nil {
            return string(b)
        }
        return fmt.Sprintf("%q", fmt.Sprintf("%v", v))
    }
    jsonErrs := func(errs []error) string {
        if errs == nil {
            return "null"
        }
        var s strings.Builder
        s.WriteString("[")
        for i, er := range errs {
            if i > 0 {
                s.WriteString(",")
            }
            if er == nil {
                s.WriteString("null")
            } else {
                s.WriteString(fmt.Sprintf("%q", er.Error()))
            }
        }
        s.WriteString("]")
        return s.String()
    }
    return fmt.Sprintf(
        `{
            "id":%s,
            "code":%s,
            "data":%s,
            "description":%s,
            "detail":%s,
            "errors":%s,
            "information":%s,
            "message":%q,
            "profile":%s,
            "reason":%s,
            "status":%s,
            "type":%s
        }`,
        jsonVal(e.ID),
        jsonVal(e.Code),
        jsonVal(e.Data),
        jsonVal(e.Description),
        jsonVal(e.Detail),
        jsonErrs(e.LoadErrors()),
        jsonVal(e.Information),
        e.Message,
        jsonVal(e.Profile),
        jsonVal(e.Reason),
        jsonVal(e.Status),
        jsonVal(e.Type),
    )
}

func (e *uniformException) String() string {
    if e == nil {
        return ""
    }
    return e.Error()
}

func (e *uniformException) Is(target error) bool {
    if e == nil {
        return e == target
    }
    if e.equationFunction != nil {
        return e.equationFunction(e, target)
    }
    var generalError GeneralError
    if As(target, &generalError) {
        for _, fn := range []func(a, b GeneralError) bool{
            func(a, b GeneralError) bool { return safeDeepEqual(a.LoadID(), b.LoadID()) },
            func(a, b GeneralError) bool { return safeDeepEqual(a.LoadCode(), b.LoadCode()) },
            func(a, b GeneralError) bool { return safeDeepEqual(a.LoadData(), b.LoadData()) },
            func(a, b GeneralError) bool { return safeDeepEqual(a.LoadDescription(), b.LoadDescription()) },
            func(a, b GeneralError) bool { return safeDeepEqual(a.LoadDetail(), b.LoadDetail()) },
            func(a, b GeneralError) bool { return errorsEquals(a.LoadErrors(), b.LoadErrors()) },
            func(a, b GeneralError) bool { return safeDeepEqual(a.LoadInformation(), b.LoadInformation()) },
            func(a, b GeneralError) bool { return safeDeepEqual(a.LoadMessage(), b.LoadMessage()) },
            func(a, b GeneralError) bool { return safeDeepEqual(a.LoadProfile(), b.LoadProfile()) },
            func(a, b GeneralError) bool { return safeDeepEqual(a.LoadReason(), b.LoadReason()) },
            func(a, b GeneralError) bool { return safeDeepEqual(a.LoadStatus(), b.LoadStatus()) },
            func(a, b GeneralError) bool { return safeDeepEqual(a.LoadType(), b.LoadType()) },
        } {
            if !fn(e, generalError) {
                return false
            }
        }
        return true
    }
    return safeDeepEqual(e, target)
}

func (e *uniformException) Unwrap() []error {
    return e.LoadErrors()
}

func (e *uniformException) LoadID() any {
    if e == nil {
        return nil
    }
    return copyAny(e.ID)
}

func (e *uniformException) LoadCode() any {
    if e == nil {
        return nil
    }
    return copyAny(e.Code)
}

func (e *uniformException) LoadData() any {
    if e == nil {
        return nil
    }
    return copyAny(e.Data)
}

func (e *uniformException) LoadDescription() ErrorDescription {
    if e == nil || e.Description == nil {
        return nil
    }
    desc := (&descriptionBuilder{}).loadFromErrorDescription(e.Description).BuildErrorDescription()
    return desc
}

func (e *uniformException) LoadDetail() any {
    if e == nil {
        return nil
    }
    return copyAny(e.Detail)
}

func (e *uniformException) LoadErrors() []error {
    if e == nil || e.Errors == nil {
        return nil
    }
    dst := make([]error, 0, len(e.Errors))
    for _, err := range e.Errors {
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
                    err: err,
                },
            )
            continue
        }
        if copied, ok := copyAny(err).(error); ok && copied != nil {
            dst = append(dst, copied)
        } else {
            dst = append(dst, err)
        }
    }
    if len(dst) == 0 {
        return nil
    }
    return dst
}

func (e *uniformException) LoadInformation() any {
    if e == nil {
        return nil
    }
    return copyAny(e.Information)
}

func (e *uniformException) LoadMessage() string {
    if e == nil {
        return ""
    }
    return e.Message
}

func (e *uniformException) LoadProfile() any {
    if e == nil {
        return nil
    }
    return copyAny(e.Profile)
}

func (e *uniformException) LoadReason() any {
    if e == nil {
        return nil
    }
    return copyAny(e.Reason)
}

func (e *uniformException) LoadStatus() any {
    if e == nil {
        return nil
    }
    return copyAny(e.Status)
}

func (e *uniformException) LoadType() any {
    if e == nil {
        return nil
    }
    return copyAny(e.Type)
}

func (e *uniformException) LoadEquationFunction() func(error, error) bool {
    if e == nil {
        return nil
    }
    return e.equationFunction
}

func (e *uniformException) WithID(id any) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.ID = copyAny(id)
    return ex
}

func (e *uniformException) WithCode(code any) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.Code = copyAny(code)
    return ex
}

func (e *uniformException) WithData(data any) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.Data = copyAny(data)
    return ex
}

func (e *uniformException) WithDescription(description ErrorDescription) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    if description == nil {
        return ex
    }
    ex.Description = (&descriptionBuilder{}).loadFromErrorDescription(description).BuildErrorDescription()
    return ex
}

func (e *uniformException) WithDetail(detail any) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.Detail = copyAny(detail)
    return ex
}

func (e *uniformException) WithErrors(errors []error) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.Errors = normalizeErrors(errors)
    return ex
}

func (e *uniformException) WithInformation(information any) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.Information = copyAny(information)
    return ex
}

func (e *uniformException) WithMessage(message string) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.Message = message
    return ex
}

func (e *uniformException) WithProfile(profile any) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.Profile = copyAny(profile)
    return ex
}

func (e *uniformException) WithReason(reason any) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.Reason = copyAny(reason)
    return ex
}

func (e *uniformException) WithStatus(status any) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.Status = copyAny(status)
    return ex
}

func (e *uniformException) WithType(t any) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.Type = copyAny(t)
    return ex
}

func (e *uniformException) WithEquationFunction(equation func(error, error) bool) GeneralError {
    if e == nil {
        return nil
    }
    ex := e.copy()
    ex.equationFunction = equation
    return ex
}

// NewUniformErrorBuilder 创建统一错误构建器
func NewUniformErrorBuilder() UniformErrorBasicBuilder {
    return &exceptionBuilder{}
}

// UniformErrorBasicBuilder 统一错误基本构建器接口
type UniformErrorBasicBuilder interface {
    WithID(any) UniformErrorBasicBuilder
    WithCode(any) UniformErrorBasicBuilder
    WithData(any) UniformErrorBasicBuilder
    WithDescription(ErrorDescription) UniformErrorBasicBuilder
    WithDetail(any) UniformErrorBasicBuilder
    WithErrors([]error) UniformErrorBasicBuilder
    WithInformation(any) UniformErrorBasicBuilder
    WithMessage(string) UniformErrorBasicBuilder
    WithProfile(any) UniformErrorBasicBuilder
    WithReason(any) UniformErrorBasicBuilder
    WithStatus(any) UniformErrorBasicBuilder
    WithType(any) UniformErrorBasicBuilder
    WithEquationFunction(func(error, error) bool) UniformErrorBasicBuilder
    CreateGeneralError() GeneralError
    CreateProjectError() ProjectError
    CreateRuntimeError() RuntimeError
    CreateServiceError() ServiceError
    CreateUniformError() UniformError
    CreateGeneralException() GeneralException
    CreateProjectException() ProjectException
    CreateRuntimeException() RuntimeException
    CreateServiceException() ServiceException
    CreateUniformException() UniformException
    MakeGeneralError() GeneralError
    MakeProjectError() ProjectError
    MakeRuntimeError() RuntimeError
    MakeServiceError() ServiceError
    MakeUniformError() UniformError
    MakeGeneralException() GeneralException
    MakeProjectException() ProjectException
    MakeRuntimeException() RuntimeException
    MakeServiceException() ServiceException
    MakeUniformException() UniformException
}

type exceptionBuilder struct {
    uniformException uniformException
}

func (builder *exceptionBuilder) WithID(id any) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithID(id)
    }
    ex := builder.uniformException.copy()
    ex.ID = copyAny(id)
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithCode(code any) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithCode(code)
    }
    ex := builder.uniformException.copy()
    ex.Code = copyAny(code)
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithData(data any) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithData(data)
    }
    ex := builder.uniformException.copy()
    ex.Data = copyAny(data)
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithDescription(description ErrorDescription) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithDescription(description)
    }
    ex := builder.uniformException.copy()
    if description != nil {
        ex.Description = (&descriptionBuilder{}).loadFromErrorDescription(description).BuildErrorDescription()
    }
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithDetail(detail any) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithDetail(detail)
    }
    ex := builder.uniformException.copy()
    ex.Detail = copyAny(detail)
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithErrors(errors []error) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithErrors(errors)
    }
    ex := builder.uniformException.copy()
    ex.Errors = normalizeErrors(errors)
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithInformation(information any) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithInformation(information)
    }
    ex := builder.uniformException.copy()
    ex.Information = copyAny(information)
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithMessage(message string) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithMessage(message)
    }
    ex := builder.uniformException.copy()
    ex.Message = message
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithProfile(profile any) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithProfile(profile)
    }
    ex := builder.uniformException.copy()
    ex.Profile = copyAny(profile)
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithReason(reason any) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithReason(reason)
    }
    ex := builder.uniformException.copy()
    ex.Reason = copyAny(reason)
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithStatus(status any) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithStatus(status)
    }
    ex := builder.uniformException.copy()
    ex.Status = copyAny(status)
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithType(t any) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithType(t)
    }
    ex := builder.uniformException.copy()
    ex.Type = copyAny(t)
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) WithEquationFunction(equation func(error, error) bool) UniformErrorBasicBuilder {
    if builder == nil {
        return NewUniformErrorBuilder().WithEquationFunction(equation)
    }
    ex := builder.uniformException.copy()
    ex.equationFunction = equation
    return &exceptionBuilder{uniformException: *ex}
}

func (builder *exceptionBuilder) CreateGeneralError() GeneralError {
    if builder == nil {
        return NewUniformErrorBuilder().CreateGeneralError()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) CreateProjectError() ProjectError {
    if builder == nil {
        return NewUniformErrorBuilder().CreateProjectError()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) CreateRuntimeError() RuntimeError {
    if builder == nil {
        return NewUniformErrorBuilder().CreateRuntimeError()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) CreateServiceError() ServiceError {
    if builder == nil {
        return NewUniformErrorBuilder().CreateServiceError()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) CreateUniformError() UniformError {
    if builder == nil {
        return NewUniformErrorBuilder().CreateUniformError()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) CreateGeneralException() GeneralException {
    if builder == nil {
        return NewUniformErrorBuilder().CreateGeneralException()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) CreateProjectException() ProjectException {
    if builder == nil {
        return NewUniformErrorBuilder().CreateProjectException()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) CreateRuntimeException() RuntimeException {
    if builder == nil {
        return NewUniformErrorBuilder().CreateRuntimeException()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) CreateServiceException() ServiceException {
    if builder == nil {
        return NewUniformErrorBuilder().CreateServiceException()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) CreateUniformException() UniformException {
    if builder == nil {
        return NewUniformErrorBuilder().CreateUniformException()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) MakeGeneralError() GeneralError {
    if builder == nil {
        return NewUniformErrorBuilder().MakeGeneralError()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) MakeProjectError() ProjectError {
    if builder == nil {
        return NewUniformErrorBuilder().MakeProjectError()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) MakeRuntimeError() RuntimeError {
    if builder == nil {
        return NewUniformErrorBuilder().MakeRuntimeError()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) MakeServiceError() ServiceError {
    if builder == nil {
        return NewUniformErrorBuilder().MakeServiceError()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) MakeUniformError() UniformError {
    if builder == nil {
        return NewUniformErrorBuilder().MakeUniformError()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) MakeGeneralException() GeneralException {
    if builder == nil {
        return NewUniformErrorBuilder().MakeGeneralException()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) MakeProjectException() ProjectException {
    if builder == nil {
        return NewUniformErrorBuilder().MakeProjectException()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) MakeRuntimeException() RuntimeException {
    if builder == nil {
        return NewUniformErrorBuilder().MakeRuntimeException()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) MakeServiceException() ServiceException {
    if builder == nil {
        return NewUniformErrorBuilder().MakeServiceException()
    }
    return builder.uniformException.copy()
}

func (builder *exceptionBuilder) MakeUniformException() UniformException {
    if builder == nil {
        return NewUniformErrorBuilder().MakeUniformException()
    }
    return builder.uniformException.copy()
}
