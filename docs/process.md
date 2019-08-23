SPDX-License-Identifier: CC-BY-4.0

General process:
1) load configuration
2) determine set of selected files to check
3) use spdx/tools-golang to get existing IDs for selected files
4) determine which IDs should be added to selected files, either:
   1) based on config per file type
   2) based on command line flag
   3) based on imported SPDX file, looking at concluded license
   4) based on imported SPDX file, looking at concluded license for file if
      present and declared package license if not present
5) prepare dry run action to take:
   1) skip: if desired ID matches already-present ID
   2) error: if desired ID does not match already-present ID
   3) error: if no read permission or no write permission for file
   4) add: if no ID is currently present
6) stop if dry run only and/or if any errors detected
7) walk through and update files; for each:
   1) open file for reading; read all lines into slice of strings
   2) determine where to add new line:
      - if skipIfFirstPrefix is set and prefix matches, line = 1
      - otherwise, line = 0
   3) reopen file for writing
   4) walk back through slice; when hitting line number, write ID based on
      desired format for that extension
   5) close file

To be added:
- as part of step 1, load Git configuration
- between steps 4 and 5, create new branch
- after step 7, add to staging and commit

Assumptions:
- contents of file tree will not otherwise change while spdxify is running
