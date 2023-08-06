
    // Функция создания всплывающего окна
    function createPopup() {
        const popup = document.createElement('div');
        popup.id = 'popup';
        popup.style.display = 'none';
        popup.style.position = 'fixed'; // Изменили значение на 'fixed'
        popup.style.width = '300px';
        popup.style.backgroundColor = '#fff';
        popup.style.border = '1px solid #ccc';
        popup.style.padding = '20px';
        popup.style.borderRadius = '5px';
        popup.style.boxShadow = '0 0 10px rgba(0, 0, 0, 0.3)';
        popup.style.top = '50%'; // Установили значение 'top' на '50%'
        popup.style.left = '50%'; // Установили значение 'left' на '50%'
        popup.style.transform = 'translate(-50%, -50%)'; // Сдвинули попап на половину его ширины и высоты

    const heading = document.createElement('h2');
    heading.textContent = 'Выберите тип:';
    popup.appendChild(heading);

    const option1 = document.createElement('label');
    option1.innerHTML = '<input type="radio" name="type" value="Доработка"> Доработка';
    popup.appendChild(option1);

    const option2 = document.createElement('label');
    option2.innerHTML = '<input type="radio" name="type" value="Баг"> Баг';
    popup.appendChild(option2);

    const description = document.createElement('textarea');
    description.id = 'description';
    description.placeholder = 'Введите ваш комментарий';
    description.style.width = '100%';
    description.style.height = '100px';
    description.style.resize = 'none';
    description.style.marginBottom = '10px';
    popup.appendChild(description);

    const sendButton = document.createElement('button');
    sendButton.id = 'sendButton';
    sendButton.textContent = 'Отправить';
    popup.appendChild(sendButton);

    const successMessage = document.createElement('p');
    successMessage.id = 'successMessage';
    successMessage.textContent = 'Успешно отправлено';
    successMessage.style.display = 'none';
    successMessage.style.color = 'green';
    popup.appendChild(successMessage);

    const closeButton = document.createElement('button');
    closeButton.id = 'closeButton';
    closeButton.textContent = 'Закрыть';
    popup.appendChild(closeButton);

    document.body.appendChild(popup);

    // Назначаем обработчики событий
    const openButton = document.getElementById('feedback');
    const closeButtonElem = document.getElementById('closeButton');
    const sendButtonElem = document.getElementById('sendButton');
    const successMessageElem = document.getElementById('successMessage');

    openButton.addEventListener('click', () => {
        console.log("OK")
        event.preventDefault()
    popup.style.display = 'block';
});

    closeButtonElem.addEventListener('click', () => {
    popup.style.display = 'none';
    successMessageElem.style.display = 'none';
});

    sendButtonElem.addEventListener('click', () => {
    const selectedType = document.querySelector('input[name="type"]:checked');
    const descriptionValue = description.value.trim();

    if (selectedType && descriptionValue !== '') {
    // Здесь можно выполнить POST-запрос с помощью XMLHttpRequest или fetch
    // После успешной отправки можно показать сообщение об успехе и очистить поля
    successMessageElem.style.display = 'block';
    description.value = '';
}
});
}
