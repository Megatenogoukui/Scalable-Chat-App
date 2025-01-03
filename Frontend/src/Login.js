import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [statusMessage, setStatusMessage] = useState("");
  const navigate = useNavigate();

  const handleLogin = async () => {
    const apiUrl = "http://localhost:8080/login"; // Adjust to your backend API URL

    // Basic validation
    if (!email || !password) {
      setStatusMessage("Please fill out all fields.");
      return;
    }

    try {
      const response = await fetch(apiUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email,
          password,
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Login failed");
      }

      const data = await response.json();

      // Assuming the token and _id are returned in the response
      if (data.token ) {
        localStorage.setItem("token", data.token);
        // Assuming the token is returned in the response
        if (data.token) {
          localStorage.setItem("token", data.token);
          try {
            const decodedToken = jwtDecode(data.token);
            const userId = decodedToken._id; // Adjust the key based on your token structure
            if (userId) {
              localStorage.setItem("_id", userId);
              setStatusMessage("Login successful!");
              // Redirect to /chat after successful login
              navigate("/users");
            } else {
              throw new Error("User ID not found in token");
            }
          } catch (error) {
            throw new Error("Error decoding token: " + error.message);
          }
        } else {
          throw new Error("Token not received");
        }
       
      } else {
        throw new Error("Token or user ID not received");
      }
    } catch (error) {
      console.error("Login error:", error);
      setStatusMessage(error.message || "Login failed. Please try again.");
    }
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh",
        fontFamily: "Arial, sans-serif",
        backgroundColor: "#f9f9f9",
      }}
    >
      <div
        style={{
          width: "300px",
          padding: "20px",
          border: "1px solid #ddd",
          borderRadius: "10px",
          backgroundColor: "#fff",
        }}
      >
        <h2 style={{ textAlign: "center", marginBottom: "20px" }}>Login</h2>

        <div style={{ marginBottom: "15px" }}>
          <label style={{ display: "block", marginBottom: "5px" }}>Email</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            style={{
              width: "100%",
              padding: "10px",
              border: "1px solid #ddd",
              borderRadius: "5px",
            }}
          />
        </div>

        <div style={{ marginBottom: "15px" }}>
          <label style={{ display: "block", marginBottom: "5px" }}>
            Password
          </label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            style={{
              width: "100%",
              padding: "10px",
              border: "1px solid #ddd",
              borderRadius: "5px",
            }}
          />
        </div>

        <button
          onClick={handleLogin}
          style={{
            width: "100%",
            padding: "10px",
            border: "none",
            backgroundColor: "#0288d1",
            color: "white",
            borderRadius: "5px",
            cursor: "pointer",
          }}
        >
          Login
        </button>

        <p style={{ marginTop: "10px", textAlign: "center" }}>
          Don't have an account?{" "}
          <span
            onClick={() => navigate("/signup")}
            style={{ color: "#0288d1", cursor: "pointer" }}
          >
            Sign Up
          </span>
        </p>

        {statusMessage && (
          <p
            style={{ marginTop: "10px", color: "#FF5722", textAlign: "center" }}
          >
            {statusMessage}
          </p>
        )}
      </div>
    </div>
  );
}

export default Login;
