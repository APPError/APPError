package apperror

import (
    "encoding/json"
    "fmt"
    "reflect"
    "strings"
)

// ErrorDescription 错误描述
type ErrorDescription interface {
    // LoadID 返回错误描述的 ID 标识
    LoadID() string
    // LoadCategory 返回错误描述的 Category 类别
    LoadCategory() string
    // LoadCode 返回错误描述的 Code 编码
    LoadCode() string
    // LoadData 返回错误描述的 Data 数据
    LoadData() any
    // LoadDetail 返回错误描述的 Detail 详情
    LoadDetail() any
    // LoadDescription 返回错误描述的 Description 描述
    LoadDescription() ErrorDescription
    // LoadDescriptions 返回错误描述的 Descriptions 描述集合
    LoadDescriptions() []ErrorDescription
    // LoadInformation 返回错误描述的 Information 信息
    LoadInformation() any
    // LoadMessage 返回错误描述的 Message 消息
    LoadMessage() string
    // LoadName 返回错误描述的 Name 名称
    LoadName() string
    // LoadOrganization 返回错误描述的 Organization 所属组织
    LoadOrganization() string
    // LoadProject 返回错误描述所属的 Project 项目
    LoadProject() string
    // LoadService 返回错误描述所属的 Service 服务
    LoadService() string
    // LoadStatus 返回错误描述的 Status 状态
    LoadStatus() any
    // LoadType 返回错误描述的 Type 类型
    LoadType() any
    // LoadEquationFunction 返回错误描述的 equationFunction 等式函数
    LoadEquationFunction() func(error, error) bool
    // WithID 设置错误描述的 ID 标识
    WithID(string) ErrorDescription
    // WithCategory 设置错误描述的 Category 类别
    WithCategory(string) ErrorDescription
    // WithCode 设置错误描述的 Code 编码
    WithCode(string) ErrorDescription
    // WithData 设置错误描述的 Data 数据
    WithData(any) ErrorDescription
    // WithDetail 设置错误描述的 Detail 详情
    WithDetail(any) ErrorDescription
    // WithDescription 设置错误描述的 Description 描述
    WithDescription(ErrorDescription) ErrorDescription
    // WithDescriptions 设置错误描述的 Descriptions 描述集合
    WithDescriptions([]ErrorDescription) ErrorDescription
    // WithInformation 设置错误描述的 Information 信息
    WithInformation(any) ErrorDescription
    // WithMessage 设置错误描述的 Message 消息
    WithMessage(string) ErrorDescription
    // WithName 设置错误描述的 Name 名称
    WithName(string) ErrorDescription
    // WithOrganization 设置错误描述的 Organization 所属组织
    WithOrganization(string) ErrorDescription
    // WithProject 设置错误描述的 Project 项目
    WithProject(string) ErrorDescription
    // WithService 设置错误描述的 Service 服务
    WithService(string) ErrorDescription
    // WithStatus 设置错误描述的 Status 状态
    WithStatus(any) ErrorDescription
    // WithType 设置错误描述的 Type 类型
    WithType(any) ErrorDescription
    // WithEquationFunction 设置错误描述的 equationFunction 等式函数
    WithEquationFunction(func(error, error) bool) ErrorDescription
}

// SimpleErrorDescription 错误描述。
type SimpleErrorDescription struct {
    ID               string                  `json:"id"`           // ID 标识
    Category         string                  `json:"category"`     // Category 类别
    Code             string                  `json:"code"`         // Code 编码
    Data             any                     `json:"data"`         // Data 数据
    Description      ErrorDescription        `json:"description"`  // Description 描述
    Descriptions     []ErrorDescription      `json:"descriptions"` // Descriptions 描述集合
    Detail           any                     `json:"detail"`       // Detail 详情
    Information      any                     `json:"information"`  // Information 信息
    Message          string                  `json:"message"`      // Message 消息
    Name             string                  `json:"name"`         // Name 名称
    Organization     string                  `json:"organization"` // Organization 所属组织
    Project          string                  `json:"project"`      // Project 所属项目
    Service          string                  `json:"service"`      // Service 所属服务
    Status           any                     `json:"status"`       // Status 状态
    Type             any                     `json:"type"`         // Type 类型
    equationFunction func(error, error) bool // equationFunction 等式函数
}

// clone 深拷贝 SimpleErrorDescription，返回副本（保留所有字段，包括未导出的 equationFunction）。
//
// 实现 interface{ clone() any }，使 copyAny 在处理嵌套的 SimpleErrorDescription 时
// 优先调用此方法而非反射深拷贝。反射路径会跳过未导出字段（如 equationFunction），
// 导致嵌套描述经 LoadDescription/LoadDescriptions/WithDescription/WithDescriptions/copy
// 等操作后丢失自定义比较函数，Is() 回退到全字段默认比较，造成错误匹配行为不一致。
// 委派给 copy() 可正确保留 equationFunction，因为 copy() 通过 *desc = *description
// 浅拷贝在顶层保留了该未导出字段，且对嵌套的 Description/Descriptions 同样经 copyAny
// 递归走 clone() 路径。
//
// 本方法未导出，仅作为包内私有约定，避免将内部适配方法暴露为公开 API。
// 外部类型若需被 copyAny 识别，仍应实现导出的 Clone() any 或 Copy() any。
func (description *SimpleErrorDescription) clone() any {
    return description.copy()
}

// copy 深拷贝 SimpleErrorDescription，确保可变引用类型字段不与原对象共享底层数据。
// 通过反射统一处理以下类型的导出字段，结构体新增同类型导出字段时自动覆盖：
//   - any（interface）类型：Data、Detail、Information、Status、Type
//   - 切片类型：如 Descriptions（[]ErrorDescription）
//
// 注意：未导出的 any/pointer/slice/map 字段无法通过反射安全读写，将被跳过。
// 若结构体包含需要深拷贝的未导出字段，应实现 Clone() any 或 Copy() any 接口。
func (description *SimpleErrorDescription) copy() *SimpleErrorDescription {
    if description == nil {
        return nil
    }
    desc := new(SimpleErrorDescription)
    *desc = *description
    {
        srcVal := reflect.ValueOf(description).Elem()
        dstVal := reflect.ValueOf(desc).Elem()
        for i := 0; i < srcVal.NumField(); i++ {
            ft := srcVal.Type().Field(i)
            fv := srcVal.Field(i)
            switch ft.Type.Kind() {
            case reflect.Interface:
                // any 类型字段：导出字段使用 copyAny 深拷贝
                if !ft.IsExported() {
                    continue
                }
                if !fv.IsNil() {
                    // 当字段是 typed-nil（如 (*T)(nil)）时，copyAny 会返回 untyped nil，
                    // 而 reflect.ValueOf(nil) 得到的是无效 Value，直接 Set 会 panic。
                    // 因此先用 copied != nil 过滤掉这种情况：
                    //   - 非 nil：正常深拷贝并 Set；
                    //   - nil：跳过 Set，保留上面 *desc = *description 已复制的原始 typed-nil，
                    //          这样类型信息不丢失，后续 Is/Unwrap 等判断依然正确。
                    if copied := copyAny(fv.Interface()); copied != nil {
                        dstVal.Field(i).Set(reflect.ValueOf(copied))
                    }
                }
            case reflect.Slice:
                // 切片类型字段：导出切片执行逐元素深拷贝，避免共享底层数组
                if !ft.IsExported() {
                    continue
                }
                if fv.IsNil() {
                    continue
                }
                newSlice := reflect.MakeSlice(ft.Type, fv.Len(), fv.Len())
                for j := 0; j < fv.Len(); j++ {
                    if copied := copyAny(fv.Index(j).Interface()); copied != nil {
                        newSlice.Index(j).Set(reflect.ValueOf(copied))
                    }
                }
                dstVal.Field(i).Set(newSlice)
            case reflect.Pointer:
                if fv.IsNil() {
                    continue
                }
                // 其他指针类型对导出字段使用 copyAny 深拷贝
                if ft.IsExported() {
                    if copied := copyAny(fv.Interface()); copied != nil {
                        dstVal.Field(i).Set(reflect.ValueOf(copied))
                    }
                }
            default:
                // 值类型（Bool、Int、String、Struct、Array、Func 等）已由 *desc = *description 正确复制，无需处理
                // 注：函数值不可变，浅拷贝安全
            }
        }
    }
    return desc
}

// Error 返回错误描述的字符串表示
func (description *SimpleErrorDescription) Error() string {
    if description == nil {
        return ""
    }
    if bytes, err := json.Marshal(description); err == nil {
        return string(bytes)
    }
    // json.Marshal 失败时，逐字段序列化以产生合法 JSON
    jsonVal := func(v any) string {
        if v == nil {
            return "null"
        }
        if b, err := json.Marshal(v); err == nil {
            return string(b)
        }
        return fmt.Sprintf("%q", fmt.Sprintf("%v", v))
    }
    jsonDescs := func(descs []ErrorDescription) string {
        if descs == nil {
            return "null"
        }
        var s strings.Builder
        s.WriteString("[")
        for i, d := range descs {
            if i > 0 {
                s.WriteString(",")
            }
            if d == nil {
                s.WriteString("null")
            } else if b, err := json.Marshal(d); err == nil {
                s.WriteString(string(b))
            } else {
                s.WriteString(fmt.Sprintf("%q", fmt.Sprintf("%v", d)))
            }
        }
        s.WriteString("]")
        return s.String()
    }
    return fmt.Sprintf(
        `{
        "id":%q,
        "category":%q,
        "code":%q,
        "data":%s,
        "detail":%s,
        "description":%s,
        "descriptions":%s,
        "information":%s,
        "message":%q,
        "name":%q,
        "organization":%q,
        "project":%q,
        "service":%q,
        "status":%s,
        "type":%s
        }`,
        description.ID,
        description.Category,
        description.Code,
        jsonVal(description.Data),
        jsonVal(description.Detail),
        jsonVal(description.Description),
        jsonDescs(description.Descriptions),
        jsonVal(description.Information),
        description.Message,
        description.Name,
        description.Organization,
        description.Project,
        description.Service,
        jsonVal(description.Status),
        jsonVal(description.Type),
    )
}

// String 返回错误描述的字符串表示
func (description *SimpleErrorDescription) String() string {
    return description.Error()
}

func (description *SimpleErrorDescription) Is(err error) (valid bool) {
    defer func() {
        if recover() != nil {
            valid = false
        }
    }()
    if description == nil || err == nil {
        return description == err
    }
    if description.equationFunction != nil {
        return description.equationFunction(description, err)
    }
    var desc ErrorDescription
    if As(err, &desc) {
        for _, get := range []func(ErrorDescription) any{
            func(d ErrorDescription) any { return d.LoadID() },
            func(d ErrorDescription) any { return d.LoadCategory() },
            func(d ErrorDescription) any { return d.LoadCode() },
            func(d ErrorDescription) any { return d.LoadData() },
            func(d ErrorDescription) any { return d.LoadDetail() },
            func(d ErrorDescription) any { return d.LoadDescription() },
            func(d ErrorDescription) any { return d.LoadDescriptions() },
            func(d ErrorDescription) any { return d.LoadInformation() },
            func(d ErrorDescription) any { return d.LoadName() },
            func(d ErrorDescription) any { return d.LoadOrganization() },
            func(d ErrorDescription) any { return d.LoadProject() },
            func(d ErrorDescription) any { return d.LoadService() },
            func(d ErrorDescription) any { return d.LoadStatus() },
            func(d ErrorDescription) any { return d.LoadType() },
        } {
            if !safeDeepEqual(get(description), get(desc)) {
                valid = false
                return false
            }
        }
        return true
    }
    return safeDeepEqual(description, err)
}

func (description *SimpleErrorDescription) Unwrap() []error {
    if description == nil {
        return nil
    }
    if len(description.Descriptions) > 0 {
        descList := make([]error, 0, len(description.Descriptions))
        for _, desc := range description.Descriptions {
            if err, ok := desc.(error); ok && err != nil {
                descList = append(descList, err)
            } else {
                descList = append(
                    descList, &SimpleErrorDescription{
                        ID:               desc.LoadID(),
                        Category:         desc.LoadCategory(),
                        Code:             desc.LoadCode(),
                        Data:             desc.LoadData(),
                        Description:      desc.LoadDescription(),
                        Descriptions:     desc.LoadDescriptions(),
                        Detail:           desc.LoadDetail(),
                        Information:      desc.LoadInformation(),
                        Message:          desc.LoadMessage(),
                        Name:             desc.LoadName(),
                        Organization:     desc.LoadOrganization(),
                        Project:          desc.LoadProject(),
                        Service:          desc.LoadService(),
                        Status:           desc.LoadStatus(),
                        Type:             desc.LoadType(),
                        equationFunction: desc.LoadEquationFunction(),
                    },
                )
            }
        }
        if len(descList) > 0 {
            result := make([]error, len(descList))
            copy(result, descList)
            return result
        }
    }
    if description.Description != nil {
        if err, ok := description.Description.(error); ok && err != nil {
            return []error{err}
        }
        return []error{
            &SimpleErrorDescription{
                ID:               description.Description.LoadID(),
                Category:         description.Description.LoadCategory(),
                Code:             description.Description.LoadCode(),
                Data:             description.Description.LoadData(),
                Description:      description.Description.LoadDescription(),
                Descriptions:     description.Description.LoadDescriptions(),
                Detail:           description.Description.LoadDetail(),
                Information:      description.Description.LoadInformation(),
                Message:          description.Description.LoadMessage(),
                Name:             description.Description.LoadName(),
                Organization:     description.Description.LoadOrganization(),
                Project:          description.Description.LoadProject(),
                Service:          description.Description.LoadService(),
                Status:           description.Description.LoadStatus(),
                Type:             description.Description.LoadType(),
                equationFunction: description.Description.LoadEquationFunction(),
            },
        }
    }
    return nil
}

// LoadID 返回错误描述的 ID 标识
func (description *SimpleErrorDescription) LoadID() string {
    if description == nil {
        return ""
    }
    return description.ID
}

// LoadCategory 返回错误描述的 Category 类别
func (description *SimpleErrorDescription) LoadCategory() string {
    if description == nil {
        return ""
    }
    return description.Category
}

// LoadCode 返回错误描述的 Code 编码
func (description *SimpleErrorDescription) LoadCode() string {
    if description == nil {
        return ""
    }
    return description.Code
}

// LoadData 返回错误描述的 Data 数据
func (description *SimpleErrorDescription) LoadData() any {
    if description == nil {
        return nil
    }
    return copyAny(description.Data)
}

// LoadDetail 返回错误描述的 Detail 详情
func (description *SimpleErrorDescription) LoadDetail() any {
    if description == nil {
        return nil
    }
    return copyAny(description.Detail)
}

// LoadDescription 返回错误描述的 Description 描述
func (description *SimpleErrorDescription) LoadDescription() ErrorDescription {
    if description == nil || description.Description == nil {
        return nil
    }
    if desc, ok := copyAny(description.Description).(ErrorDescription); ok {
        return desc
    }
    return nil
}

// LoadDescriptions 返回错误描述的 Descriptions 描述集合
func (description *SimpleErrorDescription) LoadDescriptions() []ErrorDescription {
    if description == nil || description.Descriptions == nil {
        return nil
    }
    descriptions := make([]ErrorDescription, 0, len(description.Descriptions))
    for _, desc := range description.Descriptions {
        if d, ok := copyAny(desc).(ErrorDescription); ok && d != nil {
            descriptions = append(descriptions, d)
        }
    }
    result := make([]ErrorDescription, len(descriptions))
    copy(result, descriptions)
    return result
}

// LoadInformation 返回错误描述的 Information 信息
func (description *SimpleErrorDescription) LoadInformation() any {
    if description == nil {
        return nil
    }
    return copyAny(description.Information)
}

// LoadMessage 返回错误描述的 Message 消息
func (description *SimpleErrorDescription) LoadMessage() string {
    if description == nil {
        return ""
    }
    return description.Message
}

// LoadName 返回错误描述的 Name 名称
func (description *SimpleErrorDescription) LoadName() string {
    if description == nil {
        return ""
    }
    return description.Name
}

// LoadOrganization 返回错误描述的 Organization 所属组织
func (description *SimpleErrorDescription) LoadOrganization() string {
    if description == nil {
        return ""
    }
    return description.Organization
}

// LoadProject 返回错误描述所属的 Project 项目
func (description *SimpleErrorDescription) LoadProject() string {
    if description == nil {
        return ""
    }
    return description.Project
}

// LoadService 返回错误描述所属的 Service 服务
func (description *SimpleErrorDescription) LoadService() string {
    if description == nil {
        return ""
    }
    return description.Service
}

// LoadStatus 返回错误描述的 Status 状态
func (description *SimpleErrorDescription) LoadStatus() any {
    if description == nil {
        return nil
    }
    return copyAny(description.Status)
}

// LoadType 返回错误描述的 Type 类型
func (description *SimpleErrorDescription) LoadType() any {
    if description == nil {
        return nil
    }
    return copyAny(description.Type)
}

// LoadEquationFunction 返回错误描述的 equationFunction 等式函数
func (description *SimpleErrorDescription) LoadEquationFunction() func(error, error) bool {
    if description == nil || description.equationFunction == nil {
        return nil
    }
    return description.equationFunction
}

// WithID 设置错误描述的 ID 标识
func (description *SimpleErrorDescription) WithID(id string) ErrorDescription {
    duplicate := description.copy()
    if duplicate == nil {
        duplicate = new(SimpleErrorDescription)
    }
    duplicate.ID = id
    return duplicate
}

// WithCategory 设置错误描述的 Category 类别
func (description *SimpleErrorDescription) WithCategory(category string) ErrorDescription {
    duplicate := description.copy()
    if duplicate == nil {
        duplicate = new(SimpleErrorDescription)
    }
    duplicate.Category = category
    return duplicate
}

// WithCode 设置错误描述的 Code 编码
func (description *SimpleErrorDescription) WithCode(code string) ErrorDescription {
    duplicate := description.copy()
    if duplicate == nil {
        duplicate = new(SimpleErrorDescription)
    }
    duplicate.Code = code
    return duplicate
}

// WithData 设置错误描述的 Data 数据。
func (description *SimpleErrorDescription) WithData(data any) ErrorDescription {
    duplicate := description.copy()
    if duplicate == nil {
        duplicate = new(SimpleErrorDescription)
    }
    duplicate.Data = copyAny(data)
    return duplicate
}

// WithDescription 设置错误描述的 Description 描述
func (description *SimpleErrorDescription) WithDescription(desc ErrorDescription) ErrorDescription {
    duplicate := description.copy()
    if duplicate == nil {
        duplicate = new(SimpleErrorDescription)
    }
    if dup, ok := copyAny(desc).(ErrorDescription); ok && dup != nil {
        duplicate.Description = dup
    }
    return duplicate
}

// WithDescriptions 设置错误描述的 Descriptions 描述集合
func (description *SimpleErrorDescription) WithDescriptions(descriptions []ErrorDescription) ErrorDescription {
    duplicate := description.copy()
    if duplicate == nil {
        duplicate = new(SimpleErrorDescription)
    }
    if descriptions == nil {
        duplicate.Descriptions = nil
        return duplicate
    }
    descList := make([]ErrorDescription, 0, len(descriptions))
    for _, original := range descriptions {
        if dup, ok := copyAny(original).(ErrorDescription); ok && dup != nil {
            descList = append(descList, dup)
        }
    }
    list := make([]ErrorDescription, len(descList))
    copy(list, descList)
    duplicate.Descriptions = list
    return duplicate
}

// WithDetail 设置错误描述的 Detail 详情。
func (description *SimpleErrorDescription) WithDetail(detail any) ErrorDescription {
    desc := description.copy()
    if desc == nil {
        desc = new(SimpleErrorDescription)
    }
    desc.Detail = copyAny(detail)
    return desc
}

// WithInformation 设置错误描述的 Information 信息。
func (description *SimpleErrorDescription) WithInformation(information any) ErrorDescription {
    desc := description.copy()
    if desc == nil {
        desc = new(SimpleErrorDescription)
    }
    desc.Information = copyAny(information)
    return desc
}

// WithMessage 设置错误描述的 Message 消息
func (description *SimpleErrorDescription) WithMessage(message string) ErrorDescription {
    desc := description.copy()
    if desc == nil {
        desc = new(SimpleErrorDescription)
    }
    desc.Message = message
    return desc
}

// WithName 设置错误描述的 Name 名称
func (description *SimpleErrorDescription) WithName(name string) ErrorDescription {
    desc := description.copy()
    if desc == nil {
        desc = new(SimpleErrorDescription)
    }
    desc.Name = name
    return desc
}

// WithOrganization 设置错误描述的 Organization 所属组织
func (description *SimpleErrorDescription) WithOrganization(organization string) ErrorDescription {
    desc := description.copy()
    if desc == nil {
        desc = new(SimpleErrorDescription)
    }
    desc.Organization = organization
    return desc
}

// WithProject 设置错误描述的 Project 项目
func (description *SimpleErrorDescription) WithProject(project string) ErrorDescription {
    desc := description.copy()
    if desc == nil {
        desc = new(SimpleErrorDescription)
    }
    desc.Project = project
    return desc
}

// WithService 设置错误描述的 Service 服务
func (description *SimpleErrorDescription) WithService(service string) ErrorDescription {
    desc := description.copy()
    if desc == nil {
        desc = new(SimpleErrorDescription)
    }
    desc.Service = service
    return desc
}

// WithStatus 设置错误描述的 Status 状态。
func (description *SimpleErrorDescription) WithStatus(status any) ErrorDescription {
    desc := description.copy()
    if desc == nil {
        desc = new(SimpleErrorDescription)
    }
    desc.Status = copyAny(status)
    return desc
}

// WithType 设置错误描述的 Type 类型。
func (description *SimpleErrorDescription) WithType(t any) ErrorDescription {
    desc := description.copy()
    if desc == nil {
        desc = new(SimpleErrorDescription)
    }
    desc.Type = copyAny(t)
    return desc
}

// WithEquationFunction 设置错误描述的 equationFunction 等式函数
func (description *SimpleErrorDescription) WithEquationFunction(equationFunction func(error, error) bool) ErrorDescription {
    desc := description.copy()
    if desc == nil {
        desc = new(SimpleErrorDescription)
    }
    desc.equationFunction = equationFunction
    return desc
}

// ErrorDescriptionBasicBuilder 错误描述基本构建器
type ErrorDescriptionBasicBuilder interface {
    // WithID 设置错误描述的 ID 标识
    WithID(string) ErrorDescriptionBasicBuilder
    // WithCategory 设置错误描述的 Category 类别
    WithCategory(string) ErrorDescriptionBasicBuilder
    // WithCode 设置错误描述的 Code 编码
    WithCode(string) ErrorDescriptionBasicBuilder
    // WithData 设置错误描述的 Data 数据
    WithData(any) ErrorDescriptionBasicBuilder
    // WithDetail 设置错误描述的 Detail 详情
    WithDetail(any) ErrorDescriptionBasicBuilder
    // WithDescription 设置错误描述的 Description 描述
    WithDescription(ErrorDescription) ErrorDescriptionBasicBuilder
    // WithDescriptions 设置错误描述的 Descriptions 描述集合
    WithDescriptions([]ErrorDescription) ErrorDescriptionBasicBuilder
    // WithInformation 设置错误描述的 Information 信息
    WithInformation(any) ErrorDescriptionBasicBuilder
    // WithMessage 设置错误描述的 Message 消息
    WithMessage(string) ErrorDescriptionBasicBuilder
    // WithName 设置错误描述的 Name 名称
    WithName(string) ErrorDescriptionBasicBuilder
    // WithOrganization 设置错误描述的 Organization 所属组织
    WithOrganization(string) ErrorDescriptionBasicBuilder
    // WithProject 设置错误描述的 Project 项目
    WithProject(string) ErrorDescriptionBasicBuilder
    // WithService 设置错误描述的 Service 服务
    WithService(string) ErrorDescriptionBasicBuilder
    // WithStatus 设置错误描述的 Status 状态
    WithStatus(any) ErrorDescriptionBasicBuilder
    // WithType 设置错误描述的 Type 类型
    WithType(any) ErrorDescriptionBasicBuilder
    // WithEquationFunction 设置错误描述的 equationFunction 等式函数
    WithEquationFunction(func(error, error) bool) ErrorDescriptionBasicBuilder
    // BuildErrorDescription 构建并返回错误描述
    BuildErrorDescription() ErrorDescription
}

// NewErrorDescriptionBuilder 创建一个新的错误描述构建器
func NewErrorDescriptionBuilder() ErrorDescriptionBasicBuilder {
    return &descriptionBuilder{}
}

type descriptionBuilder struct {
    description SimpleErrorDescription
}

// 编译期确保 descriptionBuilder 实现了 ErrorDescriptionBasicBuilder 接口
var _ ErrorDescriptionBasicBuilder = (*descriptionBuilder)(nil)

// loadFromErrorDescription 以当前 builder 的状态为基础，从已有的 ErrorDescription 加载所有字段并构建一个新的
// ErrorDescriptionBasicBuilder，常用于基于已有描述创建变体或派生描述。
//
// 处理逻辑：
//  1. 若 builder 非 nil，则深拷贝其内部的 SimpleErrorDescription 作为初始值；
//     若 builder 为 nil，则使用零值 SimpleErrorDescription。
//  2. 以初始值构造一个新的 descriptionBuilder。
//  3. 若参数 description 非 nil，则依次调用 With* 方法，用 description 的每一个字段覆盖初始值中的对应字段。
//  4. 返回最终的 ErrorDescriptionBasicBuilder 实例。
//
// 该方法本身是 nil-safe 的：builder 为 nil 或 description 为 nil 均可正常处理，不会 panic。
func (builder *descriptionBuilder) loadFromErrorDescription(description ErrorDescription) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    var b ErrorDescriptionBasicBuilder = &descriptionBuilder{desc}
    if description != nil {
        b = b.WithID(description.LoadID())
        b = b.WithCategory(description.LoadCategory())
        b = b.WithCode(description.LoadCode())
        b = b.WithData(description.LoadData())
        b = b.WithDetail(description.LoadDetail())
        b = b.WithDescription(description.LoadDescription())
        b = b.WithDescriptions(description.LoadDescriptions())
        b = b.WithInformation(description.LoadInformation())
        b = b.WithMessage(description.LoadMessage())
        b = b.WithName(description.LoadName())
        b = b.WithOrganization(description.LoadOrganization())
        b = b.WithProject(description.LoadProject())
        b = b.WithService(description.LoadService())
        b = b.WithStatus(description.LoadStatus())
        b = b.WithType(description.LoadType())
        b = b.WithEquationFunction(description.LoadEquationFunction())
    }
    return b
}

func (builder *descriptionBuilder) WithID(id string) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.ID = id
    return &descriptionBuilder{desc}
}

func (builder *descriptionBuilder) WithCategory(category string) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Category = category
    return &descriptionBuilder{desc}
}

func (builder *descriptionBuilder) WithCode(code string) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Code = code
    return &descriptionBuilder{desc}
}

// WithData 设置错误描述的 Data 数据。
func (builder *descriptionBuilder) WithData(data any) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Data = copyAny(data)
    return &descriptionBuilder{desc}
}

// WithDetail 设置错误描述的 Detail 详情。
func (builder *descriptionBuilder) WithDetail(detail any) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Detail = copyAny(detail)
    return &descriptionBuilder{desc}
}

func (builder *descriptionBuilder) WithDescription(description ErrorDescription) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    if d, ok := copyAny(description).(ErrorDescription); ok && d != nil {
        desc.Description = d
    }
    return &descriptionBuilder{desc}
}

func (builder *descriptionBuilder) WithDescriptions(descriptions []ErrorDescription) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    if descriptions == nil {
        desc.Descriptions = nil
        return &descriptionBuilder{desc}
    }
    list := make([]ErrorDescription, 0, len(descriptions))
    for _, original := range descriptions {
        if d, ok := copyAny(original).(ErrorDescription); ok && d != nil {
            list = append(list, d)
        }
    }
    list2 := make([]ErrorDescription, len(list))
    copy(list2, list)
    desc.Descriptions = list2
    return &descriptionBuilder{desc}
}

// WithInformation 设置错误描述的 Information 信息。
func (builder *descriptionBuilder) WithInformation(information any) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Information = copyAny(information)
    return &descriptionBuilder{desc}
}

func (builder *descriptionBuilder) WithMessage(message string) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Message = message
    return &descriptionBuilder{desc}
}

func (builder *descriptionBuilder) WithName(name string) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Name = name
    return &descriptionBuilder{desc}
}

func (builder *descriptionBuilder) WithOrganization(organization string) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Organization = organization
    return &descriptionBuilder{desc}
}

func (builder *descriptionBuilder) WithProject(project string) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Project = project
    return &descriptionBuilder{desc}
}

func (builder *descriptionBuilder) WithService(service string) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Service = service
    return &descriptionBuilder{desc}
}

// WithStatus 设置错误描述的 Status 状态。
func (builder *descriptionBuilder) WithStatus(status any) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Status = copyAny(status)
    return &descriptionBuilder{desc}
}

// WithType 设置错误描述的 Type 类型。
func (builder *descriptionBuilder) WithType(t any) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.Type = copyAny(t)
    return &descriptionBuilder{desc}
}

func (builder *descriptionBuilder) WithEquationFunction(equationFunction func(error, error) bool) ErrorDescriptionBasicBuilder {
    var desc SimpleErrorDescription
    if builder != nil {
        if d := builder.description.copy(); d != nil {
            desc = *d
        }
    }
    desc.equationFunction = equationFunction
    return &descriptionBuilder{desc}
}

func (builder *descriptionBuilder) BuildErrorDescription() ErrorDescription {
    if builder == nil {
        return nil
    }
    return builder.description.copy()
}
