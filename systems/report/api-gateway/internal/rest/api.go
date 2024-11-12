/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type GetReportRequest struct {
	ReportId string `example:"{{ReportUUID}}" path:"report_id" validate:"required"`
	AsPdf    bool   `form:"as_pdf" json:"as_pdf" query:"as_pdf" binding:"required"`
}
