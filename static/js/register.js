document.getElementById('registerForm').addEventListener('submit', function (event) {
    event.preventDefault()
    const username = document.getElementById('username').value
    const password = document.getElementById('password').value
    const email = document.getElementById('email').value

    // fetch API to register user
    fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            username: username,
            password: password,
            email: email
        })
    }).then(response => {
        if (response.redirected) {
            window.location.href = response.url
        } else {
            response.json().then(data => {
                alert(data.message)
            })
        }
    }).catch(_ => {
        alert('Registration failed. Please try again later.')
    })
})
