import setuptools

with open("validation/README.md", "r") as fh:
    long_description = fh.read()

setuptools.setup(
    name="stackhead-project_validator",
    version="1.0.1",
    author="Mario Lubenka",
    author_email="me@saitho.me",
    description="Validate your StackHead project definition files",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/getstackhead/project-validator",
    packages=['stackhead_project_validator'],
    package_data={'stackhead_project_validator': ['validation/bin/*', 'validation/schema/*']},
    classifiers=[
        "Environment :: Console",
        "Intended Audience :: Developers",
        "Intended Audience :: System Administrators",
        "License :: OSI Approved :: GNU General Public License v2 or later (GPLv2+)",
        "Operating System :: OS Independent",
        "Topic :: Software Development :: Testing",
        "Topic :: Software Development :: Quality Assurance",
    ]
)
