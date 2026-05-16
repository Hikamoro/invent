const API_BASE_URL = `${window.location.origin}`;

// Глобальные переменные
let currentCategoryId = null;
let currentCategoryName = null;

// Элементы DOM
const categoryTitle = document.getElementById('categoryTitle');
const backToCategories = document.getElementById('backToCategories');
const createGoodsBtn = document.getElementById('createGoodsBtn');
const deleteGoodsBtn = document.getElementById('deleteGoodsBtn');
const editGoodsBtn = document.getElementById('editGoodsBtn');
const goodsTableBody = document.getElementById('goodsTableBody');

// Модальные окна
const createGoodsModal = document.getElementById('createGoodsModal');
const deleteGoodsModal = document.getElementById('deleteGoodsModal');
const editGoodsModal = document.getElementById('editGoodsModal');
const closeModals = document.querySelectorAll('.close-modal');

// Формы
const createGoodsForm = document.getElementById('createGoodsForm');
const deleteGoodsSelect = document.getElementById('deleteGoodsSelect');
const editGoodsForm = document.getElementById('editGoodsForm');

// Кнопки управления модалками
const createGoodsCancel = document.getElementById('createGoodsCancel');
const deleteGoodsCancel = document.getElementById('deleteGoodsCancel');
const editGoodsCancel = document.getElementById('editGoodsCancel');
const deleteGoodsSubmit = document.getElementById('deleteGoodsSubmit');

// Поля редактирования
const editGoodsId = document.getElementById('editGoodsId');
const editGoodsName = document.getElementById('editGoodsName');
const editGoodsUnit = document.getElementById('editGoodsUnit');
const editGoodsQuantity = document.getElementById('editGoodsQuantity');
const editGoodsDescription = document.getElementById('editGoodsDescription');

// Уведомление
function showNotification(message, type = 'info') {
    const notification = document.getElementById('notification');
    const notificationText = document.getElementById('notificationText');
    notificationText.textContent = message;
    notification.className = `notification ${type}`;
    notification.style.display = 'block';
    setTimeout(() => notification.style.display = 'none', 3000);
}

// Проверка, что мы на странице товаров
function isOnGoodsPage() {
    return window.location.pathname.endsWith('/goods');
}

// Загрузка категории из sessionStorage
function loadCategoryFromStorage() {
    const categoryData = sessionStorage.getItem('selectedCategory');
    if (!categoryData) {
        showNotification('Ошибка: не указана категория', 'error');
        setTimeout(() => window.location.href = '/', 2000);
        return false;
    }
    const category = JSON.parse(categoryData);
    currentCategoryId = Number(category.id);
    currentCategoryName = category.name;
    categoryTitle.textContent = `Товары категории: ${currentCategoryName}`;
    return true;
}

// Загрузка товаров (фильтрация на клиенте)
async function loadGoods() {
    if (!currentCategoryId) return;

    try {
        const response = await fetch(`${API_BASE_URL}/api/goods`);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const allGoods = await response.json();
        const filtered = allGoods.filter(good => good.CategoryID === currentCategoryId);
        displayGoods(filtered);
    } catch (error) {
        console.error(error);
        showNotification('Ошибка загрузки товаров', 'error');
    }
}

// Отрисовка таблицы товаров с кнопками "Редактировать"
function displayGoods(goods) {
    goodsTableBody.innerHTML = '';

    if (goods.length === 0) {
        const row = goodsTableBody.insertRow();
        const cell = row.insertCell(0);
        cell.colSpan = 6;
        cell.textContent = 'В этой категории пока нет товаров';
        cell.style.textAlign = 'center';
        cell.style.padding = '40px';
        cell.style.color = '#7f8c8d';
        return;
    }

    goods.forEach(good => {
        const row = goodsTableBody.insertRow();

        // Ячейка ID
        row.insertCell(0).textContent = good.ID;
        // Название
        row.insertCell(1).textContent = good.Name;
        // Единица измерения
        row.insertCell(2).textContent = good.Unit;
        // Количество
        row.insertCell(3).textContent = good.Quantity;
        // Описание
        row.insertCell(4).textContent = good.Description || '-';

        // Ячейка с кнопкой редактирования
        const actionCell = row.insertCell(5);
        const editBtn = document.createElement('button');
        editBtn.textContent = 'Редактировать';
        editBtn.className = 'btn small primary';
        editBtn.onclick = () => editGood(good.ID);
        actionCell.appendChild(editBtn);
    });
}

// Создание товара
async function createGood(event) {
    event.preventDefault();
    const formData = new FormData(createGoodsForm);
    const goodsData = {
        name: formData.get('name'),
        unit: formData.get('unit'),
        quantity: parseInt(formData.get('quantity'), 10),
        CategoryID: currentCategoryId,
        description: formData.get('description') || ''
    };

    try {
        const response = await fetch(`${API_BASE_URL}/api/goods`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(goodsData)
        });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        showNotification('Товар создан', 'success');
        createGoodsForm.reset();
        createGoodsModal.style.display = 'none';
        await loadGoods();
    } catch (error) {
        console.error(error);
        showNotification('Ошибка создания товара', 'error');
    }
}

// Редактирование – загружаем данные товара и открываем модалку
async function editGood(goodId) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/goods/${goodId}`);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const good = await response.json();

        editGoodsId.value = good.ID;
        editGoodsName.value = good.Name;
        editGoodsUnit.value = good.Unit;
        editGoodsQuantity.value = good.Quantity;
        editGoodsDescription.value = good.Description || '';

        editGoodsModal.style.display = 'block';
    } catch (error) {
        console.error(error);
        showNotification('Ошибка загрузки товара', 'error');
    }
}

// Обновление товара
async function updateGood(event) {
    event.preventDefault();
    const goodsData = {
        name: editGoodsName.value,
        unit: editGoodsUnit.value,
        quantity: parseInt(editGoodsQuantity.value, 10),
        CategoryID: currentCategoryId,
        description: editGoodsDescription.value
    };

    try {
        const response = await fetch(`${API_BASE_URL}/api/goods/${editGoodsId.value}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(goodsData)
        });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        showNotification('Товар обновлён', 'success');
        editGoodsModal.style.display = 'none';
        await loadGoods();
    } catch (error) {
        console.error(error);
        showNotification('Ошибка обновления товара', 'error');
    }
}

// Заполнение выпадающего списка для удаления
async function populateDeleteGoodsSelect() {
    try {
        const response = await fetch(`${API_BASE_URL}/api/goods`);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const allGoods = await response.json();
        const filtered = allGoods.filter(good => good.CategoryID === currentCategoryId);

        deleteGoodsSelect.innerHTML = '<option value="">-- Выберите товар --</option>';
        filtered.forEach(good => {
            const option = document.createElement('option');
            option.value = good.ID;
            option.textContent = `${good.Name} (#${good.ID})`;
            deleteGoodsSelect.appendChild(option);
        });
    } catch (error) {
        console.error(error);
        showNotification('Ошибка загрузки списка товаров', 'error');
    }
}

// Удаление товара по ID
async function deleteGood(goodId) {
    if (!confirm('Удалить товар?')) return;
    try {
        const response = await fetch(`${API_BASE_URL}/api/goods/${goodId}`, { method: 'DELETE' });
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        showNotification('Товар удалён', 'success');
        deleteGoodsModal.style.display = 'none';
        await loadGoods();
    } catch (error) {
        console.error(error);
        showNotification('Ошибка удаления товара', 'error');
    }
}

// ----- Обработчики событий -----
backToCategories.addEventListener('click', () => window.location.href = '/');

createGoodsBtn.addEventListener('click', () => createGoodsModal.style.display = 'block');

deleteGoodsBtn.addEventListener('click', async () => {
    await populateDeleteGoodsSelect();
    deleteGoodsModal.style.display = 'block';
});

editGoodsBtn.addEventListener('click', () => {
    showNotification('Нажмите "Редактировать" рядом с товаром', 'info');
});

createGoodsForm.addEventListener('submit', createGood);
editGoodsForm.addEventListener('submit', updateGood);

deleteGoodsSubmit.addEventListener('click', () => {
    const goodsId = parseInt(deleteGoodsSelect.value, 10);
    if (!goodsId) {
        showNotification('Выберите товар', 'error');
        return;
    }
    deleteGood(goodsId);
});

// Закрытие модальных окон по крестику
closeModals.forEach(btn => {
    btn.addEventListener('click', (e) => {
        e.target.closest('.modal').style.display = 'none';
    });
});

// Отдельные кнопки "Отмена"
createGoodsCancel.addEventListener('click', () => createGoodsModal.style.display = 'none');
deleteGoodsCancel.addEventListener('click', () => deleteGoodsModal.style.display = 'none');
editGoodsCancel.addEventListener('click', () => editGoodsModal.style.display = 'none');

// Закрытие по клику вне модалки
window.addEventListener('click', (event) => {
    if (event.target === createGoodsModal) createGoodsModal.style.display = 'none';
    if (event.target === deleteGoodsModal) deleteGoodsModal.style.display = 'none';
    if (event.target === editGoodsModal) editGoodsModal.style.display = 'none';
});

// Закрытие по Escape
document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape') {
        createGoodsModal.style.display = 'none';
        deleteGoodsModal.style.display = 'none';
        editGoodsModal.style.display = 'none';
    }
});

// Инициализация страницы
document.addEventListener('DOMContentLoaded', async () => {
    if (!isOnGoodsPage()) {
        window.location.href = '/';
        return;
    }
    if (!loadCategoryFromStorage()) return;
    await loadGoods();
    showNotification(`Категория: ${currentCategoryName}`, 'info');
});