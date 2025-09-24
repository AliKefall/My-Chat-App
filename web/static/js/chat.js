
document.addEventListener("DOMContentLoaded", () => {
    const token = localStorage.getItem("token");
    const username = localStorage.getItem("username");

    if (!token || !username) {
        alert("Önce giriş yapmalısınız!");
        window.location.href = "/login.html";
        return;
    }

    // 🔗 WebSocket bağlantısı
    const socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

    socket.addEventListener("open", () => {
        console.log("✅ WebSocket bağlantısı kuruldu.");
    });

    socket.addEventListener("message", (event) => {
        const data = JSON.parse(event.data);

        // Gelen mesajları ekrana yaz
        const messages = document.getElementById("messages");
        const msgDiv = document.createElement("div");

        msgDiv.className = "p-2 my-1 rounded-lg " +
            (data.username === username
                ? "bg-blue-500 text-white self-end text-right"
                : "bg-gray-200 text-black self-start text-left");

        msgDiv.textContent = `${data.username}: ${data.content}`;
        messages.appendChild(msgDiv);

        // Otomatik scroll
        messages.scrollTop = messages.scrollHeight;
    });

    socket.addEventListener("close", () => {
        console.log("❌ WebSocket bağlantısı kapandı.");
    });

    // 📨 Mesaj gönderme
    const form = document.getElementById("chat-form");
    const input = document.getElementById("message");

    form.addEventListener("submit", (e) => {
        e.preventDefault();
        const message = input.value.trim();
        if (!message) return;

        const payload = {
            type: "message",
            content: message,
            username: username,
        };

        socket.send(JSON.stringify(payload));
        input.value = "";
    });
});

