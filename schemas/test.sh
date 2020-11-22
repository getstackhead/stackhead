#!/bin/sh

TEST_FAILED=0

validate_files()
{
  SCHEMAFILE="./${1}.schema.json"
  EXAMPLE_DIR="examples/${1}"
  echo "### TESTING ${SCHEMAFILE} ... START\n"
	printf "Valid %s files should be valid:\n" EXAMPLE_DIR
	for filename in "${EXAMPLE_DIR}"/valid/*; do
		if bin/jsonschema-validator "${SCHEMAFILE}" "${filename}"; then
		: # file is valid
		else
			TEST_FAILED=1
			echo "File ${filename} is invalid." > /dev/stderr
		fi
	done
	printf "\n\n"
	# Invalid files should be invalid
	printf "Invalid %s files should be invalid:\n" EXAMPLE_DIR
	for filename in "${EXAMPLE_DIR}"/invalid/*; do
		if bin/jsonschema-validator "${SCHEMAFILE}" "${filename}"; then
			TEST_FAILED=1
			echo "File ${filename} is valid." > /dev/stderr
		fi
	done
  echo "### TESTING ${SCHEMAFILE} ... END\n"
}

validate_files "project-definition"

if [ $TEST_FAILED -eq 0 ]; then
  echo 'All tests succeeded.'
else
  echo 'Some tests failed!'
fi

exit $TEST_FAILED
