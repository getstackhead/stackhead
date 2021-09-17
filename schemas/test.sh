#!/bin/sh

TEST_FAILED=0

validate_files() {
	SCHEMAFILE="./${1}.schema.json"
	EXAMPLE_DIR="examples/${1}"
	printf "### TESTING %s ... START\n" "${SCHEMAFILE}"
	printf "Valid %s files should be valid:\n" "${EXAMPLE_DIR}"
	for filename in "${EXAMPLE_DIR}"/valid/*; do
		if bin/jsonschema-validator "${SCHEMAFILE}" "${filename}" 1>/dev/null; then
			: # file is valid
		else
			TEST_FAILED=1
			printf "ERROR: File %s is invalid." "${filename}" >/dev/stderr
		fi
	done
	printf "\n\n"
	# Invalid files should be invalid
	printf "Invalid %s files should be invalid.\n" "${EXAMPLE_DIR}"
	for filename in "${EXAMPLE_DIR}"/invalid/*; do
		if bin/jsonschema-validator "${SCHEMAFILE}" "${filename}" 2>/dev/null; then
			TEST_FAILED=1
			printf "ERROR: %s is valid." "${filename}" >/dev/stderr
		fi
	done
	printf "### TESTING %s ... END\n" "${SCHEMAFILE}"
}

validate_files "cli-config"

if [ $TEST_FAILED -eq 0 ]; then
	echo 'All tests succeeded.'
else
	echo 'Some tests failed!'
fi

exit $TEST_FAILED
