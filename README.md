# Financing Application Aggregator

A backend service that aggregates financing applications from multiple banks and provides offers to customers.

## Running the Application

1. Set up environment variables:
```bash
cp .env.example .env
```

2. Edit `.env` and add your bank API URLs:
```bash
# Replace with actual API endpoints
FASTBANK_BASE_URL=https://your-actual-fastbank-url
SOLIDBANK_BASE_URL=https://your-actual-solidbank-url
```

3. Start the application using Docker Compose:
```bash
docker-compose up -d
```

The service will be available at `http://localhost:8080`

## Example Usage

### Submit an Application

```bash
curl -X POST http://localhost:8080/api/v1/applications \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+37126000000",
    "email": "john.doe@example.com",
    "monthlyIncome": 3000,
    "monthlyExpenses": 1200,
    "maritalStatus": "MARRIED",
    "agreeToBeScored": true,
    "amount": 8000,
    "dependents": 1
  }'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "PENDING"
}
```

### Check Application Status

```bash
curl http://localhost:8080/api/v1/applications/550e8400-e29b-41d4-a716-446655440000
```

Response (after processing):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "COMPLETED",
  "offers": [
    {
      "bank": "FastBank",
      "status": "APPROVED",
      "amount": 8000,
      "interestRate": 5.5,
      "monthlyPayment": 152.33
    },
    {
      "bank": "SolidBank",
      "status": "DECLINED",
      "reason": "Insufficient income"
    }
  ]
}
```

## Running Tests

### Unit Tests

```bash
go test ./...
```

### API Tests with Bruno

1. Install Bruno desktop app:
   - Download from [bruno.usebruno.com](https://www.usebruno.com/downloads)
   - Or install via package manager for your system

2. Open Bruno and import the collection:
   - Click "Open Collection"
   - Navigate to the `.bru/aggregator` folder
   - The collection includes test flows for successful and high-risk applications

3. Set the environment:
   - Select "Docker" environment
   - Verify `base_url` is set to `http://localhost:8080`

4. Run tests:
   - Use "02-successful-flow" for applications that should get approved
   - Use "03-high-risk-flow" for high-risk profile testing

## API Endpoints

- `POST /api/v1/applications` - Submit application
- `GET /api/v1/applications/{id}` - Get application status
- `GET /health` - Health check

## Application Processing

Applications are processed asynchronously:
1. Submit application → Returns ID and PENDING status
2. System processes application with partner banks (5-30 seconds)
3. Check status → Returns complete results with offers

## Further considerations

For a production ready solution:
- **Code quality**: Clean up handler and improve error handling.
- **Monitoring**: Set up monitoring tools to track system performance and detect issues early. Add correlation IDs for tracing.
- **Tests**: Expand API tests to cover other cases.

## Note 06.07.2025

Added a proper submission processing flow that runs every 5 minutes by default.
