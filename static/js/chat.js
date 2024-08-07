document.addEventListener('DOMContentLoaded', function () {
    const messageInput = document.getElementById('messageInput');
    const sendButton = document.getElementById('sendButton');
    const chatMessages = document.getElementById('chatMessages');
    const username = window.username; // 从页面脚本中获取用户名

    // 事件监听器
    sendButton.addEventListener('click', sendMessage);
    messageInput.addEventListener('keypress', function (event) {
        if (event.key === 'Enter') {
            sendMessage();
        }
    });

    const socket = new WebSocket(`ws://${window.location.host}/ws/chat/${username}`);
    socket.addEventListener('open', handleSocketOpen);
    socket.addEventListener('message', handleSocketMessage);
    socket.addEventListener('error', handleSocketError);

    // WebSocket 消息处理
    function handleSocketOpen() {
        console.log('WebSocket connection established');
    }

    function handleSocketMessage(event) {
        try {
            const message = JSON.parse(event.data);
            appendMessage(message.username, message.message);
        } catch (e) {
            console.error('Failed to parse message:', e);
        }
    }

    function handleSocketError(error) {
        console.error('WebSocket error:', error);
    }

    // 消息发送处理
    function sendMessage() {
        const message = messageInput.value.trim();
        if (message === "") return;

        // 将消息通过ws发送到服务器
        socket.send(JSON.stringify({
            username: username,
            message: message
        }));

        // 清空输入框
        messageInput.value = "";
    }

    // 消息显示处理
    function appendMessage(sender, message) {
        const messageElement = document.createElement('div');
        messageElement.classList.add('chat-message');
        messageElement.innerHTML = `<p class="sender">${sender}</p><p>${message}</p>`;
        chatMessages.appendChild(messageElement);
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }
});
