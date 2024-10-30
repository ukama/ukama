/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Role_Type } from '@/client/graphql/generated/subscriptions';
import { cookies } from 'next/headers';
import type { NextRequest } from 'next/server';
import { NextResponse } from 'next/server';

type User = {
  id: string;
  role: string;
  name: string;
  email: string;
  orgId: string;
  token: string;
  orgName: string;
  isShowWelcome: boolean;
  isEmailVerified: boolean;
};

const USER_INIT = {
  id: '',
  name: '',
  token: '',
  email: '',
  orgId: '',
  orgName: '',
  isShowWelcome: false,
  isEmailVerified: false,
  role: Role_Type.RoleInvalid,
};

const whoami = async (session: string) => {
  return await fetch(`${process.env.NEXT_PUBLIC_API_GW_4SS}/get-user`, {
    method: 'GET',
    cache: 'no-store',
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
  try {
    const parseToken = decodeBase64Token(token);
    const parts = parseToken.split(';');

    if (parts.length < 8) {
      return USER_INIT;
    }

    const [
      orgId,
      orgName,
      id,
      name,
      email,
      role,
      isEmailVerified,
      isShowWelcome,
    ] = parts;
    return {
      id,
      role,
      name,
      email,
      orgId,
      token,
      orgName,
      isShowWelcome: isShowWelcome.includes('true'),
      isEmailVerified: isEmailVerified.includes('true'),
    };
  } catch (error) {
    return USER_INIT;
  }
}

const getUserObject = async (session: string, cookieToken: string) => {
  if (cookieToken) {
    return getUserFromToken(cookieToken);
  } else {
    const res = await whoami(`ukama_session=${session}`);
    if (!res.ok) {
      throw new Error('Unauthorized');
    }
    const jsonRes = await res.json();
    return {
      id: jsonRes.userId,
      role: jsonRes.role,
      name: jsonRes.name,
      email: jsonRes.email,
      orgId: jsonRes.orgId,
      token: jsonRes.token,
      orgName: jsonRes.orgName,
      isShowWelcome: jsonRes.isShowWelcome,
      isEmailVerified: jsonRes.isEmailVerified,
    };
  }
};

const middleware = async (request: NextRequest) => {
  const response = NextResponse.next();
  const cookieStore = cookies();
  const { pathname } = request.nextUrl;

  if (pathname.includes('/ping')) {
    return response;
  }

  if (request.url.includes('logout')) {
    // cookieStore.set('token', '', {
    //   path: '/',
    //   name: 'token',
    //   secure: false,
    //   httpOnly: true,
    //   sameSite: 'lax',
    //   value: '',
    //   domain: process.env.NEXT_PUBLIC_APP_DOMAIN,
    //   expires: new Date(Date.now() + 1000 * 60 * 60 * 24),
    // });
    cookieStore.delete('token');
    response.cookies.delete('token');
    return response;
  }

  const session = cookieStore.get('ukama_session');
  const cookieToken = cookieStore.get('token')?.value ?? '';

  if (!session) {
    return NextResponse.redirect(
      new URL(
        `/auth/login?return_to=${request.url}`,
        process.env.NEXT_PUBLIC_AUTH_APP_URL,
      ),
    );
  }

  if (pathname.includes('/refresh')) {
    response.cookies.delete('token');
    return response;
  }

  let userObj: User = USER_INIT;
  try {
    userObj = await getUserObject(session.value, cookieToken);
  } catch (error) {
    return NextResponse.rewrite(
      new URL('/unauthorized', process.env.NEXT_PUBLIC_APP_URL),
    );
  }

  if (!userObj?.isEmailVerified) {
    return NextResponse.redirect(
      new URL('/user/verification', process.env.NEXT_PUBLIC_AUTH_APP_URL),
    );
  }

  if (userObj.token && !cookieToken) {
    response.cookies.set('token', userObj.token, {
      path: '/',
      name: 'token',
      secure: false,
      httpOnly: true,
      sameSite: 'lax',
      value: userObj.token,
      domain: process.env.NEXT_PUBLIC_APP_DOMAIN,
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

  if (userObj.isShowWelcome) {
    console.log("Redirecting to '/welcome'");
    return NextResponse.redirect(
      new URL('/welcome', process.env.NEXT_PUBLIC_APP_URL),
    );
  }

  if (
    (pathname.includes('/console') || pathname === '/') &&
    !isUserHaveOrg(userObj)
  ) {
    console.log("Redirecting to '/onboarding' ");
    return NextResponse.redirect(
      new URL('/onboarding', process.env.NEXT_PUBLIC_APP_URL),
    );
  }

  if (
    pathname.includes('/console') &&
    (userObj.role === Role_Type.RoleInvalid ||
      userObj.role === Role_Type.RoleUser)
  ) {
    console.log("Redirecting to '/403' ");
    return NextResponse.redirect(
      new URL('/403', process.env.NEXT_PUBLIC_APP_URL),
    );
  }

  if (pathname.includes('/welcome') && userObj.role !== Role_Type.RoleOwner) {
    console.log("Redirecting to '/' ");
    return NextResponse.redirect(new URL('/', process.env.NEXT_PUBLIC_APP_URL));
  }

  if (pathname.includes('/manage') && userObj.role !== Role_Type.RoleOwner) {
    console.log("Redirecting to '/unauthorized' ");
    return NextResponse.redirect(
      new URL('/unauthorized', process.env.NEXT_PUBLIC_APP_URL),
    );
  }

  if (
    (pathname.includes('/console/nodes') ||
      pathname.includes('/console/sites')) &&
    userObj.role !== Role_Type.RoleOwner &&
    userObj.role !== Role_Type.RoleAdmin
  ) {
    console.log("Redirecting to '/unauthorized' ");
    return NextResponse.redirect(
      new URL('/unauthorized', process.env.NEXT_PUBLIC_APP_URL),
    );
  }

  if (pathname === '/' && isUserHaveOrg(userObj) && isValidUser(userObj)) {
    console.log("Redirecting to '/console/home'");
    return NextResponse.redirect(
      new URL('/console/home', process.env.NEXT_PUBLIC_APP_URL),
    );
  }

  return response;
};

export { middleware };
