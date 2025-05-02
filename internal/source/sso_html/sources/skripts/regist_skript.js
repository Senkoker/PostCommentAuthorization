const login = document.querySelector('.login');
const regist = document.querySelector('.password')
const act = document.querySelector('.act')
login.addEventListener("input",cheking_login)
regist.addEventListener("input", cheking_pass);
act.addEventListener('submit', function(event) {
    event.preventDefault();
    inputer();
});

function isValidPassword(password) {
    // Регулярное выражение для проверки наличия заглавной и строчной буквы, а также специальных символов
    const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[.!$]).+$/;
    return passwordRegex.test(password);
}

function cheking_pass() {
    let text = regist.value;
    let warning = document.querySelector('.warning_pass');
    if (text.length > 9 && isValidPassword(text)) {
        console.log('Все ок с паролем');
        warning.removeAttribute("id","warning_ok");
    } else {
        warning.setAttribute("id","warning_ok");
    }
}
function cheking_login(){
    let text = login.value;
    let warning = document.querySelector('.warning_log')
    if (text.includes('@')){
    console.log('Все ok с логином')
    warning.removeAttribute("id","warning_ok")
    }else{
        warning.setAttribute("id","warning_ok")
    }
}   
function inputer() {
    const login = document.querySelector('.login');
    const regist = document.querySelector('.password')
    const act = document.querySelector('.act')
    const registerData = {
        email: login.value,
        password: regist.value,
    };

    console.log(JSON.stringify(registerData));

    fetch('http://localhost:8082/auth/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(registerData)
    })
    .then(function(response) {
        if (!response.ok) {
            throw new Error(response.statusText);
        }
        return response.json(); // Преобразуем ответ в объект JavaScript
    })
    .then(function(result) {
         })
    .catch(function(error) {
        alert(`error: ${error}`)
    });
}
