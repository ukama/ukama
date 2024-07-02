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
            file_path = os.path.join(root, file)
            with open(file_path, 'r') as f:
                for line in f:
                    match = re.search(r'message\s+(\w+)\s*{', line)
                    if match:
                        message_name = match.group(1)
                        message_names.add(message_name)
    return message_names

def generate_go_code(name):
    print(f"func Unmarshal{name}(msg *anypb.Any, emsg string) (*{name}, error) {{")
    print(f"  p := &{name}" + "{}")
    print("  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})")
    print("  if err != nil {")
    print("    log.Errorf(\"%s : %+v. Error %s.\", emsg, msg, err.Error())")
    print("    return nil, err")
    print("  }")
    print("  return p, nil")
    print("}")
    print("")
    
def clean_file(file_path):
    with open(file_path, 'w') as f:
        f.truncate(0)

path = "./events"
clean_file("./gen/events/unmarshals.go")

message_names = find_message_names(path)
print("""
/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
""")
print("package events")
print("import (")
print("\"google.golang.org/protobuf/types/known/anypb\"")
print("\"google.golang.org/protobuf/proto\"")
print("log \"github.com/sirupsen/logrus\"")
print(")")

for name in message_names:
    generate_go_code(name)