# Postman Collections

This folder contains Postman collections for testing the Astras API services locally using AWS SAM CLI.

## Collections

### 📧 `caregiver_service.json`
Complete collection for the Caregiver Service API including:

**CRUD Operations:**
- `GET /caregivers` - Get all caregivers
- `GET /caregivers/{id}` - Get caregiver by ID  
- `POST /caregivers` - Create new caregiver
- `PUT /caregivers/{id}` - Update caregiver
- `DELETE /caregivers/{id}` - Delete caregiver

**Validation Endpoints:**
- `POST /validate/email` - Validate email addresses (with valid/invalid examples)
- `POST /validate/relationship` - Validate relationships (with valid options)

### 👶 `kid_service.json` 
Complete collection for the Kid Service API including:

**CRUD Operations:**
- `GET /kids` - Get all kids
- `GET /kids/{id}` - Get kid by ID
- `POST /kids` - Create new kid  
- `PUT /kids/{id}` - Update kid
- `DELETE /kids/{id}` - Delete kid

### ⭐ `star_service.json`
Complete collection for the Star Transaction Service API including:

**CRUD Operations:**
- `GET /transactions` - Get all star transactions
- `GET /transactions/{id}` - Get transaction by ID
- `POST /transactions` - Create new transaction (earn/spend stars)
- `PUT /transactions/{id}` - Update transaction
- `DELETE /transactions/{id}` - Delete transaction

**Validation Endpoints:**
- `POST /validate/type` - Validate transaction type (earn/spend with case-insensitive examples)
- `POST /validate/amount` - Validate star amounts (1-100 range validation)

## Usage

### Prerequisites
1. **Start SAM Local API:**
   ```bash
   sam local start-api
   ```
   
2. **Import Collections:**
   - Open Postman
   - Click "Import"
   - Select both JSON files from this folder
   - Collections will appear in your Postman workspace

### Testing
- All requests are pre-configured to use `http://127.0.0.1:3000`
- Collections include both success and error test cases
- Validation endpoints include examples of valid and invalid data

### Features
- **Environment Variables:** Base URL configured as collection variable
- **Test Examples:** Multiple request examples for different scenarios  
- **Validation Testing:** Real-time validation endpoints for frontend integration
- **Error Cases:** Examples of invalid requests and expected error responses

## Local Development
Make sure you have:
1. Built the services: `mage build:kid && mage build:caregiver && mage build:star`
2. Started SAM local: `sam local start-api`
3. Services running on http://127.0.0.1:3000

## API Documentation
For detailed API documentation, see the main project README and service-specific files in `cmd/` directories.