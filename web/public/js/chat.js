
document.addEventListener("DOMContentLoaded", () => {
	const token = localStorage.getItem("token");

	if (!token) {
		alert("You must be logged in first!");
		window.location.href = "/pages/login.html";
		return;
	}

	const socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
	const messagesDiv = document.getElementById("messages");
	const input = document.getElementById("messageInput");
	const sendBtn = document.getElementById("sendBtn");

	socket.onopen = () => {
		console.log("Connected to chat server");
	};

	socket.onmessage = (event) => {
		const msg = JSON.parse(event.data);
		const messageElement = document.createElement("div");
		messageElement.textContent = `${msg.user}: ${msg.message}`;
		messagesDiv.appendChild(messageElement);
	};

	socket.onclose = () => {
		console.log("Disconnected from chat server");
	};

	sendBtn.addEventListener("click", () => {
		const message = input.value.trim();
		if (message !== "") {
			socket.send(JSON.stringify({ message }));
			input.value = "";
		}
	});

	input.addEventListener("keypress", (e) => {
		if (e.key === "Enter") {
			sendBtn.click();
		}
	});
});

