// Логика формы
const login = document.querySelector('.login');
const regiserBtn = document.querySelector('.btn_regist');
const btnLogin = document.querySelector('.btn_go')
regiserBtn.addEventListener('click', function(event){
    event.preventDefault()
    console.log("hello")
    window.location.href = 'file:///C:/Golang_social_project/VK_posts/internal/source/sso_html/register.html'
})
login.addEventListener("input",cheking_warn)
function cheking_warn(){
    let text = login.value;
    let warning = document.querySelector('.warning')
    warning.textContent = 'Ваш логин не соответсвует стандарту'
    if (isValidEmail(text)){
    console.log('Все ок')
    warning.removeAttribute("id","warning_ok")
    }else{
        warning.setAttribute("id","warning_ok")
    }
}

function isValidEmail(email) {
    // Регулярное выражение для проверки email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

btnLogin.addEventListener('click',function(event){
    event.preventDefault();
    const login = document.querySelector('.login')
    const password = document.querySelector('.password')
    fetch('http://localhost:8082/auth/login',{
        method: 'post',
        headers: {
            'Content-type': 'application/json'},
        body: JSON.stringify({
            email: login.value,
            password: password.value,
            appid: "1",
        })}).then(function(response){
            console.log(response)
            if (!response.ok){
                let warning = document.querySelector('.warning')
                warning.textContent = response.message
                warning.setAttribute("id","warning_ok")
                throw new Error(response.statusText)
            }
            return response.json()
        }).then(function(result){
            console.log(result.access_token)
            const token = result.access_token
            if (token !== "" ){
                localStorage.setItem('authToken', token);
                window.location.href = "file:///d%3A/sso_html/filluser.html"
            }
        }).catch(function(error) {
            alert(error)
        });
        
          
})
