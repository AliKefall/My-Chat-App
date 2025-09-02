
document.addEventListener("DOMContentLoaded", () => {
	// LocalStorage’dan JWT’yi al
	const token = localStorage.getItem("token");
	if (!token) {
		alert("Lütfen önce giriş yapın!");
		window.location.href = "/login.html";
		return;
	}

	// ✅ WebSocket bağlantısını başlat
	const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

	const chatBox = document.getElementById("chat-box");
	const messageInput = document.getElementById("message-input");
	const sendButton = document.getElementById("send-btn");

	// Mesaj geldiğinde ekrana yaz
	ws.onmessage = (event) => {
		const msg = JSON.parse(event.data);
		const el = document.createElement("div");
		el.innerHTML = `<b>${msg.user}:</b> ${msg.message}`;
		chatBox.appendChild(el);
		chatBox.scrollTop = chatBox.scrollHeight; // Auto-scroll
	};

	// Bağlantı açıldığında
	ws.onopen = () => {
		console.log("✅ WebSocket connected!");
	};

	ws.onerror = (err) => {
		console.error("❌ WebSocket error:", err);
	};

	ws.onclose = () => {
		console.log("❌ WebSocket closed!");
	};

	// Mesaj gönder
	sendButton.addEventListener("click", () => {
		const msg = messageInput.value.trim();
		if (msg) {
			ws.send(JSON.stringify({ message: msg }));
			messageInput.value = "";
		}
	});

	// Enter tuşu ile gönder
	messageInput.addEventListener("keypress", (e) => {
		if (e.key === "Enter") {
			sendButton.click();
		}
	});
});


