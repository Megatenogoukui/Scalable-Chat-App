import React, { useState, useEffect, useRef, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import SignUp from "./Signup.js";
import Chat from "./Chat.js";

function App() {
  const [messages, setMessages] = useState([]);
  const [message, setMessage] = useState("");
  const [connectionStatus, setConnectionStatus] = useState("Connecting...");
  const [clientId, setClientId] = useState(localStorage.getItem("_id") || null);
  const queryParams = new URLSearchParams(window.location.search);

  const socketRef = useRef(null);
  const reconnectTimeout = useRef(null);

  const socketUrl =
    process.env.REACT_APP_SOCKET_URL ||
    "ws://localhost:8080/websocket_connection";
    // Construct the WebSocket URL with query parameters
  const wsUrlWithQuery = `${socketUrl}?userId=${encodeURIComponent(clientId)}`;

  const apiEndpoint =
    process.env.REACT_APP_API_ENDPOINT || "http://localhost:8080/get_messages";

  const navigate = useNavigate();

  const addMessage = useCallback((text, type) => {
    setMessages((prevMessages) => [
      ...prevMessages,
      {
        id: `${Date.now()}-${Math.random()}`,
        text,
        type,
      },
    ]);
  }, []);

  const fetchPreviousMessages = useCallback(async () => {
    try {
      const channelName = "Messages";
      const queryParams = new URLSearchParams(window.location.search);

      const apiUrl = `${apiEndpoint}?channel_name=${encodeURIComponent(
        channelName
      )}&userId1=${localStorage.getItem("_id")}&userId2=${queryParams.get(
        "userId1"
      )}`;

      const response = await fetch(apiUrl);
      if (!response.ok) {
        throw new Error("Failed to fetch previous messages");
      }
      const data = await response.json();

      const storedClientId = localStorage.getItem("_id");
      if (!storedClientId) {
        console.warn("No client ID found in local storage");
        navigate("/login");
        return;
      }

      const newMessages =
        data?.messages?.map((rawMessage) => ({
          id: rawMessage.messageId,
          text: rawMessage.message,
          type: rawMessage.sender_id === storedClientId ? "sent" : "received",
        })) || [];

      setMessages(newMessages);
    } catch (error) {
      console.error("Error fetching previous messages:", error);
    }
  }, [apiEndpoint, navigate]);

  const connectSocket = useCallback(() => {
    if (socketRef.current) {
      socketRef.current.close();
    }

    socketRef.current = new WebSocket(wsUrlWithQuery);

    socketRef.current.onopen = () => {
      console.log("WebSocket connection opened");
      setConnectionStatus("Connected");
    };

    socketRef.current.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.clientId) {
          localStorage.setItem("_id", data.clientId);
          setClientId(data.clientId);
        } else {
          addMessage(data.message || "Unknown message format", "received");
        }
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
        addMessage(event.data || "Invalid message format", "received");
      }
    };

    socketRef.current.onerror = (error) => {
      console.error("WebSocket error:", error);
      setConnectionStatus("Error");
    };

    socketRef.current.onclose = (event) => {
      if (!event.wasClean) {
        console.log(
          "WebSocket connection closed unexpectedly, attempting to reconnect..."
        );
        setConnectionStatus("Disconnected");
        reconnectTimeout.current = setTimeout(connectSocket, 3000);
      }
    };
  }, [wsUrlWithQuery, addMessage]);

  const sendMessage = useCallback(() => {
    const trimmedMessage = message.trim();
    if (!trimmedMessage) return;
  
    const userId = localStorage.getItem("_id");
    const queryParams = new URLSearchParams(window.location.search);
    const recipientId = queryParams.get("userId1");
    const conversationId = queryParams.get("conversationId"); // Extract conversationId here
  
    if (!userId || !recipientId || !conversationId) {
      navigate("/"); // Redirect to home or login if necessary data is missing
      return;
    }
  
    const messagePayload = JSON.stringify({
      message: trimmedMessage,
      sender_id: userId,
      recipient_id: recipientId,
      conversationId, // Include conversationId in the payload
    });
  
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      try {
        socketRef.current.send(messagePayload);
        addMessage(trimmedMessage, "sent");
        setMessage(""); // Clear input after sending
      } catch (error) {
        console.error("Failed to send message:", error);
      }
    }
  }, [message, addMessage, navigate]);

  const handleKeyUp = (event) => {
    if (event.key === "Enter") {
      sendMessage();
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("_id");
    localStorage.removeItem("token");
    navigate("/login");
  };

  useEffect(() => {
    fetchPreviousMessages();
    connectSocket();

    return () => {
      clearTimeout(reconnectTimeout.current);
      if (socketRef.current) {
        socketRef.current.close();
      }
    };
  }, [connectSocket, fetchPreviousMessages]);

  return (
    <div
      style={{
        fontFamily: "Arial, sans-serif",
        display: "flex",
        flexDirection: "column",
        height: "92vh",
        backgroundColor: "#f9f9f9",
      }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          padding: "10px",
          backgroundColor:
            connectionStatus === "Connected" ? "#4CAF50" : "#FF5722",
          color: "white",
        }}
      >
        <span>
          Connection Status: {connectionStatus}{" "}
          {clientId && `(Client ID: ${clientId})`}
        </span>
        <div>
          <button
            onClick={() => navigate("/users")}
            style={{
              padding: "5px 10px",
              marginRight: "10px",
              backgroundColor: "#0288d1",
              color: "white",
              border: "none",
              borderRadius: "5px",
              cursor: "pointer",
            }}
          >
            Back
          </button>
          <button
            onClick={handleLogout}
            style={{
              padding: "5px 10px",
              backgroundColor: "#FF5722",
              color: "white",
              border: "none",
              borderRadius: "5px",
              cursor: "pointer",
            }}
          >
            Logout
          </button>
        </div>
      </div>

      <div
        style={{
          flex: 1,
          display: "flex",
          flexDirection: "column",
          padding: "10px",
          overflowY: "auto",
          backgroundColor: "#fff",
          border: "1px solid #ddd",
          margin: "10px",
          borderRadius: "5px",
        }}
      >
        {messages.map((msg) => (
          <p
            key={msg.id}
            style={{
              margin: "5px 0",
              padding: "8px",
              borderRadius: "5px",
              backgroundColor:
                msg.type === "received"
                  ? "#e1f5fe"
                  : msg.type === "sent"
                  ? "#c8e6c9"
                  : "#f0f0f0",
              alignSelf: msg.type === "received" ? "flex-start" : "flex-end",
              maxWidth: "80%",
              wordWrap: "break-word",
            }}
          >
            {msg.text}
          </p>
        ))}
      </div>

      <div
        style={{
          display: "flex",
          padding: "10px",
          backgroundColor: "#fff",
          borderTop: "1px solid #ddd",
        }}
      >
        <input
          type="text"
          placeholder="Type your message here..."
          style={{
            flex: 1,
            padding: "10px",
            border: "1px solid #ddd",
            borderRadius: "5px",
          }}
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyUp={handleKeyUp}
          disabled={connectionStatus !== "Connected"}
        />
        <button
          onClick={sendMessage}
          style={{
            padding: "10px 20px",
            marginLeft: "10px",
            border: "none",
            backgroundColor:
              connectionStatus === "Connected" ? "#0288d1" : "#9E9E9E",
            color: "white",
            borderRadius: "5px",
            cursor: connectionStatus === "Connected" ? "pointer" : "not-allowed",
          }}
          disabled={connectionStatus !== "Connected"}
        >
          Send
        </button>
      </div>
    </div>
  );
}

export default App;
