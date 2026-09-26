import { createServerClient } from '@supabase/ssr'
import { NextResponse, type NextRequest } from 'next/server'
import { getUserData } from './user'
import { 
  UserRoleEnum, 
  ADMIN_WHITELISTED_PAGES, 
  PUBLIC_PAGES, 
  STUDENT_WHITELISTED_PAGES, 
  TUTOR_WHITELISTED_PAGES 
} from '@/app/(admin)/role-authorization'


export async function updateSession(request: NextRequest) {
  let supabaseResponse = NextResponse.next({
    request,
  })

  const supabase = createServerClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
    {
      cookies: {
        getAll() {
          return request.cookies.getAll()
        },
        setAll(cookiesToSet) {
          cookiesToSet.forEach(({ name, value, options }) => request.cookies.set(name, value))
          supabaseResponse = NextResponse.next({
            request,
          })
          cookiesToSet.forEach(({ name, value, options }) =>
            supabaseResponse.cookies.set(name, value, options)
          )
        },
      },
    }
  )

  const {
    data: { user },
  } = await supabase.auth.getUser()


  if (
    !user &&
    !PUBLIC_PAGES.some(pageRegex => request.nextUrl.pathname.match(pageRegex))
  ) {
    // no user, potentially respond by redirecting the user to the login page
    const url = request.nextUrl.clone()
    url.pathname = '/sign-in'
    return NextResponse.redirect(url)
  } else if (!user) {
    return supabaseResponse
  }

  const userData = await getUserData(user.id)
  
  let redirect: NextResponse<unknown> | null = null;
  switch (userData.role) {
    case UserRoleEnum.STUDENT:
      redirect = checkAndRedirect(request, STUDENT_WHITELISTED_PAGES, '/')
      if (redirect) {
        return redirect
      }
      break
    case UserRoleEnum.TUTOR:
      redirect = checkAndRedirect(request, TUTOR_WHITELISTED_PAGES, '/admin/')
      if (redirect) {
        return redirect
      }
      break
    case UserRoleEnum.ADMIN:
      redirect = checkAndRedirect(request, ADMIN_WHITELISTED_PAGES, '/admin/')
      if (redirect) {
        return redirect
      }
      break
  }

  return supabaseResponse
}

function checkAndRedirect(request: NextRequest, whitelistedPages: RegExp[], redirectUrl: string) {
  if (!whitelistedPages.some(pageRegex => request.nextUrl.pathname.match(pageRegex))) {
    const url = request.nextUrl.clone()
    url.pathname = redirectUrl
    return NextResponse.redirect(url)
  }
  return null;
}
