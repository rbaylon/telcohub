
function openModal(lat, lng, marker = null) {
  const modal = document.getElementById("markerModal");
  const form = modal.querySelector("form");

  modal.classList.remove("hidden");

  if (marker) {
    document.getElementById("modalTitle").textContent = "Edit Marker";
    form.title.value = marker.Title;
    form.description.value = marker.Description;
    form.latitude.value = marker.Latitude;
    form.longitude.value = marker.Longitude;
    form.marker_id.value = marker.ID;
  } else {
    document.getElementById("modalTitle").textContent = "Add Marker";
    form.reset();
    form.latitude.value = lat;
    form.longitude.value = lng;
    form.marker_id.value = "";
  }
}


function closeModal() {
  document.getElementById("markerModal").classList.add("hidden");
}

document.getElementById("markerForm")?.addEventListener("submit", function (e) {
  e.preventDefault();
  const form = e.target;
  const data = new URLSearchParams(new FormData(form));

  const markerId = form.marker_id.value;
  const endpoint = markerId ? `/marker/edit/${markerId}` : `/marker/create`;

  fetch(endpoint, {
    method: "POST",
    body: data,
  })
    .then(() => {
      closeModal();
      location.reload();
    })
    .catch(err => {
      alert("Failed to submit marker");
      console.error(err);
    });
});

// Fill table with markers
function populateTable(markers) {
  const table = document.getElementById("marker-table");
  table.innerHTML = markers.map(m =>
    `<tr>
      <td class="px-4 py-2">${m.Title}</td>
      <td class="px-4 py-2 hidden md:table-cell">${m.Description}</td>
      <td class="px-4 py-2">${m.Latitude}</td>
      <td class="px-4 py-2">${m.Longitude}</td>
      <td class="px-4 py-2">${m.User.Username}</td>
      <td class="px-4 py-2 space-x-2">
        <button onclick="editMarker(${m.ID})" title="Edit">
          <i class="fas fa-edit text-blue-500"></i>
        </button>
        <button onclick="deleteMarker(${m.ID})" title="Delete">
          <i class="fas fa-trash text-red-500"></i>
        </button>
      </td>
    </tr>`
  ).join('');
}

function editMarker(id) {
  fetch("/marker/list")
    .then(res => res.json())
    .then(data => {
      const marker = data.find(m => m.ID === id);
      openModal(marker.Latitude, marker.Longitude, marker);
    });
}

function deleteMarker(id) {
  if (confirm("Are you sure you want to delete this marker?")) {
    fetch(`/marker/delete/${id}`, { method: "POST" })
      .then(() => location.reload());
  }
}

function toggleUserMenu() {
  const menu = document.getElementById("userMenu");
  menu.classList.toggle("hidden");
}

document.addEventListener("click", function (e) {
  const toggleBtn = e.target.closest("button[onclick='toggleUserMenu()']");
  const menu = document.getElementById("userMenu");

  if (!toggleBtn && !e.target.closest("#userMenu")) {
    menu.classList.add("hidden");
  }
});

