SPDX-License-Identifier: CC-BY-4.0

Modify files in a directory, that share a common file extension, to add an SPDX
short-form license identifier:
- based on identifier passed from command line
- based on identifier indicated in a separate SPDX file

Prior pass, scan for short-form identifiers using github.com/spdx/tools-golang
and _do not_ add to files where the same identifier is already present.

Fail with error (and _do not_ write anything) where some files already have an
identifier, but it differs from the requested new license ID.

Optionally, also create and file Git commit in a new branch with those
added license IDs.

Depending on file type, use the appropriate comment syntax!

DO NOT add on first line where that would cause functional issues.
- e.g. shell scripts, or really any file starting with `#!`, skip to next line
- will modifying the first line cause linting problems in Golang files? others?

Handle binary files correctly (most likely, just refuse to handle them unless
called with a -f force flag)

Allow configuration (global and/or per-run) to select which comment styles will
be used for which types of files, e.g.:
```
{
    "filetypes": {
        ".c": {
            "comment": "/* SPDX */"
        },
        ".py": {
            "comment": "# SPDX"
        },
    },
    "skip": [".jpg", ".png", ...]
}
```

Eventually, for binary or non-commentable files with structured metadata
formats, insert into those as well
- into appropriate fields for image files with a header field for license or comments
- into the license field of a package.json file
    - NOTE: this wouldn't be an SPDX-License-Identifier: tag, so maybe don't do this
