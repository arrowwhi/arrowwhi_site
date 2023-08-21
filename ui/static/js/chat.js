const input = document.getElementById("chat_input");
const output = document.getElementById("message_list");
let ws;
let pendingResponses = {};


let single_message = {
    "id": 0,
    "m_type": "message",
    "message": "",
    "user_to": "",
    "user_from": "You",
    "is_read": false,
    "create_date": ""
}

// Функция для создания сообщения
function createMessage(author, messageText, dateTime, read) {
    const newRow = document.createElement("tr");
    if (!read) {
        newRow.classList.add("unread")
    }
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
const print = function (elem) {
    const d = createMessage(elem.user_from, elem.message, elem.create_date, elem.is_read)
    output.appendChild(d);
    if (elem.id !== 0) {
        d.setAttribute('data-message-id', elem.id);
        ws.send(JSON.stringify({
            m_type: "read_message",
            ids: [elem.id]
        }));
    }
    output.scroll(0, output.scrollHeight);
    const msgs = document.getElementById('message_list')
    msgs.scrollTop = msgs.scrollHeight;
    return d
};

// функция для печати сообщений сверху
const print_forward = function (elem) {
    const dateTime = new Date(elem.create_date)
    const d = createMessage(elem.user_from, elem.message, dateTime, elem.is_read);
    d.setAttribute('data-message-id', elem.id);
    if (elem.user_from === single_message.user_to && !elem.is_read) {
        ws.send(JSON.stringify({
            m_type: "read_message",
            ids: [elem.id]
        }));
    }
    output.prepend(d);
    output.scroll(0, output.scrollHeight);
};

//проверка эвента на нажатую клаишу enter           - ?
function checkEnter(evt) {
    if (event.key === "Enter") {
        event.preventDefault();
        press_send(evt)
    }
}

// отправка сообщения на сервер
function press_send() {
    if (!ws) {
        console.log("No connection");
        return;
    }
    single_message.message = input.value;
    if (single_message.message.trim() === "" || single_message.user_to.trim() === "") {
        console.log("Empty message");
        return;
    }
    console.log(`user user_to: ${single_message.user_to}`)


    single_message.create_date = new Date().toISOString();
    // TODO поправить формированиие локал ид
    single_message.local_id = Math.random()
    const d = print(single_message);
    console.log("SEND: " + JSON.stringify(single_message));
    ws.send(JSON.stringify(single_message));
    pendingResponses[single_message.local_id] = d;
    GetNewMessageOnLoginList(single_message, true)
    input.value = ''
    const msgs = document.getElementById('message_list')
    msgs.scrollTop = msgs.scrollHeight;
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
            const loginPlace = document.getElementById('MessageListProfileName')
            loginPlace.textContent = data.name
            const photoPlace = document.getElementById('MessageListProfilePhoto')
            photoPlace.src = data.photo
            let read_now = []
            for (let elem of data.messages) {

                lastMessage = elem.id
                if (!elem.is_read && elem.user_from === usr) {
                    read_now.push(elem.id)
                }
                print_forward(elem)
            }
            if (read_now.length > 0) {
                ws.send(JSON.stringify({
                    m_type: "read_message",
                    ids: read_now
                }));
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

    // показываем "загрузить ранние"
    const q = document.getElementById('early_href_div')
    if (q) {
        q.classList.remove('invisible')
    }

    // снимаем активность с других кнопок и даем нужной
    const allLinks = document.querySelectorAll('.login_click');
    for (let i = 0; i < allLinks.length; i++) {
        allLinks[i].classList.remove("active")
    }
    clickedElement.classList.add("active")

    // чистим старые сообщения
    while (output.firstChild) {
        output.removeChild(output.firstChild);
    }

    // показываем блок с сообщениями
    let all_messages = document.getElementById('message_zero')
    if (all_messages) {
        all_messages.classList.remove("d-none")
    }

    // скрываем блок с логинами на мобилке
    hideLogins()

    const strongElement = clickedElement.querySelector(".mb-1"); // Выбираем внутренний элемент с классом "mb-1"
    const user = strongElement.textContent; // Получаем текстовое содержимое элемента

    single_message.user_to = user;
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
        setTimeout(function () {
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
        handleLoginClick(event, message)
    }
}

// функция для добавления логина в список логинов при поиске
function add_login(login) {
    // Создаем новый элемент "a"
    const newLink = document.createElement('a');
    // newLink.href = '/';
    newLink.className = 'list-group-item list-group-item-action py-3 lh-sm';
    newLink.addEventListener("click", function (event) {
        event.preventDefault();
        HandleLoginSearchClick(login)
    })
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
    aTag.addEventListener("click", function (event) {
        event.preventDefault();
        handleLoginClick(event, this);
    });

    const div1 = document.createElement('div');
    div1.classList.add('d-flex', 'w-100', 'align-items-center', 'justify-content-between');

    const unread = document.createElement('span');
    unread.classList.add('top-100', 'start-50', 'badge', 'translate-middle', 'bg-danger', 'rounded-pill', 'd-none');
    unread.textContent = '0';

    const strong = document.createElement('strong');
    strong.classList.add('mb-1');
    strong.textContent = data.user;

    const small = document.createElement('small');
    small.classList.add('text-body-secondary');
    small.textContent = data.create_date.toString();

    div1.appendChild(strong);
    div1.appendChild(unread);
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

// функция для управления количеством непрочитанных
function changeUnreadCount(login, count) {
    const loginElement = getFromCurrentLogin(login);
    const smallElement = loginElement.querySelector('span');
    let c = parseInt(smallElement.textContent)
    smallElement.textContent = String(c + count)
    if (c + count === 0) {
        smallElement.classList.add('d-none')
    } else {
        smallElement.classList.remove('d-none')
    }

}

function dateTimeValidation(dateTime) {
    dateTime = new Date(dateTime)
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
        formatDate = `${day.padStart(2, '0')} ${months[month - 1]}`;
    }
    return formatDate
}


function GetNewMessageOnLoginList(struct, to=false) {
    let elem;
    if (!to) {
        elem = getFromCurrentLogin(struct.user_from)
    } else {
        elem = getFromCurrentLogin(struct.user_to)
    }
    const unread = parseInt(elem.querySelector("span").innerText);
    const is_active = elem.classList.contains("active")
    if (elem) {
        elem.remove()
    }
    const messageData = {
        user: struct.user_from,
        create_date: dateTimeValidation(new Date()),
        message: struct.message
    };
    const newElem = createLoginMessage(messageData)
    if (is_active) {
        newElem.classList.add("active")
    }
    const q = document.getElementById("messagesList");
    q.prepend(newElem)
    if (!to) {
        changeUnreadCount(struct.user_from, unread + 1)
    }
}

function MarkAsRead(id) {
    const elem = document.querySelector(`[data-message-id="${id}"]`)
    if (elem) {
        elem.classList.remove("unread")
    }
}

const links = document.querySelectorAll(".login_click");
links.forEach((link) => {
    link.addEventListener("click", function (event) {
        event.preventDefault();
        handleLoginClick(event, this);
    });
});

const myInput = document.getElementById("loginsSearch");
myInput.addEventListener("focus", function () {
    console.log("Поле ввода получило фокус!");
    get_logins_messages(true)
    take_logins()
});

myInput.addEventListener("blur", function () {
    console.log("Поле ввода потеряло фокус!");
    get_logins_messages(false)
});

const toggleButton = document.getElementById('login_button');
toggleButton.addEventListener('click', function () {
    hideLogins()
});

const earlySearch = document.getElementById("early_href");
earlySearch.addEventListener("click", function (event) {
    event.preventDefault();
    take_messages(single_message.user_to, lastMessage);
})
