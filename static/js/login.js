document.getElementById('loginForm').addEventListener('submit', function (event) {
    event.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    // 将用户名字传递到后端,说明该用户要登录
    fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            username: username,
            password: password
        })
    }).then(response => {
        if (response.redirected) {
            window.location.href = response.url;
        } else {
            response.json().then(data => {
                alert(data.message);
            });
        }
    })
});
