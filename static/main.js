let socket;

function addSocketEvents(socket) {
  const messagesList = document.getElementById("messages");

  socket.onmessage = (event) => {
    const li = document.createElement("li");
    if (event.data === "") {
      console.log("Invalid message");
    }
    li.textContent = event.data;
    messagesList.appendChild(li);
  };

  socket.onclose = (event) => {
    console.log("Closing Connection", event);
  };

  socket.onerror = (error) => {
    console.log("Socket Error:", error);
  };
}

function handleConnect(event) {
  const messagingDiv = document.getElementById("message-div");
  event.preventDefault();

  const formData = new FormData(event.target);
  const name = formData.get("name");
  const connectFormDiv = document.getElementById("connect-form");

  socket = new WebSocket("ws://localhost:8080/connect");

  socket.onopen = () => {
    console.log("Successfully connected");
    socket.send(JSON.stringify({ name: name }));
  };
  console.log("Attempting websocket connection");
  addSocketEvents(socket);

  connectFormDiv.remove();
  messagingDiv.style.display = "block";
}

function sendMessage(event) {
  event.preventDefault();
  const formData = new FormData(event.target);
  const message = formData.get("message");
  const input = document.getElementById("messagefield");

  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({ message: message }));
    console.log("Message sent");
  } else {
    console.log("Socket is not open");
  }
  input.value = "";
}
