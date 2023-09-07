// SPDX-License-Identifier: GPL-2.0+
/*
 *  EFI image loader
 *
 *  based partly on wine code
 *
 *  Copyright (c) 2016 Alexander Graf
 */

#include <common.h>
#include <efi_loader.h>
#include <pe.h>

const efi_guid_t efi_global_variable_guid = EFI_GLOBAL_VARIABLE_GUID;
const efi_guid_t efi_guid_device_path = DEVICE_PATH_GUID;
const efi_guid_t efi_guid_loaded_image = LOADED_IMAGE_GUID;
const efi_guid_t efi_simple_file_system_protocol_guid =
		EFI_SIMPLE_FILE_SYSTEM_PROTOCOL_GUID;
const efi_guid_t efi_file_info_guid = EFI_FILE_INFO_GUID;

static int machines[] = {
#if defined(CONFIG_ARM64)
	IMAGE_FILE_MACHINE_ARM64,
#elif defined(CONFIG_ARM)
	IMAGE_FILE_MACHINE_ARM,
	IMAGE_FILE_MACHINE_THUMB,
	IMAGE_FILE_MACHINE_ARMNT,
#endif

#if defined(CONFIG_X86_64)
	IMAGE_FILE_MACHINE_AMD64,
#elif defined(CONFIG_X86)
	IMAGE_FILE_MACHINE_I386,
#endif

#if defined(CONFIG_CPU_RISCV_32)
	IMAGE_FILE_MACHINE_RISCV32,
#endif

#if defined(CONFIG_CPU_RISCV_64)
	IMAGE_FILE_MACHINE_RISCV64,
#endif
	0 };

/*
 * Print information about a loaded image.
 *
 * If the program counter is located within the image the offset to the base
 * address is shown.
 *
 * @image:	loaded image
 * @pc:		program counter (use NULL to suppress offset output)
 * @return:	status code
 */
efi_status_t efi_print_image_info(struct efi_loaded_image *image, void *pc)
{
	if (!image)
		return EFI_INVALID_PARAMETER;
	printf("UEFI image");
	printf(" [0x%p:0x%p]",
	       image->reloc_base, image->reloc_base + image->reloc_size - 1);
	if (pc && pc >= image->reloc_base &&
	    pc < image->reloc_base + image->reloc_size)
		printf(" pc=0x%zx", pc - image->reloc_base);
	if (image->file_path)
		printf(" '%pD'", image->file_path);
	printf("\n");
	return EFI_SUCCESS;
}

/*
 * Print information about all loaded images.
 *
 * @pc:		program counter (use NULL to suppress offset output)
 */
void efi_print_image_infos(void *pc)
{
	struct efi_object *efiobj;
	struct efi_handler *handler;

	list_for_each_entry(efiobj, &efi_obj_list, link) {
		list_for_each_entry(handler, &efiobj->protocols, link) {
			if (!guidcmp(handler->guid, &efi_guid_loaded_image)) {
				efi_print_image_info(
					handler->protocol_interface, pc);
			}
		}
	}
}

static efi_status_t efi_loader_relocate(const IMAGE_BASE_RELOCATION *rel,
			unsigned long rel_size, void *efi_reloc,
			unsigned long pref_address)
{
	unsigned long delta = (unsigned long)efi_reloc - pref_address;
	const IMAGE_BASE_RELOCATION *end;
	int i;

	if (delta == 0)
		return EFI_SUCCESS;

	end = (const IMAGE_BASE_RELOCATION *)((const char *)rel + rel_size);
	while (rel < end - 1 && rel->SizeOfBlock) {
		const uint16_t *relocs = (const uint16_t *)(rel + 1);
		i = (rel->SizeOfBlock - sizeof(*rel)) / sizeof(uint16_t);
		while (i--) {
			uint32_t offset = (uint32_t)(*relocs & 0xfff) +
					  rel->VirtualAddress;
			int type = *relocs >> EFI_PAGE_SHIFT;
			uint64_t *x64 = efi_reloc + offset;
			uint32_t *x32 = efi_reloc + offset;
			uint16_t *x16 = efi_reloc + offset;

			switch (type) {
			case IMAGE_REL_BASED_ABSOLUTE:
				break;
			case IMAGE_REL_BASED_HIGH:
				*x16 += ((uint32_t)delta) >> 16;
				break;
			case IMAGE_REL_BASED_LOW:
				*x16 += (uint16_t)delta;
				break;
			case IMAGE_REL_BASED_HIGHLOW:
				*x32 += (uint32_t)delta;
				break;
			case IMAGE_REL_BASED_DIR64:
				*x64 += (uint64_t)delta;
				break;
#ifdef __riscv
			case IMAGE_REL_BASED_RISCV_HI20:
				*x32 = ((*x32 & 0xfffff000) + (uint32_t)delta) |
					(*x32 & 0x00000fff);
				break;
			case IMAGE_REL_BASED_RISCV_LOW12I:
			case IMAGE_REL_BASED_RISCV_LOW12S:
				/* We know that we're 4k aligned */
				if (delta & 0xfff) {
					printf("Unsupported reloc offset\n");
					return EFI_LOAD_ERROR;
				}
				break;
#endif
			default:
				printf("Unknown Relocation off %x type %x\n",
				       offset, type);
				return EFI_LOAD_ERROR;
			}
			relocs++;
		}
		rel = (const IMAGE_BASE_RELOCATION *)relocs;
	}
	return EFI_SUCCESS;
}

void __weak invalidate_icache_all(void)
{
	/* If the system doesn't support icache_all flush, cross our fingers */
}

/*
 * Determine the memory types to be used for code and data.
 *
 * @loaded_image_info	image descriptor
 * @image_type		field Subsystem of the optional header for
 *			Windows specific field
 */
static void efi_set_code_and_data_type(
			struct efi_loaded_image *loaded_image_info,
			uint16_t image_type)
{
	switch (image_type) {
	case IMAGE_SUBSYSTEM_EFI_APPLICATION:
		loaded_image_info->image_code_type = EFI_LOADER_CODE;
		loaded_image_info->image_data_type = EFI_LOADER_DATA;
		break;
	case IMAGE_SUBSYSTEM_EFI_BOOT_SERVICE_DRIVER:
		loaded_image_info->image_code_type = EFI_BOOT_SERVICES_CODE;
		loaded_image_info->image_data_type = EFI_BOOT_SERVICES_DATA;
		break;
	case IMAGE_SUBSYSTEM_EFI_RUNTIME_DRIVER:
	case IMAGE_SUBSYSTEM_EFI_ROM:
		loaded_image_info->image_code_type = EFI_RUNTIME_SERVICES_CODE;
		loaded_image_info->image_data_type = EFI_RUNTIME_SERVICES_DATA;
		break;
	default:
		printf("%s: invalid image type: %u\n", __func__, image_type);
		/* Let's assume it is an application */
		loaded_image_info->image_code_type = EFI_LOADER_CODE;
		loaded_image_info->image_data_type = EFI_LOADER_DATA;
		break;
	}
}

/*
 * This function loads all sections from a PE binary into a newly reserved
 * piece of memory. On successful load it then returns the entry point for
 * the binary. Otherwise NULL.
 */
void *efi_load_pe(void *efi, struct efi_loaded_image *loaded_image_info)
{
	IMAGE_NT_HEADERS32 *nt;
	IMAGE_DOS_HEADER *dos;
	IMAGE_SECTION_HEADER *sections;
	int num_sections;
	void *efi_reloc;
	int i;
	const IMAGE_BASE_RELOCATION *rel;
	unsigned long rel_size;
	int rel_idx = IMAGE_DIRECTORY_ENTRY_BASERELOC;
	void *entry;
	uint64_t image_base;
	uint64_t image_size;
	unsigned long virt_size = 0;
	int supported = 0;

	dos = efi;
	if (dos->e_magic != IMAGE_DOS_SIGNATURE) {
		printf("%s: Invalid DOS Signature\n", __func__);
		return NULL;
	}

	nt = (void *) ((char *)efi + dos->e_lfanew);
	if (nt->Signature != IMAGE_NT_SIGNATURE) {
		printf("%s: Invalid NT Signature\n", __func__);
		return NULL;
	}

	for (i = 0; machines[i]; i++)
		if (machines[i] == nt->FileHeader.Machine) {
			supported = 1;
			break;
		}

	if (!supported) {
		printf("%s: Machine type 0x%04x is not supported\n",
		       __func__, nt->FileHeader.Machine);
		return NULL;
	}

	/* Calculate upper virtual address boundary */
	num_sections = nt->FileHeader.NumberOfSections;
	sections = (void *)&nt->OptionalHeader +
			    nt->FileHeader.SizeOfOptionalHeader;

	for (i = num_sections - 1; i >= 0; i--) {
		IMAGE_SECTION_HEADER *sec = &sections[i];
		virt_size = max_t(unsigned long, virt_size,
				  sec->VirtualAddress + sec->Misc.VirtualSize);
	}

	/* Read 32/64bit specific header bits */
	if (nt->OptionalHeader.Magic == IMAGE_NT_OPTIONAL_HDR64_MAGIC) {
		IMAGE_NT_HEADERS64 *nt64 = (void *)nt;
		IMAGE_OPTIONAL_HEADER64 *opt = &nt64->OptionalHeader;
		image_base = opt->ImageBase;
		image_size = opt->SizeOfImage;
		efi_set_code_and_data_type(loaded_image_info, opt->Subsystem);
		efi_reloc = efi_alloc(virt_size,
				      loaded_image_info->image_code_type);
		if (!efi_reloc) {
			printf("%s: Could not allocate %lu bytes\n",
			       __func__, virt_size);
			return NULL;
		}
		entry = efi_reloc + opt->AddressOfEntryPoint;
		rel_size = opt->DataDirectory[rel_idx].Size;
		rel = efi_reloc + opt->DataDirectory[rel_idx].VirtualAddress;
		virt_size = ALIGN(virt_size, opt->SectionAlignment);
	} else if (nt->OptionalHeader.Magic == IMAGE_NT_OPTIONAL_HDR32_MAGIC) {
		IMAGE_OPTIONAL_HEADER32 *opt = &nt->OptionalHeader;
		image_base = opt->ImageBase;
		image_size = opt->SizeOfImage;
		efi_set_code_and_data_type(loaded_image_info, opt->Subsystem);
		efi_reloc = efi_alloc(virt_size,
				      loaded_image_info->image_code_type);
		if (!efi_reloc) {
			printf("%s: Could not allocate %lu bytes\n",
			       __func__, virt_size);
			return NULL;
		}
		entry = efi_reloc + opt->AddressOfEntryPoint;
		rel_size = opt->DataDirectory[rel_idx].Size;
		rel = efi_reloc + opt->DataDirectory[rel_idx].VirtualAddress;
		virt_size = ALIGN(virt_size, opt->SectionAlignment);
	} else {
		printf("%s: Invalid optional header magic %x\n", __func__,
		       nt->OptionalHeader.Magic);
		return NULL;
	}

	/* Load sections into RAM */
	for (i = num_sections - 1; i >= 0; i--) {
		IMAGE_SECTION_HEADER *sec = &sections[i];
		memset(efi_reloc + sec->VirtualAddress, 0,
		       sec->Misc.VirtualSize);
		memcpy(efi_reloc + sec->VirtualAddress,
		       efi + sec->PointerToRawData,
		       sec->SizeOfRawData);
	}

	/* Run through relocations */
	if (efi_loader_relocate(rel, rel_size, efi_reloc,
				(unsigned long)image_base) != EFI_SUCCESS) {
		efi_free_pages((uintptr_t) efi_reloc,
			       (virt_size + EFI_PAGE_MASK) >> EFI_PAGE_SHIFT);
		return NULL;
	}

	/* Flush cache */
	flush_cache((ulong)efi_reloc,
		    ALIGN(virt_size, EFI_CACHELINE_SIZE));
	invalidate_icache_all();

	/* Populate the loaded image interface bits */
	loaded_image_info->image_base = efi;
	loaded_image_info->image_size = image_size;
	loaded_image_info->reloc_base = efi_reloc;
	loaded_image_info->reloc_size = virt_size;

	return entry;
}
