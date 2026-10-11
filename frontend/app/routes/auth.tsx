import { useEffect, useRef } from "react";
import { useLocation, useNavigate } from "react-router";
import { useAuthStore } from "~/lib/auth";

export const meta = () => ([
  { title: "Resumind | Sign In" },
  { name: "description", content: "Sign in to your Resumind account" },
]);

declare global {
  interface Window {
    google?: {
      accounts: {
        id: {
          initialize: (config: object) => void;
          renderButton: (element: HTMLElement, config: object) => void;
          prompt: () => void;
        };
      };
    };
  }
}

const Auth = () => {
  const { isAuthenticated, checkAuth } = useAuthStore();
  const location = useLocation();
  const navigate = useNavigate();
  const googleBtnRef = useRef<HTMLDivElement>(null);


  const next = location.search.split("next=")[1] || "/";

  useEffect(() => {
    if (isAuthenticated) navigate(next);
  }, [isAuthenticated, next, navigate]);


  useEffect(() => {
    const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID;

    if (!clientId) {
      console.error("VITE_GOOGLE_CLIENT_ID is not set in your .env file");
      return;
    }


    const initGoogle = () => {
      if (!window.google || !googleBtnRef.current) return;

      window.google.accounts.id.initialize({
        client_id: clientId,
        callback: async (response: { credential: string }) => {

          const res = await fetch("/api/auth/google", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ idToken: response.credential }),
          });

          if (res.ok) {
            await checkAuth();
            navigate(next);
          } else {
            const err = await res.json();
            console.error("Google login failed:", err.error);
          }
        },
      });


      window.google.accounts.id.renderButton(googleBtnRef.current, {
        theme: "outline",
        size: "large",
        width: 280,
        text: "signin_with",
        shape: "rectangular",
      });
    };


    if (window.google) {
      initGoogle();
    } else {

      const script = document.querySelector(
        'script[src="https://accounts.google.com/gsi/client"]'
      );
      if (script) {
        script.addEventListener("load", initGoogle);
        return () => script.removeEventListener("load", initGoogle);
      }
    }
  }, [checkAuth, navigate, next]);

  return (
    <main className="bg-[url('/images/bg-auth.svg')] bg-cover min-h-screen flex items-center justify-center">
      <div className="gradient-border shadow-lg">
        <section className="flex flex-col gap-8 bg-white rounded-2xl p-10">
          <div className="flex flex-col items-center gap-2 text-center">
            <h1>Welcome</h1>
            <h2>Sign in to Continue Your Job Journey</h2>
          </div>

          {}
          <div className="flex justify-center">
            <div ref={googleBtnRef} id="google-signin-btn" />
          </div>

          <p className="text-xs text-center text-gray-400">
            By signing in, you agree to our terms of service.
          </p>
        </section>
      </div>
    </main>
  );
};

export default Auth;
