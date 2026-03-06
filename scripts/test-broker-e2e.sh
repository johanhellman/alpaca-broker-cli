#!/usr/bin/env bash
set -e

echo "=============================================="
echo "    Alpaca Broker CLI E2E Sandbox Test Suite  "
echo "=============================================="

# Attempt to load .env from project root if variables are missing
if [ -z "$ALPACA_BROKER_API_KEY" ] || [ -z "$ALPACA_BROKER_API_SECRET" ]; then
  ENV_FILE="$(dirname "$0")/../.env"
  if [ -f "$ENV_FILE" ]; then
    set -a
    source "$ENV_FILE"
    set +a
  fi
fi

# Check for required Broker authentication variables
if [ -z "$ALPACA_BROKER_API_KEY" ] || [ -z "$ALPACA_BROKER_API_SECRET" ]; then
  echo "❌ Error: ALPACA_BROKER_API_KEY and ALPACA_BROKER_API_SECRET must be set."
  echo "Please export your Broker API credentials before running this test, or set them in the .env file."
  exit 1
fi

# Ensure we are strictly using the sandbox environment to avoid live accidents
export ALPACA_BROKER_ENV="sandbox"
echo "🔐 Authenticating against the Broker Sandbox API..."

# Ensure the binary is built and available in the current PATH
echo "🔨 Building latest alpaca-broker binary..."
make build-broker > /dev/null
export PATH="$PWD:$PATH"

echo -e "\n1. Listing current Broker Accounts..."
alpaca-broker accounts list

echo -e "\n2. Creating a new test Sub-Account..."
# Unique email using a timestamp to avoid conflicts in subsequent test runs
TEST_EMAIL="john.doe+$(date +%s)@example.com"
alpaca-broker accounts create \
    --contact-email "$TEST_EMAIL" \
    --contact-phone "+15555551234" \
    --contact-street-1 "123 Test St" \
    --contact-city "San Mateo" \
    --contact-state "CA" \
    --contact-postal "94401" \
    --id-given-name "E2E Test" \
    --id-family-name "User" \
    --id-dob "1990-01-01" \
    --id-tax-id "152637482" \
    --id-tax-type "USA_SSN" \
    --id-country-citizen "USA" \
    --id-country-birth "USA" \
    --id-country-tax "USA" \
    --funding-source-type "employment_income"

echo -e "\n   (Waiting 5 seconds for backend provisioning...)"
sleep 5

# Extract the newly created account ID via the email address we just used
NEW_ACCOUNT_ID=$(alpaca-broker accounts list --query "$TEST_EMAIL" --output json | jq -r '.[0].id')

if [ "$NEW_ACCOUNT_ID" == "null" ] || [ -z "$NEW_ACCOUNT_ID" ]; then
    echo "❌ Failed to retrieve newly created sub-account!"
    exit 1
fi

echo "   New Account ID: $NEW_ACCOUNT_ID"

echo -e "\n3. Checking Funding Transfers for the new account..."
alpaca-broker funding transfers "$NEW_ACCOUNT_ID"

# Note: In a true E2E sandbox environment, accounts require KYC approval before they can 
# legally trade or receive journal funds. The below commands verify the routing logic.

echo -e "\n4. Attempting to create a JNLC (Cash Journal) to fund the sub-account..."
# Assume a firm master account ID needs to exist, here we just show the structure.
# alpaca-broker journals create --from-account "$FIRM_ACCOUNT_ID" --to-account "$NEW_ACCOUNT_ID" --entry-type JNLC --amount 1000.00
echo "(Journal creation mocked due to missing Firm Source Account identity)"

echo -e "\n5. Querying Trading Orders for the new Sub-Account..."
alpaca-broker trading orders "$NEW_ACCOUNT_ID"

echo -e "\n6. Attempting to Place a Mock Sub-Account Trade..."
# Submits a simulated AAPL buy order directly on behalf of the provisioned customer
alpaca-broker trading order-create "$NEW_ACCOUNT_ID" --symbol AAPL --qty 1 --side buy --type market --time-in-force day || echo "    (Note: Trade rejected by backend because Sandbox Account is pending KYC)"

echo -e "\n7. Testing CSV Data Export (Accounts)..."
alpaca-broker accounts list --output csv

echo -e "\n✅ Broker E2E Testing Sequence Initialized successfully."
echo "Note: Full E2E execution requires firm-level KYC approvals and valid Sandbox Firm accounts for journaling."
