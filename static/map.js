document.addEventListener('DOMContentLoaded', function () {
  var mapEl = document.getElementById('hero-map');
  if (!mapEl) return;
  var lat = parseFloat(mapEl.dataset.lat);
  var lon = parseFloat(mapEl.dataset.lon);
  if (isNaN(lat) || isNaN(lon)) return;

  var map = L.map('hero-map').setView([lat, lon], 13);
  mapTileLayer(L).addTo(map);
  L.marker([lat, lon]).addTo(map).bindPopup('Photo taken here').openPopup();
});
