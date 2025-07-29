import { NextResponse } from 'next/server';
import { execSync } from 'child_process';

export async function GET() {
  try {
    const output = execSync('../n8n-ops status --json', {
      cwd: process.cwd(),
      encoding: 'utf-8'
    });
    const data = JSON.parse(output);
    return NextResponse.json(data);
  } catch {
    return NextResponse.json({ status: 'unavailable' });
  }
}
