## ADDED Requirements

### Requirement: 容器启动时自动执行 SQL migration
后端容器启动前 SHALL 执行 `cmd/migrate` 工具，运行所有未执行的 SQL migration 文件。migration 成功后才启动 server 进程。

#### Scenario: 首次部署空数据库
- **WHEN** 数据库为空，`schema_migrations` 表不存在
- **THEN** 容器创建 `schema_migrations` 表，按文件名排序执行 `migrations/sql/` 下所有 `.sql` 文件，记录到 `schema_migrations`，然后启动 server

#### Scenario: 增量部署有新 migration
- **WHEN** 数据库已有部分 migration 执行记录，`migrations/sql/` 中有新增的 `.sql` 文件
- **THEN** 容器跳过已执行的 migration，只执行新增的，记录到 `schema_migrations`，然后启动 server

#### Scenario: 所有 migration 已是最新
- **WHEN** 所有 `.sql` 文件都已在 `schema_migrations` 中有记录
- **THEN** 容器输出"所有 migration 已是最新"，直接启动 server

### Requirement: migration 失败时容器不启动
如果任何 migration 执行失败，容器 SHALL 以非零退出码退出，不启动 server 进程。

#### Scenario: migration SQL 语法错误
- **WHEN** 某个 `.sql` 文件包含语法错误的 SQL
- **THEN** 容器输出错误信息并 `exit 1`，server 不启动。配合 `restart: unless-stopped`，容器会自动重试。

### Requirement: 数据库未就绪时自动重试
MySQL 容器启动后需要时间初始化，migration 工具 SHALL 在连接失败时退出，由 Docker restart 策略自动重试。

#### Scenario: MySQL 容器刚启动还未就绪
- **WHEN** migration 工具无法连接到数据库
- **THEN** 容器退出，Docker `restart: unless-stopped` 策略自动重启容器，直到 MySQL 就绪后 migration 成功执行

### Requirement: Dockerfile 包含 migrate 二进制
后端 Dockerfile SHALL 在构建阶段编译 `cmd/migrate`，最终镜像同时包含 `./server` 和 `./migrate` 两个二进制。

#### Scenario: 镜像构建产物
- **WHEN** 执行 `docker build`
- **THEN** 最终镜像中包含 `/app/server` 和 `/app/migrate` 两个可执行文件，以及 `migrations/sql/` 目录

### Requirement: 通过环境变量配置数据库连接
docker-compose.yml 中 SHALL 配置 migration 所需的数据库连接环境变量（`DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASSWORD`、`DB_NAME`），Docker 环境中 host 为 `db`。

#### Scenario: 容器内连接数据库
- **WHEN** 后端容器在 docker-compose 环境中启动
- **THEN** `DB_HOST` 为 `db`（docker-compose 服务名），migration 工具能通过该地址连接到 MySQL 容器
