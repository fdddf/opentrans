#!/bin/bash

# env key


# Test script for opentrans

go build -o opentrans

echo "=== opentrans Test Script ==="
echo

# Check if binary exists
if [ ! -f "opentrans" ]; then
    echo "Error: opentrans binary not found. Please build it first."
    exit 1
fi

# Test help command
echo "1. Testing help command..."
./opentrans --help
echo

# Test version (if implemented)
echo "2. Testing version command..."
./opentrans -v 2>/dev/null || echo "Version command not implemented"
echo

# Test provider list
echo "3. Testing provider commands..."
for provider in google deepl baidu openai; do
    echo "Testing $provider provider help..."
    ./opentrans $provider --help | head -10
    echo "----------------------------------------"
done

# Test with example file
echo "4. Testing with example file..."
echo "Example file content:"
cat example.xcstrings | jq '.strings | keys'
echo

# Test dry run (without actual translation)
echo "5. Testing dry run with Google provider (simulated)..."
./opentrans google \
    --api-key "test-key" \
    --input "example.xcstrings" \
    --output "example_translated.xcstrings" \
    --target-languages "zh-Hans" \
    --verbose 2>&1 | grep -E "(Loading|Found|strings|Exiting)"

echo
echo "=== Test completed ==="