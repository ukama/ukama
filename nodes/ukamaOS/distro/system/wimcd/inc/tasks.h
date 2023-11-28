/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef WIMC_TASKS_H
#define WIMC_TASKS_H

void clear_tasks(WTasks **tasks);
void add_to_tasks(WTasks **tasks, WimcReq *req);
void delete_from_tasks(WTasks **tasks, WTasks *target);
char *process_task_request(WTasks *task);
char *process_cli_response(WRespType type, char *path, char *idStr,
                           WTasks *task, char *errStr); 

#endif /* WIMC_TASKS_H */
