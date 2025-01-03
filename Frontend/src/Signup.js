import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom'; // Import useNavigate for navigation
import { jwtDecode } from "jwt-decode";


function SignUp() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [userName, setUserName] = useState('');
    const [statusMessage, setStatusMessage] = useState('');
    const navigate = useNavigate(); // Hook for navigation

    const handleSignUp = async () => {
        const apiUrl = "http://localhost:8080/sign_up"; // Adjust to your backend API URL

        // Basic validation
        if (!email || !password || !userName) {
            setStatusMessage("Please fill out all fields.");
            return;
        }

        try {
            const response = await fetch(apiUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    email,
                    password,
                    user_name: userName, // Match the expected key from the backend
                }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || "Signup failed");
            }

            const data = await response.json();

            // Assuming the token is returned as "token" in the response
            if (data.token) {
                localStorage.setItem('token', data.token);
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
                throw new Error("Token not received");
            }
        } catch (error) {
            console.error("Signup error:", error);
            setStatusMessage(error.message || "Signup failed. Please try again.");
        }
    };

    return (
        <div style={{
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh',
            fontFamily: 'Arial, sans-serif',
            backgroundColor: '#f9f9f9',
        }}>
            <div style={{
                width: '300px',
                padding: '20px',
                border: '1px solid #ddd',
                borderRadius: '10px',
                backgroundColor: '#fff',
            }}>
                <h2 style={{ textAlign: 'center', marginBottom: '20px' }}>Sign Up</h2>

                <div style={{ marginBottom: '15px' }}>
                    <label style={{ display: 'block', marginBottom: '5px' }}>Email</label>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        style={{
                            width: '100%',
                            padding: '10px',
                            border: '1px solid #ddd',
                            borderRadius: '5px',
                        }}
                    />
                </div>

                <div style={{ marginBottom: '15px' }}>
                    <label style={{ display: 'block', marginBottom: '5px' }}>Password</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        style={{
                            width: '100%',
                            padding: '10px',
                            border: '1px solid #ddd',
                            borderRadius: '5px',
                        }}
                    />
                </div>

                <div style={{ marginBottom: '15px' }}>
                    <label style={{ display: 'block', marginBottom: '5px' }}>Username</label>
                    <input
                        type="text"
                        value={userName}
                        onChange={(e) => setUserName(e.target.value)}
                        style={{
                            width: '100%',
                            padding: '10px',
                            border: '1px solid #ddd',
                            borderRadius: '5px',
                        }}
                    />
                </div>

                <button
                    onClick={handleSignUp}
                    style={{
                        width: '100%',
                        padding: '10px',
                        border: 'none',
                        backgroundColor: '#0288d1',
                        color: 'white',
                        borderRadius: '5px',
                        cursor: 'pointer',
                    }}
                >
                    Sign Up
                </button>

                {statusMessage && (
                    <p style={{ marginTop: '10px', color: '#FF5722', textAlign: 'center' }}>
                        {statusMessage}
                    </p>
                )}
            </div>
        </div>
    );
}

export default SignUp;
