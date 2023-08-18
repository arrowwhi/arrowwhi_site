const input = document.getElementById("chat_input");
const output = document.getElementById("message_list");
let ws;

let single_message = {
    "m_type": "message",
    "message": "",
    "recipient": "",
}

// Функция для создания структуры тегов для сообщения
function createMessage(author, messageText, dateTime) {
    const newRow = document.createElement("tr");
    // Создание ячеек <th> и <td> для новой строки
    const thCell = document.createElement("th");
    thCell.textContent = author; // Здесь можно указать имя
    const tdCell1 = document.createElement("td");
    tdCell1.textContent = messageText; // Здесь можно указать сообщение
    const tdCell2 = document.createElement("td");
    tdCell2.textContent = dateTimeValidation(dateTime);

    newRow.appendChild(thCell);
    newRow.appendChild(tdCell1);
    newRow.appendChild(tdCell2);

    return newRow;
}

// функция для печати сообщений снизу
const print = function (author, message, time) {
    const d = createMessage(author, message, time)
    output.appendChild(d);
    output.scroll(0, output.scrollHeight);
};

// функция для печати сообщений сверху
const print_forward = function (author, message, time) {
    const d = createMessage(author, message, time);
    output.prepend(d);
    output.scroll(0, output.scrollHeight);
};

//проверка евента на нажатую клаишу enter           - ?
function checkEnter(evt) {
    if (event.key === "Enter") {
        event.preventDefault();
        press_send(evt)
    }
}

// отправка сообщения на сервер
function press_send() {
    if (!ws) {
        return;
    }
    single_message.message = input.value;
    if (single_message.message.trim() === "") {
        return;
    }
    console.log(`user recipient: ${single_message.recipient}`)

    print("Вы:", input.value, new Date());
    console.log("SEND: " + JSON.stringify(single_message));
    ws.send(JSON.stringify(single_message));
    input.value = ''
}

let lastMessage = 0

// функция для получения сообщений
async function take_messages(usr, last) {
    const url = "/database/get_messages"
    const data = {
        username: usr,
        lastId: last,
        count: 10
    };
    await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json' // Заголовок для указания типа данных JSON
        },
        body: JSON.stringify(data) // Преобразование данных в JSON и отправка в теле запроса
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            // Обработка полученных данных
            for (let elem of data.messages) {
                const dateTime = new Date(elem.create_date)
                lastMessage = elem.id
                print_forward(elem.user_from, elem.message, dateTime)
            }
            if (data.messages.length < 10) {
                const q = document.getElementById('early_href_div')
                q.classList.add('invisible')
            }
        })
        .catch(error => console.error('Ошибка:', error));
}

// функция для обработки выбора чата
function handleLoginClick(event, clickedElement) {
    event.preventDefault();
    const q = document.getElementById('early_href_div')
    if (q) {
        q.classList.remove('invisible')
    }
    const allLinks = document.querySelectorAll('.login_click');
    for (let i = 0; i < allLinks.length; i++) {
        allLinks[i].classList.remove("active")
    }
    clickedElement.classList.add("active")
    while (output.firstChild) {
        output.removeChild(output.firstChild);
    }

    let all_messages = document.getElementById('message_zero')
    if (all_messages) {
        all_messages.classList.remove("d-none")
    }
    hideLogins()

    const strongElement = clickedElement.querySelector(".mb-1"); // Выбираем внутренний элемент с классом "mb-1"
    const user = strongElement.textContent; // Получаем текстовое содержимое элемента

    single_message.recipient = user;
    take_messages(user, 0)
}

// функция для выделения поиска логина
function get_logins_messages(logins_on) {
    const mgs_div = document.getElementById('messagesList')
    const logins_div = document.getElementById('loginsList')
    if (logins_on) {
        mgs_div.classList.add('invisible')
        logins_div.classList.remove('invisible')
    } else {
        setTimeout(function() {
        mgs_div.classList.remove('invisible')
        logins_div.classList.add('invisible')
        }, 100)
    }
}

// функция для обработки выбора логина из поиска
function HandleLoginSearchClick(login) {
    let current_login = getFromCurrentLogin(login)
    if (current_login) {
        handleLoginClick(event, current_login)
    } else {
        const messageData = {
            user: login,
            create_date: '',
            message: ''
        };
        const container = document.getElementById('messagesList');
        const message = createLoginMessage(messageData);
        container.prepend(message);
    }
}

// функция для добавления логина в список логинов при поиске
function add_login(login) {
    // Создаем новый элемент "a"
    const newLink = document.createElement('a');
    // newLink.href = '/';
    newLink.className = 'list-group-item list-group-item-action py-3 lh-sm';
    newLink.addEventListener("click", function(event) {
        event.preventDefault();
        HandleLoginSearchClick(login)})
    // Создаем новый элемент "strong"
    const strongElement = document.createElement('strong');
    strongElement.textContent = login; // Задаем содержимое для "strong"
    newLink.appendChild(strongElement);
    const loginsList = document.getElementById('loginsList');
    // Добавляем новый "a" внутрь блока с id "loginsList"
    loginsList.appendChild(newLink);
}


// функция для создания структуры тегов для логина с сообщениями
function createLoginMessage(data) {
    const aTag = document.createElement('a');
    aTag.href = '#';
    aTag.classList.add('list-group-item', 'list-group-item-action', 'py-3', 'lh-sm', 'login_click');
    aTag.addEventListener("click", function(event) {
        event.preventDefault();
        handleLoginClick(event, this);
    });

    const div1 = document.createElement('div');
    div1.classList.add('d-flex', 'w-100', 'align-items-center', 'justify-content-between');

    const strong = document.createElement('strong');
    strong.classList.add('mb-1');
    strong.textContent = data.user;

    const small = document.createElement('small');
    small.classList.add('text-body-secondary');
    small.textContent = data.create_date.toString();

    div1.appendChild(strong);
    div1.appendChild(small);

    const div2 = document.createElement('div');
    div2.classList.add('col-10', 'mb-1', 'small');
    div2.textContent = data.message;

    aTag.appendChild(div1);
    aTag.appendChild(div2);

    return aTag;
}


// функция для получения логинов из базы данных и формирования списка
function take_logins() {
    const block = document.getElementById('loginsList')
    while (block.firstChild) {
        block.removeChild(block.firstChild);
    }
    fetch("/database/get_logins", {
        method: 'POST',
    }).then(response => response.json()).then(data => {
        for (let elem of data.logins) {
            add_login(elem)
        }
    }).catch(error => console.error('Ошибка:', error));
}

// функция для поиска логинов
function getFromCurrentLogin(login) {
    const links = document.querySelectorAll(".login_click");
    for (let i = 0; i < links.length; ++i) {
        const mb1Element = links[i].querySelector('.mb-1');
        if (mb1Element.textContent.trim() === login) {
            return links[i]
        }
    }
    return null
}

// функция для скрытия и показа блоков с логинами
function hideLogins() {
    const logins = document.getElementById('lgns')
    if (logins.classList.contains('d-none')) {
        logins.classList.remove('d-none')
        const all_messages = document.getElementById('message_zero')
        all_messages.classList.add('d-none')
    } else {
        logins.classList.add('d-none')

    }
}

function dateTimeValidation(dateTime) {
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    const month = String(dateTime.getMonth() + 1);
    const day = String(dateTime.getDate());
    const hours = String(dateTime.getHours());
    const minutes = String(dateTime.getMinutes());
    const currentDate = new Date();
    let formatDate = ""
    if (
        currentDate.getFullYear() === dateTime.getFullYear() &&
        currentDate.getMonth() === dateTime.getMonth() &&
        currentDate.getDate() === dateTime.getDate()
    ) {
        formatDate = `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}`;
    } else {
        formatDate = `${day.padStart(2, '0')} ${months[month-1]}`;
    }
    console.log(formatDate)
    return formatDate
}

const links = document.querySelectorAll(".login_click");
links.forEach((link) => {
    link.addEventListener("click", function(event) {
        event.preventDefault();
        handleLoginClick(event, this);
    });
});

const myInput = document.getElementById("loginsSearch");
myInput.addEventListener("focus", function() {
    console.log("Поле ввода получило фокус!");
    get_logins_messages(true)
    take_logins()
});
myInput.addEventListener("blur", function() {
    console.log("Поле ввода потеряло фокус!");
    get_logins_messages(false)
});

const toggleButton = document.getElementById('login_button');
toggleButton.addEventListener('click', function () {
    hideLogins()
});



const earlySearch = document.getElementById("early_href");
earlySearch.addEventListener("click", function(event) {
    event.preventDefault();
    take_messages(single_message.recipient, lastMessage);
})


