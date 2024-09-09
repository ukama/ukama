# #  This Source Code Form is subject to the terms of the Mozilla Public
# #  License, v. 2.0. If a copy of the MPL was not distributed with this
# #  file, You can obtain one at https://mozilla.org/MPL/2.0/.
# # 
# #  Copyright (c) 2023-present, Ukama Inc.
 
# # This script generates Go code for unmarshalling protobuf messages for events

# # To run this script, execute the following command:
# # python3 generate_unmarshall_func.py > ./gen/events/unmarshals.go
import os
import re
from io import StringIO

LICENSE_TEXT = """
/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
"""

PACKAGE_IMPORTS = """
package events

import (
    "google.golang.org/protobuf/types/known/anypb"
    "google.golang.org/protobuf/proto"
    log "github.com/sirupsen/logrus"
)
"""

def find_message_names(directory):
    message_names = set()
    for root, dirs, files in os.walk(directory):
        for file in files:
            file_path = os.path.join(root, file)
            try:
                with open(file_path, 'r') as f:
                    for line in f:
                        match = re.search(r'message\s+(\w+)\s*{', line)
                        if match:
                            message_name = match.group(1)
                            message_names.add(message_name)
            except (IOError, OSError) as e:
                print(f"Error reading file {file_path}: {e}")
    return message_names

def generate_go_code(name):
    return f"""func Unmarshal{name}(msg *anypb.Any, emsg string) (*{name}, error) {{
    p := &{name}{{
    }}
    err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{{AllowPartial: true, DiscardUnknown: true}})
    if err != nil {{
        log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
        return nil, err
    }}
    return p, nil
}}
"""

def clean_file(file_path):
    try:
        with open(file_path, 'w') as f:
            f.truncate(0)
    except (IOError, OSError) as e:
        print(f"Error cleaning file {file_path}: {e}")

def main():
    path = "./pb/events"
    output_file_path = "./pb/gen/events/unmarshals.go"

    clean_file(output_file_path)

    message_names = find_message_names(path)

    output = StringIO()
    output.write(LICENSE_TEXT)
    output.write(PACKAGE_IMPORTS)

    for name in message_names:
        output.write(generate_go_code(name))

    try:
        with open(output_file_path, 'w') as f:
            f.write(output.getvalue())
    except (IOError, OSError) as e:
        print(f"Error writing to file {output_file_path}: {e}")

if __name__ == "__main__":
    main()