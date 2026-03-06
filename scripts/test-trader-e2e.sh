#!/usr/bin/env bash
set -e

echo "=============================================="
echo "    Alpaca Trader CLI E2E Paper Test Suite    "
echo "=============================================="

# Attempt to load .env from project root if variables are missing
if [ -z "$APCA_API_KEY_ID" ] || [ -z "$APCA_API_SECRET_KEY" ]; then
  ENV_FILE="$(dirname "$0")/../.env"
  if [ -f "$ENV_FILE" ]; then
    set -a
    source "$ENV_FILE"
    set +a
  fi
fi

# Check for required authentication variables
if [ -z "$APCA_API_KEY_ID" ] || [ -z "$APCA_API_SECRET_KEY" ]; then
  echo "❌ Error: APCA_API_KEY_ID and APCA_API_SECRET_KEY must be set."
  echo "Please export your Paper API credentials before running this test, or set them in the .env file."
  exit 1
fi

# Ensure we are strictly using the paper environment to avoid live trading accidents
export APCA_ENV="paper"
echo "🔐 Authenticating against the Paper API..."

# Ensure the binary is built and available in the current PATH
echo "🔨 Building latest alpaca-trader binary..."
make build-trader > /dev/null
export PATH="$PWD:$PATH"

echo -e "\n1. Fetching Account Details..."
alpaca-trader account get

echo -e "\n2. Verifying Asset Tradability (AAPL)..."
alpaca-trader assets get AAPL

echo -e "\n3. Fetching Market Data Quotes (AAPL)..."
alpaca-trader market-data quotes AAPL --total-limit 1

echo -e "\n4. Creating a Mock Watchlist..."
# Generating a random suffix to avoid unique constraint collisions
LIST_NAME="E2E-Test-List-$RANDOM"
alpaca-trader watchlists create --name "$LIST_NAME" --symbols AAPL,GOOG
# Parse the ID of the new watchlist for later cleanup
WATCHLIST_ID=$(alpaca-trader watchlists list --query '0.id')

echo -e "\n5. Submitting Orders: Market (AAPL), Limit (NVDA), and Fractional (MSFT)..."
# Buy 1 share of AAPL at Market
alpaca-trader orders create --symbol AAPL --qty 1 --side buy --type market --time-in-force day
# Buy 1 share of NVDA at a strict Limit
alpaca-trader orders create --symbol NVDA --qty 1 --side buy --type limit --limit-price 10.00 --time-in-force day 
# Buy $50 USD of MSFT Fractionally
alpaca-trader orders create --symbol MSFT --notional 50 --side buy --type market --time-in-force day

echo -e "\n6. Querying Open Orders..."
alpaca-trader orders list

echo -e "\n7. Waiting 10 seconds for fills..."
sleep 10

echo -e "\n8. Querying Current Positions..."
alpaca-trader positions list

echo -e "\n9. Cleaning up test artifacts (Watchlists, Open Limit Orders)..."
alpaca-trader watchlists delete "$WATCHLIST_ID"
# Using a complex jq-like query to find and cancel our unrealistic $10.00 NVDA limit order
NVDA_ORDER=$(alpaca-trader orders list --query '#(symbol=="NVDA").id')
if [ "$NVDA_ORDER" != "null" ]; then
  # The SDK doesn't natively expose order cancellation via CLI yet, so we just acknowledge it
  echo "    (Note: NVDA $10 limit order remains Open in paper ecosystem)"
fi

echo -e "\n10. Viewing Account Portfolio Value change..."
alpaca-trader account get --query "PortfolioValue"

echo -e "\n11. Testing CSV Data Export (Positions)..."
alpaca-trader positions list --output csv

echo -e "\n✅ E2E Testing Complete! "
echo "Note: 'positions close-all' is bypassed by default to prevent triggering Alpaca's Paper API Pattern-Day-Trader (PDT) wash-trade protections."
