
wx.ready(function () {

  // 2. 分享接口
  // 2.1 监听“分享给朋友”，按钮点击、自定义分享内容及分享结果接口
  

	wx.onMenuShareAppMessage({
	  title: '云喇叭快递通知助手，关注送话费啦！',
	  desc: '关注云喇叭，第一时间获悉快递到达通知。快递员帮手，您的贴心小秘书。',
	  link: 'http://wechat2.shenbianvip.com/h5/r',
	  imgUrl: 'http://ylb-pub.oss-cn-hangzhou.aliyuncs.com/file-bag/ylb-hb.jpeg',
	  trigger: function (res) {
	    // 不要尝试在trigger中使用ajax异步请求修改本次分享的内容，因为客户端分享操作是一个同步操作，这时候使用ajax的回包会还没有返回
	  },
	  success: function (res) {
	    //alert('已分享');
	  },
	  cancel: function (res) {
	    //alert('已取消');
	  },
	  fail: function (res) {
	    //alert(JSON.stringify(res));
	  }
	});

    // 点击发送 验证码 按钮 
    document.querySelector('#sendSnCode').onclick = function () {
    	var m = $('#userPhoneNumber')[0].value;

      var x = $('#imgcode')[0].value;

    	var $cellSnCode = $('#cellSnCode');  // 
      var $cellImgCode = $('#cellImgCode');  // 
    	var $sendSnCodeBtn = $('#sendSnCodeBtn'); 
    	var $showQrCodeBtn = $('#showQrCodeBtn');

    	$.getJSON('/api/v1/mobile/sn/send?m='+m+'&x='+x, function (d) {
    		if (d.msg == 'success') {
    			$cellSnCode.show();   // 现实 验证码 输入框 
          $cellImgCode.hide(); 

    			$sendSnCodeBtn.hide();
    			$showQrCodeBtn.show(); 
    		} else {
          alert("图像验证码或手机号码错误");
        }
		});
    	
    };

  	document.querySelector('#showQrCode').onclick = function () {
  		var m = $('#userPhoneNumber')[0].value;
  		var sn = $('#snkey')[0].value; // 验证码 

  		$.getJSON('/api/v1/mobile/sn/check?m='+m+'&x='+sn, function (d) {
    		if (d.msg == 'success') {

    			var $loadingToast = $('#loadingToast');
		      $loadingToast.show();
          $('#qrImg').attr('src', "/img/qr_sence_123.jpeg"); // qr_sence_123.jpeg
				  $loadingToast.hide();

				  var $dialog = $('#dialog2');
          $dialog.show();
          $dialog.find('.weui_btn_dialog').one('click', function () {$dialog.hide();});

    		} else {
    			alert("验证码错误");
    			return;
    		}
		});
  		
			
  	};

});


wx.error(function (res) {
  alert(res.errMsg);
});
