import {
  ComponentAPIDto,
  ComponentDto,
  ComponentsAPIResDto,
  ComponentsResDto,
} from "../resolvers/types";

export const dtoToComponentsDto = (
  res: ComponentsAPIResDto
): ComponentsResDto => {
  const components: ComponentDto[] = [];
  res.components.forEach(component => {
    components.push({
      id: component.id,
      specification: component.specification,
      inventoryId: component.inventory_id,
      category: component.category,
      type: component.type,
      userId: component.user_id,
      description: component.description,
      datasheetUrl: component.datasheet_url,
      imageUrl: component.images_url,
      partNumber: component.part_number,
      manufacturer: component.manufacturer,
      managed: component.managed,
      warranty: component.warranty,
    });
  });
  return {
    components: components,
  };
};

export const dtoToComponentDto = (res: ComponentAPIDto): ComponentDto => {
  return {
    id: res.id,
    specification: res.specification,
    inventoryId: res.inventory_id,
    category: res.category,
    type: res.type,
    userId: res.user_id,
    description: res.description,
    datasheetUrl: res.datasheet_url,
    imageUrl: res.images_url,
    partNumber: res.part_number,
    manufacturer: res.manufacturer,
    managed: res.managed,
    warranty: res.warranty,
  };
};
