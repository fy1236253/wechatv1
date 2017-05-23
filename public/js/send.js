$(document).ready(function() {

	function getRandom(n){
		return Math.floor(Math.random()*n+1)
	}
	//var imgs = [" ","#5cb85c","#5bc0de","#f0ad4e","#d9534f","#f9690e","#5bc0de"];
	var arr = new Array();
	for (var i = 0; i < 10; i++) {
		//var n = getRandom(5);
		var ss = $("h4").eq(i).html()
		switch(ss){
		case "韵达速递":
		  $("li img").eq(i).attr("src","/img/yunda.png");		
		  break;
		case "顺丰速运":
		  $("li img").eq(i).attr("src","/img/shunfeng.png");
		  break;
		case "京东快递":
		  $("li img").eq(i).attr("src","/img/jingdong.png");
		  break;
		case "圆通速递":
		  $("li img").eq(i).attr("src","/img/yuantong.png");		
		  break;
		case "申通快递":
		  $("li img").eq(i).attr("src","/img/shentong.png");
		  break;
		case "百世快递":
		  $("li img").eq(i).attr("src","/img/huitong.png");
		  break;
		case "天天快递":
		  $("li img").eq(i).attr("src","/img/tiantian.png");		
		  break;
		case "中国邮政":
		  $("li img").eq(i).attr("src","/img/ems.png");		
		  break;
		case "中通快递":
		  $("li img").eq(i).attr("src","/img/zhongtong.png");		
		  break;
		default:
		  $("li img").eq(i).attr("src","/img/logo@2x.png");
		}
		//$("li img").attr("src","/img/yunda.png");
		arr.push(ss)	
	}
});