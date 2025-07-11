window.onload = () => {
  fetch("/admin/dashboard")
    .then(res => res.json())
    .then(data => populateAdminTable(data))
    .catch(err => console.error("Failed to load admin data:", err));
};

function populateAdminTable(users) {
  const table = document.getElementById("admin-user-table");

  table.innerHTML = users.map(u => `
    <tr>
      <td class="px-4 py-2">${u.username}</td>
      <td class="px-4 py-2">
        <select onchange="updateRole(${u.id}, this.value)" class="border px-2 py-1 rounded">
          <option value="user" ${u.role === "user" ? "selected" : ""}>User</option>
          <option value="admin" ${u.role === "admin" ? "selected" : ""}>Admin</option>
        </select>
      </td>
      <td class="px-4 py-2">${u.markerCount}</td>
      <td class="px-4 py-2 space-x-2">
        <button onclick="showHistory(${u.id})" title="View History">
          <i class="fas fa-history text-indigo-500"></i>
        </button>
        <button onclick="deleteUser(${u.id})" title="Delete User">
          <i class="fas fa-user-slash text-red-500"></i>
        </button>
      </td>
    </tr>
  `).join('');
}

function updateRole(userId, role) {
  fetch(`/admin/user/${userId}/role`, {
    method: "POST",
    body: new URLSearchParams({ role }),
  })
    .then(() => alert("✅ Role updated successfully"))
    .catch(err => {
      alert("❌ Failed to update role");
      console.error(err);
    });
}

function deleteUser(userId) {
  if (!confirm("Are you sure you want to permanently delete this user?")) return;

  fetch(`/admin/user/${userId}/delete`, {
    method: "DELETE",
  })
    .then(() => location.reload())
    .catch(err => {
      alert("❌ Failed to delete user");
      console.error(err);
    });
}

function showHistory(userId) {
  fetch("/admin/dashboard")
    .then(res => res.json())
    .then(data => {
      const user = data.find(u => u.id === userId);
      const modal = document.getElementById("historyModal");
      const list = user.markerHistory.map(m => `
        <li>
          <strong>${m.title}</strong>: ${m.description}<br/>
          🌍 (${m.latitude.toFixed(4)}, ${m.longitude.toFixed(4)})<br/>
          📅 ${m.created_at}
        </li>
      `).join('');

      modal.querySelector(".modal-body").innerHTML = list;
      modal.classList.remove("hidden");
    })
    .catch(err => {
      alert("❌ Failed to load marker history");
      console.error(err);
    });
}

function closeHistory() {
  document.getElementById("historyModal").classList.add("hidden");
}
