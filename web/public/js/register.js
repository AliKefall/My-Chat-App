
document.addEventListener("DOMContentLoaded", () => {
	const form = document.getElementById("register-form");

	form.addEventListener("submit", async (e) => {
		e.preventDefault();

		const username = document.getElementById("username").value.trim();
		const password = document.getElementById("password").value.trim();

		try {
			const res = await fetch("/api/register", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ username, password }),
			});

			if (!res.ok) {
				const err = await res.text();
				alert("Hata: " + err);
				return;
			}

			const data = await res.json();
			// ✅ JWT token al ve sakla
			localStorage.setItem("token", data.token);

			// ✅ Chat sayfasına yönlendir
			window.location.href = "/chat.html";
		} catch (err) {
			console.error("Register error:", err);
			alert("Bir hata oluştu!");
		}
	});
});

