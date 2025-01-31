import NextAuth from 'next-auth';

declare module 'next-auth' {
  interface User {
    id: string;
    role: 'admin' | 'student';
  }

  interface Session {
    user: User;
  }

  interface JWT {
    id: string;
    role: 'admin' | 'student';
  }
}
