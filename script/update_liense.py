import os
import re
import sys
from pathlib import Path

# Original MIT License text (preserve original line widths)
MIT_LICENSE = """\
/**
 *MIT License
 *
 *Copyright (c) 2025 ylgeeker
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
**/

"""

# New Apache License text (preserve original line widths)
APACHE_LICENSE = """\
/**
 * Copyright 2025 saber authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/
"""

# Precompile regex pattern with raw string
MIT_PATTERN = re.compile(
    re.escape(MIT_LICENSE),
    re.MULTILINE | re.DOTALL
)

# Supported file extensions
SUPPORTED_EXTS = {'.go'}


def replace_license(file_path):
    """
    Replace MIT license with Apache 2.0 in a file.

    Args:
        file_path (str): Path to target file.
    Returns:
        bool: True if replacement occurred.
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()

        if MIT_PATTERN.search(content):
            new_content = MIT_PATTERN.sub(APACHE_LICENSE, content)
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(new_content)
            print(f"Replaced: {file_path}")
            return True
        else:
            print(f"Skipped: {file_path} (no MIT license)")
            return False
    except Exception as e:
        print(f"Error processing {file_path}: {str(e)}",
              file=sys.stderr)
        return False


def process_directory(dir_path):
    """
    Process all files in directory recursively.

    Args:
        dir_path (str): Directory to process.
    Returns:
        tuple: Total files checked, files replaced.
    """
    total = 0
    replaced = 0

    for root, _, files in os.walk(dir_path):
        for file in files:
            ext = Path(file).suffix.lower()
            if ext in SUPPORTED_EXTS:
                file_path = os.path.join(root, file)
                total += 1
                if replace_license(file_path):
                    replaced += 1

    return total, replaced


def main():
    """Main entry point for command line execution."""
    if len(sys.argv) < 2:
        print("Usage: python license_replacer.py <path>")
        sys.exit(1)

    target = sys.argv[1]

    if not os.path.exists(target):
        print(f"Error: Path '{target}' does not exist",
              file=sys.stderr)
        sys.exit(1)

    if os.path.isfile(target):
        ext = Path(target).suffix.lower()
        if ext not in SUPPORTED_EXTS:
            print(f"Error: File '{target}' is not C/C++",
                  file=sys.stderr)
            sys.exit(1)

        print(f"Processing single file: {target}")
        replaced = replace_license(target)
        status = "Replaced" if replaced else "Not replaced"
        print(f"Status: {status}")
    else:
        print(f"Processing directory: {target}")
        total, replaced = process_directory(target)
        print(f"Summary: Checked {total} files, "
              f"replaced {replaced} files")


if __name__ == "__main__":
    main()
