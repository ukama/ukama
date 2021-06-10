/*
 * parser.c
 *
 *  Created on: Jun 7, 2021
 *      Author: vishal
 */

#include "parser.h"

client_config_t lw_cfg = {0};

/* Free server cfg */
int free_server_cfg(server_cfg_t ** cfg) {
	if (*cfg) {
		if ((*cfg)->addr) {
			free((*cfg)->addr);
		}
		free(*cfg);
		*cfg = NULL;
	}
}

/* Free file store cfg */
int free_file_store_cfg(file_store_t ** cfg) {
	if (*cfg) {

		/* certs */
		if ((*cfg)->certs) {
			free((*cfg)->certs);
		}

		/* address */
		if ((*cfg)->addr) {
			free((*cfg)->addr);
		}

		free(*cfg);
		*cfg = NULL;

	}
}

/* Extract string from the Toml table */
int extract_string(const toml_table_t* table, const char* tag, char** data)
{
	int ret = -1;
	/* Extract values */
	toml_datum_t tstr = toml_string_in(table, tag);
	if (!tstr.ok)
	{
		fprintf(stderr, "Error:: Failed to read tag %s\n", tag);
		*data = NULL;

	}
	else
	{
		uint32_t len =strlen(tstr.u.s);

		*data = malloc(sizeof(char)*len + 1);
		if(*data)
		{
			memset(*data, '\0', len+1);
			memcpy(*data,tstr.u.s, len);
		}

		fprintf(stdout, "%s tag read %s value.\n", tag, *data);
		free(tstr.u.s);

		ret = 0;
	}

	return ret;

}

/* Extract integer from the Toml table */
int extract_integer(const toml_table_t* table, const char* tag, int* data)
{
	int ret = -1;
	/* Extract values */
	toml_datum_t tint = toml_int_in(table, tag);
	if (!tint.ok)
	{

		fprintf(stderr, "Error:: Failed to read tag %s\n", tag);
		*data = 0;

	}
	else
	{
		*data = tint.u.b;
		ret = 0;
	}

	return ret;

}

/* Parse Server config entry */
int parse_server_subtable(toml_table_t* conf, char* tag, server_cfg_t** cfgdata)
{
	int ret = -1;
	server_cfg_t *data = NULL;

	/* Traverse to tag in table. */
	toml_table_t* ttag = toml_table_in(conf, tag);
	if (ttag)
	{
		 data = malloc(sizeof(server_cfg_t));
		if (data) {

			/* Extract address values */
			ret = extract_string(ttag, SERVER_ADD_STR, &(data->addr));
			if (ret) {
				goto error;
			}

			/* Extract port values */
			ret = extract_integer(ttag, SERVER_PORT_STR, &(data->port));
			if (ret) {
				goto error;
			}

			ret = 0;
		}

	}
	else
	{
		fprintf(stderr, "Error:: Missing %s \n", tag);
	}

	/* error handling */
	error:
	if (ret) {
		free_server_cfg(&data);
	}

	/* Assign data to server cfg data*/
	*cfgdata = data;
	return ret;
}

int parse_file_store_subtable(toml_table_t* conf, char* tag, file_store_t** cfgdata)
{
	int ret = -1;
	file_store_t *data = NULL;

	/* Traverse to tag in table. */
	toml_table_t* ttag = toml_table_in(conf, tag);
	if (ttag)
	{
		 data = malloc(sizeof(file_store_t));
		if (data) {

			/* Extract certs path */
			ret = extract_string(ttag, CERTS_STORE_STR, &(data->certs));
			if (ret) {
				goto error;
			}

			/* Extract IP path */
			ret = extract_string(ttag, IP_STORE_STR, &(data->addr));
			if (ret) {
				goto error;
			}

			ret = 0;
		}

	}
	else
	{
		fprintf(stderr, "Error:: Missing %s \n", tag);
	}

	/* error handling */
	error:
	if (ret) {
		free_file_store_cfg(&data);
	}

	/* Assign data to file store cfg data*/
	*cfgdata = data;
	return ret;
}

/* Parse config file. */
int parse_config(char* cfgName)
{
	int ret = 0;
	FILE* fp;
	char errbuf[200];

	/* Read and parse toml file */
	fp = fopen(cfgName, "r");
	if (!fp)
	{
		fprintf(stderr, "ERROR:: cannot open %s file error %s\n", cfgName, strerror(errno));
		ret = -1;
		return ret;
	}

	/* Parse TOML config file */
	toml_table_t* conf = toml_parse_file(fp, errbuf, sizeof(errbuf));
	fclose(fp);

	if (!conf)
	{
		fprintf(stderr, "ERROR:: cannot parse %s file error %s\n", cfgName, errbuf);
		ret = -1;
		goto error;
	}

	/* Traverse to Lookup server in table. */
	ret = parse_server_subtable(conf, SERVER_STR , &lw_cfg.server);
	if (ret)
	{
		goto error;
	}

	/* Traverse to a DMR Server in table. */
	ret = parse_file_store_subtable(conf, FILE_STORE_STR , &lw_cfg.file_store);
	if (ret)
	{
		goto error;
	}

	error:
	toml_free(conf);
	return ret;
}

