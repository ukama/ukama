/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package ingest

var EventRoutes = []string{
	"event.cloud.local.{{ .Org}}.registry.site.site.create",
	"event.cloud.local.{{ .Org}}.registry.site.site.update",
	"event.cloud.local.{{ .Org}}.registry.network.network.add",
	"event.cloud.local.{{ .Org}}.registry.node.node.create",
	"event.cloud.local.{{ .Org}}.registry.node.node.update",
	"event.cloud.local.{{ .Org}}.registry.node.node.assign",
	"event.cloud.local.{{ .Org}}.registry.node.node.release",
	"event.cloud.local.{{ .Org}}.messaging.mesh.node.online",
	"event.cloud.local.{{ .Org}}.messaging.mesh.node.offline",
	"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activate",
	"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate",
	"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.addpackage",
	"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activepackage",
	"event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.create",
	"event.cloud.local.{{ .Org}}.dataplan.package.package.create",
	"event.cloud.local.{{ .Org}}.dataplan.package.package.update",
	"event.cloud.local.{{ .Org}}.inventory.component.components.sync",
	"event.cloud.local.{{ .Org}}.node.health.report.store",
	"event.cloud.local.{{ .Org}}.payments.processor.payment.success",
	"event.cloud.local.{{ .Org}}.payments.processor.payment.failed",
	"event.cloud.local.{{ .Org}}.node.state.node.transition",
	"event.cloud.local.{{ .Org}}.billing.report.invoice.generate",
	"event.cloud.local.{{ .Org}}.webhooks.dispatcher.webhook.update",
}
