(function () {
  function getCookie(n) {
    var m = document.cookie.match(new RegExp('(?:^|; )' + n + '=([^;]*)'));
    return m ? decodeURIComponent(m[1]) : '';
  }
  var t = getCookie('theme');
  if (t === 'dark') document.documentElement.setAttribute('data-theme', 'dark');
  else if (t === 'light') document.documentElement.setAttribute('data-theme', 'light');
  else if (t === 'auto') document.documentElement.removeAttribute('data-theme');
  else document.documentElement.setAttribute('data-theme', 'dark'); // no preference set yet — default to dark
})();

function setTheme(t) {
  document.documentElement.setAttribute('data-theme', t === 'auto' ? '' : t);
  document.cookie = 'theme=' + t + ';path=/;max-age=31536000;samesite=lax';
  document.querySelectorAll('.theme-btn').forEach(function (btn) {
    btn.classList.toggle('active', btn.dataset.theme === t);
  });
}

function isDarkTheme() {
  var attr = document.documentElement.getAttribute('data-theme');
  if (attr === 'dark') return true;
  if (attr === 'light') return false;
  return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
}

// CARTO's free anonymous dark_all raster tiles now bounce unauthenticated
// requests to a "get an API key" placeholder — an account-scoped key query
// param is required (confirmed against a real CARTO account key: the param
// must be named "key", not "api_key" — the latter silently returns the same
// watermarked placeholder as no key at all, no error). If none is configured
// server-side, fall back to plain OpenStreetMap tiles rather than ever
// hitting CARTO unauthenticated.
function mapTileLayer(L, cartoKey) {
  var dark = isDarkTheme() && cartoKey;
  return L.tileLayer(
    dark
      ? 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png?key=' + encodeURIComponent(cartoKey)
      : 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
    {
      maxZoom: 19,
      attribution: dark
        ? '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
        : '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }
  );
}
