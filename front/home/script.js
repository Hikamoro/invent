const API_BASE_URL = `${window.location.origin}`;

// Элементы DOM
const selectCategoryBtn = document.getElementById('selectCategoryBtn');
const createCategoryBtn = document.getElementById('createCategoryBtn');
const deleteCategoryBtn = document.getElementById('deleteCategoryBtn');
const categoryDropdown = document.getElementById('categoryDropdown');
const categoryList = document.getElementById('categoryList');
const closeDropdown = document.getElementById('closeDropdown');

// Модальные окна
const createModal = document.getElementById('createModal');
const deleteModal = document.getElementById('deleteModal');
const closeModal = document.querySelector('.close-modal');
const createCategorySubmit = document.getElementById('createCategorySubmit');
const createCategoryCancel = document.getElementById('createCategoryCancel');
const deleteCategorySubmit = document.getElementById('deleteCategorySubmit');
const deleteCategoryCancel = document.getElementById('deleteCategoryCancel');

// Формы
const categoryNameInput = document.getElementById('categoryName');
const deleteCategorySelect = document.getElementById('deleteCategorySelect');

// Уведомления
const notification = document.getElementById('notification');
const notificationText = document.getElementById('notificationText');

// Показать уведомление
function showNotification(message, type = 'info') {
    notificationText.textContent = message;
    notification.className = `notification ${type}`;
    notification.style.display = 'block';
    
    setTimeout(() => {
        notification.style.display = 'none';
    }, 3000);
}

// Загрузить категории
async function loadCategories() {
    try {
        const response = await fetch(`${API_BASE_URL}/api/categories`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const categories = await response.json();
        return categories;
    } catch (error) {
        console.error('Error loading categories:', error);
        showNotification('Ошибка загрузки категорий', 'error');
        return [];
    }
}

// Отобразить категории в выпадающем списке
async function showCategoryDropdown() {
    const categories = await loadCategories();
    categoryList.innerHTML = '';
    
    if (categories.length === 0) {
        categoryList.innerHTML = '<div class="category-item">Нет категорий</div>';
    } else {
        categories.forEach(category => {
            const categoryItem = document.createElement('div');
            categoryItem.className = 'category-item';
            categoryItem.innerHTML = `
                <span class="category-name">${category.name}</span>
                <span class="category-id">#${category.id}</span>
            `;
            categoryItem.addEventListener('click', () => {
                showNotification(`Выбрана категория: ${category.name}`, 'success');
                categoryDropdown.style.display = 'none';

                const categoryData = { id: category.id, name: category.name };
                sessionStorage.setItem('selectedCategory', JSON.stringify(categoryData));
                window.location.href = '/goods';

            });
            categoryList.appendChild(categoryItem);
        });
    }
    
    categoryDropdown.style.display = 'block';
}

// Отобразить категории в выпадающем списке для удаления
async function populateDeleteCategories() {
    const categories = await loadCategories();
    deleteCategorySelect.innerHTML = '<option value="">-- Выберите категорию --</option>';
    
    categories.forEach(category => {
        const option = document.createElement('option');
        option.value = category.id;
        option.textContent = `${category.name} (#${category.id})`;
        deleteCategorySelect.appendChild(option);
    });
}

// Создать новую категорию
async function createCategory() {
    const name = categoryNameInput.value.trim();
    
    if (!name) {
        showNotification('Введите название категории', 'error');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/api/categories`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name }),
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const newCategory = await response.json();
        showNotification(`Категория "${name}" создана успешно`, 'success');
        categoryNameInput.value = '';
        createModal.style.display = 'none';
        
        // Обновляем выпадающие списки
        await showCategoryDropdown();
        await populateDeleteCategories();
        
    } catch (error) {
        console.error('Error creating category:', error);
        showNotification('Ошибка создания категории', 'error');
    }
}

// Удалить категорию
async function deleteCategory() {
    const categoryId = parseInt(deleteCategorySelect.value);
    
    if (!categoryId) {
        showNotification('Выберите категорию для удаления', 'error');
        return;
    }
    
    const category = deleteCategorySelect.options[deleteCategorySelect.selectedIndex].text;
    
    if (!confirm(`Вы уверены, что хотите удалить категорию "${category}"?`)) {
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/api/categories/${categoryId}`, {
            method: 'DELETE',
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        showNotification(`Категория удалена успешно`, 'success');
        deleteCategorySelect.value = '';
        deleteModal.style.display = 'none';
        
        // Обновляем выпадающие списки
        await showCategoryDropdown();
        await populateDeleteCategories();
        
    } catch (error) {
        console.error('Error deleting category:', error);
        showNotification('Ошибка удаления категории', 'error');
    }
}

// События
selectCategoryBtn.addEventListener('click', showCategoryDropdown);
createCategoryBtn.addEventListener('click', () => {
    createModal.style.display = 'block';
});
deleteCategoryBtn.addEventListener('click', async () => {
    await populateDeleteCategories();
    deleteModal.style.display = 'block';
});

closeDropdown.addEventListener('click', () => {
    categoryDropdown.style.display = 'none';
});

closeModal.addEventListener('click', () => {
    createModal.style.display = 'none';
    deleteModal.style.display = 'none';
});

createCategoryCancel.addEventListener('click', () => {
    createModal.style.display = 'none';
    categoryNameInput.value = '';
});

deleteCategoryCancel.addEventListener('click', () => {
    deleteModal.style.display = 'none';
    deleteCategorySelect.value = '';
});

createCategorySubmit.addEventListener('click', createCategory);
deleteCategorySubmit.addEventListener('click', deleteCategory);

// Закрытие модальных окон при клике вне их
window.addEventListener('click', (event) => {
    if (event.target === createModal) {
        createModal.style.display = 'none';
        categoryNameInput.value = '';
    }
    if (event.target === deleteModal) {
        deleteModal.style.display = 'none';
        deleteCategorySelect.value = '';
    }
});

// Обработка нажатия Escape
document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape') {
        categoryDropdown.style.display = 'none';
        createModal.style.display = 'none';
        deleteModal.style.display = 'none';
        categoryNameInput.value = '';
        deleteCategorySelect.value = '';
    }
});

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', async () => {
    showNotification('Система инвентаризации загружена', 'info');
});