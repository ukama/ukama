# #  This Source Code Form is subject to the terms of the Mozilla Public
# #  License, v. 2.0. If a copy of the MPL was not distributed with this
# #  file, You can obtain one at https://mozilla.org/MPL/2.0/.
# # 
# #  Copyright (c) 2023-present, Ukama Inc.
 
# # This script generates Go code for unmarshalling protobuf messages for events

# # To run this script, execute the following command:
# # python3 generate_unmarshall_func.py > ./gen/events/unmarshals.go

# import os
# import re
# import sys

# # Saving the reference of the standard output
# original_stdout = sys.stdout 

# def find_message_names(directory):
#     message_names = set()
#     for root, dirs, files in os.walk(directory):
#         for file in files:
#             file_path = os.path.join(root, file)
#             with open(file_path, 'r') as f:
#                 for line in f:
#                     match = re.search(r'message\s+(\w+)\s*{', line)
#                     if match:
#                         message_name = match.group(1)
#                         message_names.add(message_name)
#     return message_names

# def generate_go_code(name):
#     print(f"func Unmarshal{name}(msg *anypb.Any, emsg string) (*{name}, error) {{")
#     print(f"  p := &{name}" + "{}")
#     print("  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})")
#     print("  if err != nil {")
#     print("    log.Errorf(\"%s : %+v. Error %s.\", emsg, msg, err.Error())")
#     print("    return nil, err")
#     print("  }")
#     print("  return p, nil")
#     print("}")
#     print("")
    
# def clean_file(file_path):
#     with open(file_path, 'w') as f:
#         f.truncate(0)

# path = "./pb/events"
# gfile ="./pb/gen/events/unmarshals.go"
# clean_file(gfile)

# message_names = find_message_names(path)
# with open(gfile, 'w') as of:
#         sys.stdout = of
#         print("""
#         /*
#         * This Source Code Form is subject to the terms of the Mozilla Public
#         * License, v. 2.0. If a copy of the MPL was not distributed with this
#         * file, You can obtain one at https://mozilla.org/MPL/2.0/.
#         *
#         * Copyright (c) 2023-present, Ukama Inc.
#         */
#         """)
#         print("package events")
#         print("import (")
#         print("\"google.golang.org/protobuf/types/known/anypb\"")
#         print("\"google.golang.org/protobuf/proto\"")
#         print("log \"github.com/sirupsen/logrus\"")
#         print(")")

#         for name in message_names:
#             generate_go_code(name)
            
#         # Reset the standard output
#         sys.stdout = original_stdout  

#  This Source Code Form is subject to the terms of the Mozilla Public
#  License, v. 2.0. If a copy of the MPL was not distributed with this
#  file, You can obtain one at https://mozilla.org/MPL/2.0/.
# 
#  Copyright (c) 2023-present, Ukama Inc.
 
# This script generates Go code for unmarshalling protobuf messages for events

# To run this script, execute the following command:
# python3 generate_unmarshall_func.py > ./gen/events/unmarshals.go

import os
import re


def find_message_names(directory):
    message_names = set()
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith('.proto'):
                file_path = os.path.join(root, file)
                with open(file_path, 'r') as f:
                    for line in f:
                        match = re.search(r'message\s+(\w+)\s*{', line)
                        if match:
                            message_name = match.group(1)
                            message_names.add(message_name)
    return message_names

def generate_go_code(name, output_file):
    output_file.write(f"func Unmarshal{name}(msg *anypb.Any, emsg string) (*{name}, error) {{\n")
    output_file.write(f"\tp := &{name}" + "{}\n")
    output_file.write("\terr := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})\n")
    output_file.write("\tif err != nil {\n")
    output_file.write("\t\tlog.Errorf(\"%s : %+v. Error %s.\", emsg, msg, err.Error())\n")
    output_file.write("\t\treturn nil, err\n")
    output_file.write("\t}\n")
    output_file.write("\treturn p, nil\n")
    output_file.write("}\n\n")

# Get the script's directory to determine correct paths
script_dir = os.path.dirname(os.path.abspath(__file__))
pb_dir = os.path.dirname(script_dir) if os.path.basename(script_dir) == 'pb' else script_dir

# Paths relative to where script is located (pb directory)
events_path = os.path.join(script_dir, "events")
output_file_path = os.path.join(script_dir, "gen", "events", "unmarshals.go")

# Ensure the output directory exists
output_dir = os.path.dirname(output_file_path)
os.makedirs(output_dir, exist_ok=True)

# Find all message names from proto files
message_names = find_message_names(events_path)

# Write the generated code to the file
with open(output_file_path, 'w') as f:
    f.write("""/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
""")
    f.write("package events\n\n")
    f.write("import (\n")
    f.write("\t\"google.golang.org/protobuf/types/known/anypb\"\n")
    f.write("\t\"google.golang.org/protobuf/proto\"\n")
    f.write("\tlog \"github.com/sirupsen/logrus\"\n")
    f.write(")\n\n")
    
    for name in sorted(message_names):
        generate_go_code(name, f)