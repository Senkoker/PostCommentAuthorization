function openPostMenu () {
    const createPostMenu = document.querySelector('.createPostMenu')
    createPostMenu.style.display = 'flex'
}
window.onclick = function(event) {
    const createPostMenu = document.querySelector('.createPostMenu')
    const closeMenu = document.querySelector('.closeMenu')
    if (event.target == closeMenu) {
        createPostMenu.style.display = "none";
    }
}
async function createPost(){

}


async function uploadPhoto() {
    const fileInput = document.getElementById('photoInput');
    const content = document.getElementById('content')
    const hashtags = document.getElementById('hashtags');
    const private = document.getElementById('private');

    const file = fileInput.files[0];
    const token = getToken()
    const contentValue = content.value;
    const hashtagsValue = hashtags.value
    const privateValue = private.value


    const formData = new FormData();
    formData.append('hashtags', hashtagsValue);
    formData.append('content', contentValue);
    formData.append('img', file);
    formData.append('private',privateValue)

    try {
        const response = await fetch('http://localhost:8082/feed/create_post', {
            method: 'POST',
            body: formData,
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        const result = await response.json();
        if (response.ok) {
            alert(`Upload successful! ${result.id}`);
            const createPostMenu = document.querySelector('.createPostMenu');
            createPostMenu.style.display = "none";
        } else {
            alert(`Error: ${response.status}`);
        }
    } catch (error) {
        alert(`Network error: ${error.message}`);
    }

}

function getToken() {
    return localStorage.getItem('authToken');
}