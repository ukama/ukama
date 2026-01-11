#include <assert.h>
#include <arpa/inet.h>
#include <resolv.h>
#include <netdb.h>
#include <stdio.h>
#include <poll.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>

#include "initClient.h"
#include "log.h"

enum
{
	POLL_TIMEOUT  = 5000,
	MAXREVLEN     = 73,
};

UpdateIpCallback register_fxn_cb = NULL;

static struct addrinfo *resolve_server(const char *server)
{
	/* translate the server name to an address */
	struct addrinfo hints = {
			.ai_family = AF_UNSPEC,
			.ai_socktype = SOCK_DGRAM,
			.ai_flags = AI_PASSIVE|AI_CANONNAME|AI_NUMERICSERV,
	};
	struct addrinfo *res;

	int rv = getaddrinfo(server, "53", &hints, &res);
	if (rv != 0 || res == NULL) {
		log_error("cannot resolve %s\n", server);
		return NULL;
	}

	/* print the address of the server */
	char str[INET6_ADDRSTRLEN];
	log_info("%-10s %s", "server:", server);
	log_info("%-10s %s", "Address 1:",
			inet_ntop(res->ai_family, res->ai_addr, str, sizeof(str)));

	return res;
}

static int res_ssend(struct addrinfo *srv, const unsigned char *msg,
		int msglen, unsigned char *answer, int anslen)
{
	struct sockaddr_storage src = {
			.ss_family = srv->ai_family,
	};

	int fd = socket(srv->ai_family,
			SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, 0);
	if (fd < 0) {
		log_error("cannot create the socket");
		goto err;
	}

	if (bind(fd, (void *)&src, srv->ai_addrlen) < 0) {
		log_error("bind failure");
		goto err_close;
	}

	/* send the query */
	if (sendto(fd, msg, msglen, MSG_NOSIGNAL,
			srv->ai_addr, srv->ai_addrlen) < 0)
	{
		log_error("sendto failure");
		goto err_close;
	}

	/* wait for the answer */
	struct pollfd pfd;
	pfd.fd = fd;
	pfd.events = POLLIN;
	if (poll(&pfd, 1, POLL_TIMEOUT) <= 0) {
		log_error("poll timeout");
		goto err_close;
	}

	/* receive the data */
	size_t alen = anslen;
	ssize_t rlen = recvfrom(fd, answer, alen, 0, NULL, NULL);
	if (rlen < 0) {
		log_error("recvfrom error");
		goto err_close;
	}

	close(fd);
	return rlen;

	err_close:
	close(fd);
	err:
	return -1;
}

static int get_ip(const uint8_t *data, char **ip) {
	*ip = (char*)calloc(1, (sizeof(char)*INET6_ADDRSTRLEN));
	if (inet_ntop(AF_INET,data, *ip, (sizeof(char)*INET6_ADDRSTRLEN))){
		log_info("IPV4 address associated is %s", *ip);
		return TRUE;
	} else {
		log_info("IPV4 address not found", *ip);
		return FALSE;
	}
}

/* modified from MUSL libc code */
static int dns_parse(const unsigned char *r, int rlen, void *ctx, char **ip)
{
	/* return if we didn't even get the header */
	if (rlen < 12)
		return FALSE;

	/* return in case of errors */
	if ((r[3] & 15))
		return FALSE;

	int qdcount = r[4]*256 + r[5];
	int ancount = r[6]*256 + r[7];

	if (qdcount + ancount > 64)
		return FALSE;

	const unsigned char *p = r+12;
	while (qdcount--) {
		while (p-r < rlen && *p-1U < 127)
			p++;

		if (*p>193 || (*p==193 && p[1]>254) || p>r+rlen-6)
			return FALSE;

		p += 5 + !!*p;
	}

	while (ancount--) {
		while (p-r < rlen && *p-1U < 127)
			p++;

		if (*p>193 || (*p==193 && p[1]>254) || p>r+rlen-6)
			return FALSE;

		p += 1 + !!*p;
		size_t len = p[8]*256U + p[9];
		if (p+len > r+rlen)
			return FALSE;

		if (!get_ip(p+10, ip)) {
			return FALSE;
		}
		p += 10 + len;
	}
	return TRUE;
}

static const char *dns_strerror(const unsigned char *r, int rlen)
{
	switch (rlen > 3 ? (r[3] & 15) : 0)
	{
	case 0: return "NoAnswer";
	case 1: return "FormErr";
	case 2: return "SerFail";
	case 3: return "NXDomain";
	case 4: return "NotImp";
	case 5: return "Refused";
	case 6: return "YXDomain";
	case 7: return "YXRRSet";
	case 8: return "NXRRSet";
	case 9:
	case 10: return "NotAuth";
	case 11: return "NotZone";
	default: return "Unassigned";
	}
}


char* nslookup(char* name, char *server)
{

	char *ip = NULL;
	unsigned char query[280];
	unsigned char response[1024];

	/* query types */
	ns_type qtype = ns_t_a;

	struct addrinfo *srv = (!server)
										? resolve_server("127.0.0.1")
												: resolve_server(server);
	if (!srv)
	{
		log_error("cannot resolve the nameserver");
		return FALSE;
	}

	int qlen = res_mkquery(0, name, ns_c_in, qtype, 0, 0, 0,
			query, sizeof(query));
	if (qlen < 0)
	{
		log_error("dns query build failed");
		goto clean;
	}

	/* send the query to the server */
	int rlen = (!server) ? res_send(query, qlen, response, sizeof(response))
						: res_ssend(srv, query, qlen, response, sizeof(response));
	if (rlen < 0)
	{
		log_error("failed to execute query.");
		goto clean;
	}

	/* check if query and response id match */
	if (memcmp(query, response, 2))
	{
		log_error("query don't match to response");
		goto clean;
	}

	/* decode the response */
	if (!dns_parse(response, rlen, NULL, &ip))
	{
		log_error("dns error: %s", dns_strerror(response, rlen));
		goto clean;
	}

	log_info("IP for %s is %s", name, ip);

	clean:
	freeaddrinfo(srv);

	return ip;
}

void register_callback(UpdateIpCallback cb) {
	if (cb) {
		register_fxn_cb = cb;
	}
}

void* refresh_lookup(void* args) {
	Config *c = (Config*) args;
	char* rIp = NULL;
	while(TRUE) {
		rIp = nslookup(c->systemDNS, c->nameServer);
		if (rIp) {
			if (strcmp(rIp, c->systemAddr) != 0) {
				/* update IP */
				free(c->systemAddr);
				c->systemAddr = strdup(rIp);
				free(rIp);

				/* callback function */
				if (register_fxn_cb) {
				  int status = register_fxn_cb(c);
				  if (status) {
					  log_info("Failed to update IP for the %s to %s.", c->systemName, c->systemAddr);
				  }
				}
			}
		}
		sleep(c->timePeriod);
	}

	pthread_exit (NULL);
}

