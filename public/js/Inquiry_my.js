
(function(own) {

	own.setCompany = function(){
		var companySelect = document.getElementById('companySelect');
		var com_select = document.createElement('select');
		com_select.className= "weui_select"
		
		var companyList = ['选择快递公司(可不选择)','EMS','申通快递','顺丰速运','圆通速递','韵达快递','中通快递','百世快递','天天快递','宅急送','京东','德邦物流','邮政包裹','联邦快递','优速物流','aae全球专递','安捷快递','安信达快递','彪记快递','bht','dpex','d速快递','递四方','fedex（国外）','飞康达物流','凤凰快递','飞快达','国通快递','港中能达物流','广东邮政物流','共速达','汇通快运','恒路物流','华夏龙物流','海红','海外环球','佳怡物流','京广速递','急先达','佳吉物流','加运美物流','金大物流','嘉里大通','晋越快递','快捷速递','联昊通物流','龙邦物流','立即送','乐捷递','民航快递','美国快递','门对门','OCS','配思货运','全晨快递','全峰快递','全际通物流','全日通快递','全一快递','如风达','三态速递','盛辉物流','速尔物流','盛丰物流','赛澳递','天地华宇','tnt','ups','万家物流','伍圆','万象物流','新邦物流','信丰物流','亚风速递','一邦速递','远成物流','源伟丰快递','元智捷诚快递','越丰物流','运通快递','源安达','银捷速递','中铁快运','中邮物流','忠信达','芝麻开门'];
		var html = '';
		for (var i = 0; i < companyList.length; i++) {
			html = html +'<option value="'+companyList[i]+'">'+companyList[i]+'</option>'
		}
		com_select.innerHTML = html;
		companySelect.appendChild(com_select);
		com_select.style.color = 'lightgray'
		
	}

})(window);