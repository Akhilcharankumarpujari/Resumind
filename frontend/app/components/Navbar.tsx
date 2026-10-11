import { Link, useNavigate } from "react-router";
import { useAuthStore } from "~/lib/auth";

const Navbar = () => {
  const { user, isAuthenticated, signOut } = useAuthStore();
  const navigate = useNavigate();

  const handleSignOut = async () => {
    await signOut();
    navigate("/auth");
  };

  return (
    <nav className="navbar">
      <Link to="/">
        <p className="text-2xl font-bold text-gradient">RESUMIND</p>
      </Link>

      <div className="flex items-center gap-3">
        {isAuthenticated && user ? (
          <>
            {}
            {user.picture && (
              <img
                src={user.picture}
                alt={user.name}
                className="w-8 h-8 rounded-full object-cover border-2 border-white shadow"
                referrerPolicy="no-referrer"
              />
            )}
            <span className="text-sm font-medium text-gray-700 hidden sm:block">
              {user.name}
            </span>
            <button
              onClick={handleSignOut}
              className="text-sm text-gray-500 hover:text-red-500 transition-colors"
            >
              Sign Out
            </button>
          </>
        ) : null}

        <Link to="/upload" className="primary-button w-fit">
          Upload Resume
        </Link>
      </div>
    </nav>
  );
};

export default Navbar;
