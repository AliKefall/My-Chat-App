document.getElementById("registerForm").addEventListener("submit", async (e) => {
	e.preventDefault();
	const username = document.getElementById("username").value.trim();
	const password = document.getElementById("password").value.trim();

	const res = await fetch("/api/register", {
		method: "POST",
		headers: { "Content-type": "application/json" },
		body: JSON.stringify({ username, password })
	});

	if (res.ok) {
		const data = await res.json();
		localStorage.setItem("token", data.token);
		window.location.href = "/pages/chat.html"

	} else {
		alert("Register failed")
	}

})
