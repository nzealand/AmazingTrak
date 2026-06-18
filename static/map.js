document.addEventListener('DOMContentLoaded', function () {
  var mapEl = document.getElementById('hero-map');
  if (!mapEl) return;
  var lat = parseFloat(mapEl.dataset.lat);
  var lon = parseFloat(mapEl.dataset.lon);
  if (isNaN(lat) || isNaN(lon)) return;

  var map = L.map('hero-map').setView([lat, lon], 13);
  L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
  }).addTo(map);
  L.marker([lat, lon]).addTo(map).bindPopup('Photo taken here').openPopup();
});
