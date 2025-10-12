---
title: 命名映射
description: 学习 Pie 的命名映射功能，自定义字段和表名转换
---

# 命名映射

Pie 提供了灵活的命名映射功能，允许您自定义结构体字段名和数据库字段名之间的转换规则。

## 基本用法

### 内置映射器

```go
// 蛇形命名映射
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(&pie.SnakeMapper{}),
)

// 驼峰命名映射
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(&pie.CamelMapper{}),
)

// 相同命名映射（不转换）
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(&pie.SameMapper{}),
)
```

### 映射效果

```go
type User struct {
    ID        bson.ObjectID `bson:"_id,omitempty"`
    FirstName string             `bson:"first_name"`
    LastName  string             `bson:"last_name"`
    Email     string             `bson:"email"`
    CreatedAt time.Time          `bson:"created_at"`
}

// 使用 SnakeMapper
// 结构体字段 -> 数据库字段
// FirstName -> first_name
// LastName -> last_name
// CreatedAt -> created_at

// 使用 CamelMapper
// 结构体字段 -> 数据库字段
// first_name -> FirstName
// last_name -> LastName
// created_at -> CreatedAt

// 使用 SameMapper
// 结构体字段 -> 数据库字段
// FirstName -> FirstName
// LastName -> LastName
// CreatedAt -> CreatedAt
```

## 自定义映射器

### 实现映射器接口

```go
// 映射器接口
type Mapper interface {
    TableName(structName string) string
    FieldName(fieldName string) string
}

// 自定义映射器
type CustomMapper struct{}

func (m CustomMapper) TableName(structName string) string {
    return "t_" + strings.ToLower(structName)
}

func (m CustomMapper) FieldName(fieldName string) string {
    return strings.ToLower(fieldName)
}

// 使用自定义映射器
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(CustomMapper{}),
)
```

### 前缀映射器

```go
type PrefixMapper struct {
    TablePrefix string
    FieldPrefix string
}

func (m PrefixMapper) TableName(structName string) string {
    return m.TablePrefix + strings.ToLower(structName)
}

func (m PrefixMapper) FieldName(fieldName string) string {
    return m.FieldPrefix + strings.ToLower(fieldName)
}

// 使用前缀映射器
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(PrefixMapper{
        TablePrefix: "app_",
        FieldPrefix: "f_",
    }),
)
```

### 条件映射器

```go
type ConditionalMapper struct {
    SnakeMapper pie.SnakeMapper
    CamelMapper pie.CamelMapper
    UseSnake    bool
}

func (m ConditionalMapper) TableName(structName string) string {
    if m.UseSnake {
        return m.SnakeMapper.TableName(structName)
    }
    return m.CamelMapper.TableName(structName)
}

func (m ConditionalMapper) FieldName(fieldName string) string {
    if m.UseSnake {
        return m.SnakeMapper.FieldName(fieldName)
    }
    return m.CamelMapper.FieldName(fieldName)
}

// 使用条件映射器
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(ConditionalMapper{
        UseSnake: true,
    }),
)
```

## 高级映射器

### 缓存映射器

```go
type CachedMapper struct {
    mapper    pie.Mapper
    tableCache map[string]string
    fieldCache map[string]string
    mutex     sync.RWMutex
}

func NewCachedMapper(mapper pie.Mapper) *CachedMapper {
    return &CachedMapper{
        mapper:     mapper,
        tableCache: make(map[string]string),
        fieldCache: make(map[string]string),
    }
}

func (m *CachedMapper) TableName(structName string) string {
    m.mutex.RLock()
    if name, ok := m.tableCache[structName]; ok {
        m.mutex.RUnlock()
        return name
    }
    m.mutex.RUnlock()
    
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    if name, ok := m.tableCache[structName]; ok {
        return name
    }
    
    name := m.mapper.TableName(structName)
    m.tableCache[structName] = name
    return name
}

func (m *CachedMapper) FieldName(fieldName string) string {
    m.mutex.RLock()
    if name, ok := m.fieldCache[fieldName]; ok {
        m.mutex.RUnlock()
        return name
    }
    m.mutex.RUnlock()
    
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    if name, ok := m.fieldCache[fieldName]; ok {
        return name
    }
    
    name := m.mapper.FieldName(fieldName)
    m.fieldCache[fieldName] = name
    return name
}

// 使用缓存映射器
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(NewCachedMapper(&pie.SnakeMapper{})),
)
```

### 链式映射器

```go
type ChainMapper struct {
    mappers []pie.Mapper
}

func NewChainMapper(mappers ...pie.Mapper) *ChainMapper {
    return &ChainMapper{mappers: mappers}
}

func (m *ChainMapper) TableName(structName string) string {
    name := structName
    for _, mapper := range m.mappers {
        name = mapper.TableName(name)
    }
    return name
}

func (m *ChainMapper) FieldName(fieldName string) string {
    name := fieldName
    for _, mapper := range m.mappers {
        name = mapper.FieldName(name)
    }
    return name
}

// 使用链式映射器
engine, err := pie.NewEngine(ctx, "mydb",
    pie.WithMapper(NewChainMapper(
        &pie.SnakeMapper{},
        &PrefixMapper{TablePrefix: "app_"},
    )),
)
```

## 实际应用场景

### 多租户映射器

```go
type TenantMapper struct {
    tenantID string
    baseMapper pie.Mapper
}

func NewTenantMapper(tenantID string, baseMapper pie.Mapper) *TenantMapper {
    return &TenantMapper{
        tenantID:   tenantID,
        baseMapper: baseMapper,
    }
}

func (m *TenantMapper) TableName(structName string) string {
    baseName := m.baseMapper.TableName(structName)
    return fmt.Sprintf("tenant_%s_%s", m.tenantID, baseName)
}

func (m *TenantMapper) FieldName(fieldName string) string {
    return m.baseMapper.FieldName(fieldName)
}

// 使用多租户映射器
func createTenantEngine(tenantID string) (*pie.Engine, error) {
    return pie.NewEngine(ctx, "mydb",
        pie.WithMapper(NewTenantMapper(tenantID, &pie.SnakeMapper{})),
    )
}
```

### 版本化映射器

```go
type VersionedMapper struct {
    version    string
    baseMapper pie.Mapper
}

func NewVersionedMapper(version string, baseMapper pie.Mapper) *VersionedMapper {
    return &VersionedMapper{
        version:    version,
        baseMapper: baseMapper,
    }
}

func (m *VersionedMapper) TableName(structName string) string {
    baseName := m.baseMapper.TableName(structName)
    return fmt.Sprintf("v%s_%s", m.version, baseName)
}

func (m *VersionedMapper) FieldName(fieldName string) string {
    return m.baseMapper.FieldName(fieldName)
}

// 使用版本化映射器
func createVersionedEngine(version string) (*pie.Engine, error) {
    return pie.NewEngine(ctx, "mydb",
        pie.WithMapper(NewVersionedMapper(version, &pie.SnakeMapper{})),
    )
}
```

### 环境特定映射器

```go
type EnvironmentMapper struct {
    environment string
    baseMapper  pie.Mapper
}

func NewEnvironmentMapper(environment string, baseMapper pie.Mapper) *EnvironmentMapper {
    return &EnvironmentMapper{
        environment: environment,
        baseMapper:  baseMapper,
    }
}

func (m *EnvironmentMapper) TableName(structName string) string {
    baseName := m.baseMapper.TableName(structName)
    return fmt.Sprintf("%s_%s", m.environment, baseName)
}

func (m *EnvironmentMapper) FieldName(fieldName string) string {
    return m.baseMapper.FieldName(fieldName)
}

// 使用环境特定映射器
func createEnvironmentEngine(environment string) (*pie.Engine, error) {
    return pie.NewEngine(ctx, "mydb",
        pie.WithMapper(NewEnvironmentMapper(environment, &pie.SnakeMapper{})),
    )
}
```

## 映射器测试

### 单元测试

```go
func TestCustomMapper(t *testing.T) {
    mapper := CustomMapper{}
    
    // 测试表名映射
    tableName := mapper.TableName("User")
    assert.Equal(t, "t_user", tableName)
    
    // 测试字段名映射
    fieldName := mapper.FieldName("FirstName")
    assert.Equal(t, "firstname", fieldName)
}

func TestPrefixMapper(t *testing.T) {
    mapper := PrefixMapper{
        TablePrefix: "app_",
        FieldPrefix: "f_",
    }
    
    // 测试表名映射
    tableName := mapper.TableName("User")
    assert.Equal(t, "app_user", tableName)
    
    // 测试字段名映射
    fieldName := mapper.FieldName("FirstName")
    assert.Equal(t, "f_firstname", fieldName)
}
```

### 集成测试

```go
func TestMapperIntegration(t *testing.T) {
    engine, err := pie.NewEngine(ctx, "test_db",
        pie.WithMapper(CustomMapper{}),
    )
    require.NoError(t, err)
    defer engine.Disconnect(ctx)
    
    session := pie.Table[User](engine)
    
    // 测试插入
    user := &User{
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john@example.com",
    }
    
    result, err := session.Insert(ctx, user)
    require.NoError(t, err)
    require.NotNil(t, result.InsertedID)
    
    // 测试查询
    var users []User
    err = session.Find(ctx, &users)
    require.NoError(t, err)
    require.Len(t, users, 1)
    assert.Equal(t, "John", users[0].FirstName)
}
```

## 性能优化

### 映射器缓存

```go
type OptimizedMapper struct {
    tableCache sync.Map
    fieldCache sync.Map
    baseMapper pie.Mapper
}

func NewOptimizedMapper(baseMapper pie.Mapper) *OptimizedMapper {
    return &OptimizedMapper{
        baseMapper: baseMapper,
    }
}

func (m *OptimizedMapper) TableName(structName string) string {
    if cached, ok := m.tableCache.Load(structName); ok {
        return cached.(string)
    }
    
    name := m.baseMapper.TableName(structName)
    m.tableCache.Store(structName, name)
    return name
}

func (m *OptimizedMapper) FieldName(fieldName string) string {
    if cached, ok := m.fieldCache.Load(fieldName); ok {
        return cached.(string)
    }
    
    name := m.baseMapper.FieldName(fieldName)
    m.fieldCache.Store(fieldName, name)
    return name
}
```

### 预编译映射

```go
type PrecompiledMapper struct {
    tableNames map[string]string
    fieldNames map[string]string
}

func NewPrecompiledMapper(structs ...any) *PrecompiledMapper {
    mapper := &PrecompiledMapper{
        tableNames: make(map[string]string),
        fieldNames: make(map[string]string),
    }
    
    // 预编译所有结构体的映射
    for _, s := range structs {
        structType := reflect.TypeOf(s)
        if structType.Kind() == reflect.Ptr {
            structType = structType.Elem()
        }
        
        // 预编译表名
        tableName := strings.ToLower(structType.Name())
        mapper.tableNames[structType.Name()] = tableName
        
        // 预编译字段名
        for i := 0; i < structType.NumField(); i++ {
            field := structType.Field(i)
            fieldName := strings.ToLower(field.Name)
            mapper.fieldNames[field.Name] = fieldName
        }
    }
    
    return mapper
}

func (m *PrecompiledMapper) TableName(structName string) string {
    if name, ok := m.tableNames[structName]; ok {
        return name
    }
    return strings.ToLower(structName)
}

func (m *PrecompiledMapper) FieldName(fieldName string) string {
    if name, ok := m.fieldNames[fieldName]; ok {
        return name
    }
    return strings.ToLower(fieldName)
}
```

## 最佳实践

### 1. 映射器设计原则

```go
// 好的映射器设计
type GoodMapper struct {
    // 简单、可预测的转换规则
    // 支持缓存
    // 线程安全
}

// 避免的映射器设计
type BadMapper struct {
    // 复杂的转换逻辑
    // 不可预测的结果
    // 性能问题
}
```

### 2. 映射器文档化

```go
// UserMapper 提供用户相关的命名映射
// 表名映射: User -> users
// 字段映射: FirstName -> first_name, LastName -> last_name
type UserMapper struct {
    baseMapper pie.Mapper
}

func (m *UserMapper) TableName(structName string) string {
    // 将结构体名转换为复数形式
    return m.baseMapper.TableName(structName) + "s"
}

func (m *UserMapper) FieldName(fieldName string) string {
    // 将驼峰命名转换为蛇形命名
    return m.baseMapper.FieldName(fieldName)
}
```

### 3. 映射器测试

```go
func TestMapperComprehensive(t *testing.T) {
    testCases := []struct {
        name        string
        mapper      pie.Mapper
        structName  string
        fieldName   string
        expectedTable string
        expectedField  string
    }{
        {
            name:        "SnakeMapper",
            mapper:      &pie.SnakeMapper{},
            structName:  "User",
            fieldName:   "FirstName",
            expectedTable: "user",
            expectedField:  "first_name",
        },
        {
            name:        "CustomMapper",
            mapper:      CustomMapper{},
            structName:  "User",
            fieldName:   "FirstName",
            expectedTable: "t_user",
            expectedField:  "firstname",
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            tableName := tc.mapper.TableName(tc.structName)
            fieldName := tc.mapper.FieldName(tc.fieldName)
            
            assert.Equal(t, tc.expectedTable, tableName)
            assert.Equal(t, tc.expectedField, fieldName)
        })
    }
}
```

## 下一步

- [错误处理](/reference/error-handling/) - 学习错误处理
- [性能优化](/reference/performance/) - 了解性能优化
- [最佳实践](/best-practices/) - 掌握开发最佳实践
