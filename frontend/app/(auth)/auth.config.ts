import type { NextAuthConfig } from 'next-auth';

export const authConfig = {
  pages: {
    signIn: '/login',
    newUser: '/',
  },
  providers: [
    // added later in auth.ts since it requires bcrypt which is only compatible with Node.js
    // while this file is also used in non-Node.js environments
  ],
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.id = user.id;
        token.role = user.role;
      }
      return token;
    },
    async session({ session, token }) {
      if (session.user) {
        session.user.id = token.id as string;
        session.user.role = token.role as 'student' | 'admin';
      }
      return session;
    },
    authorized({ auth, request: { nextUrl } }) {
      const isLoggedIn = !!auth?.user;
      const isAdmin = auth?.user?.role === 'admin';
      const isOnAdmin = nextUrl.pathname.startsWith('/admin');
      const isOnChat = nextUrl.pathname.startsWith('/');
      const isOnRegister = nextUrl.pathname.startsWith('/register');
      const isOnLogin = nextUrl.pathname.startsWith('/login');

      if (isLoggedIn && isOnLogin) {
        return Response.redirect(new URL(isAdmin ? '/admin' : '/', nextUrl as unknown as URL));
      }

      if (isOnAdmin) {
        return isAdmin; // Protect admin routes
      }

      if (isOnRegister || isOnLogin) {
        return true; // Always allow access to register and login pages
      }

      if (isOnChat) {
        if (isLoggedIn) {
          // Redirect admins to admin dashboard
          if (isAdmin) {
            return Response.redirect(new URL('/admin', nextUrl));
          }
          return true;
        }
        return false; // Redirect unauthenticated users to login
      }

      return true;
    },
  },
} satisfies NextAuthConfig;
