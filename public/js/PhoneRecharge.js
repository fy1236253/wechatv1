(function(own) {
	var rechargeItems = [];
	var curChooseIndex = -1;
	var currentChooseName = ''; 
	var currentChoosePrice = ''; 
	var currentChoosePriceY = ''; 
	var prePhone = '';
	var currentPhone = ''; 

	//设置recharegeItems,将之前选中的数据全部设置为初始
	own.initSomeRechargeItem = function(itemArr) {
		curChooseIndex = -1;
		var ul = document.getElementById('recharge_item');
		ul.innerHTML = '';
		rechargeItems = itemArr;
		for (var i = 0; i < rechargeItems.length; i++) {
			var item_recharge = rechargeItems[i];
			var li = document.createElement('li');
			li.id = 'item' + i;
			html = '<div class="li_style">\
			<div class="li_border">\
			<p class="p1">'+item_recharge.name+'</p>\
			<p class="p2">(原价'+item_recharge.real_money+'元)</p>\
			</div>\
			</div>';
			li.innerHTML = html;
			ul.appendChild(li);
			li.addEventListener('tap', function() {
				var li_pause = document.getElementById('recharge_item').getElementsByTagName('li');
				for (var j = 0; j < li_pause.length; j++) {
					var styleDiv_pause = li_pause[j].getElementsByClassName('li_style')[0];
					var borederDiv_pause = styleDiv_pause.getElementsByClassName('li_border')[0];
					borederDiv_pause.style.border = '1px solid white';
					var liP_pause = borederDiv_pause.getElementsByTagName('p');
					liP_pause[0].style.color = 'black';
					liP_pause[1].style.color = 'darkgray'
				}
				var curLi = this;
				var styleDiv = curLi.getElementsByClassName('li_style')[0];
				var borderDiv = styleDiv.getElementsByClassName('li_border')[0];
				borderDiv.style.border = '1px solid #56BCED';
				var liP = borderDiv.getElementsByTagName('p');
				liP[0].style.color = '#56BCED';
				liP[1].style.color = '#56BCED';
				var liidstr = curLi.id;
				var idStr = liidstr.substr(4, liidstr.length - 4);
				curChooseIndex = parseInt(idStr);
				setChooseItem();
			});
		}
	}

	//设置选中
	own.setChooseItem = function() {
		if (curChooseIndex >= 0 && rechargeItems.length > 0) {
			//从数组中去获取
			var real_money_value = document.getElementById('real_money_value');
			var pref_money_value = document.getElementById('pref_money_value');
			real_money_value.innerText = '' + (rechargeItems[curChooseIndex].real_money - rechargeItems[curChooseIndex].prel_money);
			pref_money_value.innerText = '' + rechargeItems[curChooseIndex].prel_money;

			currentChooseName = rechargeItems[curChooseIndex].name
			currentChoosePriceY = rechargeItems[curChooseIndex].real_money
			currentChoosePrice = real_money_value.innerText

		}
	}

	//监听input变化
	own.listenInput = function() {
		var num = document.getElementById('num');
		num.oninput = function() {
			if (this.value.length >= 11) {
				if (prePhone == this.value.substr(0, 11)) {
					this.value = this.value.substr(0, 11);
					return;
				} else {
					prePhone = this.value.substr(0, 11);
					this.value = this.value.substr(0, 11);
					own.getPhoneDetailAndRechargeItemsDetail(this.value);
				}
			} else {
				prePhone = '';
			}
		}
		num.OnPropChanged = function() {
			if (this.value.length >= 11) {
				if (prePhone == this.value.substr(0, 11)) {
					this.value = this.value.substr(0, 11);
					return;
				} else {
					prePhone = this.value.substr(0, 11);
					this.value = this.value.substr(0, 11);
					own.getPhoneDetailAndRechargeItemsDetail(this.value);
				}
			} else {
				prePhone = '';
			}
		}
	}

	own.setPhoneNumDetail = function(numDetail) {
		var phone_detail = document.getElementById('phone_detail');
		phone_detail.innerText = numDetail;
	}

	own.getPhoneDetailAndRechargeItemsDetail = function(phoneNum) {
		//$.get('http://api.fnbird.com/mobile/index.php?act=index&op=commonInfo', function(data, status, xhr) {})	

		// {10M 3} {30M 5} {100M 10} {300M 20}  {500M 30}   // {1G 50} {2G 70} {3G 100} {4G 130} {6G 180} {11G 280} 
		// {20M 3} {50M 6} {100M 10} {200M 15}  {500M 30}
		// {5M 1} {10M 2} {30M 5} {50M 7} {100M 10} {200M 15} {500M 30} //  {1G 50}

        var chinaMobile = "|134|135|136|137|138|139|147|150|151|152|157|158|159|178|182|183|184|187|188";
        var chinaUnicom = "|130|131|132|145|155|156|175|176|185|186";
        var chinaTelecom= "|133|153|177|180|181|189";

        if (chinaMobile.indexOf(phoneNum.substr(0,3)) > 0) {
			var itemArray = [{
				"name": "10M",
				"real_money": 3.00,
				"prel_money": 0.30
			},{
				"name": "30M",
				"real_money": 5.00,
				"prel_money": 0.50
			},{
				"name": "100M",
				"real_money": 10.00,
				"prel_money": 1.00
			},{
				"name": "300M",
				"real_money": 20.00,
				"prel_money": 2.00
			},{
				"name": "500M",
				"real_money": 30.00,
				"prel_money": 3.00
			}];

			initSomeRechargeItem(itemArray);
			setPhoneNumDetail('中国移动');
		}

		// {20M 3} {50M 6} {100M 10} {200M 15}  {500M 30}
		if (chinaUnicom.indexOf(phoneNum.substr(0,3)) > 0) {
			var itemArray = [{
				"name": "20M",
				"real_money": 3.00,
				"prel_money": 0.30
			},{
				"name": "50M",
				"real_money": 6.00,
				"prel_money": 0.60
			},{
				"name": "100M",
				"real_money": 10.00,
				"prel_money": 1.00
			},{
				"name": "200M",
				"real_money": 15.00,
				"prel_money": 1.50
			},{
				"name": "500M",
				"real_money": 30.00,
				"prel_money": 3.00
			}];

			initSomeRechargeItem(itemArray);
			setPhoneNumDetail('中国联通');
		}


		// {5M 1} {10M 2} {30M 5} {50M 7} {100M 10} {200M 15} {500M 30}
		if (chinaTelecom.indexOf(phoneNum.substr(0,3)) > 0) {
			var itemArray = [{
				"name": "5M",
				"real_money": 1.00,
				"prel_money": 0.10
			},{
				"name": "10M",
				"real_money": 2.00,
				"prel_money": 0.20
			},{
				"name": "30M",
				"real_money": 5.00,
				"prel_money": 0.50
			},{
				"name": "50M",
				"real_money": 7.00,
				"prel_money": 0.70
			},{
				"name": "100M",
				"real_money": 10.00,
				"prel_money": 1.00
			},{
				"name": "200M",
				"real_money": 15.00,
				"prel_money": 1.50
			},{
				"name": "500M",
				"real_money": 30.00,
				"prel_money": 3.00
			}];

			initSomeRechargeItem(itemArray);
			setPhoneNumDetail('中国电信');
		}
	}

	own.getcurChooseIndex = function() {
		var num = document.getElementById('num');
		if (num.value.length < 11) {
			alert("请输入正确的手机号码");
			return -1;
		}
		if (num.value.substr(0, 1) != '1') {
			alert("请输入正确的手机号码");
			return -1;
		}
		if (curChooseIndex == -1) {
			alert("请选择要充值的流量包");
			return -1;
		}
		
		return curChooseIndex;
	}

	own.getcurrentPhone = function() {
		var num = document.getElementById('num');
		return num.value;
	}

	own.getcurChooseName = function() {
		return currentChooseName;
	}

	own.getcurChoosePrice = function() {
		return currentChoosePrice;
	}

	own.getcurChoosePriceY = function() {
		return currentChoosePriceY;
	}
})(window);