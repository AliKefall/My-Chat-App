document.getElementById("loginForm").addEventListener("submit", async (e) => {
	e.preventDefault()

	const username = document.getElementById("username").value.trim();
	const password = document.getElementById("password").value.trim();

	const res = await fetch("/api/login", {
		method: "POST",
		header: { "Content-type": "application/json" },
		body: JSON.stringify({ username, password })
	});

	if (res.ok) {
		const data = await res.json();
		localStorage.setItem("token", data.token);
		window.location.href = "/pages/chat.html"
	} else {
		alert("Login failed")
	}

})
