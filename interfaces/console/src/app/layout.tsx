import AppContextWrapper from '@/context';
import '@/styles/global.css';
import AppThemeProvider from '@/theme/AppThemeProvider';
import { ApolloWrapper } from '@/wrappers/apolloWrapper';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { cookies, headers } from 'next/headers';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Ukama Console',
  description: 'Ukama Conosle app to manage your network',
  icons: {
    icon: [
      {
        url: '/svg/ulogo.svg',
        href: '/svg/ulogo.svg',
      },
    ],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const headersList = headers();
  const cookieStore = cookies();
  const cookieTheme = cookieStore.get('theme') || {
    name: 'theme',
    value: 'light',
  };
  const role = headersList.get('role');
  const name = headersList.get('name');
  const email = headersList.get('email');
  const orgId = headersList.get('org-id');
  const userId = headersList.get('user-id');
  const orgName = headersList.get('org-name');
  const tokenStr = cookieStore.get('token') ?? {
    name: 'token',
    value: '',
  };
  return (
    <html lang="en">
      <body className={inter.className}>
        <ApolloWrapper>
          <AppContextWrapper
            token={tokenStr.value}
            initalUserValues={{
              id: userId ?? '',
              name: name ?? '',
              role: role ?? '',
              email: email ?? '',
              orgId: orgId ?? '',
              orgName: orgName ?? '',
            }}
          >
            <AppThemeProvider themeCookie={cookieTheme?.value ?? 'dark'}>
              {children}
            </AppThemeProvider>
          </AppContextWrapper>
        </ApolloWrapper>
      </body>
    </html>
  );
}
