import type { Metadata } from "next";
import { Toaster } from "sonner";

import { ThemeProvider } from "@/components/theme-provider";
import { Providers } from "./providers";

import "./globals.css";
import { getUser } from "@/lib/supabase/user";

export const metadata: Metadata = {
  metadataBase: new URL("https://segp-frontend.vercel.app"),
  title: "Imperial College Tutor Agent",
  description: "A web app for Imperial College students to message an AI agent " +
    "as a first call before their personal tutor.",
};

export const viewport = {
  maximumScale: 1, // Disable auto-zoom on mobile Safari
};

const LIGHT_THEME_COLOR = "hsl(0 0% 100%)";
const DARK_THEME_COLOR = "hsl(255 6% 13%)";
const THEME_COLOR_SCRIPT = `\
(function() {
  var html = document.documentElement;
  var meta = document.querySelector('meta[name="theme-color"]');
  if (!meta) {
    meta = document.createElement('meta');
    meta.setAttribute('name', 'theme-color');
    document.head.appendChild(meta);
  }
  function updateThemeColor() {
    var isDark = html.classList.contains('dark');
    meta.setAttribute('content', isDark ? '${DARK_THEME_COLOR}' : '${LIGHT_THEME_COLOR}');
  }
  var observer = new MutationObserver(updateThemeColor);
  observer.observe(html, { attributes: true, attributeFilter: ['class'] });
  updateThemeColor();
})();`;

export default async function RootLayout({children}: {children: React.ReactNode}) {
  const user = await getUser();

  return (
    <html
      lang="en"
      // `next-themes` injects an extra classname to the body element to avoid
      // visual flicker before hydration. Hence the `suppressHydrationWarning`
      // prop is necessary to avoid the React hydration mismatch warning.
      // https://github.com/pacocoursey/next-themes?tab=readme-ov-file#with-app
      suppressHydrationWarning
    >
      <head>
        <script
          dangerouslySetInnerHTML={{
            __html: THEME_COLOR_SCRIPT,
          }}
        />
      </head>
      <body 
        className="antialiased vsc-initialized"
        data-new-gr-c-s-check-loaded="14.1224.0"
        data-gr-ext-installed=""
      >
        <Providers user={user}>
          <ThemeProvider
            attribute="class"
            defaultTheme="system"
            enableSystem
            disableTransitionOnChange
          >
            <Toaster richColors position="top-center" />
            {children}
          </ThemeProvider>
        </Providers>
      </body>
    </html>
  );
}
