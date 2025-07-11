window.onload = () => {
  fetch("/admin/category/list")
    .then(res => res.json())
    .then(data => populateCategoryTable(data))
    .catch(err => console.error("Failed to load admin data:", err));
};

function updateColorPreview() {
  const input = document.getElementById("colorInput");
  const preview = document.getElementById("colorPreview");

  // Remove any previous color classes
  preview.className = "w-8 h-8 rounded border";

  // Add new class dynamically
  const colorClass = input.value.trim();
  if (colorClass) preview.classList.add(colorClass);
}

function populateCategoryTable(cats) {
  const table = document.getElementById("category-table");

  table.innerHTML = cats.map(c => `
    <tr>
      <td class="px-4 py-2">${c.name}</td>
      <td class="px-4 py-2">${c.color}</td>
      <td class="px-4 py-2 space-x-2">
        <button title="Edit Category">
          <i class="fas fa-edit text-indigo-500"></i>
        </button>
        <button onclick="deleteCategory(${c.id})" title="Delete Category">
          <i class="fas fa-trash text-red-500"></i>
        </button>
      </td>
    </tr>
  `).join('');
}

function deleteCategory(id) {
  if (confirm("Are you sure you want to delete this category?")) {
    fetch(`/admin/category/delete/${id}`, { method: "POST" })
      .then(() => location.reload());
  }
}