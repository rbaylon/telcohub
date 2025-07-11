let map;

function initMap() {
  map = new google.maps.Map(document.getElementById("map"), {
    center: { lat: 10.3157, lng: 123.8854 },
    zoom: 12,
  });

  map.addListener("click", (e) => {
    placeMarker(e.latLng);
    openModal(e.latLng.lat(), e.latLng.lng());
  });

  fetchMarkers();
}


function placeMarker(latLng) {
  new google.maps.Marker({
    position: latLng,
    map: map,
  });
}

let markerObjects = [];

function filterMarkers() {
  const query = document.getElementById("markerSearch").value.toLowerCase();
  const filtered = markerObjects.filter(m =>
    m.Title.toLowerCase().includes(query) ||
    m.Description.toLowerCase().includes(query) ||
    m.User.Username.toLowerCase().includes(query)
  );
  populateTable(filtered);
}

function fetchMarkers() {
  fetch("/marker/list")
    .then(res => res.json())
    .then(data => {
      markerObjects.forEach(m => m.setMap(null)); // Clear old markers
      markerObjects = [];
      data.forEach(m => {
        const marker = new google.maps.Marker({
          position: { lat: m.Latitude, lng: m.Longitude },
          map: map,
          title: m.Title,
        });

        const infoWindow = new google.maps.InfoWindow({
          content: `
            <div class="text-sm ${m.Category.color}">
              <h3 class="font-semibold text-blue-600">${m.Title}</h3>
              <p class="text-gray-700">${m.Description}</p>
            </div>
          `,
        });
        // 👇 Show the window immediately
        infoWindow.open(map, marker);
        marker.addListener("click", () => {
          infoWindow.open(map, marker);
        });
        markerObjects.push(m);
      });
      console.log(data);
      populateTable(data);
    })
    .catch(err => console.error("Failed to load markers:", err));
}


