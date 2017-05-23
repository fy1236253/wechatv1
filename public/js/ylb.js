

wx.ready(function () {
  // 1 判断当前版本是否支持指定 JS 接口，支持批量判断
  document.querySelector('#give-me-hand').onclick = function () {
    window.open('./give-me-hand/','');
  };

});

wx.error(function (res) {
  alert(res.errMsg);
});
