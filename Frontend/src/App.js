import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import SignUp from './Signup.js'; // Ensure the correct path to your SignUp component
import Chat from './Chat'; // Ensure the correct path to your Chat component
import Login from './Login.js';
import UsersPage from './Users.js';

function App() {
    return (
        <Router>
            <div style={{ fontFamily: 'Arial, sans-serif', padding: '20px' }}>
                <Routes>
                    <Route path="/signup" element={<SignUp />} />
                    <Route path="/chat" element={<Chat />} /> {/* The chat component for WebSocket messages */}
                    <Route path="/login" element={<Login />} /> {/* The chat component for WebSocket messages */}
                    <Route path="/users" element={<UsersPage />} /> {/* The chat component for WebSocket messages */}
                </Routes>
            </div>
        </Router>
    );
}

export default App;
