# match-engine 学习路线图

参考项目：`../matching-engine`（Go 写的高性能内存撮合引擎 SDK，`github.com/0x5487/matching-engine`）

本文件是**进度追踪表**，不是代码。代码由你自己逐课敲入。

---

## 为什么是这个顺序

参考项目的真实历史是 2021→2026 分 8 个版本长出来的（见 `../matching-engine/CHANGELOG.md`），中间有大量返工：
API 从 `Submit` 统一入口拆回类型化请求、`Metadata` 字段加了又删、`AggregatedBook` 写完就废弃。

**按历史顺序学是错的。** 本路线图按**依赖关系自底向上**重排，保证每一课结束时项目都能
`go build ./...` 通过。核心依赖链是：

```
protocol（谁都不依赖）
   ↓
error / models / order_book_log（领域模型）
   ↓
structure（跳表）→ queue（价格档位）
   ↓
order_book（撮合核心）
   ↓
disruptor（并发原语）→ engine（编排）
   ↓
snapshot / 线格式 / aggregated_book
```

---

## 环境

```bash
cd ~/Projects/match-engine
source ./goenv.sh          # 提供 go 1.26.5 + gcc 15.3.0 + 工作区内的缓存目录
cd match-engine
```

`goenv.sh` 存在的两个原因：
1. 本机没有全局 Go/gcc（NixOS，且 go 不在 profile 里）；
2. 沙箱只允许写工作区，默认的 `~/.cache/go-build` 和 `~/go/pkg/mod` 不可写，必须重定向。

（`devenv.nix` 里的 `languages.go.enable` 也能提供 Go，二选一即可。）

---

## 进度表

阶段图例：`[ ]` 未开始 · `[~]` 进行中 · `[x]` 完成

### A 地基

- [ ] **L1 项目骨架** — `go.mod`、目录布局、`Makefile`
      产出：`go build ./...` 通过
- [ ] **L2 `protocol/types.go`** — 枚举、日志类型、拒绝原因、查询响应
      产出：175 行纯声明
- [ ] **L3 `protocol/command.go`（上）+ `protocol/query.go`** — `BaseCommand` + 8 个类型化请求 + `GetRequestBase`
      产出：命令信封与类型化负载
      ⚠️ `query.go`（18 行）必须在这里一起写，`order_book.go` 的 `processQuery` 依赖 `GetDepthQuery`/`GetStatsQuery`
- [ ] **L4 `error.go`** — 11 个哨兵错误值

### B 领域模型与事件日志

- [ ] **L5 `models.go` + `snapshot.go`** — `Order` / `DepthChange` / `InputEvent` + 快照 DTO
      ⚠️ 本课 `models.go` **不含** `Future[T]`（它引用 `*MatchingEngine`），L17 再加回来
      ⚠️ `snapshot.go` 的 `OrderBookSnapshot` 必须在这里写，`order_book.go` 的 `Restore`/`createSnapshot` 依赖它
- [ ] **L6 `order_book_log.go`** — `OrderBookLog`、`LogBatch`、`Publisher`、7 个日志构造器、`sync.Pool`

### C 数据结构与价格队列

- [ ] **L7 `structure/pooled_skiplist.go`** — arena 池化跳表（本路线图最难的数据结构）
      ⚠️ 需自行声明 `NullIndex`，见下方「与参考项目的偏差」
- [ ] **L8 `queue.go`** — `priceUnit` + 侵入式 FIFO 链表 + 价格索引
      产出：买卖两侧队列，最优价 O(1) 查询

### D 撮合核心

> **分课策略（已实测验证）**：`order_book.go` 有 1146 行，一次敲完不现实，而且文件中间状态无法编译。
> 实测确认的做法是**增量式**：L9 先写一个 **517 行可编译可撮合的子集**，之后每课往里**追加**一个
> handler + 一个 `switch` case。每一课结束都能 `go build` 并跑通冒烟测试。
> （另一种选择：L9 一次性贴完整 1146 行——实测也能编译，但后续每课变成"改已有函数"，手感差。）

- [ ] **L9 `order_book.go`(1)** — `OrderBook` 骨架、校验分层（`validateBasic`/`validateState`）、**Limit 撮合主循环**、`handlePlaceOrder`、`depth`/`stats`/`createSnapshot`
      产出：**517 行可撮合子集**。`checkReplenish` 先留 TODO 桩
- [ ] **L10 `order_book.go`(2)** — 追加 `Cancel` + `Amend`（优先级规则：改价/增量失去优先级）
- [ ] **L11 `order_book.go`(3)** — 追加 `Market`（Size / QuoteSize 双模式 + `lotSize` 防死循环）
- [ ] **L12 `order_book.go`(4)** — 追加 `IOC` / `FOK`（两阶段预检）/ `PostOnly`
- [ ] **L13 `order_book.go`(5)** — 把 `checkReplenish` 桩换成真实现（Iceberg 自动补量 + 优先级重置）
- [ ] **L14 `order_book.go`(6)** — 追加 `Suspend`/`Resume`/`UpdateConfig` handler + `snapshotQuery`
      ⚠️ `snapshotQuery` 在参考项目里声明于 `engine.go`，我们会声明在 `order_book.go`（见偏差表 #5）

### E 并发与引擎

- [ ] **L15 `disruptor.go`** — MPSC 环形缓冲、缓存行填充、`IdleStrategy`
- [ ] **L16 `engine.go` 核心** — `MatchingEngine`、多市场路由、`CreateMarket`、`UserEvent`
- [ ] **L17 `Future[T]` + 响应通道池** — 回到 `models.go` 补 `Future`，`engine.go` 补 `responsePool`

### F 持久化与线格式

- [ ] **L18 `snapshot.go` + 快照 IO** — CRC32、footer、段校验、临时目录原子替换
- [ ] **L19 `protocol/command.go`（下）** — 手写二进制线格式 `MarshalRequest`/`UnmarshalRequest`
- [ ] **L20 `aggregated_book.go`** — 下游聚合簿（参考项目此文件是 **TODO 空实现**，我们补完）

### G 质量

- [ ] **L21 测试与 benchmark** — testify 风格、`assert.Eventually`、benchmark、`golangci-lint`

---

## 与参考项目的偏差（重要）

这些是**有意识的偏离**，不是抄错：

| # | 偏差 | 原因 |
|---|---|---|
| 1 | 不写 `structure/llrb_tree.go`（423 行） | 参考项目同一 commit 里做了**两个**价格索引候选（LLRB 红黑树 + 池化跳表），最终**只用跳表**，LLRB 是死代码。省掉它 |
| 2 | 因此需自己声明 `const NullIndex int32 = -1` | 参考项目的 `NullIndex` 声明在 `llrb_tree.go:21`，跳表文件反而依赖它。跳过 LLRB 后必须自己声明。建议放在 `pooled_skiplist.go` 的 import 块之后 |
| 3 | `L5` 的 `models.go` 不含 `Future[T]` | `Future` 有 `engine *MatchingEngine` 字段并调用 `engine.releaseResponseChannel`，在 L17 引擎就绪前无法编译。先省略，L17 补回（这本身就是真实的增量开发方式） |
| 4 | module path 是 **`github.com/miliu2cc/match-engine`**（你已填好 `go.mod`） | 参考项目是 `github.com/0x5487/matching-engine`。**所有 `import` 行必须用你的路径**，例如 `"github.com/miliu2cc/match-engine/protocol"` |
| 4b | `go 1.26.5`（你写的是带补丁号的版本） | 语法合法（Go 1.21+ 允许），语义是"要求工具链 ≥1.26.5"。若想让项目在稍旧工具链上也能编译，可改成 `go 1.25`（**语言版本**）。两种都能跑，不影响本项目功能 |
| 5 | `type snapshotQuery struct{}` 声明在 `order_book.go` | 参考项目把它放在 `engine.go`，但 `order_book.go` 的 `processQuery` 要用它。放在 `order_book.go` 可让 Phase D 不依赖 Phase E |
| 6 | 提前写 `protocol/query.go`（L3） | 它是 18 行的纯声明，但 `order_book.go` 的 `processQuery` 依赖它。实测：缺了它 `order_book.go` 直接编译失败 |
| 7 | 提前写 `snapshot.go` 的 DTO（L5） | `OrderBookSnapshot` 被 `order_book.go` 的 `Restore()`/`createSnapshot()` 引用。实测：缺了它 `order_book.go` 编译失败 |

### ⚠️ 建议加固的地方（参考项目的潜在隐患）

这些**不是抄错**，是参考项目本身的薄弱点。实测确认过，你的实现可以做得更稳。

| # | 隐患 | 实测证据 | 建议 |
|---|---|---|---|
| A | **`map[udecimal.Decimal]` 的键用 `==` 比较，而 `Decimal.Equal()` 用数值比较——两者不等价** | `MustParse("10") == MustParse("10.0")` → **false**（但 `Equal()` → true，且两者 `String()` 都打印 `10`）。所以 `m[MustParse("10")]` 用 `MustParse("10.0")` **查不到**；两个都赋值会产生 **2 个键** → **同一价格出现两个档位** | 在订单入口处把价格规范化一次，例如 `price, _ = udecimal.Parse(price.String())`。走二进制线格式（L19 的 `String()`→`Parse()` 往返）**会自动规范化**，所以只有"进程内直接构造请求"（README 的 Quick Start 用法）才会踩到 |
| B | **`handleSuspendMarket` 的 `LogTypeAdmin.EventType` 传的是 `req.Reason`，而文档规定是 `market_suspended`** | 实测：`l14Suspend(book, "maintenance")` 产出的 `EventType` 是 **`"maintenance"`**，不是 `"market_suspended"`。而 `handleResumeMarket` 传字面量 `"market_resumed"`、`handleUpdateConfig` 传 `"market_config_updated"`——**三个 handler 只有 suspend 不一致**，且与 `docs/design/management.md` 第 92–97 行的规定冲突 | 下游若按文档用 `EventType == "market_suspended"` 过滤，**会漏掉所有停牌事件**。建议传固定字面量 `"market_suspended"`，把 `Reason` 放到别的字段（如 `Data`）。**你的实现可以修正这一点** |

**参考项目里其他可以跳过的死代码**（了解即可，不必实现）：
- `aggregated_book.go` 的 `Replay()` 是 TODO 空实现，`aggregated_book_test.go` 只有一行 `package match`
- Makefile 的 `prof-mem` / `prof-cpu` 目标引用的 `BenchmarkPlaceOrders` / `BenchmarkMatching` **已不存在**
- `docs/benchmark.md` 报告的 385ns/op 也来自那个已删除的 benchmark，数字已失效

---

## 参考项目的关键约定（照抄，否则会踩坑）

| 约定 | 说明 |
|---|---|
| `go 1.25` | 用到 Go 1.22+ 的 `for i := range 整数`、1.21+ 的 `min()` 内建函数 |
| 根包名 `match` | 不是 `main`。SDK 是给别人 import 的库 |
| 行宽 120 | 由 **formatter `golines`** 强制，**不是** linter `lll`（`lll` 是关闭的） |
| `//nolint:x // 理由` | 注释与 `//` 之间**无空格**（`revive` 的 `comment-spacings: [nolint]` 允许） |
| import 分三组 | 标准库 / 第三方 / 本项目，组间空行（`gci` 的 `standard, default, localmodule`） |
| 禁止 `float64` 存价格 | 一律 `udecimal.Decimal`；对外 JSON 序列化成**字符串** |
| 引擎内禁止 `time.Now()` 打业务时间戳 | 破坏重放确定性，时间戳必须由上游命令携带 |
| 每个导出标识符都要有 doc 注释 | `revive` 强制，且注释须以标识符名开头 |

---

## 每课的自检命令

```bash
go build ./...     # 无输出 = 成功
go vet ./...       # 静态检查
gofmt -l .         # 列出格式不合规的文件；无输出 = 干净
```

参考项目实测：上述三条在**全部**源码上均通过（`gofmt -l` 输出为空）。

---

## 已验证的里程碑

下面这些状态是在**独立的临时 module** 里真实编译/运行过的，不是推测。
你敲出来的代码如果和这些里程碑对不上，就是抄写出了问题。

| 里程碑 | 内容 | 验证结果 |
|---|---|---|
| **L2** | `protocol/types.go`（175 行） | BUILD OK / VET OK / FMT CLEAN |
| **L3** | `command.go` 类型部分 + `udecimal` 依赖拉取 | BUILD OK / VET OK / FMT CLEAN |
| **L4–L8 累积** | `error.go` + `models.go`(无 Future) + `order_book_log.go` + `pooled_skiplist.go` + `queue.go` | BUILD OK / VET OK / FMT CLEAN |
| **L9 子集** | `order_book.go` 517 行（OrderBook + 校验 + Limit 撮合） | BUILD OK，**且冒烟测试通过**：sell@110 挂单 → buy@115 吃单，产出 `open` + `match`（maker 价 110、size 1、amount 110、tradeID 1），撮合后订单簿清空 |
| **完整 Phase D** | 参考项目原始 `order_book.go`（1146 行，未裁剪） | BUILD OK / VET OK / FMT CLEAN —— 证明 L9–L14 的最终形态可行 |

> **注意**：`snapshotQuery` 和 `protocol/query.go` 是上面两次构建的关键依赖，见偏差表 #5/#6/#7。
