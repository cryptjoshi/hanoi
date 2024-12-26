import { NextResponse } from 'next/server';
import { getIronSession } from 'iron-session';
import { sessionOptions } from '@/lib';
import { cookies } from 'next/headers';

export async function GET() {
  const session = await getIronSession(cookies(), sessionOptions);

  return NextResponse.json({ session: session || null });
}