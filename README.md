## 廣告投放服務

### 介紹

廣告投放服務提供了創建廣告、檢索廣告等功能。該系統使用了 Go 語言編寫後端服務，並使用 Gin 框架處理 HTTP 請求，同時利用 PostgreSQL 存儲廣告信息，並使用 Redis 緩存來提高查詢性能。另外，也將所有服務打包成一個docker-compose檔方便佈署。

### 技術

- **後端開發語言**: Go
- **Web 框架**: Gin
- **資料庫**: PostgreSQL
- **緩存**: Redis

### 結構

專案分為三個主要模塊：`api`、`database`、`redisDB`。

- **api**: 包含處理 HTTP 請求的路由和處理函數。提供了創建廣告和檢索廣告的接口。
- **database**: 提供了與 PostgreSQL 數據庫交互的功能。包含了創建廣告、檢索廣告。
- **redisDB**: 提供了與 Redis 緩存交互的功能。初始化 Redis 客戶端並提供了設置和獲取緩存數據的函數。

### 設計選擇

1. **Gin 框架**: 選擇 Gin 框架作為 Web 框架，因為它具有輕量級、快速和簡單易用的特點，適合快速開發 RESTful API。
   
2. **PostgreSQL 資料庫**: 使用 PostgreSQL 儲存廣告信息。選擇 PostgreSQL 的原因包括其穩定性、可靠性和功能強大的 JSONB 數據類型，能夠存儲靈活的廣告條件數據。
   
3. **Redis 緩存**: 利用 Redis 緩存來提高廣告檢索接口的性能。通過緩存最近的查詢結果，減少對資料庫的頻繁查詢。

4. **分表分區**: 針對廣告資料表，根據 `end_at` 字段進行分區，以提高查詢效率。
5. **索引**：創建`end_at`的索引，提高檢索廣告的排序效率。
6. **上下文傳遞**: 使用 Gin 框架的上下文功能，將資料庫連接和其他相關對象傳遞給處理函數，使得處理函數可以方便地訪問資料庫和其他資源。

### Getting Started

1. Clone the repository to your go path.
2. Set up a database PostgreSQL and configure the connection details in the application.
3. Run the application.
   ```bash
   docker-compose up --build
   ```
4. Access the APIs using the provided endpoints and methods.

### Request範例

1.  建立廣告 API

**請求方式**

```
POST /api/v1/ad
```

**請求參數**

| 參數        | 類型   | 描述                 |
|-------------|--------|----------------------|
| title       | 字符串 | 廣告標題             |
| startAt     | 時間   | 廣告開始時間         |
| endAt       | 時間   | 廣告結束時間         |
| conditions  | JSON   | 廣告條件             |

**Request**

```json
{
    "title": "Summer Sale",
    "startAt": "2024-06-01T00:00:00Z",
    "endAt": "2024-06-30T23:59:59Z",
    "conditions": {
        "ageStart": 18,
        "ageEnd": 50,
        "gender": "M",
        "country": ["US", "CA"],
        "platform": ["ios", "android"]
    }
}
```

**Response**

- **成功**:
  - 狀態碼: 201 Created
  - 內容: `{"message": "Advertisement created successfully"}`
- **失敗**:
  - 狀態碼: 錯誤狀態碼
  - 內容: `{"error": "錯誤消息"}`

---

2. 獲取廣告 API

**請求方式**

```
GET http://localhost:8080/api/v1/ad
```

**請求參數(Optional)**

| 參數        | 類型   | 描述                  |
|-------------|--------|-----------------------|
| age         | 字符串 | 年齡條件（Optional）      |
| gender      | 字符串 | 性別條件（Optional）      |
| country     | 字符串 | 國家條件（Optional）      |
| platform    | 字符串 | 平台條件（Optional）      |
| offset      | 數字   | 資料開始值（Optional）    |
| limit       | 數字   | 返回資料限制（Optional）  |

**Request**

```
GET /api/v1/ad?age=25&gender=M&country=US&platform=ios&offset=0&limit=10
```

**Response**

- **成功**:
  - 狀態碼: 200 OK
  - 內容: 返回廣告列表的 JSON 資料
- **失敗**:
  - 狀態碼: 適當的錯誤狀態碼
  - 內容: `{"error": "錯誤消息"}`

**注意事項**

- 如果沒有提供任何參數，則將返回所有廣告。
- 參數 `offset` 和 `limit` 用於分頁，`offset` 指定返回結果的起始位置，`limit` 指定返回的數據量。
- 其他參數用於篩選廣告，可以根據年齡、性別、國家和平台進行篩選。



### 參考資料

- [Gin 文檔](https://github.com/gin-gonic/gin)
- [PostgreSQL 文檔](https://www.postgresql.org/docs/)
- [Redis 文檔](https://redis.io/documentation)

