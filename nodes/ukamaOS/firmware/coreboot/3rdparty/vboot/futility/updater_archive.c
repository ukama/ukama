/* Copyright 2018 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Accessing updater resources from an archive.
 */

#include <assert.h>
#include <ctype.h>
#include <fts.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <unistd.h>

#ifdef HAVE_LIBZIP
#include <zip.h>
#endif

#include "host_misc.h"
#include "updater.h"
#include "util_misc.h"
#include "vb2_common.h"

/*
 * A firmware update package (archive) is a file packed by either shar(1) or
 * zip(1). See https://chromium.googlesource.com/chromiumos/platform/firmware/
 * for more information.
 *
 * A package for single board (i.e., not Unified Build) will have all the image
 * files in top folder:
 *  - host: 'image.bin' (or 'bios.bin' as legacy name before CL:1318712)
 *  - ec: 'ec.bin'
 *  - pd: 'pd.bin'
 * If white label is supported, a 'keyset/' folder will be available, with key
 * files in it:
 *  - rootkey.$WLTAG
 *  - vblock_A.$WLTAG
 *  - vblock_B.$WLTAG
 * The $WLTAG should come from VPD value 'whitelabel_tag', or the
 * 'customization_id'. Note 'customization_id' is in format LOEM[-VARIANT] and
 * we can only take LOEM as $WLTAG, for example A-B => $WLTAG=A.
 *
 * A package for Unified Build is more complicated. There will be a models/
 * folder, and each model (by $(mosys platform model) ) should appear as a sub
 * folder, with a 'setvars.sh' file inside. The 'setvars.sh' is a shell script
 * describing what files should be used and the signature ID ($SIGID) to use.
 *
 * Similar to write label in non-Unified-Build, the keys and vblock files will
 * be in 'keyset/' folder:
 *  - rootkey.$SIGID
 *  - vblock_A.$SIGID
 *  - vblock_B.$SIGID
 * If $SIGID starts with 'sig-id-in-*' then we have to replace it by VPD value
 * 'whitelabel_tag' as '$MODEL-$WLTAG'.
 */

static const char * const SETVARS_IMAGE_MAIN = "IMAGE_MAIN",
		  * const SETVARS_IMAGE_EC = "IMAGE_EC",
		  * const SETVARS_IMAGE_PD = "IMAGE_PD",
		  * const SETVARS_SIGNATURE_ID = "SIGNATURE_ID",
		  * const SIG_ID_IN_VPD_PREFIX = "sig-id-in",
		  * const DIR_KEYSET = "keyset",
		  * const DIR_MODELS = "models",
		  * const DEFAULT_MODEL_NAME = "default",
		  * const VPD_WHITELABEL_TAG = "whitelabel_tag",
		  * const VPD_CUSTOMIZATION_ID = "customization_id",
		  * const ENV_VAR_MODEL_DIR = "${MODEL_DIR}",
		  * const PATH_STARTSWITH_KEYSET = "keyset/",
		  * const PATH_ENDSWITH_SERVARS = "/setvars.sh";

struct archive {
	void *handle;

	void * (*open)(const char *name);
	int (*close)(void *handle);

	int (*walk)(void *handle, void *arg,
		    int (*callback)(const char *path, void *arg));
	int (*has_entry)(void *handle, const char *name);
	int (*read_file)(void *handle, const char *fname,
			 uint8_t **data, uint32_t *size);
	int (*write_file)(void *handle, const char *fname,
			  uint8_t *data, uint32_t size);
};

/*
 * -- Begin of archive implementations --
 */

/* Callback for archive_open on a general file system. */
static void *archive_fallback_open(const char *name)
{
	assert(name && *name);
	return strdup(name);
}

/* Callback for archive_close on a general file system. */
static int archive_fallback_close(void *handle)
{
	free(handle);
	return 0;
}

/* Callback for archive_walk on a general file system. */
static int archive_fallback_walk(
		void *handle, void *arg,
		int (*callback)(const char *path, void *arg))
{
	FTS *fts_handle;
	FTSENT *ent;
	char *fts_argv[2] = {};
	char default_path[] = ".";
	char *root = default_path;
	size_t root_len;

	if (handle)
		root = (char *)handle;
	root_len = strlen(root);
	fts_argv[0] = root;

	fts_handle = fts_open(fts_argv, FTS_NOCHDIR, NULL);
	if (!fts_handle)
		return -1;

	while ((ent = fts_read(fts_handle)) != NULL) {
		char *path = ent->fts_path + root_len;
		if (ent->fts_info != FTS_F && ent->fts_info != FTS_SL)
			continue;
		while (*path == '/')
			path++;
		if (!*path)
			continue;
		if (callback(path, arg))
			break;
	}
	return 0;
}

/* Callback for fallback drivers to get full path easily. */
static const char *archive_fallback_get_path(void *handle, const char *fname,
					     char **temp_path)
{
	if (handle && *fname != '/') {
		ASPRINTF(temp_path, "%s/%s", (char *)handle, fname);
		return *temp_path;
	}
	return fname;
}

/* Callback for archive_has_entry on a general file system. */
static int archive_fallback_has_entry(void *handle, const char *fname)
{
	int r;
	char *temp_path = NULL;
	const char *path = archive_fallback_get_path(handle, fname, &temp_path);

	VB2_DEBUG("Checking %s\n", path);
	r = access(path, R_OK);
	free(temp_path);
	return r == 0;
}

/* Callback for archive_read_file on a general file system. */
static int archive_fallback_read_file(void *handle, const char *fname,
				      uint8_t **data, uint32_t *size)
{
	int r;
	char *temp_path = NULL;
	const char *path = archive_fallback_get_path(handle, fname, &temp_path);

	VB2_DEBUG("Reading %s\n", path);
	*data = NULL;
	*size = 0;
	r = vb2_read_file(path, data, size) != VB2_SUCCESS;
	free(temp_path);
	return r;
}

/* Callback for archive_write_file on a general file system. */
static int archive_fallback_write_file(void *handle, const char *fname,
				       uint8_t *data, uint32_t size)
{
	int r;
	char *temp_path = NULL;
	const char *path = archive_fallback_get_path(handle, fname, &temp_path);

	VB2_DEBUG("Writing %s\n", path);
	if (strchr(path, '/')) {
		char *dirname = strdup(path);
		*strrchr(dirname, '/') = '\0';
		/* TODO(hungte): call mkdir(2) instead of shell invocation. */
		if (access(dirname, W_OK) != 0) {
			char *command;
			ASPRINTF(&command, "mkdir -p %s", dirname);
			free(host_shell(command));
			free(command);
		}
		free(dirname);
	}
	r = vb2_write_file(path, data, size) != VB2_SUCCESS;
	free(temp_path);
	return r;
}

#ifdef HAVE_LIBZIP

/* Callback for archive_open on a ZIP file. */
static void *archive_zip_open(const char *name)
{
	return zip_open(name, 0, NULL);
}

/* Callback for archive_close on a ZIP file. */
static int archive_zip_close(void *handle)
{
	struct zip *zip = (struct zip *)handle;

	if (zip)
		return zip_close(zip);
	return 0;
}

/* Callback for archive_has_entry on a ZIP file. */
static int archive_zip_has_entry(void *handle, const char *fname)
{
	struct zip *zip = (struct zip *)handle;
	assert(zip);
	return zip_name_locate(zip, fname, 0) != -1;
}

/* Callback for archive_walk on a ZIP file. */
static int archive_zip_walk(
		void *handle, void *arg,
		int (*callback)(const char *name, void *arg))
{
	zip_int64_t num, i;
	struct zip *zip = (struct zip *)handle;
	assert(zip);

	num = zip_get_num_entries(zip, 0);
	if (num < 0)
		return 1;
	for (i = 0; i < num; i++) {
		const char *name = zip_get_name(zip, i, 0);
		if (*name && name[strlen(name) - 1] == '/')
			continue;
		if (callback(name, arg))
			break;
	}
	return 0;
}

/* Callback for archive_zip_read_file on a ZIP file. */
static int archive_zip_read_file(void *handle, const char *fname,
			     uint8_t **data, uint32_t *size)
{
	struct zip *zip = (struct zip *)handle;
	struct zip_file *fp;
	struct zip_stat stat;

	assert(zip);
	*data = NULL;
	*size = 0;
	zip_stat_init(&stat);
	if (zip_stat(zip, fname, 0, &stat)) {
		ERROR("Fail to stat entry in ZIP: %s\n", fname);
		return 1;
	}
	fp = zip_fopen(zip, fname, 0);
	if (!fp) {
		ERROR("Failed to open entry in ZIP: %s\n", fname);
		return 1;
	}
	*data = (uint8_t *)malloc(stat.size);
	if (*data) {
		if (zip_fread(fp, *data, stat.size) == stat.size) {
			*size = stat.size;
		} else {
			ERROR("Failed to read entry in zip: %s\n", fname);
			free(*data);
			*data = NULL;
		}
	}
	zip_fclose(fp);
	return *data == NULL;
}

/* Callback for archive_zip_write_file on a ZIP file. */
static int archive_zip_write_file(void *handle, const char *fname,
				  uint8_t *data, uint32_t size)
{
	struct zip *zip = (struct zip *)handle;
	struct zip_source *src;

	VB2_DEBUG("Writing %s\n", fname);
	assert(zip);
	src = zip_source_buffer(zip, data, size, 0);
	if (!src) {
		ERROR("Internal error: cannot allocate buffer: %s\n", fname);
		return 1;
	}

	if (zip_file_add(zip, fname, src, ZIP_FL_OVERWRITE) < 0) {
		zip_source_free(src);
		ERROR("Internal error: failed to add: %s\n", fname);
		return 1;
	}
	/* zip_source_free is not needed if zip_file_add success. */
#if LIBZIP_VERSION_MAJOR >= 1
	zip_file_set_mtime(zip, zip_name_locate(zip, fname, 0), 0, 0);
#endif
	return 0;
}
#endif

/*
 * Opens an archive from given path.
 * The type of archive will be determined automatically.
 * Returns a pointer to reference to archive (must be released by archive_close
 * when not used), otherwise NULL on error.
 */
struct archive *archive_open(const char *path)
{
	struct stat path_stat;
	struct archive *ar;

	if (stat(path, &path_stat) != 0) {
		ERROR("Cannot identify type of path: %s\n", path);
		return NULL;
	}

	ar = (struct archive *)malloc(sizeof(*ar));
	if (!ar) {
		ERROR("Internal error: allocation failure.\n");
		return NULL;
	}

	if (S_ISDIR(path_stat.st_mode)) {
		VB2_DEBUG("Found directory, use fallback (fs) driver: %s\n",
			  path);
		/* Regular file system. */
		ar->open = archive_fallback_open;
		ar->close = archive_fallback_close;
		ar->walk = archive_fallback_walk;
		ar->has_entry = archive_fallback_has_entry;
		ar->read_file = archive_fallback_read_file;
		ar->write_file = archive_fallback_write_file;
	} else {
#ifdef HAVE_LIBZIP
		VB2_DEBUG("Found file, use ZIP driver: %s\n", path);
		ar->open = archive_zip_open;
		ar->close = archive_zip_close;
		ar->walk = archive_zip_walk;
		ar->has_entry = archive_zip_has_entry;
		ar->read_file = archive_zip_read_file;
		ar->write_file = archive_zip_write_file;
#else
		ERROR("Found file, but no drivers were enabled: %s\n", path);
		free(ar);
		return NULL;
#endif
	}
	ar->handle = ar->open(path);
	if (!ar->handle) {
		ERROR("Failed to open archive: %s\n", path);
		free(ar);
		return NULL;
	}
	return ar;
}

/*
 * Closes an archive reference.
 * Returns 0 on success, otherwise non-zero as failure.
 */
int archive_close(struct archive *ar)
{
	int r = ar->close(ar->handle);
	free(ar);
	return r;
}

/*
 * Checks if an entry (either file or directory) exists in archive.
 * If entry name (fname) is an absolute path (/file), always check
 * with real file system.
 * Returns 1 if exists, otherwise 0
 */
int archive_has_entry(struct archive *ar, const char *name)
{
	if (!ar || *name == '/')
		return archive_fallback_has_entry(NULL, name);
	return ar->has_entry(ar->handle, name);
}

/*
 * Traverses all files within archive (directories are ignored).
 * For every entry, the path (relative the archive root) will be passed to
 * callback function, until the callback returns non-zero.
 * The arg argument will also be passed to callback.
 * Returns 0 on success otherwise non-zero as failure.
 */
static int archive_walk(struct archive *ar, void *arg,
			int (*callback)(const char *path, void *arg))
{
	if (!ar)
		return archive_fallback_walk(NULL, arg, callback);
	return ar->walk(ar->handle, arg, callback);
}

/*
 * Reads a file from archive.
 * If entry name (fname) is an absolute path (/file), always read
 * from real file system.
 * Returns 0 on success (data and size reflects the file content),
 * otherwise non-zero as failure.
 */
int archive_read_file(struct archive *ar, const char *fname,
		      uint8_t **data, uint32_t *size)
{
	if (!ar || *fname == '/')
		return archive_fallback_read_file(NULL, fname, data, size);
	return ar->read_file(ar->handle, fname, data, size);
}

/*
 * Writes a file into archive.
 * If entry name (fname) is an absolute path (/file), always write into real
 * file system.
 * Returns 0 on success, otherwise non-zero as failure.
 */
int archive_write_file(struct archive *ar, const char *fname,
		       uint8_t *data, uint32_t size)
{
	if (!ar || *fname == '/')
		return archive_fallback_write_file(NULL, fname, data, size);
	return ar->write_file(ar->handle, fname, data, size);
}

struct _copy_arg {
	struct archive *from, *to;
};

/* Callback for archive_copy. */
static int archive_copy_callback(const char *path, void *_arg)
{
	const struct _copy_arg *arg = (const struct _copy_arg*)_arg;
	uint32_t size;
	uint8_t *data;
	int r;

	INFO("Copying: %s\n", path);
	if (archive_read_file(arg->from, path, &data, &size)) {
		ERROR("Failed reading: %s\n", path);
		return 1;
	}
	r = archive_write_file(arg->to, path, data, size);
	VB2_DEBUG("result=%d\n", r);
	free(data);
	return r;
}

/*
 * Copies all entries from one archive to another.
 * Returns 0 on success, otherwise non-zero as failure.
 */
int archive_copy(struct archive *from, struct archive *to)
{
	struct _copy_arg arg = { .from = from, .to = to };
	return archive_walk(from, &arg, archive_copy_callback);
}

/*
 * -- End of archive implementations --
 */

/* Utility function to convert a string. */
static void str_convert(char *s, int (*convert)(int c))
{
	int c;

	for (; *s; s++) {
		c = *s;
		if (!isascii(c))
			continue;
		*s = convert(c);
	}
}

/* Returns 1 if name ends by given pattern, otherwise 0. */
static int str_endswith(const char *name, const char *pattern)
{
	size_t name_len = strlen(name), pattern_len = strlen(pattern);
	if (name_len < pattern_len)
		return 0;
	return strcmp(name + name_len - pattern_len, pattern) == 0;
}

/* Returns 1 if name starts by given pattern, otherwise 0. */
static int str_startswith(const char *name, const char *pattern)
{
	return strncmp(name, pattern, strlen(pattern)) == 0;
}

/* Returns the VPD value by given key name, or NULL on error (or no value). */
static char *vpd_get_value(const char *fpath, const char *key)
{
	char *command, *result;

	assert(fpath);
	ASPRINTF(&command, "vpd -g %s -f %s 2>/dev/null", key, fpath);
	result = host_shell(command);
	free(command);

	if (result && !*result) {
		free(result);
		result = NULL;
	}
	return result;
}

/*
 * Reads and parses a setvars type file from archive, then stores into config.
 * Returns 0 on success (at least one entry found), otherwise failure.
 */
static int model_config_parse_setvars_file(
		struct model_config *cfg, struct archive *archive,
		const char *fpath)
{
	uint8_t *data;
	uint32_t len;

	char *ptr_line, *ptr_token;
	char *line, *k, *v;
	int valid = 0;

	if (archive_read_file(archive, fpath, &data, &len) != 0) {
		ERROR("Failed reading: %s\n", fpath);
		return -1;
	}

	/* Valid content should end with \n, or \"; ensure ASCIIZ for parsing */
	if (len)
		data[len - 1] = '\0';

	for (line = strtok_r((char *)data, "\n\r", &ptr_line); line;
	     line = strtok_r(NULL, "\n\r", &ptr_line)) {
		char *expand_path = NULL;
		int found_valid = 1;

		/* Format: KEY="value" */
		k = strtok_r(line, "=", &ptr_token);
		if (!k)
			continue;
		v = strtok_r(NULL, "\"", &ptr_token);
		if (!v)
			continue;

		/* Some legacy updaters may be still using ${MODEL_DIR}. */
		if (str_startswith(v, ENV_VAR_MODEL_DIR)) {
			ASPRINTF(&expand_path, "%s/%s%s", DIR_MODELS, cfg->name,
				 v + strlen(ENV_VAR_MODEL_DIR));
		}

		if (strcmp(k, SETVARS_IMAGE_MAIN) == 0)
			cfg->image = strdup(v);
		else if (strcmp(k, SETVARS_IMAGE_EC) == 0)
			cfg->ec_image = strdup(v);
		else if (strcmp(k, SETVARS_IMAGE_PD) == 0)
			cfg->pd_image = strdup(v);
		else if (strcmp(k, SETVARS_SIGNATURE_ID) == 0) {
			cfg->signature_id = strdup(v);
			if (str_startswith(v, SIG_ID_IN_VPD_PREFIX))
				cfg->is_white_label = 1;
		} else
			found_valid = 0;
		free(expand_path);
		valid += found_valid;
	}
	free(data);
	return valid == 0;
}

/*
 * Changes the rootkey in firmware GBB to given new key.
 * Returns 0 on success, otherwise failure.
 */
static int change_gbb_rootkey(struct firmware_image *image,
			      const char *section_name,
			      const uint8_t *rootkey, uint32_t rootkey_len)
{
	const struct vb2_gbb_header *gbb = find_gbb(image);
	uint8_t *gbb_rootkey;
	if (!gbb) {
		ERROR("Cannot find GBB in image %s.\n", image->file_name);
		return -1;
	}
	if (gbb->rootkey_size < rootkey_len) {
		ERROR("New root key (%u bytes) larger than GBB (%u bytes).\n",
		      rootkey_len, gbb->rootkey_size);
		return -1;
	}

	gbb_rootkey = (uint8_t *)gbb + gbb->rootkey_offset;
	/* See cmd_gbb_utility: root key must be first cleared with zero. */
	memset(gbb_rootkey, 0, gbb->rootkey_size);
	memcpy(gbb_rootkey, rootkey, rootkey_len);
	return 0;
}

/*
 * Changes the VBlock in firmware section to new data.
 * Returns 0 on success, otherwise failure.
 */
static int change_vblock(struct firmware_image *image, const char *section_name,
			 const uint8_t *vblock, uint32_t vblock_len)
{
	struct firmware_section section;

	find_firmware_section(&section, image, section_name);
	if (!section.data) {
		ERROR("Need section %s in image %s.\n", section_name,
		      image->file_name);
		return -1;
	}
	if (section.size < vblock_len) {
		ERROR("Section %s too small (%zu bytes) for vblock (%u bytes).\n",
		      section_name, section.size, vblock_len);
		return -1;
	}
	memcpy(section.data, vblock, vblock_len);
	return 0;
}

/*
 * Applies a key file to firmware image.
 * Returns 0 on success, otherwise failure.
 */
static int apply_key_file(
		struct firmware_image *image, const char *path,
		struct archive *archive, const char *section_name,
		int (*apply)(struct firmware_image *image, const char *section,
			     const uint8_t *data, uint32_t len))
{
	int r = 0;
	uint8_t *data = NULL;
	uint32_t len;

	r = archive_read_file(archive, path, &data, &len);
	if (r == 0) {
		VB2_DEBUG("Loaded file: %s\n", path);
		r = apply(image, section_name, data, len);
		if (r)
			ERROR("Failed applying %s to %s\n", path, section_name);
	} else {
		ERROR("Failed reading: %s\n", path);
	}
	free(data);
	return r;
}

/*
 * Modifies a firmware image from patch information specified in model config.
 * Returns 0 on success, otherwise number of failures.
 */
int patch_image_by_model(
		struct firmware_image *image, const struct model_config *model,
		struct archive *archive)
{
	int err = 0;
	if (model->patches.rootkey)
		err += !!apply_key_file(
				image, model->patches.rootkey, archive,
				FMAP_RO_GBB, change_gbb_rootkey);
	if (model->patches.vblock_a)
		err += !!apply_key_file(
				image, model->patches.vblock_a, archive,
				FMAP_RW_VBLOCK_A, change_vblock);
	if (model->patches.vblock_b)
		err += !!apply_key_file(
				image, model->patches.vblock_b, archive,
				FMAP_RW_VBLOCK_B, change_vblock);
	return err;
}

/*
 * Finds available patch files by given model.
 * Updates `model` argument with path of patch files.
 */
static void find_patches_for_model(struct model_config *model,
				   struct archive *archive,
				   const char *signature_id)
{
	char *path;
	int i;

	const char *names[] = {
		"rootkey",
		"vblock_A",
		"vblock_B",
	};

	char **targets[] = {
		&model->patches.rootkey,
		&model->patches.vblock_a,
		&model->patches.vblock_b,
	};

	assert(ARRAY_SIZE(names) == ARRAY_SIZE(targets));
	for (i = 0; i < ARRAY_SIZE(names); i++) {
		ASPRINTF(&path, "%s/%s.%s", DIR_KEYSET, names[i], signature_id);
		if (archive_has_entry(archive, path))
			*targets[i] = path;
		else
			free(path);
	}
}

/*
 * Adds and copies one new model config to the existing list of given manifest.
 * Returns a pointer to the newly allocated config, or NULL on failure.
 */
static struct model_config *manifest_add_model(
		struct manifest *manifest,
		const struct model_config *cfg)
{
	struct model_config *model;
	manifest->num++;
	manifest->models = (struct model_config *)realloc(
			manifest->models, manifest->num * sizeof(*model));
	if (!manifest->models) {
		ERROR("Internal error: failed to allocate buffer.\n");
		return NULL;
	}
	model = &manifest->models[manifest->num - 1];
	memcpy(model, cfg, sizeof(*model));
	return model;
}

/*
 * A callback function for manifest to scan files in archive.
 * Returns 0 to keep scanning, or non-zero to stop.
 */
static int manifest_scan_entries(const char *name, void *arg)
{
	struct manifest *manifest = (struct manifest *)arg;
	struct archive *archive = manifest->archive;
	struct model_config model = {0};
	char *slash;

	if (str_startswith(name, PATH_STARTSWITH_KEYSET))
		manifest->has_keyset = 1;
	if (!str_endswith(name, PATH_ENDSWITH_SERVARS))
		return 0;

	/* name: models/$MODEL/setvars.sh */
	model.name = strdup(strchr(name, '/') + 1);
	slash = strchr(model.name, '/');
	if (slash)
		*slash = '\0';

	VB2_DEBUG("Found model <%s> setvars: %s\n", model.name, name);
	if (model_config_parse_setvars_file(&model, archive, name)) {
		ERROR("Invalid setvars file: %s\n", name);
		return 0;
	}

	/* In legacy setvars.sh, the ec_image and pd_image may not exist. */
	if (model.ec_image && !archive_has_entry(archive, model.ec_image)) {
		VB2_DEBUG("Ignore non-exist EC image: %s\n", model.ec_image);
		free(model.ec_image);
		model.ec_image = NULL;
	}
	if (model.pd_image && !archive_has_entry(archive, model.pd_image)) {
		VB2_DEBUG("Ignore non-exist PD image: %s\n", model.pd_image);
		free(model.pd_image);
		model.pd_image = NULL;
	}

	/* Find patch files. */
	if (model.signature_id)
		find_patches_for_model(&model, archive, model.signature_id);

	return !manifest_add_model(manifest, &model);
}

/*
 * Finds the existing model_config from manifest that best matches current
 * system (as defined by model_name).
 * Returns a model_config from manifest, or NULL if not found.
 */
const struct model_config *manifest_find_model(const struct manifest *manifest,
					       const char *model_name)
{
	char *sys_model_name = NULL;
	const struct model_config *model = NULL;
	int i;

	/*
	 * For manifest with single model defined, we should just return because
	 * there are other mechanisms like platform name check to double confirm
	 * if the firmware is valid.
	 */
	if (manifest->num == 1)
		return &manifest->models[0];

	if (!model_name) {
		sys_model_name = host_shell("mosys platform model");
		VB2_DEBUG("System model name: '%s'\n", sys_model_name);
		model_name = sys_model_name;
	}

	for (i = 0; !model && i < manifest->num; i++) {
		if (strcmp(model_name, manifest->models[i].name) == 0)
			model = &manifest->models[i];
	}
	if (!model) {
		if (!*model_name)
			ERROR("Cannot get model name.\n");
		else
			ERROR("Unsupported model: '%s'.\n", model_name);

		fprintf(stderr,
			"You are probably running an image for wrong board, or "
			"a device in early stage that 'mosys' command is not "
			"ready, or image from old (or factory) branches that "
			"Unified Build config is not updated yet for 'mosys'.\n"
			"Please check command 'mosys platform model', "
			"which should output one of the supported models below:"
			"\n");

		for (i = 0; i < manifest->num; i++)
			fprintf(stderr, " %s", manifest->models[i].name);
		fprintf(stderr, "\n");
	}


	free(sys_model_name);
	return model;
}

/*
 * Determines the signature ID to use for white label.
 * Returns the signature ID for looking up rootkey and vblock files.
 * Caller must free the returned string.
 */
static char *resolve_signature_id(struct model_config *model, const char *image)
{
	int is_unibuild = model->signature_id ? 1 : 0;
	char *wl_tag = vpd_get_value(image, VPD_WHITELABEL_TAG);
	char *sig_id = NULL;

	/* Unified build: $model.$wl_tag, or $model (b/126800200). */
	if (is_unibuild) {
		if (!wl_tag) {
			WARN("No VPD '%s' set for white label - use model name "
			     "'%s' as default.\n", VPD_WHITELABEL_TAG,
			     model->name);
			return strdup(model->name);
		}

		ASPRINTF(&sig_id, "%s-%s", model->name, wl_tag);
		free(wl_tag);
		return sig_id;
	}

	/* Non-Unibuild: Upper($wl_tag), or Upper(${cid%%-*}). */
	if (!wl_tag) {
		char *cid = vpd_get_value(image, VPD_CUSTOMIZATION_ID);
		if (cid) {
			/* customization_id in format LOEM[-VARIANT]. */
			char *dash = strchr(cid, '-');
			if (dash)
				*dash = '\0';
			wl_tag = cid;
		}
	}
	if (wl_tag)
		str_convert(wl_tag, toupper);
	return wl_tag;
}

/*
 * Applies white label information to an existing model configuration.
 * Collects signature ID information from either parameter signature_id or
 * image file (via VPD) and updates model.patches for key files.
 * Returns 0 on success, otherwise failure.
 */
int model_apply_white_label(
		struct model_config *model,
		struct archive *archive,
		const char *signature_id,
		const char *image)
{
	char *sig_id = NULL;
	int r = 0;

	if (!signature_id) {
		sig_id = resolve_signature_id(model, image);
		signature_id = sig_id;
	}

	if (signature_id) {
		VB2_DEBUG("Find white label patches by signature ID: '%s'.\n",
		      signature_id);
		find_patches_for_model(model, archive, signature_id);
	} else {
		signature_id = "";
		WARN("No VPD '%s' set for white label - use default keys.\n",
		     VPD_WHITELABEL_TAG);
	}
	if (!model->patches.rootkey) {
		ERROR("No keys found for signature_id: '%s'\n", signature_id);
		r = 1;
	} else {
		INFO("Applied for white label: %s\n", signature_id);
	}
	free(sig_id);
	return r;
}

/*
 * Creates a new manifest object by scanning files in archive.
 * Returns the manifest on success, otherwise NULL for failure.
 */
struct manifest *new_manifest_from_archive(struct archive *archive)
{
	struct manifest manifest = {0}, *new_manifest;
	struct model_config model = {0};
	const char * const host_image_name = "image.bin",
		   * const old_host_image_name = "bios.bin",
	           * const ec_name = "ec.bin",
		   * const pd_name = "pd.bin";

	manifest.archive = archive;
	manifest.default_model = -1;
	archive_walk(archive, &manifest, manifest_scan_entries);
	if (manifest.num == 0) {
		const char *image_name = NULL;
		struct firmware_image image = {0};

		/* Try to load from current folder. */
		if (archive_has_entry(archive, old_host_image_name))
			image_name = old_host_image_name;
		else if (archive_has_entry(archive, host_image_name))
			image_name = host_image_name;
		else
			return 0;

		model.image = strdup(image_name);
		if (archive_has_entry(archive, ec_name))
			model.ec_image = strdup(ec_name);
		if (archive_has_entry(archive, pd_name))
			model.pd_image = strdup(pd_name);
		/* Extract model name from FWID: $Vendor_$Platform.$Version */
		if (!load_firmware_image(&image, image_name, archive)) {
			char *token = NULL;
			if (strtok(image.ro_version, "_"))
				token = strtok(NULL, ".");
			if (token && *token) {
				str_convert(token, tolower);
				model.name = strdup(token);
			}
			free_firmware_image(&image);
		}
		if (!model.name)
			model.name = strdup(DEFAULT_MODEL_NAME);
		if (manifest.has_keyset)
			model.is_white_label = 1;
		manifest_add_model(&manifest, &model);
		manifest.default_model = manifest.num - 1;
	}
	VB2_DEBUG("%d model(s) loaded.\n", manifest.num);
	if (!manifest.num) {
		ERROR("No valid configurations found from archive.\n");
		return NULL;
	}

	new_manifest = (struct manifest *)malloc(sizeof(manifest));
	if (!new_manifest) {
		ERROR("Internal error: memory allocation error.\n");
		return NULL;
	}
	memcpy(new_manifest, &manifest, sizeof(manifest));
	return new_manifest;
}

/* Releases all resources allocated by given manifest object. */
void delete_manifest(struct manifest *manifest)
{
	int i;
	assert(manifest);
	for (i = 0; i < manifest->num; i++) {
		struct model_config *model = &manifest->models[i];
		free(model->name);
		free(model->signature_id);
		free(model->image);
		free(model->ec_image);
		free(model->pd_image);
		free(model->patches.rootkey);
		free(model->patches.vblock_a);
		free(model->patches.vblock_b);
	}
	free(manifest->models);
	free(manifest);
}

static const char *get_gbb_key_hash(const struct vb2_gbb_header *gbb,
				    int32_t offset, int32_t size)
{
	struct vb2_packed_key *key;

	if (!gbb)
		return "<No GBB>";
	key = (struct vb2_packed_key *)((uint8_t *)gbb + offset);
	if (!packed_key_looks_ok(key, size))
	    return "<Invalid key>";
	return packed_key_sha1_string(key);
}

/* Prints the information of given image file in JSON format. */
static void print_json_image(
		const char *name, const char *fpath, struct model_config *m,
		struct archive *archive, int indent, int is_host)
{
	struct firmware_image image = {0};
	const struct vb2_gbb_header *gbb = NULL;
	if (!fpath)
		return;
	if (load_firmware_image(&image, fpath, archive))
		return;
	if (is_host)
		gbb = find_gbb(&image);
	else
		printf(",\n");
	printf("%*s\"%s\": { \"versions\": { \"ro\": \"%s\", \"rw\": \"%s\" },",
	       indent, "", name, image.ro_version, image.rw_version_a);
	indent += 2;
	if (is_host && patch_image_by_model(&image, m, archive) != 0) {
		ERROR("Failed to patch images by model: %s\n", m->name);
	} else if (gbb) {
		printf("\n%*s\"keys\": { \"root\": \"%s\", ",
		       indent, "",
		       get_gbb_key_hash(gbb, gbb->rootkey_offset,
					gbb->rootkey_size));
		printf("\"recovery\": \"%s\" },",
		       get_gbb_key_hash(gbb, gbb->recovery_key_offset,
					gbb->recovery_key_size));
	}
	printf("\n%*s\"image\": \"%s\" }", indent, "", fpath);
	free_firmware_image(&image);
}

/* Prints the information of objects in manifest (models and images) in JSON. */
void print_json_manifest(const struct manifest *manifest)
{
	int i, indent;
	struct archive *ar = manifest->archive;

	printf("{\n");
	for (i = 0, indent = 2; i < manifest->num; i++) {
		struct model_config *m = &manifest->models[i];
		printf("%s%*s\"%s\": {\n", i ? ",\n" : "", indent, "", m->name);
		indent += 2;
		print_json_image("host", m->image, m, ar, indent, 1);
		print_json_image("ec", m->ec_image, m, ar, indent, 0);
		print_json_image("pd", m->pd_image, m, ar, indent, 0);
		if (m->patches.rootkey) {
			struct patch_config *p = &m->patches;
			printf(",\n%*s\"patches\": { \"rootkey\": \"%s\", "
			       "\"vblock_a\": \"%s\", \"vblock_b\": \"%s\" }",
			       indent, "", p->rootkey, p->vblock_a,
			       p->vblock_b);
		}
		if (m->signature_id)
			printf(",\n%*s\"signature_id\": \"%s\"", indent, "",
			       m->signature_id);
		printf("\n  }");
		indent -= 2;
		assert(indent == 2);
	}
	printf("\n}\n");
}
