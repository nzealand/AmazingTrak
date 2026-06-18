(function () {
  function getCookie(n) {
    var m = document.cookie.match(new RegExp('(?:^|; )' + n + '=([^;]*)'));
    return m ? decodeURIComponent(m[1]) : '';
  }
  var t = getCookie('theme');
  if (t === 'dark') document.documentElement.setAttribute('data-theme', 'dark');
  else if (t === 'light') document.documentElement.setAttribute('data-theme', 'light');
})();

function setTheme(t) {
  document.documentElement.setAttribute('data-theme', t === 'auto' ? '' : t);
  document.cookie = 'theme=' + t + ';path=/;max-age=31536000;samesite=lax';
  document.querySelectorAll('.theme-btn').forEach(function (btn) {
    btn.classList.toggle('active', btn.dataset.theme === t);
  });
}
