'use client';

import { createContext } from 'react';
import { User } from './supabase/user';

export const UserContext = createContext<User | null>(null);