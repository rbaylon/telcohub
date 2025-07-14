let map;

function initMap() {
  map = new google.maps.Map(document.getElementById("map"), {
    center: { lat: 10.3157, lng: 123.8854 },
    zoom: 12,
    mapId: "AIzaSyAOAkWiVJ_XseIUPwN3Xtn0Nm74oa01SQI",
  });

  map.addListener("click", (e) => {
    placeMarker(e.latLng);
    openModal(e.latLng.lat(), e.latLng.lng());
  });

  fetchMarkers();
}


function placeMarker(latLng) {
  new google.maps.marker.AdvancedMarkerElement({
    map,
    position: latLng,
  });
}

let markerObjects = [];
let filteredmarkersObjects = [];

function fetchMarkers() {
  fetch("/marker/list")
    .then(res => res.json())
    .then(data => {
      markerObjects.forEach(m => m.marker.setMap(null)); // Clear old markers
      markerObjects = [];
      data.forEach(m => {
        const pinBackground = document.createElement('div');
        pinBackground.className = "marker-class";
        pinBackground.textContent = `${m.Title} ${m.Description}`;
        const marker = new google.maps.marker.AdvancedMarkerElement({
          map,
          position: { lat: m.Latitude, lng: m.Longitude },
          title: m.Title,
          content: pinBackground,
        });

        marker.content.classList.add("drop");
        if (m.Category.color){
          marker.content.classList.add(`${m.Category.color}`);
        } else {
          marker.content.classList.add("bg-blue-500");
        }
        marker.content.classList.add("bg-opacity-75");
        marker.content.classList.add("text-white");
        const infoWindow = new google.maps.InfoWindow({
          headerContent: `${m.Category.name}`,
          content: `
            <div class="text-sm text-gray-700">
              <p class="font-bold text-blue-500">${m.Title} ${m.Description}</p>
              <p>Lat: ${m.Latitude} Long: ${m.Longitude}</p>
            </div>
          `,
        });
        marker.addListener("click", () => {
          infoWindow.open(map, marker);
        });
        const markerdata = {
          "data": m,
          "marker": marker,
          "infowindow": infoWindow,
        }; 
        markerObjects.push(markerdata);
      });
      console.log(markerObjects);
      populateTable(markerObjects);
      filteredmarkersObjects = markerObjects;
    })
    .catch(err => console.error("Failed to load markers:", err));
}

function focusMarker(index) {
  const marker = filteredmarkersObjects[index];
  const infoWindow = marker.infowindow;
  window.scrollTo(0, 0);
  map.panTo(marker.marker.position);
  map.setZoom(17); // or adjust as needed

  infoWindow.open(map, marker.marker);
}

function filterMarkers() {
  const query = document.getElementById("markerSearch").value.toLowerCase();
  const filtered = markerObjects.filter(m =>
    m.data.Title.toLowerCase().includes(query) ||
    m.data.Description.toLowerCase().includes(query) ||
    m.data.User.Username.toLowerCase().includes(query)
  );
  populateTable(filtered);
  filteredmarkersObjects = filtered;
}