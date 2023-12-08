// script.js

// Создадим объект для управления шахматной доской
const chessboard = document.querySelector('.chessboard');
const yourColor = 1;  //white
let IsChoice = "";
let TurnArray = []



const InitSetup = [
    "rnbqkbnr",
    "pppppppp",
    "........",
    "........",
    "........",
    "........",
    "PPPPPPPP",
    "RNBQKBNR"
];

// Добавим клетки на доску
function createCells() {
    for (let i = 0; i < 64; ++i) {
        const cell = document.createElement('div');
        cell.classList.add('cell');
        cell.draggable = false;
        cell.id = `cell-${i}`;
        chessboard.appendChild(cell);
    }
}

// Добавим фигуры на доску
function createPieces(initialSetup) {
    const cells = document.querySelectorAll('.cell');

    for (let i = 0; i < cells.length; i++) {
        const row = Math.floor(i / 8);
        const col = i % 8;
        const piece = initialSetup[row][col];
        if (piece !== '.') {
            const pieceElement = document.createElement('div');
            pieceElement.classList.add('figure');
            switch (piece) {
                case 'p':
                    pieceElement.style.backgroundImage = 'url(/images/chess/piece_black.png';
                    break;
                case 'P':
                    pieceElement.style.backgroundImage = 'url(/images/chess/piece_white.png';
                    break;
                case 'r':
                    pieceElement.style.backgroundImage = 'url(/images/chess/rook_black.png';
                    break;
                case 'R':
                    pieceElement.style.backgroundImage = 'url(/images/chess/rook_white.png';
                    break;
                case 'n':
                    pieceElement.style.backgroundImage = 'url(/images/chess/knight_black.png';
                    break;
                case 'N':
                    pieceElement.style.backgroundImage = 'url(/images/chess/knight_white.png';
                    break;
                case 'b':
                    pieceElement.style.backgroundImage = 'url(/images/chess/bishop_black.png';
                    break;
                case 'B':
                    pieceElement.style.backgroundImage = 'url(/images/chess/bishop_white.png';
                    break;
                case 'q':
                    pieceElement.style.backgroundImage = 'url(/images/chess/queen_black.png';
                    break;
                case 'Q':
                    pieceElement.style.backgroundImage = 'url(/images/chess/queen_white.png';
                    break;
                case 'k':
                    pieceElement.style.backgroundImage = 'url(/images/chess/king_black.png';
                    break;
                case 'K':
                    pieceElement.style.backgroundImage = 'url(/images/chess/king_white.png';
                    break;
            }
            pieceElement.dataset.piece = piece;
            cells[i].appendChild(pieceElement);
        }
    }
}

function isUpperCase(character) {
    return character === character.toUpperCase() && character !== character.toLowerCase();
}


function block_is_empty(block_id) {
    const block = document.getElementById(block_id);
    const innerDiv = block.querySelector("div");
    if (!innerDiv) {
        return 0;
    } else if (isUpperCase(innerDiv.dataset.piece)) {
        return 1;   // фигура белая
    } else {
        return -1;  // фигура черная
    }
}


function init_commands() {
    const blocks = document.querySelectorAll(".cell");
    blocks.forEach(function (block) {
        block.addEventListener("click", function () {
            if (!IsChoice) {
                figure_click(block);
                TurnArray = take_turn(block)
                for (let i = 0; i < TurnArray.length; ++i) {
                    const cell = document.getElementById(TurnArray[i]);
                    cell.style.backgroundColor = "rgba(50, 200, 50, 0.5)";
                }
            } else {
                console.log("OK")
                figure_click_2(block)
            }
        });
    });
}

// Добавить обработчик события "click" на родительский див
function figure_click(block) {
    // Проверить, есть ли текст внутри внутреннего дива
    if (block_is_empty(block.id) !== yourColor) {
        return;
    }
    block.style.backgroundColor = "rgba(255, 0, 0, 0.5)";
    IsChoice = block.id;
}

function clear_background() {
    for (let i = 0; i < 64; ++i) {
        const cell = document.getElementById(`cell-${i}`);
        cell.style.backgroundColor = null;
    }
}


function figure_click_2(block) {
    const outDiv = document.getElementById(IsChoice);
    const innerDiv = outDiv.querySelector("div");
    if (block === outDiv) {
        IsChoice = ""
        TurnArray = []
        clear_background()
        return
    }
    if (!TurnArray.includes(block.id)) {
        return
    }
    if (block_is_empty(block.id) === -yourColor) {
        const delDiv = block.querySelector("div");
        delDiv.remove();
    }
    block.appendChild(innerDiv);
    IsChoice = ""
    TurnArray = []
    clear_background()
}

function piece_turn(cell_num) {
    let available =[]
    const row = Math.floor(cell_num / 8)
    const col = cell_num % 8
    let block_id = ""
    if (row > 0) {
        block_id = `cell-${(row - 1) * 8 + col}`
        if (block_is_empty(block_id) === yourColor) {
            return available
        }
        available.push(block_id)
    }
    if (row === 6) {
        block_id = `cell-${(row-2) * 8 + col}`
        if (block_is_empty(block_id) === yourColor) {
            return available
        }
        available.push(block_id)
    }
    return available
}

function rook_turn(cell_num) {
    let available =[]
    const row = Math.floor(cell_num / 8)
    const col = cell_num % 8
    let block_id = ""
    for (let i = row - 1; i >= 0; i--) {
        block_id = `cell-${i * 8 + col}`
        if (block_is_empty(block_id) === yourColor) {
            break;
        } else if (block_is_empty(block_id) === -yourColor) {
            available.push(block_id)
            break;
        }
        available.push(block_id)
    }
    for (let i = row + 1; i < 8; i++) {
        block_id = `cell-${i * 8 + col}`
        if (block_is_empty(block_id) === yourColor) {
            break;
        } else if (block_is_empty(block_id) === -yourColor) {
            available.push(block_id)
            break;
        }
        available.push(block_id)
    }
    for (let i = col - 1; i >= 0; i--) {
        block_id = `cell-${row * 8 + i}`
        if (block_is_empty(block_id) === yourColor) {
            break;
        } else if (block_is_empty(block_id) === -yourColor) {
            available.push(block_id)
            break;
        }
        available.push(block_id)
    }
    for (let i = col + 1; i < 8; i++) {
        block_id= `cell-${row * 8 + i}`
        if (block_is_empty(block_id) === yourColor) {
            break;
        } else if (block_is_empty(block_id) === -yourColor) {
            available.push(block_id)
            break;
        }
        available.push(block_id)
    }
    return available
}

function knight_turn(cell_num) {
    let available =[]
    let block_id = ""
    const row = Math.floor(cell_num / 8)
    const col = cell_num % 8
    if (row - 2 >= 0) {
        if (col - 1 >= 0) {
            block_id = `cell-${(row - 2) * 8 + col - 1}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
        if (col + 1 < 8) {
            block_id = `cell-${(row - 2) * 8 + col + 1}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
    }
    if (row + 2 < 8) {
        if (col - 1 >= 0) {
            block_id = `cell-${(row + 2) * 8 + col - 1}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
        if (col + 1 < 8) {
            block_id = `cell-${(row + 2) * 8 + col + 1}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
    }
    if (col - 2 >= 0) {
        if (row - 1 >= 0) {
            block_id = `cell-${(row - 1) * 8 + col - 2}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
        if (row + 1 < 8) {
            block_id = `cell-${(row + 1) * 8 + col - 2}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
    }
    if (col + 2 < 8) {
        if (row - 1 >= 0) {
            block_id = `cell-${(row - 1) * 8 + col + 2}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
        if (row + 1 < 8) {
            block_id = `cell-${(row + 1) * 8 + col + 2}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
    }
    return available
}

function bishop_turn(cell_num) {
    let available =[]
    let block_id = ""
    const row = Math.floor(cell_num / 8)
    const col = cell_num % 8
    for (let i = 1; i < 8; i++) {
        if (row - i >= 0 && col - i >= 0) {
            block_id = `cell-${(row - i) * 8 + col - i}`
            if (block_is_empty(block_id) === yourColor) {
                break;
            } else if (block_is_empty(block_id) === -yourColor) {
                available.push(block_id)
                break;
            }
            available.push(block_id)
        }
    }
    for (let i = 1; i < 8; i++) {

        if (row - i >= 0 && col + i < 8) {
            block_id = `cell-${(row - i) * 8 + col + i}`
            if (block_is_empty(block_id) === yourColor) {
                break;
            } else if (block_is_empty(block_id) === -yourColor) {
                available.push(block_id)
                break;
            }
            available.push(block_id)
        }
    }
    for (let i = 1; i < 8; i++) {
        if (row + i < 8 && col - i >= 0) {
            block_id = `cell-${(row + i) * 8 + col - i}`
            if (block_is_empty(block_id) === yourColor) {
                break;
            } else if (block_is_empty(block_id) === -yourColor) {
                available.push(block_id)
                break;
            }
            available.push(block_id)
        }
    }
    for (let i = 1; i < 8; i++) {
        if (row + i < 8 && col + i < 8) {
            block_id = `cell-${(row + i) * 8 + col + i}`
            if (block_is_empty(block_id) === yourColor) {
                break;
            } else if (block_is_empty(block_id) === -yourColor) {
                available.push(block_id)
                break;
            }
            available.push(block_id)
        }
    }
    return available
}

function queen_turn(cell_num) {
    let available =[]
    available = rook_turn(cell_num)
    available = available.concat(bishop_turn(cell_num))
    return available
}

function king_turn(cell_num) {
    let available =[]
    let block_id = ""
    const row = Math.floor(cell_num / 8)
    const col = cell_num % 8
    if (row - 1 >= 0) {
        if (col - 1 >= 0) {
            block_id = `cell-${(row - 1) * 8 + col - 1}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
        if (col + 1 < 8) {
            block_id = `cell-${(row - 1) * 8 + col + 1}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
        block_id = `cell-${(row - 1) * 8 + col}`
        if (block_is_empty(block_id) !== yourColor) {
            available.push(block_id)
        }
    }
    if (row + 1 < 8) {
        if (col - 1 >= 0) {
            block_id = `cell-${(row + 1) * 8 + col - 1}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
        if (col + 1 < 8) {
            block_id = `cell-${(row + 1) * 8 + col + 1}`
            if (block_is_empty(block_id) !== yourColor) {
                available.push(block_id)
            }
        }
        block_id = `cell-${(row + 1) * 8 + col}`
        if (block_is_empty(block_id) !== yourColor) {
            available.push(block_id)
        }
    }
    if (col - 1 >= 0) {
        block_id = `cell-${row * 8 + col - 1}`
        if (block_is_empty(block_id) !== yourColor) {
            available.push(block_id)
        }
    }
    if (col + 1 < 8) {
        block_id = `cell-${row * 8 + col + 1}`
        if (block_is_empty(block_id) !== yourColor) {
            available.push(block_id)
        }
    }
    return available
}

function take_turn(block) {
    const num = Number(block.id.substring(5))
    const innerDiv = block.querySelector("div");
    switch (innerDiv.getAttribute("data-piece").toLowerCase()) {
        case 'p':
            return piece_turn(num)
        case 'r':
            return rook_turn(num)
        case 'n':
            return knight_turn(num)
        case 'b':
            return bishop_turn(num)
        case 'q':
            return queen_turn(num)
        case 'k':
            return king_turn(num)
    }
}

// Инициализация
createCells();
createPieces(InitSetup);
init_commands();

