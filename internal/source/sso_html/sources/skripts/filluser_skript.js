document.addEventListener('DOMContentLoaded', function() {
    const secondName = document.querySelector('.second_name');
    const firstName = document.querySelector('.first_name');
    const imgFile = document.querySelector('.img_file');
    const birth = document.querySelector('.birth_date');
    const country = document.querySelector('.country');
    const city = document.querySelector('.city');
    const education = document.querySelector('.education');
    const btn = document.querySelector('.go_btn');

    function getToken() {
        return localStorage.getItem('authToken');
    }

    function unputer() {
        const token = getToken();
        console.log(token);
        const formData = new FormData();
        const date = new Date(birth.value);
        const formattedDate = date.toISOString().split('T')[0];
        formData.append("first_name", firstName.value);
        formData.append("second_name", secondName.value);
        formData.append("img", imgFile.files[0]);
        formData.append("birth_date", formattedDate);
        formData.append("country", country.value);
        formData.append("city", city.value);
        formData.append("education", education.value);

        fetch('http://localhost:8082/profile/fill_data', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            },
            body: formData
        }).then(function(response) {
            if (!response.ok) {
                throw new Error(response.status);
            }
            return response.json();
        }).then(function(result) {
            alert("Successful", result.id);
        }).catch(function(error) {
            alert(error);
        });
    }

    if (btn) {
        btn.addEventListener('click', function(event) {
            event.preventDefault();
            unputer();
        });
    } else {
        console.error("Element with class 'go_btn' not found.");
    }
});
