package apperror

import (
    "encoding/json"
    "fmt"
)

// ErrorDescription 错误描述
type ErrorDescription struct {
    ID           string `json:"id"`           // ID 标识
    Category     string `json:"category"`     // Category 类别
    Code         string `json:"code"`         // Code 编码
    Data         any    `json:"data"`         // Data 数据
    Detail       any    `json:"detail"`       // Detail 详情
    Identifier   string `json:"identifier"`   // Identifier 标识
    Information  any    `json:"information"`  // Information 信息
    Organization string `json:"organization"` // Organization 所属组织
    Project      string `json:"project"`      // Project 所属项目
    Service      string `json:"service"`      // Service 所属服务
    Status       any    `json:"status"`       // Status 状态
    Type         any    `json:"type"`         // Type 类型
}

func (description *ErrorDescription) Error() string {
    if description == nil {
        return ""
    }
    if bytes, err := json.Marshal(description); err == nil {
        return string(bytes)
    }
    return fmt.Sprintf(
        `ErrorDescription{
        ID:%s
        Category:%s
        Code:%s
        Data:%v
        Detail:%v
        Identifier:%s
        Information:%v
        Organization:%s
        Project:%s
        Service:%s
        Status:%v
        Type:%v
    }`,
        description.ID, description.Category, description.Code, description.Data, description.Detail,
        description.Identifier, description.Information, description.Organization, description.Project,
        description.Service, description.Status, description.Type,
    )
}

func (description *ErrorDescription) String() string {
    return description.Error()
}

// LoadID 返回错误描述的 ID 标识。
func (description *ErrorDescription) LoadID() string {
    if description != nil {
        return description.ID
    }
    return ""
}

// LoadCategory 返回错误描述的类别。
func (description *ErrorDescription) LoadCategory() string {
    if description != nil {
        return description.Category
    }
    return ""
}

// LoadCode 返回错误描述的编码。
func (description *ErrorDescription) LoadCode() string {
    if description != nil {
        return description.Code
    }
    return ""
}

// LoadData 返回错误描述的自定义数据。
func (description *ErrorDescription) LoadData() any {
    if description != nil {
        return description.Data
    }
    return nil
}

// LoadDetail 返回错误描述的详情。
func (description *ErrorDescription) LoadDetail() any {
    if description != nil {
        return description.Detail
    }
    return nil
}

// LoadInformation 返回错误描述的附加信息。
func (description *ErrorDescription) LoadInformation() any {
    if description != nil {
        return description.Information
    }
    return nil
}

// LoadIdentifier 返回错误描述的标识符。
func (description *ErrorDescription) LoadIdentifier() string {
    if description != nil {
        return description.Identifier
    }
    return ""
}

// LoadOrganization 返回错误描述所属的组织。
func (description *ErrorDescription) LoadOrganization() string {
    if description != nil {
        return description.Organization
    }
    return ""
}

// LoadProject 返回错误描述所属的项目。
func (description *ErrorDescription) LoadProject() string {
    if description != nil {
        return description.Project
    }
    return ""
}

// LoadService 返回错误描述所属的服务。
func (description *ErrorDescription) LoadService() string {
    if description != nil {
        return description.Service
    }
    return ""
}

// LoadStatus 返回错误描述的状态。
func (description *ErrorDescription) LoadStatus() any {
    if description != nil {
        return description.Status
    }
    return nil
}

// LoadType 返回错误描述的类型。
func (description *ErrorDescription) LoadType() any {
    if description != nil {
        return description.Type
    }
    return nil
}

// WithID 设置错误描述的 ID 标识，返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithID(id string) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.ID = id
    return desc
}

// WithCategory 设置错误描述的类别，返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithCategory(category string) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.Category = category
    return desc
}

// WithCode 设置错误描述的编码，返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithCode(code string) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.Code = code
    return desc
}

// WithData 设置错误描述的自定义数据（内部深拷贝），返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithData(data any) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.Data = copyAny(data)
    return desc
}

// WithDetail 设置错误描述的详情（内部深拷贝），返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithDetail(detail any) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.Detail = copyAny(detail)
    return desc
}

// WithIdentifier 设置错误描述的标识符，返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithIdentifier(identifier string) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.Identifier = identifier
    return desc
}

// WithInformation 设置错误描述的附加信息（内部深拷贝），返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithInformation(information any) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.Information = copyAny(information)
    return desc
}

// WithOrganization 设置错误描述所属的组织，返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithOrganization(organization string) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.Organization = organization
    return desc
}

// WithProject 设置错误描述所属的项目，返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithProject(project string) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.Project = project
    return desc
}

// WithService 设置错误描述所属的服务，返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithService(service string) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.Service = service
    return desc
}

// WithStatus 设置错误描述的状态（内部深拷贝），返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithStatus(status any) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.Status = copyAny(status)
    return desc
}

// WithType 设置错误描述的类型（内部深拷贝），返回新的 ErrorDescription（不可变风格）。
func (description *ErrorDescription) WithType(typeValue any) ErrorDescription {
    desc := ErrorDescription{}
    if description != nil {
        desc = *description
    }
    desc.Type = copyAny(typeValue)
    return desc
}

type ErrorDescriptionBasicBuilder interface {
    ID(string) ErrorDescriptionBasicBuilder
    Category(string) ErrorDescriptionBasicBuilder
    Code(string) ErrorDescriptionBasicBuilder
    Data(any) ErrorDescriptionBasicBuilder
    Detail(any) ErrorDescriptionBasicBuilder
    Identifier(string) ErrorDescriptionBasicBuilder
    Information(any) ErrorDescriptionBasicBuilder
    Organization(string) ErrorDescriptionBasicBuilder
    Project(string) ErrorDescriptionBasicBuilder
    Service(string) ErrorDescriptionBasicBuilder
    Status(any) ErrorDescriptionBasicBuilder
    Type(any) ErrorDescriptionBasicBuilder
    BuildErrorDescription() ErrorDescription
}

func NewErrorDescriptionBuilder() ErrorDescriptionBasicBuilder {
    return &descriptionBuilder{}
}

type descriptionBuilder struct {
    description ErrorDescription
}

func (b *descriptionBuilder) ID(id string) ErrorDescriptionBasicBuilder {
    b.description.ID = id
    return b
}

func (b *descriptionBuilder) Category(category string) ErrorDescriptionBasicBuilder {
    b.description.Category = category
    return b
}

func (b *descriptionBuilder) Code(code string) ErrorDescriptionBasicBuilder {
    b.description.Code = code
    return b
}

func (b *descriptionBuilder) Data(data any) ErrorDescriptionBasicBuilder {
    b.description.Data = copyAny(data)
    return b
}

func (b *descriptionBuilder) Detail(detail any) ErrorDescriptionBasicBuilder {
    b.description.Detail = copyAny(detail)
    return b
}

func (b *descriptionBuilder) Identifier(identifier string) ErrorDescriptionBasicBuilder {
    b.description.Identifier = identifier
    return b
}

func (b *descriptionBuilder) Information(information any) ErrorDescriptionBasicBuilder {
    b.description.Information = copyAny(information)
    return b
}

func (b *descriptionBuilder) Organization(organization string) ErrorDescriptionBasicBuilder {
    b.description.Organization = organization
    return b
}

func (b *descriptionBuilder) Project(project string) ErrorDescriptionBasicBuilder {
    b.description.Project = project
    return b
}

func (b *descriptionBuilder) Service(service string) ErrorDescriptionBasicBuilder {
    b.description.Service = service
    return b
}

func (b *descriptionBuilder) Status(status any) ErrorDescriptionBasicBuilder {
    b.description.Status = copyAny(status)
    return b
}

func (b *descriptionBuilder) Type(typeValue any) ErrorDescriptionBasicBuilder {
    b.description.Type = copyAny(typeValue)
    return b
}

func (b *descriptionBuilder) BuildErrorDescription() ErrorDescription {
    desc := b.description
    // 深拷贝引用类型字段，防止返回后通过 builder 继续修改或外部持有者间接篡改
    desc.Data = copyAny(desc.Data)
    desc.Detail = copyAny(desc.Detail)
    desc.Information = copyAny(desc.Information)
    desc.Status = copyAny(desc.Status)
    desc.Type = copyAny(desc.Type)
    return desc
}
