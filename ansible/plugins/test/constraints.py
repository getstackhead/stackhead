import re
import semantic_version


# Checks if the given `constraints` are met by comparing it with the StackHead configuration
def constraints_fulfilled(constraints, stackhead_config):
    extractFromConstraint = re.compile(r"^([\w.]+)([+><=\s]+(.*))?$")

    for constraint in constraints:
        match = extractFromConstraint.match(constraint)
        (constraint_name, constraint_range) = match.group(1, 2)
        if constraint_name != "stackhead":  # for now only support stackhead
            continue

        if constraint_range is None:
            # todo when extending this for module constraints: check if package exists
            continue

        # remove whitespace
        # @see https://stackoverflow.com/questions/8270092/remove-all-whitespace-in-a-string
        constraint_range = re.sub(r"\s+", "", constraint_range, flags=re.UNICODE)

        constraint_version = ""
        if constraint_name == "stackhead":
            constraint_version = stackhead_config["version"]["current"]

        if semantic_version.Version(
            constraint_version
        ) not in semantic_version.SimpleSpec(constraint_range):
            return False
    return True


class TestModule(object):
    """
    Jinja2 tests for version compare
    """

    def tests(self):
        return {"constraintsFulfilled": constraints_fulfilled}
