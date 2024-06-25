/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { cookies } from 'next/headers';
import type { NextRequest } from 'next/server';
import { NextResponse } from 'next/server';
import { Role_Type } from './client/graphql/generated';

type User = {
  id: string;
  role: string;
  name: string;
  email: string;
  orgId: string;
  token: string;
  orgName: string;
  isEmailVerified: boolean;
};

const whoami = async (session: string) => {
  return await fetch(`${process.env.NEXT_PUBLIC_API_GW}/get-user`, {
    method: 'GET',
    cache: 'force-cache',
    credentials: 'include',
    headers: {
      cookie: session,
      'Content-Type': 'application/json',
    },
  });
};

function isUserHaveOrg(userObj: { orgId: string; orgName: string }) {
  return userObj.orgId !== '' || userObj.orgName !== '';
}

function isValidUser(userObj: {
  id: string;
  name: string;
  email: string;
  token: string;
  isEmailVerified: boolean;
}) {
  return (
    userObj.id !== '' &&
    userObj.name !== '' &&
    userObj.email !== '' &&
    userObj.token !== '' &&
    userObj.isEmailVerified
  );
}

function decodeBase64Token(token: string): string {
  return Buffer.from(token, 'base64').toString('utf8');
}

function getUserFromToken(token: string): User {
  const parseToken = decodeBase64Token(token);
  const s = parseToken.split(';');
  return {
    id: s[2] ?? '',
    role: s[5] ?? '',
    name: s[3] ?? '',
    email: s[4] ?? '',
    orgId: s[0] ?? '',
    token: token ?? '',
    orgName: s[1] ?? '',
    isEmailVerified: s[6] === 'true' ? true : false ?? false,
  };
}

const middleware = async (request: NextRequest) => {
  const { pathname } = request.nextUrl;
  const cookieStore = cookies();
  const session = cookieStore.get('ukama_session');
  let cookieToken = cookieStore.get('token')?.value || '';
  const response = NextResponse.next();
  let userObj: User = {
    id: '',
    role: Role_Type.RoleInvalid,
    name: '',
    email: '',
    orgId: '',
    token: '',
    orgName: '',
    isEmailVerified: false,
  };

  if (!session) {
    return NextResponse.redirect(
      new URL('/auth/login', process.env.NEXT_PUBLIC_AUTH_APP_URL),
    );
  } else if (session && cookieToken) {
    userObj = getUserFromToken(cookieToken);
  } else if (session && !cookieToken) {
    const res = await whoami(`ukama_session=${session.value}`);
    if (!res.ok) {
      return NextResponse.rewrite(
        new URL('/unauthorized', process.env.NEXT_PUBLIC_APP_URL),
      );
    } else {
      const jsonRes = await res.json();
      userObj.id = jsonRes.userId;
      userObj.role = jsonRes.role;
      userObj.name = jsonRes.name;
      userObj.email = jsonRes.email;
      userObj.orgId = jsonRes.orgId;
      userObj.token = jsonRes.token;
      userObj.orgName = jsonRes.orgName;
      userObj.isEmailVerified = jsonRes.isEmailVerified;
    }
  }

  // TODO: ADD CHECK TO CROSS MATCH cookieToken and whoami resposne token

  if (!userObj.isEmailVerified) {
    return NextResponse.redirect(
      new URL('/user/verification', process.env.NEXT_PUBLIC_AUTH_APP_URL),
    );
  } else if (userObj.token && !cookieToken) {
    response.cookies.set('token', userObj.token, {
      path: '/',
      name: 'token',
      httpOnly: true,
      sameSite: 'lax',
      value: userObj.token,
      domain: process.env.NEXT_PUBLIC_APP_DOMAIN,
      secure: process.env.NODE_ENV === 'production',
      expires: new Date(Date.now() + 1000 * 60 * 60 * 24),
    });
  } else if (!userObj.token) {
    return NextResponse.rewrite(
      new URL('/unauthorized', process.env.NEXT_PUBLIC_APP_URL),
    );
  }

  response.headers.set('role', userObj.role);
  response.headers.set('name', userObj.name);
  response.headers.set('user-id', userObj.id);
  response.headers.set('email', userObj.email);
  response.headers.set('org-id', userObj.orgId);
  response.headers.set('org-name', userObj.orgName);

  if (
    (pathname.includes('/console') || pathname === '/') &&
    !isUserHaveOrg(userObj)
  ) {
    return NextResponse.redirect(
      new URL('/onboarding', process.env.NEXT_PUBLIC_APP_URL),
    );
  } else if (
    pathname.includes('/console') &&
    (userObj.role === Role_Type.RoleInvalid ||
      userObj.role === Role_Type.RoleUser)
  ) {
    return NextResponse.redirect(
      new URL('/403', process.env.NEXT_PUBLIC_APP_URL),
    );
  } else if (
    pathname.includes('/manage') &&
    userObj.role !== Role_Type.RoleOwner &&
    userObj.role !== Role_Type.RoleAdmin
  ) {
    return NextResponse.redirect(
      new URL('/unauthorized', process.env.NEXT_PUBLIC_APP_URL),
    );
  } else if (
    pathname === '/' &&
    isUserHaveOrg(userObj) &&
    isValidUser(userObj)
  ) {
    return NextResponse.redirect(
      new URL('/console/home', process.env.NEXT_PUBLIC_APP_URL),
    );
  }
  return response;
};

export { middleware };
