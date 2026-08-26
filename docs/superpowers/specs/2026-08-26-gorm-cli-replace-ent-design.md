# 用 GORM CLI 替换 Ent

日期：2026-08-26  
状态：待 review

## 背景

仓库里业务表（User / Role / Permission / Product）走 Ent schema + 生成代码；GORM 已经接好连接、读写分离、日志、软删除 `baseModel` 和 Casbin `gorm-adapter`。迁移一直是 `sql-migrate` 的 SQL 文件，不是 Ent schema migrate。

目标：用官方 GORM CLI（`gorm.io/cli`，命令 `gorm gen`）替换 Ent，并删除全部 Ent 相关代码与依赖。usecase / controller / HTTP / gRPC 对外接口不变。

## 目标与非目标

**目标**

- repository 对外接口、domain 实体、`toEntity` 映射保持不变
- 手写 GORM model + CLI 生成 field helper；通用读查询用 typed SQL interface
- 拼条件 / 排序 / 更新字段禁止手写列名，一律走 field helper
- 完整移除 Ent（含 `casbin/ent-adapter`、`pkg/log/ent`、Taskfile 安装项、`go.mod` 依赖）
- 迁移仍只走 `sql-migrate`，不为这次切换改表结构

**非目标**

- 不用 GORM AutoMigrate，也不等 CLI 的 migration 子命令
- 不改 HTTP / gRPC / usecase 行为
- 不为 4 张示例表补全新测试套件（仓库里目前没有 repository 单测）
- 不把 Filter、Create、Update、Delete 写成 SQL interface（留给以后需要时再加）

## 架构

```
handler / grpc → usecase → *RepositoryInterface
                              ↓
                         repository 实现
                              ↓
              ┌───────────────┴───────────────┐
              │                               │
     queries.Query[T]              gorm.G[models.T] + field helper
     (FindOne / Exist)             (Filter、实体特有查询、写)
              │                               │
              └───────────────┬───────────────┘
                              ↓
                         *gorm.DB（已有 internal/pkg/gorm）
                              ↓
                         domain 实体
```

**不变**

- `internal/pkg/gorm` 连接、resolver、日志
- `sql-migrate`、`migrations/`
- Casbin 的 gorm-adapter 与 file adapter
- domain 包

**变**

- repository 构造函数从 `*ent.Client` 改为 `*gorm.DefaultDB`
- 去掉 `ent.ProvideDefault` 后，启动时不再额外开一条 Ent 连接（只保留 GORM 那条）

## 目录

```
internal/app/repository/
  repository.go                 # wire、handleError（只映射 gorm）
  user.go role.go permission.go product.go
  policy.go                     # 不动
    models/
    base.go                     # 原 baseModel
    user.go role.go permission.go product.go
  queries/
    generate.go
    gengen.go                   # go:generate 入口：gen 后改名为 *.gen.go
    src/query.go                # 手写泛型 Query[T]（仅作 generate 输入）
    *.gen.go                    # CLI 输出（field helper + Query 实现）
```

### 生成物放哪（与「同目录」的约束）

GORM CLI 给每个 model 生成 `var User` 这种 field helper。Go 不允许同一包里同时存在 `type User` 和 `var User`。因此 **field helper 不能生成进 `models/`**。

落地约定：

- `models/`：**只放手写 struct**，不放生成代码
- `queries/`：手写 `Query[T]` + 生成的 Query 实现 + 生成的 field helper（`*.gen.go`）
- 调用方写 `queries.User.Username.Eq(...)`、`gorm.G[models.User](db)`

这是「源文件与生成代码同目录」在 Go 类型规则下能编译的最小调整：生成物和 query 源文件在一起；model 源文件单独一包。

## models

从现有 Ent schema 和 `migrations/*/default` 对齐，不改列。

`models.Base`（从 `repository.baseModel` 挪过来）：

```go
type Base struct {
    ID        int64                 `gorm:"primaryKey"`
    CreatedAt int64                 `gorm:"autoCreateTime;not null"`
    UpdatedAt int64                 `gorm:"autoUpdateTime;not null;default:0"`
    DeletedAt soft_delete.DeletedAt `gorm:"index;not null;default:0"`
}
```

`DeletedAt` 与现有 `repository.baseModel` 一致：unix **秒**、`0` 表示未删除。不要加 `softDelete:milli` / `nano`，否则会与表里的 `bigint` 秒时间戳对不上。

四个实体嵌入 `Base`，`TableName()` 分别为 `users` / `roles` / `permissions` / `products`。`Product.Desc` 映射列名 `` `desc` ``（SQL 保留字）。索引、unique、size 与现有 migration 一致：

| 表 | 额外约束 |
| --- | --- |
| users | unique username；index phone |
| roles | unique name，size 32 |
| permissions | unique key，size 128；parent_id |
| products | index name |

ID 为数据库 `bigserial` / autoIncrement，Create 不手写 ID（与现在 Ent Create 不 SetID 一致）。

## queries

手写一份泛型 interface，四个实体共用。CLI 会生成同名 `func Query`，与 `type Query` 不能同包，因此手写 interface 放在 `queries/src`（只给 generate 扫，不参与编译到 `queries` 包）：

```go
type Query[T any] interface {
    // SELECT * FROM @@table WHERE id=@id AND deleted_at=0
    FindOne(id int64) (T, error)

    // SELECT COUNT(1) FROM @@table WHERE id=@id AND deleted_at=0
    Exist(id int64) (int64, error)
}
```

调用：`queries.Query[models.User](db).FindOne(ctx, id)`。

**必须带 `deleted_at=0`**：CLI 把 SQL 编成 `Raw()`，会绕过 GORM soft-delete plugin。

**FindOne 空结果**：`Raw().Scan` 在 0 行时可能 `err == nil` 且零值。repository 在 `ID == 0` 时映射为 `ErrRecordNotFound`（本项目 ID 从 1 起）。实现时若 CLI/GORM 已返回 `gorm.ErrRecordNotFound`，则只走 `handleError`，不再用 ID==0 兜底。

**Exist**：生成方法返回 `int64`（COUNT）。repository 的 `Exist(...) (bool, error)` 映射为 `n > 0`。

以后加复杂 SQL：在 `queries` 再加专用 interface，不要塞进这个泛型 `Query[T]`。

## repository 实现约定

对外接口一字不改。内部：

| 方法类型 | 实现 |
| --- | --- |
| `FindOne(id)` / `Exist(id)` | `queries.Query[T]` |
| `Filter`、keyword、`IDNEQ`、`IDIn`、按 username/name/key 查 | `gorm.G[T]` + `queries.Xxx` field helper |
| Create | `gorm.G[T](db).Create(ctx, &m)` 或 `Set(helper).Create` |
| Update | `Where(queries.Xxx.ID.Eq(id)).Set(queries.Xxx.Field.Set(...)).Update` |
| Delete | `Where(queries.Xxx.ID.Eq(id)).Delete`（走软删除 plugin） |
| Casbin 授权（AssignRoles 等） | 逻辑不动，只换查表客户端 |

禁止 `"username = ?"`、`Order("updated_at desc")` 这类字符串列名。

所有返回 error 的路径只调用 `handleError`，方法里不再 `WithStack`。FindOne 若 CLI `Scan` 在 0 行时 `err == nil` 且 `ID == 0`，先把 err 设为 `gorm.ErrRecordNotFound` 再交给 `handleError`。

构造函数注入 `*gorm.DefaultDB`（及现有 `*casbin.Enforcer`）。

## 生成工作流

`queries/generate.go`：

```go
//go:generate go run gengen.go
```

`gengen.go`（`//go:build ignore`）依次执行：

1. `gorm gen -i ../models -o ../queries`
2. `gorm gen -i ./src -o ../queries`
3. 把带 CLI 头的 `.go` 改名为 `*.gen.go`（不覆盖 `generate.go` / `gengen.go`）

`-o` 必须是以 `queries` 结尾的路径；`-o .` 会生成非法的 `package .`。

手写 SQL interface 固定在 `queries/src/query.go`。生成的 `func Query` 在 `query.gen.go`。

`go generate ./internal/app/repository/...` 必须可重复执行。`Taskfile.yaml` 的 `download`：删除 `entgo.io/ent/cmd/ent`，改为 `go install gorm.io/cli/gorm@latest`。`go.mod` 用 `tool` 或 generate 里 `go run gorm.io/cli/gorm` 钉版本，避免每人 CLI 不一致。

## Ent 拆除清单

删除：

- `internal/pkg/ent/`（含全部生成代码）
- `internal/app/repository/schema/`（含 mixin、types）
- `pkg/log/ent/`
- `internal/pkg/casbin/adapter/ent.go`
- `config.CasbinAdapter.Ent` / `CasbinEntAdapter`
- `etc/config.yaml` 里注释掉的 `ent: {}`
- `go.mod`：`entgo.io/ent`、`github.com/casbin/ent-adapter`（以及因此变成无用的 indirect，如 atlas）
- `Taskfile.yaml` 里的 `ent` 安装

修改：

- `internal/pkg/pkg.go`：去掉 `ent.ProvideDefault`
- `internal/pkg/casbin`：`New` / `adapter.New` / `Provide` 去掉 `*sql.DB` 和 Ent 分支（file / gorm 分支保留）。去掉 Ent 后，`Provide` 不必再为 Ent 单独 `db.New` 一条连接
- `internal/command/wire.go` 不变入口的话，重跑 `wire` 更新 `wire_gen.go`
- `repository.go`：`IsNotFound` / `handleError` 去 Ent

`go.mod` 里已有的 `gorm.io/gorm`、`gorm.io/plugin/soft_delete`、`gorm.io/plugin/dbresolver`、casbin gorm-adapter 保留。新增 `gorm.io/cli`（generate 用）。

## 错误处理

转换和堆栈都集中在 `handleError`，repository 方法只 `return handleError(err)`（或 `return nil, handleError(err)`），不再调用 `github.com/pkg/errors.WithStack`。

```go
func handleError(err error) error {
    if err == nil {
        return nil
    }
    if errors.Is(err, gorm.ErrRecordNotFound) {
        err = ErrRecordNotFound
    }
    return uerr.WithStack(err, 4)
}
```

- 用 `go-scaffold/pkg/errors.WithStack`，**不要**用 `github.com/pkg/errors.WithStack`
- `skip` 为 **4**：与 `internal/errors.Error.Wrap` 一致。`runtime.Callers` 经 `callers` → `WithStack` → `handleError`，第一帧是调用 `handleError` 的 repository 方法，堆栈里看不到 `handleError` / `WithStack`
- `WithStack` 已有栈则不再包一层
- 对外 sentinel 仍是 `repository.ErrRecordNotFound`；`errors.Is` 可识别。去掉 `ent.IsNotFound`
- Casbin、`policy.go` 解析失败等同样走 `handleError`，不在方法里单独 `WithStack`

`IsNotFound` 只判断 `ErrRecordNotFound`（以及仍可能漏网的 `gorm.ErrRecordNotFound`）。

## 测试

仓库没有 repository `*_test.go`。这次验收：

1. `go generate ./internal/app/repository/...` 成功，且不覆盖手写文件
2. `go generate` 后的 `wire` 成功
3. `go test ./...` 与 `go build` 通过
4. grep 确认不再出现 `entgo.io/ent`、`casbin/ent-adapter`、`internal/pkg/ent`

不强制新增 sqlmock 用例。若实现时 FindOne 空结果行为不好测，再给 `handleError` / FindOne 映射加一个小纯函数测试。

## 以后加表

1. `migrations/` 加 SQL（与现在一样）
2. `models/` 加 struct，嵌入 `Base`
3. 需要通用 FindOne/Exist：已有 `Query[T]`，不用改 interface
4. `go generate` 更新 `queries/*.gen.go`
5. 新建 `xxx.go` repository：动态查询用 `queries.Xxx` helper，写操作用 `gorm.G`

## 范围外明确不做

- 不迁移历史 Ent 生成代码「对照保留」
- 不把 Casbin 规则表纳入这次的 models/queries
- 不引入 gorm.io/gen（旧 GEN）
