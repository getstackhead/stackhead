import re

# Regular expression taken from
# https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
SEMVER_RE = re.compile(
    r"""
    ^
        (?P<major>0|[1-9]\d*)
        \.
        (?P<minor>0|[1-9]\d*)
        \.
        (?P<patch>0|[1-9]\d*)
        (?:
            -
            (?P<prerelease>
                (?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)
                (?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*
            )
        )?
        (?:
            \+
            (?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*)
        )?
    $
    """,
    flags=re.X,
)


# Check if currentVersion is within targetVersion
def within_major(currentVersion, targetVersion):
    matchCurrent = SEMVER_RE.match(currentVersion)
    if not matchCurrent:
        raise ValueError("invalid semantic version '%s'" % currentVersion)
    matchTarget = SEMVER_RE.match(targetVersion)
    if not matchTarget:
        raise ValueError("invalid semantic version '%s'" % targetVersion)

    (major, minor, patch, _, _) = matchCurrent.group(1, 2, 3, 4, 5)
    (majorTarget, minorTarget, patchTarget, _, _) = matchTarget.group(1, 2, 3, 4, 5)

    return int(major) == int(majorTarget)


class TestModule(object):
    """
    Jinja2 tests for version compare
    """

    def tests(self):
        return {"withinMajorVersion": within_major}
