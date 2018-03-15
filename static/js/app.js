function accessibilityToggle(e, whichtoggle) {
    var keynum;
    if (window.event) { // IE
        keynum = e.keyCode;
    } else if (e.which) { // Netscape/Firefox/Opera
        keynum = e.which;
    }

    if (keynum == 13) {
        if (!$("#" + whichtoggle).hasClass("in")) {
            $("#" + whichtoggle).addClass("in");
            $("#" + whichtoggle).removeClass("collapse");
        } else {
            $("#" + whichtoggle).removeClass("in");
            $("#" + whichtoggle).addClass("collapse");
        }
    }
}


function goBack() {
    window.history.back();
}