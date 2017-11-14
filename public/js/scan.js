

$(function () {
    $("#reload").click(function () {
        $("#localImg").attr("src", "")
        $("#scanner").show();
        $(".scan-ok").hide();
    });
    $("#submit").click(function () {

        window.location.href = "/credits"
    })
    $("#reload").click(function(){
        window.location.href = "/scanner"
    })
})