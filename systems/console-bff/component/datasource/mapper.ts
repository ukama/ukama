/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  ComponentAPIResDto,
  ComponentDto,
  ComponentsAPIResDto,
  ComponentsResDto,
} from "../resolvers/types";

export const dtoTocomponentsDto = (
  res: ComponentsAPIResDto
): ComponentsResDto => {
  const components: ComponentDto[] = [];
  res.components.forEach(component => {
    components.push({
      id: component.id,
      inventory: component.inventory,
      userId: component.user_id,
      type: component.type,
      description: component.description,
      datasheetUrl: component.datasheet_url,
      imageUrl: component.images_url,
      partNumber: component.part_number,
      manufacturer: component.manufacturer,
      managed: component.managed,
      warranty: component.warranty,
      specification: component.specification,
    });
  });
  return {
    components: components,
  };
};

export const dtoTocomponentDto = (res: ComponentAPIResDto): ComponentDto => {
  return {
    id: res.component.id,
    inventory: res.component.inventory,
    userId: res.component.user_id,
    type: res.component.type,
    description: res.component.description,
    datasheetUrl: res.component.datasheet_url,
    imageUrl: res.component.images_url,
    partNumber: res.component.part_number,
    manufacturer: res.component.manufacturer,
    managed: res.component.managed,
    warranty: res.component.warranty,
    specification: res.component.specification,
  };
};
