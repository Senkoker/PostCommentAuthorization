document.addEventListener('DOMContentLoaded', function(){
    let searchBtn = document.querySelector('.searchButton');
    var template = document.querySelector('.post')
    template.style.display = 'none'
    if (searchBtn) {
        searchBtn.addEventListener("click", function(event) {
            event.preventDefault();
            NewPost();
        });
    } else {
        console.error('Элемент с классом "searchButton" не найден.');
    }
});

function splitStringByPattern(inputString, pattern) {
    const regex = new RegExp(pattern);
    const parts = inputString.split(regex);
    return parts.slice(0, 2);
}

function getToken() {
    return localStorage.getItem('authToken');
}

function NewComments(){
    fetch('http://localhost:8082/',{
        method: 'post',
        headers: {
            // Заполните заголовки, если необходимо
        },
        body: JSON.stringify({
            // Заполните тело запроса, если необходимо
        })
    }).then(function(response){
        if (!response.ok){
            throw new Error('Network response was not ok');
        }
        return response.json();
    }).then(function(result){
        // Обработка результата
    }).catch(function(error){
        console.error('Ошибка при выполнении запроса:', error);
    });
}

function NewPost(){
    const inputIDs = document.querySelector('.searchItems');
    const postHeader = document.querySelector('.postHeader');
    const userName = document.querySelector('.userName');
    const token = getToken();

    fetch('http://localhost:8082/feed/get_posts?limit=5&offset=0&redis=false',{
        method: 'post',
        headers: {
            'Content-type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
            "posts_ids": [inputIDs.value],
            "hashtags": ""
        })
    }).then(function(response){
        if (!response.ok){
            throw new Error('Network response was not ok');
        }
        return response.json();
    }).then(function(result){
        console.log(result);
        result.posts.forEach((post) => {
            var template = document.querySelector('.post');
            var newPost = template.cloneNode(true);
            newPost.style.display = 'block';

            // Устанавливаем src для изображения
            newPost.querySelector('.imgUser').src = post.img_person_url;
            newPost.querySelector('.userName').textContent = post.author;

            const content = splitStringByPattern(post.content, '&');
            newPost.querySelector('.content').textContent = content[0];

            // Устанавливаем src для изображения поста
            newPost.querySelector('.postImage').src = content[1];

            newPost.querySelector('.numberLikes').textContent = post.likes;
            newPost.querySelector('.numberWathing').textContent = post.watched;
            newPost.querySelector('.datePublication').textContent = String(post.created_at).slice(0, 10);
            newPost.querySelector('.commentID').textContent = post.post_id;
            newPost.querySelector('.commentID').style.display = 'none';

            document.querySelector('.center').appendChild(newPost);
        });
    }).catch(function(error){
        console.error('Ошибка при выполнении запроса:', error);
    });
}