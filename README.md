# Advertisement Delivery Service

This is a simplified advertisement delivery service designed and implemented in Golang, PostgreSQL and Docker. The service provides two main APIs: one for creating advertisements and another for listing advertisements. Each advertisement can have conditions associated with it, such as age, gender, country, and platform of the user.

## Features

- **Admin API**: Used for managing advertisement resources. It supports creating advertisements with specified conditions. The API includes:
  - Title: Title of the advertisement.
  - StartAt and EndAt: Time period for displaying the advertisement.
  - Conditions: Optional criteria such as age, gender, country, and platform.

- **Public API**: Used for fetching advertisements based on user-specified criteria. The API provides:
  - Listing active advertisements within the specified conditions.
  - Pagination support using offset and limit parameters.
  - Query parameters for filtering advertisements based on age, gender, country, and platform.

## APIs

### Admin API

- Endpoint: `/api/v1/ad`
- Method: `POST`
- Usage:
```bash
curl -X POST -H "Content-Type: application/json" \
"http://<host>/api/v1/ad" \
--data '{
  "title": "AD 55",
  "startAt": "2023-12-10T03:00:00.000Z",
  "endAt": "2023-12-31T16:00:00.000Z",
  "conditions": [
    {
      "ageStart": 20,
      "ageEnd": 30,
      "country": ["TW", "JP"],
      "platform": ["android", "ios"]
    }
  ]
}'
```

### Public API
- Endpoint: `/api/v1/ad`
- Method: `GET`
- Usage:
```bash
curl -X GET -H "Content-Type: application/json" \
"http://<host>/api/v1/ad?offset=10&limit=3&age=24&gender=F&country=TW&platform=ios"

```
- Response:
```json
{
  "items": [
    {
    "title": "AD 1",
    "endAt": "2023-12-22T01:00:00.000Z"
    },
    {
    "title": "AD 31",
    "endAt": "2023-12-30T12:00:00.000Z"
    },
    {
    "title": "AD 10",
    "endAt": "2023-12-31T16:00:00.000Z"
    }
  ]
}
```
## Prerequirements
- Install Golang if not already installed.
- Install Docker if not already installed.


## How to Start
1. Clone the repository to your go path.
2. Set up a database PostgreSQL and configure the connection details in the application.
3. Run the application.
   ```bash
   docker-compose up --build
   ```
4. Access the APIs using the provided endpoints and methods.
