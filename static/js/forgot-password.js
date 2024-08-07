document.getElementById('forgotPasswordForm').addEventListener('submit', function (event) {
    event.preventDefault();
    const email = document.getElementById('email').value;

    // fetch API to send email to user
    fetch('/forgot-password', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
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
        alert('send email failed. Please try again later.')
    })
})