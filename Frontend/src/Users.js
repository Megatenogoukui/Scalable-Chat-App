import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

function UsersPage() {
  const [users, setUsers] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const currentUserId = localStorage.getItem("_id"); // Get current user's ID from localStorage

  // If no user is logged in, navigate to the login page
  useEffect(() => {
    if (!currentUserId) {
      navigate("/login"); // Redirect to login page if no current user ID found
    }
  }, [currentUserId, navigate]);

  const apiUrl = `http://localhost:8080/GetAllUsers?userId=${currentUserId}`;

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await fetch(apiUrl, {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        });

        if (!response.ok) {
          throw new Error("Failed to fetch users");
        }

        const data = await response.json();

        if (!data || !Array.isArray(data?.result)) {
          throw new Error("Invalid response format: Expected an array");
        }

        setUsers(data?.result || []);
      } catch (err) {
        setError(err.message || "Error fetching users");
      } finally {
        setLoading(false);
      }
    };

    if (currentUserId) {
      fetchUsers();
    }
  }, [apiUrl, currentUserId]);

  const handleConnect = (userId, currentUserId, con_id) => {
    navigate(`/chat?userId1=${userId}&userId2=${currentUserId}&conversationId=${con_id}`);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg">Loading users...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen text-red-500">
        Error: {error}
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6 px-4">All Users</h1>
      {users.length === 0 ? (
        <p className="text-center text-gray-500 py-4">No users found.</p>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full border-collapse">
              <thead className="bg-gray-100">
                <tr>
                  <th className="text-left px-6 py-3 border-b border-gray-200 bg-gray-50 text-xs font-medium text-gray-500 uppercase tracking-wider">
                    ID
                  </th>
                  <th className="text-left px-6 py-3 border-b border-gray-200 bg-gray-50 text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Email
                  </th>
                  <th className="text-left px-6 py-3 border-b border-gray-200 bg-gray-50 text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Username
                  </th>
                  <th className="text-left px-6 py-3 border-b border-gray-200 bg-gray-50 text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {users.map((user) => {
                  const isCurrentUser = user?._id === currentUserId;
                  return (
                    <tr key={user?._id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 text-sm font-mono text-gray-600 truncate max-w-xs">
                        {user?._id || "N/A"}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-900 truncate max-w-xs">
                        {user?.email || "N/A"}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-900">
                        {user?.user_name || "N/A"}
                      </td>
                      <td className="px-6 py-4 text-sm">
                        <button
                          onClick={() => handleConnect(user?._id, currentUserId, user?.conversation_id)}
                          disabled={isCurrentUser}
                          className={`px-4 py-2 rounded-md text-sm font-medium transition-colors duration-150 ease-in-out ${
                            isCurrentUser
                              ? "bg-gray-300 cursor-not-allowed"
                              : "bg-blue-500 hover:bg-blue-600 text-white"
                          }`}
                          title={isCurrentUser ? "Cannot connect to yourself" : "Connect with user"}
                        >
                          {isCurrentUser ? "You" : "Connect"}
                        </button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}

export default UsersPage;
