let socket;

function addSocketEvents(socket) {
  const messagesList = document.getElementById("message-list");

  socket.onmessage = (event) => {
    try {
      const wrapper = JSON.parse(event.data);
      const sender = wrapper.sender;

      const decodedBodyAsString = atob(wrapper.body);

      console.log("decoded:", decodedBodyAsString);
      const decodedBodyAsJSON = JSON.parse(decodedBodyAsString);
      const message = decodedBodyAsJSON.message;

      // console.log("Sender: ", data.sender);
      // console.log("Body (raw): ", data.body);

      // console.log("Decoded Body: ", text);
      const li = document.createElement("li");
      const div = document.createElement("div");
      li.appendChild(div);
      li.textContent = `From ${sender}: ${message}`;
      messagesList.appendChild(li);
    } catch (err) {
      console.error("Error parsing message:", err);
    }
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
  const connectFormDiv = document.getElementById("connect-form-div");

  socket = new WebSocket("ws://152.26.89.214:8080/connect");

  socket.onopen = () => {
    console.log("Successfully connected");
    socket.send(JSON.stringify({ name: name }));
  };
  console.log("Attempting websocket connection");
  addSocketEvents(socket);

  connectFormDiv.remove();
  messagingDiv.style.display = "flex";
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
