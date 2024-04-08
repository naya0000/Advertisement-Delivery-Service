## 廣告管理系統

### 概述

廣告管理系統是一個用於管理廣告的系統，它提供了創建廣告、檢索廣告等功能。該系統使用了 Go 語言編寫後端服務，並使用 Gin 框架處理 HTTP 請求，同時利用 PostgreSQL 存儲廣告信息，並使用 Redis 緩存來提高查詢性能。另外，也將所有服務打包成一個docker-compose檔。

### 技術棧

- **後端開發語言**: Go
- **Web 框架**: Gin
- **數據庫**: PostgreSQL
- **緩存**: Redis

### 結構

項目分為三個主要模塊：`api`、`database`、`redisDB`。

- **api**: 包含處理 HTTP 請求的路由和處理函數。提供了創建廣告和檢索廣告的接口。
- **database**: 提供了與 PostgreSQL 數據庫交互的功能。包含了創建廣告、檢索廣告。
- **redisDB**: 提供了與 Redis 緩存交互的功能。初始化 Redis 客戶端並提供了設置和獲取緩存數據的函數。

### 設計選擇

1. **Gin 框架**: 選擇 Gin 框架作為 Web 框架，因為它具有輕量級、快速和簡單易用的特點，適合快速開發 RESTful API。
   
2. **PostgreSQL 數據庫**: 使用 PostgreSQL 作為主要數據庫存儲廣告信息。選擇 PostgreSQL 的原因包括其穩定性、可靠性和功能強大的 JSONB 數據類型，能夠存儲靈活的廣告條件數據。
   
3. **Redis 緩存**: 利用 Redis 緩存來提高廣告檢索接口的性能。通過緩存最近的查詢結果，減少對數據庫的頻繁查詢，提高響應速度和系統性能。

4. **分表分區**: 針對廣告表，根據 `end_at` 字段進行分區，以提高查詢效率和管理數據。
5. **索引**：創造`end_at`的索引，提高檢索廣告的排序效率。
6. **上下文傳遞**: 使用 Gin 框架的上下文功能，將數據庫連接和其他相關對象傳遞給處理函數，使得處理函數可以方便地訪問數據庫和其他資源。

### Getting Started

docker-compose up --build

### 參考資料

- [Gin 文檔](https://github.com/gin-gonic/gin)
- [PostgreSQL 文檔](https://www.postgresql.org/docs/)
- [Redis 文檔](https://redis.io/documentation)
