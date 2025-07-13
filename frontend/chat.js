const UI = (() => {
	const loginSection = document.getElementById('login');
	const chatSection = document.getElementById('chat');
	const messagesDiv = document.getElementById('messages');
	const usernameInput = document.getElementById('usernameInput');
	const messageInput = document.getElementById('messageInput');
	const loginBtn = document.getElementById('loginBtn');
	const sendBtn = document.getElementById('sendBtn');

	const bindLogin = handler => loginBtn.addEventListener('click', () => handler(usernameInput.value));
	const bindSend = handler => sendBtn.addEventListener('click', () => handler(messageInput.value));

	const showChat = () => {
		loginSection.classList.add('hidden');
		chatSection.classList.remove('hidden');
	};

	const addMessage = ({ username, message }) => {
		const msgEl = document.createElement('div');
		msgEl.classList.add('message');
		msgEl.innerHTML = `<span class="username">${username}:</span><span>${message}</span>`;
		messagesDiv.appendChild(msgEl);
		messagesDiv.scrollTop = messagesDiv.scrollHeight;
	};

	const clearComposer = () => { messageInput.value = ''; };

	return { bindLogin, bindSend, showChat, addMessage, clearComposer };
})();

const SocketClient = (() => {
	let socket;
	const connect = username => {
		socket = new WebSocket(`ws://${window.location.host}/ws`);
		socket.onopen = () => UI.showChat();
		socket.onmessage = event => UI.addMessage(JSON.parse(event.data));
		socket.onclose = () => console.log('Bağlantı kesildi');
	};

	const send = (username, message) => {
		if (!socket) return;
		socket.send(JSON.stringify({ username, message }));
		UI.clearComposer();
	};

	return { connect, send };
})();

// Uygulama başlatma
(function App() {
	let currentUser = '';
	UI.bindLogin(username => {
		if (!username.trim()) return alert('Kullanıcı adı girin');
		currentUser = username;
		SocketClient.connect(username);
	});

	UI.bindSend(msg => {
		if (!msg.trim()) return;
		SocketClient.send(currentUser, msg);
	});
})();
