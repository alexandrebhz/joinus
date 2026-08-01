import { redirect } from 'next/navigation'
import { getGoogleOAuthUrl } from '@/infrastructure/auth/oauth'

/**
 * Frontend entry for Google OAuth.
 * The browser must navigate here (not XHR) so cookies/redirects work;
 * this page then sends the user to the backend OAuth start endpoint.
 */
export default function GoogleAuthStartPage() {
  redirect(getGoogleOAuthUrl())
}
