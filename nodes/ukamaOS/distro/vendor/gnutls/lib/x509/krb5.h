/*
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Author: Nikos Mavrogiannopoulos
 *
 * This file is part of GnuTLS.
 *
 * The GnuTLS is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 2.1 of
 * the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>
 *
 */

#ifndef GNUTLS_LIB_X509_KRB5_H
#define GNUTLS_LIB_X509_KRB5_H

int _gnutls_krb5_principal_to_der(const char *name, gnutls_datum_t * der);
int _gnutls_krb5_der_to_principal(const gnutls_datum_t * der, gnutls_datum_t *name);

#endif /* GNUTLS_LIB_X509_KRB5_H */
