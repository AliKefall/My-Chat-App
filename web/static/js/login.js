
document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("login-form");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const email = document.getElementById("email").value.trim();
    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value.trim();

    try {
      const res = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, username, password }), // ✅ backend'in istediği format
      });

      if (!res.ok) {
        const err = await res.text();
        alert("❌ Hata: " + err);
        return;
      }

      const data = await res.json();
      // ✅ JWT token sakla
      localStorage.setItem("token", data.token);

      // ✅ Chat sayfasına yönlendir
      window.location.href = "/chat.html";
    } catch (err) {
      console.error("Login error:", err);
      alert("Bir hata oluştu!");
    }
  });
});

