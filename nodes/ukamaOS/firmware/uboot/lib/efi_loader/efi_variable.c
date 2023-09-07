// SPDX-License-Identifier: GPL-2.0+
/*
 *  EFI utils
 *
 *  Copyright (c) 2017 Rob Clark
 */

#include <malloc.h>
#include <charset.h>
#include <efi_loader.h>

#define READ_ONLY BIT(31)

/*
 * Mapping between EFI variables and u-boot variables:
 *
 *   efi_$guid_$varname = {attributes}(type)value
 *
 * For example:
 *
 *   efi_8be4df61-93ca-11d2-aa0d-00e098032b8c_OsIndicationsSupported=
 *      "{ro,boot,run}(blob)0000000000000000"
 *   efi_8be4df61-93ca-11d2-aa0d-00e098032b8c_BootOrder=
 *      "(blob)00010000"
 *
 * The attributes are a comma separated list of these possible
 * attributes:
 *
 *   + ro   - read-only
 *   + boot - boot-services access
 *   + run  - runtime access
 *
 * NOTE: with current implementation, no variables are available after
 * ExitBootServices, and all are persisted (if possible).
 *
 * If not specified, the attributes default to "{boot}".
 *
 * The required type is one of:
 *
 *   + utf8 - raw utf8 string
 *   + blob - arbitrary length hex string
 *
 * Maybe a utf16 type would be useful to for a string value to be auto
 * converted to utf16?
 */

#define MAX_VAR_NAME 31
#define MAX_NATIVE_VAR_NAME \
	(strlen("efi_xxxxxxxx-xxxx-xxxx-xxxxxxxxxxxxxxxx_") + \
		(MAX_VAR_NAME * MAX_UTF8_PER_UTF16))

static int hex(int ch)
{
	if (ch >= 'a' && ch <= 'f')
		return ch-'a'+10;
	if (ch >= '0' && ch <= '9')
		return ch-'0';
	if (ch >= 'A' && ch <= 'F')
		return ch-'A'+10;
	return -1;
}

static int hex2mem(u8 *mem, const char *hexstr, int size)
{
	int nibble;
	int i;

	for (i = 0; i < size; i++) {
		if (*hexstr == '\0')
			break;

		nibble = hex(*hexstr);
		if (nibble < 0)
			return -1;

		*mem = nibble;
		hexstr++;

		nibble = hex(*hexstr);
		if (nibble < 0)
			return -1;

		*mem = (*mem << 4) | nibble;
		hexstr++;
		mem++;
	}

	return i;
}

static char *mem2hex(char *hexstr, const u8 *mem, int count)
{
	static const char hexchars[] = "0123456789abcdef";

	while (count-- > 0) {
		u8 ch = *mem++;
		*hexstr++ = hexchars[ch >> 4];
		*hexstr++ = hexchars[ch & 0xf];
	}

	return hexstr;
}

static efi_status_t efi_to_native(char *native, u16 *variable_name,
				  efi_guid_t *vendor)
{
	size_t len;

	len = utf16_strlen((u16 *)variable_name);
	if (len >= MAX_VAR_NAME)
		return EFI_DEVICE_ERROR;

	native += sprintf(native, "efi_%pUl_", vendor);
	native  = (char *)utf16_to_utf8((u8 *)native, (u16 *)variable_name, len);
	*native = '\0';

	return EFI_SUCCESS;
}

static const char *prefix(const char *str, const char *prefix)
{
	size_t n = strlen(prefix);
	if (!strncmp(prefix, str, n))
		return str + n;
	return NULL;
}

/* parse attributes part of variable value, if present: */
static const char *parse_attr(const char *str, u32 *attrp)
{
	u32 attr = 0;
	char sep = '{';

	if (*str != '{') {
		*attrp = EFI_VARIABLE_BOOTSERVICE_ACCESS;
		return str;
	}

	while (*str == sep) {
		const char *s;

		str++;

		if ((s = prefix(str, "ro"))) {
			attr |= READ_ONLY;
		} else if ((s = prefix(str, "boot"))) {
			attr |= EFI_VARIABLE_BOOTSERVICE_ACCESS;
		} else if ((s = prefix(str, "run"))) {
			attr |= EFI_VARIABLE_RUNTIME_ACCESS;
		} else {
			printf("invalid attribute: %s\n", str);
			break;
		}

		str = s;
		sep = ',';
	}

	str++;

	*attrp = attr;

	return str;
}

/* http://wiki.phoenix.com/wiki/index.php/EFI_RUNTIME_SERVICES#GetVariable.28.29 */
efi_status_t EFIAPI efi_get_variable(u16 *variable_name, efi_guid_t *vendor,
				     u32 *attributes, efi_uintn_t *data_size,
				     void *data)
{
	char native_name[MAX_NATIVE_VAR_NAME + 1];
	efi_status_t ret;
	unsigned long in_size;
	const char *val, *s;
	u32 attr;

	EFI_ENTRY("\"%ls\" %pUl %p %p %p", variable_name, vendor, attributes,
		  data_size, data);

	if (!variable_name || !vendor || !data_size)
		return EFI_EXIT(EFI_INVALID_PARAMETER);

	ret = efi_to_native(native_name, variable_name, vendor);
	if (ret)
		return EFI_EXIT(ret);

	debug("%s: get '%s'\n", __func__, native_name);

	val = env_get(native_name);
	if (!val)
		return EFI_EXIT(EFI_NOT_FOUND);

	val = parse_attr(val, &attr);

	in_size = *data_size;

	if ((s = prefix(val, "(blob)"))) {
		unsigned len = strlen(s);

		/* number of hexadecimal digits must be even */
		if (len & 1)
			return EFI_EXIT(EFI_DEVICE_ERROR);

		/* two characters per byte: */
		len /= 2;
		*data_size = len;

		if (in_size < len)
			return EFI_EXIT(EFI_BUFFER_TOO_SMALL);

		if (!data)
			return EFI_EXIT(EFI_INVALID_PARAMETER);

		if (hex2mem(data, s, len) != len)
			return EFI_EXIT(EFI_DEVICE_ERROR);

		debug("%s: got value: \"%s\"\n", __func__, s);
	} else if ((s = prefix(val, "(utf8)"))) {
		unsigned len = strlen(s) + 1;

		*data_size = len;

		if (in_size < len)
			return EFI_EXIT(EFI_BUFFER_TOO_SMALL);

		if (!data)
			return EFI_EXIT(EFI_INVALID_PARAMETER);

		memcpy(data, s, len);
		((char *)data)[len] = '\0';

		debug("%s: got value: \"%s\"\n", __func__, (char *)data);
	} else {
		debug("%s: invalid value: '%s'\n", __func__, val);
		return EFI_EXIT(EFI_DEVICE_ERROR);
	}

	if (attributes)
		*attributes = attr & EFI_VARIABLE_MASK;

	return EFI_EXIT(EFI_SUCCESS);
}

/* http://wiki.phoenix.com/wiki/index.php/EFI_RUNTIME_SERVICES#GetNextVariableName.28.29 */
efi_status_t EFIAPI efi_get_next_variable_name(efi_uintn_t *variable_name_size,
					       u16 *variable_name,
					       efi_guid_t *vendor)
{
	EFI_ENTRY("%p \"%ls\" %pUl", variable_name_size, variable_name, vendor);

	return EFI_EXIT(EFI_DEVICE_ERROR);
}

/* http://wiki.phoenix.com/wiki/index.php/EFI_RUNTIME_SERVICES#SetVariable.28.29 */
efi_status_t EFIAPI efi_set_variable(u16 *variable_name, efi_guid_t *vendor,
				     u32 attributes, efi_uintn_t data_size,
				     void *data)
{
	char native_name[MAX_NATIVE_VAR_NAME + 1];
	efi_status_t ret = EFI_SUCCESS;
	char *val, *s;
	u32 attr;

	EFI_ENTRY("\"%ls\" %pUl %x %zu %p", variable_name, vendor, attributes,
		  data_size, data);

	if (!variable_name || !vendor)
		return EFI_EXIT(EFI_INVALID_PARAMETER);

	ret = efi_to_native(native_name, variable_name, vendor);
	if (ret)
		return EFI_EXIT(ret);

#define ACCESS_ATTR (EFI_VARIABLE_RUNTIME_ACCESS | EFI_VARIABLE_BOOTSERVICE_ACCESS)

	if ((data_size == 0) || !(attributes & ACCESS_ATTR)) {
		/* delete the variable: */
		env_set(native_name, NULL);
		return EFI_EXIT(EFI_SUCCESS);
	}

	val = env_get(native_name);
	if (val) {
		parse_attr(val, &attr);

		if (attr & READ_ONLY)
			return EFI_EXIT(EFI_WRITE_PROTECTED);
	}

	val = malloc(2 * data_size + strlen("{ro,run,boot}(blob)") + 1);
	if (!val)
		return EFI_EXIT(EFI_OUT_OF_RESOURCES);

	s = val;

	/* store attributes: */
	attributes &= (EFI_VARIABLE_BOOTSERVICE_ACCESS | EFI_VARIABLE_RUNTIME_ACCESS);
	s += sprintf(s, "{");
	while (attributes) {
		u32 attr = 1 << (ffs(attributes) - 1);

		if (attr == EFI_VARIABLE_BOOTSERVICE_ACCESS)
			s += sprintf(s, "boot");
		else if (attr == EFI_VARIABLE_RUNTIME_ACCESS)
			s += sprintf(s, "run");

		attributes &= ~attr;
		if (attributes)
			s += sprintf(s, ",");
	}
	s += sprintf(s, "}");

	/* store payload: */
	s += sprintf(s, "(blob)");
	s = mem2hex(s, data, data_size);
	*s = '\0';

	debug("%s: setting: %s=%s\n", __func__, native_name, val);

	if (env_set(native_name, val))
		ret = EFI_DEVICE_ERROR;

	free(val);

	return EFI_EXIT(ret);
}
