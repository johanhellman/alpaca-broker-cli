#!/usr/bin/env bash
set -e

echo "=============================================="
echo "    Alpaca Trader CLI E2E Paper Test Suite    "
echo "=============================================="

# Check for required authentication variables
if [ -z "$APCA_API_KEY_ID" ] || [ -z "$APCA_API_SECRET_KEY" ]; then
  echo "❌ Error: APCA_API_KEY_ID and APCA_API_SECRET_KEY must be set."
  echo "Please export your Paper API credentials before running this test."
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

echo -e "\n2. Submitting Market Order (1 Share AAPL)..."
alpaca-trader orders create --symbol AAPL --qty 1 --side buy --type market --time-in-force day

echo -e "\n3. Submitting Fractional Order (\$100 Notional NVDA)..."
alpaca-trader orders create --symbol NVDA --notional 100 --side buy --type market --time-in-force day

echo -e "\n4. Querying Open Orders..."
alpaca-trader orders list

echo -e "\n5. Waiting 10 seconds for fills..."
sleep 10

echo -e "\n6. Querying Current Positions..."
alpaca-trader positions list

echo -e "\n7. Viewing Account Portfolio Value change..."
alpaca-trader account get --query "PortfolioValue"

echo -e "\n✅ E2E Testing Complete! "
echo "Note: 'positions close-all' is bypassed by default to prevent triggering Alpaca's Paper API Pattern-Day-Trader (PDT) wash-trade protections."
