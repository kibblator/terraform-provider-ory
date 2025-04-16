#!/bin/bash

# Find all Go test files with "TestAcc" in the name
tests=$(go test -list . ./... 2>/dev/null | grep -E '^TestAcc')

# Run each test individually
for test in $tests; do
  echo "Running $test..."
  go test -v -run $test ./...
done
