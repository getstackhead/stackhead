#!/bin/sh

TEST_FAILED=0

# Valid files should be valid
echo "Valid files should be valid:\n"
for filename in examples/valid/*; do
  if bin/project-validator ${filename}; then
  : # file is valid
  else
    TEST_FAILED=1
    echo "File ${filename} is invalid." > /dev/stderr
  fi
done

echo "\n\n"
# Invalid files should be invalid
echo "Invalid files should be invalid:\n"
for filename in examples/invalid/*; do
  if bin/project-validator ${filename}; then
    TEST_FAILED=1
    echo "File ${filename} is valid." > /dev/stderr
  fi
done

if [ $TEST_FAILED -eq 0 ]; then
  echo 'All tests succeeded.'
else
  echo 'Some tests failed!'
fi

exit $TEST_FAILED