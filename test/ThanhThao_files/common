/* Minification failed. Returning unminified contents.
(456,26-27): run-time error JS1195: Expected expression: .
(456,39-40): run-time error JS1195: Expected expression: .
(456,54-55): run-time error JS1003: Expected ':': )
(457,21-24): run-time error JS1009: Expected '}': var
(457,21-24): run-time error JS1003: Expected ':': var
(457,63-64): run-time error JS1006: Expected ')': ;
(464,19-23): run-time error JS1006: Expected ')': else
(464,18): run-time error JS1004: Expected ';'
(464,33-34): run-time error JS1195: Expected expression: .
(464,59-60): run-time error JS1003: Expected ':': )
(465,21-25): run-time error JS1009: Expected '}': this
(465,21-25): run-time error JS1006: Expected ')': this
(466,19-23): run-time error JS1009: Expected '}': else
(166,17): run-time error JS1004: Expected ';'
(466,19-23): run-time error JS1034: Unmatched 'else'; no 'if' defined: else
(469,13-14): run-time error JS1002: Syntax error: }
(469,14-15): run-time error JS1195: Expected expression: )
(470,5-6): run-time error JS1002: Syntax error: }
(470,6-7): run-time error JS1195: Expected expression: ,
(472,18): run-time error JS1004: Expected ';'
(478,6-7): run-time error JS1195: Expected expression: ,
(483,34-35): run-time error JS1010: Expected identifier: (
(492,6-7): run-time error JS1195: Expected expression: ,
(494,33-34): run-time error JS1010: Expected identifier: (
(503,6-7): run-time error JS1195: Expected expression: ,
(505,47-48): run-time error JS1010: Expected identifier: (
(513,6-7): run-time error JS1195: Expected expression: ,
(518,28-29): run-time error JS1010: Expected identifier: (
(527,6-7): run-time error JS1195: Expected expression: ,
(528,33-34): run-time error JS1010: Expected identifier: (
(534,6-7): run-time error JS1195: Expected expression: ,
(535,32-33): run-time error JS1010: Expected identifier: (
(541,6-7): run-time error JS1195: Expected expression: ,
(542,26-27): run-time error JS1010: Expected identifier: (
(544,6-7): run-time error JS1195: Expected expression: ,
(548,28-29): run-time error JS1010: Expected identifier: (
(553,6-7): run-time error JS1195: Expected expression: ,
(554,31-32): run-time error JS1010: Expected identifier: (
(559,6-7): run-time error JS1195: Expected expression: ,
(563,43-44): run-time error JS1010: Expected identifier: (
(574,6-7): run-time error JS1195: Expected expression: ,
(575,36): run-time error JS1004: Expected ';'
(600,6-7): run-time error JS1195: Expected expression: ,
(605,22-23): run-time error JS1010: Expected identifier: (
(610,6-7): run-time error JS1195: Expected expression: ,
(614,49-50): run-time error JS1010: Expected identifier: (
(625,6-7): run-time error JS1195: Expected expression: ,
(627,31-32): run-time error JS1010: Expected identifier: (
(640,6-7): run-time error JS1195: Expected expression: ,
(642,31-32): run-time error JS1010: Expected identifier: (
(648,6-7): run-time error JS1195: Expected expression: ,
(650,37-38): run-time error JS1010: Expected identifier: (
(653,6-7): run-time error JS1195: Expected expression: ,
(655,45-46): run-time error JS1010: Expected identifier: (
(661,6-7): run-time error JS1195: Expected expression: ,
(663,41-42): run-time error JS1010: Expected identifier: (
(670,1-2): run-time error JS1002: Syntax error: }
(1402,3701-3702): run-time error JS1195: Expected expression: .
(1402,3746-3747): run-time error JS1003: Expected ':': ;
(434,5,466,18): run-time error JS1314: Implicit property name must be identifier: downloadFile(method, url, params, lang) {
        let filename = 'uname';
        return fetch(url, {
                method: method,
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: "same-origin",
                body: JSON.stringify(params)
            })
            .then(res => {
                var disposition = res.headers.get('Content-Disposition');
                if (disposition && disposition.indexOf('attachment') !== -1) {
                    const parts = disposition.split(';');
                    filename = decodeURI(parts[1].split('=')[1]).replaceAll("+", " ");

                    return res.blob();
                }
                
                return res.json();
            })
            .then(data => {
                if (data?.constructor?.name == 'Blob') {
                    var url = window.URL.createObjectURL(data);
                    var a = document.createElement('a');
                    a.href = url;
                    a.download = filename;
                    document.body.appendChild(a);
                    a.click();
                    a.remove();
                } else if (data?.Status == this.AJAX_ERROR) {
                    this.showDangerMessage(data.Message);
                }
(607,32-33): run-time error JS1013: Syntax error in regular expression: ,
(577,9,599,15): run-time error JS1018: 'return' statement outside of function: return fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: "same-origin",
            body: JSON.stringify(data)
        })
            .then(res => {
                const header = res.headers.get('Content-Disposition');
                const parts = header.split(';');
                filename = parts[1].split('=')[1];
                return res.blob();
            })
            .then(blob => {
                var url = window.URL.createObjectURL(blob);
                var a = document.createElement('a');
                a.href = url;
                a.download = filename;
                document.body.appendChild(a);
                a.click();
                a.remove();
            })
 */
// ---- tlherbal ----
var Spiner = {
    // Thực hiện mở loading
    showLoading: function (divAppendTo, level) {
        // Remove loading nếu tồn tại
        this.closeLoading();
        var over = '<div id="overlay">'
        if (typeof level !== 'undefined' && (level == 1 || level == "1")) {
            over += '<img id="loading-img" src="../../../Content/Images/loadingtc.gif">' +
            '</div>';
        } else {
            over += '<img id="loading-img" src="../../Content/Images/loadingtc.gif">' +
            '</div>';
        }

        $(over).appendTo('#' + divAppendTo);
    },
    // Thực hiện đóng loading
    closeLoading: function () {
        $('#overlay').remove();
    },
};

// ---- tlherbal ----
var Confirm = {
    // Hiển thị dialog confirm
    confirm: function (title, button1, button2, element, fn_success, fn_failure) {
        var btns = {};
        btns[button1] = function () {
            element.parents('li').hide();
            $(this).dialog("close");
            fn_success();
        };
        btns[button2] = function () {
            // Do nothing
            $(this).dialog("close");
            fn_failure();
        };
        $("<div style='margin-top:15px; font-style:italic;'>" + title + "</div>").dialog({
            autoOpen: true,
            title: 'Xác nhận',
            modal: true,
            width: 400,
            height: 210,
            resizable: false,
            buttons: btns
        });

        $(".ui-dialog-buttonset button .ui-button-text").addClass("tc-btn-1");
        $(".ui-dialog-buttonset").css("margin-right", "122px");
        $(".ui-dialog-buttonpane.ui-widget-content.ui-helper-clearfix").css("border", "none");
        $(".ui-dialog.ui-widget.ui-widget-content.ui-corner-all.ui-draggable").css("border", "none");
    },
};

var Dialog = {
    create: function (dialogID, dialogWidth, dialogTitle) {
        $('#' + dialogID).dialog({
            autoOpen: false,
            width: dialogWidth,
            resizable: false,
            title: dialogTitle,
            modal: true,
            show: {
                effect: "blind",
            }
        });
    },
    displayDataAndOpen: function (dialogID, data) {
        $('#' + dialogID).html(data);
        this.open('frm-nhap-bo-sung');
    },
    open: function (dialogID) {
        $('#' + dialogID).dialog("open");
    },
    close: function (dialogID) {
        $('#' + dialogID).dialog("close");
    }
};
// ---- /tlherbal ----



// AuNH
function elementInViewport2(el) {
    var top = el.offsetTop;
    var left = el.offsetLeft;
    var width = el.offsetWidth;
    var height = el.offsetHeight;

    while (el.offsetParent) {
        el = el.offsetParent;
        top += el.offsetTop;
        left += el.offsetLeft;
    }

    return (
      top < (window.pageYOffset + window.innerHeight) &&
      left < (window.pageXOffset + window.innerWidth) &&
      (top + height) > window.pageYOffset &&
      (left + width) > window.pageXOffset
    );
}



function elementInViewport(el) {
    var top = el.offsetTop;
    var left = el.offsetLeft;
    var width = el.offsetWidth;
    var height = el.offsetHeight;

    while (el.offsetParent) {
        el = el.offsetParent;
        top += el.offsetTop;
        left += el.offsetLeft;
    }

    return (
      top >= window.pageYOffset &&
      left >= window.pageXOffset &&
      (top + height) <= (window.pageYOffset + window.innerHeight) &&
      (left + width) <= (window.pageXOffset + window.innerWidth)
    );
}

// ---- tlherbal ----
// Make an element stay inside window with center
function centerElementInWindow(elementID) {
    $('#' + elementID).parent().addClass('centered');
    return true;
}
// ---- /tlherbal ----

(function ($) {
    $.fn.fixMe = function () {
        return this.each(function () {
            var $this = $(this),
               $t_fixed;
            function init() {
                $this.wrap('<div class="container" />');
                $t_fixed = $this.clone();
                $t_fixed.find("tbody").remove().end().addClass("fixed").insertBefore($this);
                resizeFixed();
            }
            function resizeFixed() {

            }
            function scrollFixed() {
                var offset = $(this).scrollTop(),
                tableOffsetTop = $this.offset().top,
                tableOffsetBottom = tableOffsetTop + $this.height() - $this.find("thead").height();
                if (offset < tableOffsetTop || offset > tableOffsetBottom)
                    $t_fixed.hide();
                else if (offset >= tableOffsetTop && offset <= tableOffsetBottom && $t_fixed.is(":hidden"))
                    $t_fixed.show();
            }
            $(window).resize(resizeFixed);
            $(window).scroll(scrollFixed);
            init();
        });
    };
})(jQuery);


var CommonJS = {
    AJAX_ERROR: 0,

    AJAX_SUCCESS: 1,

    DEFAULT_START_PAGE: 1,

    HEADER_HIGHT: 271,

    //hainm22 add encode decode
    htmlEncode: function (value) {
        //create a in-memory div, set it's inner text(which jQuery automatically encodes)
        //then grab the encoded contents back out.  The div never exists on the page.
        return $('<div/>').text(value).html();
    },

    htmlDecode: function (value) {
        return $('<div/>').html(value).text();
    },
    // Nhantd
    // Tim doi tuong trong mang theo property
    lookupInArray: function (array, prop, value) {
        for (var i = 0, len = array.length; i < len; i++) {
            var item = array[i];
            if (item && item.hasOwnProperty(prop)) {
                var property = item[prop];
                if ($.isNumeric(property))
                    property = property.toString();
                if (property.toLowerCase() === value.toLowerCase()) {
                    return item;
                }
            }
        }
    },
    // ---- tlherbal ----
    // Chuyển tiếng việt có dấu thành tiếng anh
    convertVietnamese: function (str) {
        str = str.toLowerCase();
        str = str.replace(/à|á|ạ|ả|ã|â|ầ|ấ|ậ|ẩ|ẫ|ă|ằ|ắ|ặ|ẳ|ẵ/g, "a");
        str = str.replace(/è|é|ẹ|ẻ|ẽ|ê|ề|ế|ệ|ể|ễ/g, "e");
        str = str.replace(/ì|í|ị|ỉ|ĩ/g, "i");
        str = str.replace(/ò|ó|ọ|ỏ|õ|ô|ồ|ố|ộ|ổ|ỗ|ơ|ờ|ớ|ợ|ở|ỡ/g, "o");
        str = str.replace(/ù|ú|ụ|ủ|ũ|ư|ừ|ứ|ự|ử|ữ/g, "u");
        str = str.replace(/ỳ|ý|ỵ|ỷ|ỹ/g, "y");
        str = str.replace(/đ/g, "d");
        //str = str.replace(/!|@|%|\^|\*|\(|\)|\+|\=|\<|\>|\?|\/|,|\.|\:|\;|\'| |\"|\&|\#|\[|\]|~|$|_/g, "-");
        //str = str.replace(/-+-/g, "-");
        //str = str.replace(/^\-+|\-+$/g, "");
        return str;
    },
    // ---- /tlherbal ----
    baseAlert: function (msg) {
        if (msg.Message != null) {
            alert(msg.Message);
            return msg;
        }
        json = {};
        if (typeof msg !== "string" || msg.trim()[0] !== "{") {
            // Not JSON, treat as plain text
            return bootbox.alert(escapeHtml(msg));
        }

        try {
            json = JSON.parse(msg);
        } catch (e) {
            return bootbox.alert("Invalid response.");
        }

        alert(json.Message);
        return json;
    },

    alert: function (msg) {
        $("#_GlobalMessage").attr("class", "");
        $("#_GlobalMessage").addClass("msg-type-SUCCESS");
        if (msg.Message != null) {
            if (msg.Message == "The custom error module does not recognize this error.") {
                $("#_GlobalMessage").addClass("msg-type-ERROR");
            }
            $("#_GlobalMessage").html(msg.Message.replace("The custom error module does not recognize this error.", "Có lỗi xảy ra hoặc do bạn đã nhập các ký tự đặc biệt"));
            $("#_GlobalMessage").fadeIn();
            setTimeout('$("#_GlobalMessage").fadeOut();', 6000);
            return msg;
        }
        json = {};
        try {
            if (typeof msg === "string"
                && (/^\s*\{[\s\S]*\}\s*$/.test(msg) || /^\s*\[[\s\S]*\]\s*$/).test(msg)) {
                json = JSON.parse(msg);   
            }
        } catch (e) {
            if (msg == "The custom error module does not recognize this error.") {
                $("#_GlobalMessage").addClass("msg-type-ERROR");
            }
            $("#_GlobalMessage").html(msg.replace("The custom error module does not recognize this error.", "Có lỗi xảy ra hoặc do bạn đã nhập các ký tự đặc biệt"));
            $("#_GlobalMessage").fadeIn();
            setTimeout('$("#_GlobalMessage").fadeOut();', 6000);
            return;
        }
        if (json == null) {
            $("#_GlobalMessage").addClass("msg-type-ERROR");
            $("#_GlobalMessage").html("null");
        } else {
            $("#_GlobalMessage").addClass("msg-type-" + json.Type);
            $("#_GlobalMessage").html(json.Message.replace("The custom error module does not recognize this error.", "Có lỗi xảy ra hoặc do bạn đã nhập các ký tự đặc biệt"));
        }
        $("#_GlobalMessage").fadeIn();
        setTimeout('$("#_GlobalMessage").fadeOut();', 6000);
        return json;
    },

    back: function () {
        history.back(1);
    },

    checkDate: function checkDate(date) {
        var minYear = 1902;
        var errorMsg = "";
        // regular expression to match required date format 
        re = /^(\d{1,2})\/(\d{1,2})\/(\d{4})$/;

        if (date != '') {
            if (regs = date.match(re)) {
                if (regs[1] < 1 || regs[1] > 31) {
                    errorMsg = "Invalid value for day: " + regs[1];
                    return false;
                }
                else if (regs[2] < 1 || regs[2] > 12) {
                    errorMsg = "Invalid value for month: " + regs[2];
                    return false;
                }
                else if (regs[3] < minYear) {
                    errorMsg = "Invalid value for year: " + regs[3] + " - must be between " + minYear;
                    return false;
                }
            }
            else {
                errorMsg = "Invalid date format: " + date;
                return false;
            }
        }
        else {
            errorMsg = "Empty date not allowed!";
            return false;
        }

        return true;
    },
    isNumber: function (n) {
        return !isNaN(parseFloat(n)) && isFinite(n);
    },

    makeSizeMatchParent: function (elementID, parentID) {
        var pheight = $("#" + parentID).height();
        var pwidth = $("#" + parentID).width();
        $("#" + elementID).css({
            width: pwidth,
            height: pheight,
        });
    },

    FullScreen: function (elementID, isFull) {
        if (isFull) {
            $("#container").css({
                width: '1px',
                height: '1px',
                overflow: 'hidden'
            });

            $("#" + elementID).css({
                position: 'absolute',
                width: $(window).width() - 3,
                height: $(window).height() - 5,
                top: 0,
                left: 0,
                'z-index': 100,
                overflow: 'hidden',
                background: '#fff'
            });

            $(".form-enter-data").css({
                width: $(window).width() - 290,
                height: $(window).height() - 100,
                'max-height': $(window).height() - 100,
                overflow: 'auto',
            });
            $("#DivGroupTree .panelScroll").width($("#DivGroupTree").width());
        } else {
            $("#container").css({
                width: '',
                height: '',
                overflow: ''
            });
            $("#" + elementID).css({
                position: '',
                width: '',
                height: '',
                top: '',
                left: '',
                'z-index': '',
                overflow: '',
                background: ''
            });

            $(".form-enter-data").css({
                width: '',
                height: '',
                overflow: '',
                'max-height': $(window).height() - CommonJS.HEADER_HIGHT,
                background: ''
            });
            $("#DivGroupTree .panelScroll").width(194);
        }


    },

    EditableFullScreen: function (elementID, isFull) {
        $(window).resize(function () {
            CommonJS.FullScreen(elementID, isFull);
        });
        CommonJS.FullScreen(elementID, isFull);
        if (isFull) {
            $("#btnFullScreen").hide();
            $("#btnUnFull").show();
            $("#TemplateHide").show();
            $("#isFullScreen").val(1);
            $(".PanelContentScroll").css({
                'max-height': $(window).height() - 90,
            });
        } else {
            $("#btnFullScreen").show();
            $("#btnUnFull").hide();
            $("#isFullScreen").val(0);
            $(".PanelContentScroll").css({
                'max-height': $(window).height() - CommonJS.HEADER_HIGHT,
            });
        }
    },

    PrepareEditableFullScreen: function () {
        if ($("#isFullScreen").val() * 1 == 1) {
            $("#btnFullScreen").click();
        } else {
            $("#btnUnFull").click();
        }
    },
    ReloadPage: function () {
        window.location.href = window.location.href;
    },
    ClearAllCache: function () {
        sessionStorage.clear();
    },
    TreeClickFirstLeave: function () {
        var cnode = $($("div.hitarea")[0])
        var _count = 1;
        while (cnode.length > 0) {
            if (!cnode.hasClass("lastCollapsable-hitarea")) {
                cnode.click();
            }
            cnode = $(cnode.parent().find("ul li div.hitarea")[0])
            _count++;
            if (_count > 20) break;
        }

        $($("li span[id^=Template]")[0]).find("a.treeitem").click();
    },
    
    downloadFile(method, url, params, lang) {
        let filename = 'uname';
        return fetch(url, {
                method: method,
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: "same-origin",
                body: JSON.stringify(params)
            })
            .then(res => {
                var disposition = res.headers.get('Content-Disposition');
                if (disposition && disposition.indexOf('attachment') !== -1) {
                    const parts = disposition.split(';');
                    filename = decodeURI(parts[1].split('=')[1]).replaceAll("+", " ");

                    return res.blob();
                }
                
                return res.json();
            })
            .then(data => {
                if (data?.constructor?.name == 'Blob') {
                    var url = window.URL.createObjectURL(data);
                    var a = document.createElement('a');
                    a.href = url;
                    a.download = filename;
                    document.body.appendChild(a);
                    a.click();
                    a.remove();
                } else if (data?.Status == this.AJAX_ERROR) {
                    this.showDangerMessage(data.Message);
                } else {
                    throw data;
                }
            });
    },

    printPDF(url) {
        bootbox.dialog({
            title: GlobalResources.IN_PDF,
            message: '<div class="row mt5"><iframe src="' + url + '" height="499" width="100%" border="0"></iframe></div>',
            size: 'large',
        });
    },

    //
    // SHOW MESSAGE
    //
    showSuccessMessage: function (msg, time) {
        jQuery.gritter.add({
            //title: data,
            text: msg,
            class_name: 'growl-success',
            sticky: false,
            time: time != null ? time : 1000,
            position: 'bottom-right'
        });
    },

    showDangerMessage: function (msg, time) {
        jQuery.gritter.add({
            //title: data,
            text: msg,
            class_name: 'growl-danger',
            sticky: false,
            time: time != null ? time : 3000,
            position: 'bottom-right'
        });
    },

    showCommonDangerMessageWithLang: function (lang, time) {
        var msg = "Something went wrong, please try again later";
        if (lang == "vi") {
            msg = "Có lỗi xảy ra, vui lòng thử lại sau";
        } else if (lang == "lo") {
            msg = "ມີຂໍ້ຜິດພາດເກີດຂຶ້ນໃນຂະບານການປະມວນຜົນ, ກະລຸນາລອງອີກຄັ້ງ";
        }
        this.showDangerMessage(msg, time);
    },

    // ------------------------------------------------
    //       Kiem tra ngay hop le
    // ------------------------------------------------
    ValidateDate: function (DateString) {
        var regex = /^(((0[1-9]|[12]\d|3[01])\/(0[13578]|1[02])\/((1[6-9]|[2-9]\d)\d{2}))|((0[1-9]|[12]\d|30)\/(0[13456789]|1[012])\/((1[6-9]|[2-9]\d)\d{2}))|((0[1-9]|1\d|2[0-8])\/02\/((1[6-9]|[2-9]\d)\d{2}))|(29\/02\/((1[6-9]|[2-9]\d)(0[48]|[2468][048]|[13579][26])|((16|[2468][048]|[3579][26])00))))$/;
        //var regex = /^(((0[1-9]|[12]\d|3[01])\/(0[13578]|1[02])\/(19[5-9][0-9]|20[0-4][0-9]|2050))|((0[1-9]|[12]\d|30)\/(0[13456789]|1[012])\/(19[5-9][0-9]|20[0-4][0-9]|2050))|((0[1-9]|1\d|2[0-8])\/02\/(19[5-9][0-9]|20[0-4][0-9]|2050))|(29\/02\/(19([68][048]|[579][26])|20(0[48]|[2468][048]|[13579][26]))))$/;
        if (!(regex.test(DateString))) {
            return false;
        }
        else {
            return true;
        }
    },
    ValidateSameMonth: function (fromDate, toDate) {
        var partsFrom = fromDate.split("/");
        var dateFrom =  new Date(partsFrom[2], partsFrom[1] - 1, partsFrom[0]);
        var partsTo = toDate.split("/");
        var dateTo =  new Date(partsTo[2], partsTo[1] - 1, partsTo[0]);
        return (dateFrom.getMonth() === dateTo.getMonth());
    },
    ValidateSameYear: function (fromDate, toDate) {
        var partsFrom = fromDate.split("/");
        var dateFrom = new Date(partsFrom[2], partsFrom[1] - 1, partsFrom[0]);
        var partsTo = toDate.split("/");
        var dateTo = new Date(partsTo[2], partsTo[1] - 1, partsTo[0]);
        return (dateFrom.getYear() === dateTo.getYear());
    },
    capitalize: function (string) {
        return string.charAt(0).toUpperCase() + string.slice(1);
    },
    // --------------------------------------------
    //       Check so sanh Ngay
    // --------------------------------------------
    ckFormatDate: function (strDateTime) {
        var dateTimeParts = strDateTime.split(' ');
        var dateParts = dateTimeParts[0].split('/');

        return new Date(dateParts[2], parseInt(dateParts[1], 10) - 1, dateParts[0]);
    },
    scrollToElement: function (el, ms) {
        var speed = (ms) ? ms : 600;
        $('html,body').animate({
            scrollTop: $(el).offset().top
        }, speed);
    },
    // ------------------------------------------------
    //       Convert Ngay Json sang String (dd/MM/yyyy)
    // ------------------------------------------------
    convertDateTimeJsonToString: function (jsonDate) {
        if (jsonDate != null && jsonDate != "") {
            var dateString = jsonDate.substr(6);
            var currentTime = new Date(parseInt(dateString));
            var month = ("0" + (currentTime.getMonth() + 1)).slice(-2);
            var day = ("0" + currentTime.getDate()).slice(-2);
            var year = currentTime.getFullYear();
            var date = day + '/' + month + '/' + year;
            return date;
        }
        return "";
    },
    downloadFile(method, url, data) {
        let filename = 'uname';
        return fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: "same-origin",
            body: JSON.stringify(data)
        })
            .then(res => {
                const header = res.headers.get('Content-Disposition');
                const parts = header.split(';');
                filename = parts[1].split('=')[1];
                return res.blob();
            })
            .then(blob => {
                var url = window.URL.createObjectURL(blob);
                var a = document.createElement('a');
                a.href = url;
                a.download = filename;
                document.body.appendChild(a);
                a.click();
                a.remove();
            });
    },

    // ------------------------------------------------
    //       Endcode du lieu hien thi len control
    // ------------------------------------------------
    encode: function (e) {
        if (e == null || e == '') return "";
        return e.replace(/[^]/g, function (e) {
            return "&#" + e.charCodeAt(0) + ";"
        });
    },
    // ------------------------------------------------
    //       Convert tieng viet co dau sang khong dau
    // ------------------------------------------------
    ConvertSignedVietnameseToUnsigned: function (alias) {
        var str = alias;
        str = str.toLowerCase();
        str = str.replace(/à|á|ạ|ả|ã|â|ầ|ấ|ậ|ẩ|ẫ|ă|ằ|ắ|ặ|ẳ|ẵ/g, "a");
        str = str.replace(/è|é|ẹ|ẻ|ẽ|ê|ề|ế|ệ|ể|ễ/g, "e");
        str = str.replace(/ì|í|ị|ỉ|ĩ/g, "i");
        str = str.replace(/ò|ó|ọ|ỏ|õ|ô|ồ|ố|ộ|ổ|ỗ|ơ|ờ|ớ|ợ|ở|ỡ/g, "o");
        str = str.replace(/ù|ú|ụ|ủ|ũ|ư|ừ|ứ|ự|ử|ữ/g, "u");
        str = str.replace(/ỳ|ý|ỵ|ỷ|ỹ/g, "y");
        str = str.replace(/đ/g, "d");
        return str;
    },

    AddSeparatorsNF: function (nStr, inD, outD, sep) {
        nStr += '';
        var dpos = nStr.indexOf(inD);
        var nStrEnd = '';
        if (dpos != -1) {
            nStrEnd = outD + nStr.substring(dpos + 1, nStr.length);
            nStr = nStr.substring(0, dpos);
        }
        var rgx = /(\d+)(\d{3})/;
        while (rgx.test(nStr)) {
            nStr = nStr.replace(rgx, '$1' + sep + '$2');
        }
        return nStr + nStrEnd;
    },

    JsonToUrlString: function (params) {
        var urlParamenterString = Object.keys(params).map(function (key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(params[key]);
        }).join('&');

        return urlParamenterString;
    },

    getDateFromDateString: function (strDate) {
        var parts = strDate.split("/");
        return new Date(parts[2], parts[1] - 1, parts[0]);
    },

    getDateTimeFromDateTimeString: function (strDateTime) {
        var dateTimeParts = strDateTime.split(' ');
        var timeParts = dateTimeParts[1].split(':');
        var dateParts = dateTimeParts[0].split('/');

        return new Date(dateParts[2], parseInt(dateParts[1], 10) - 1, dateParts[0], timeParts[0], timeParts[1]);
    },

    getDateFromDateTimeString: function (strDateTime) {
        // Loại bỏ phần giờ, chỉ lấy ngày để so sánh
        var dateTimeParts = strDateTime.split(' ');
        var dateParts = dateTimeParts[0].split('/');

        return new Date(dateParts[2], parseInt(dateParts[1], 10) - 1, dateParts[0]);
    }
}

String.prototype.format = function () {
    var formatted = this;
    for (var arg in arguments) {
        formatted = formatted.replace("{" + arg + "}", arguments[arg]);
    }
    return formatted;
}

Number.prototype.formatMoney = function (c, d, t) {
    var n = this,
        c = isNaN(c = Math.abs(c)) ? 2 : c,
        d = d == undefined ? "." : d,
        t = t == undefined ? "," : t,
        s = n < 0 ? "-" : "",
        i = String(parseInt(n = Math.abs(Number(n) || 0).toFixed(c))),
        j = (j = i.length) > 3 ? j % 3 : 0;
    return s + (j ? i.substr(0, j) + t : "") + i.substr(j).replace(/(\d{3})(?=\d)/g, "$1" + t) + (c ? d + Math.abs(n - i).toFixed(c).slice(2) : "");
};;

jQuery(window).load(function() {
   
   "use strict";
   
   // Page Preloader
   jQuery('#preloader').delay(350).fadeOut(function(){
      jQuery('body').delay(350).css({'overflow':'visible'});
   });
});

jQuery(document).ready(function() {
   
   "use strict";
   
   // Toggle Left Menu
   jQuery('.leftpanel .nav-parent > a').live('click', function() {
      
      var parent = jQuery(this).parent();
      var sub = parent.find('> ul');
      
      // Dropdown works only when leftpanel is not collapsed
      if(!jQuery('body').hasClass('leftpanel-collapsed')) {
         if(sub.is(':visible')) {
            sub.slideUp(200, function(){
               parent.removeClass('nav-active');
               jQuery('.mainpanel').css({height: ''});
               adjustmainpanelheight();
            });
         } else {
            closeVisibleSubMenu();
            parent.addClass('nav-active');
            sub.slideDown(200, function(){
               adjustmainpanelheight();
            });
         }
      }
      return false;
   });
   
   function closeVisibleSubMenu() {
      jQuery('.leftpanel .nav-parent').each(function() {
         var t = jQuery(this);
         if(t.hasClass('nav-active')) {
            t.find('> ul').slideUp(200, function(){
               t.removeClass('nav-active');
            });
         }
      });
   }
   
   function adjustmainpanelheight() {
      // Adjust mainpanel height
      var docHeight = jQuery(document).height();
      if(docHeight > jQuery('.mainpanel').height())
         jQuery('.mainpanel').height(docHeight);
   }
   adjustmainpanelheight();
   
   
   // Tooltip
   jQuery('.tooltips').tooltip({ container: 'body'});
   
   // Popover
   jQuery('.popovers').popover();
   
   // Close Button in Panels
   jQuery('.panel .panel-close').click(function(){
      jQuery(this).closest('.panel').fadeOut(200);
      return false;
   });
   
   // Form Toggles
   jQuery('.toggle').toggles({on: true});
   
   jQuery('.toggle-chat1').toggles({on: false});
   
   var scColor1 = '#428BCA';
   if (jQuery.cookie('change-skin') && jQuery.cookie('change-skin') == 'bluenav') {
      scColor1 = '#fff';
   }
   
   
   // Sparkline
   jQuery('#sidebar-chart').sparkline([4,3,3,1,4,3,2,2,3,10,9,6], {
	 type: 'bar', 
	 height:'30px',
         barColor: scColor1
   });
   
   jQuery('#sidebar-chart2').sparkline([1,3,4,5,4,10,8,5,7,6,9,3], {
	  type: 'bar', 
	  height:'30px',
         barColor: '#D9534F'
   });
   
   jQuery('#sidebar-chart3').sparkline([5,9,3,8,4,10,8,5,7,6,9,3], {
	  type: 'bar', 
	  height:'30px',
         barColor: '#1CAF9A'
   });
   
   jQuery('#sidebar-chart4').sparkline([4,3,3,1,4,3,2,2,3,10,9,6], {
	  type: 'bar', 
	  height:'30px',
         barColor: scColor1
   });
   
   jQuery('#sidebar-chart5').sparkline([1,3,4,5,4,10,8,5,7,6,9,3], {
	  type: 'bar', 
	  height:'30px',
      barColor: '#F0AD4E'
   });
   
   
   // Minimize Button in Panels
   jQuery('.minimize').click(function(){
      var t = jQuery(this);
      var p = t.closest('.panel');
      if(!jQuery(this).hasClass('maximize')) {
         p.find('.panel-body, .panel-footer').slideUp(200);
         t.addClass('maximize');
         t.html('&plus;');
      } else {
         p.find('.panel-body, .panel-footer').slideDown(200);
         t.removeClass('maximize');
         t.html('&minus;');
      }
      return false;
   });
   
   
   // Add class everytime a mouse pointer hover over it
   jQuery('.nav-bracket > li').hover(function(){
      jQuery(this).addClass('nav-hover');
   }, function(){
      jQuery(this).removeClass('nav-hover');
   });
   
   
   // Menu Toggle
   jQuery('.menutoggle').click(function(){
      
      var body = jQuery('body');
      var bodypos = body.css('position');
      
      if(bodypos != 'relative') {
         
         if(!body.hasClass('leftpanel-collapsed')) {
            body.addClass('leftpanel-collapsed');
            jQuery('.nav-bracket ul').attr('style','');
            
            jQuery(this).addClass('menu-collapsed');
            
         } else {
            body.removeClass('leftpanel-collapsed chat-view');
            jQuery('.nav-bracket li.active ul').css({display: 'block'});
            
            jQuery(this).removeClass('menu-collapsed');
            
         }
      } else {
         
         if(body.hasClass('leftpanel-show'))
            body.removeClass('leftpanel-show');
         else
            body.addClass('leftpanel-show');
         
         adjustmainpanelheight();         
      }

   });
   
   // Chat View
   jQuery('#chatview').click(function(){
      
      var body = jQuery('body');
      var bodypos = body.css('position');
      
      if(bodypos != 'relative') {
         
         if(!body.hasClass('chat-view')) {
            body.addClass('leftpanel-collapsed chat-view');
            jQuery('.nav-bracket ul').attr('style','');
            
         } else {
            
            body.removeClass('chat-view');
            
            if(!jQuery('.menutoggle').hasClass('menu-collapsed')) {
               jQuery('body').removeClass('leftpanel-collapsed');
               jQuery('.nav-bracket li.active ul').css({display: 'block'});
            } else {
               
            }
         }
         
      } else {
         
         if(!body.hasClass('chat-relative-view')) {
            
            body.addClass('chat-relative-view');
            body.css({left: ''});
         
         } else {
            body.removeClass('chat-relative-view');   
         }
      }
      
   });
   
   reposition_topnav();
   reposition_searchform();
   
   jQuery(window).resize(function(){
      
      if(jQuery('body').css('position') == 'relative') {

         jQuery('body').removeClass('leftpanel-collapsed chat-view');
         
      } else {
         
         jQuery('body').removeClass('chat-relative-view');         
         jQuery('body').css({left: '', marginRight: ''});
      }
      
      
      reposition_searchform();
      reposition_topnav();
      
   });
   
   
   
   /* This function will reposition search form to the left panel when viewed
    * in screens smaller than 767px and will return to top when viewed higher
    * than 767px
    */ 
   function reposition_searchform() {
      if(jQuery('.searchform').css('position') == 'relative') {
         jQuery('.searchform').insertBefore('.leftpanelinner .userlogged');
      } else {
         jQuery('.searchform').insertBefore('.header-right');
      }
   }
   
   

   /* This function allows top navigation menu to move to left navigation menu
    * when viewed in screens lower than 1024px and will move it back when viewed
    * higher than 1024px
    */
   function reposition_topnav() {
      if(jQuery('.nav-horizontal').length > 0) {
         
         // top navigation move to left nav
         // .nav-horizontal will set position to relative when viewed in screen below 1024
         if(jQuery('.nav-horizontal').css('position') == 'relative') {
                                  
            if(jQuery('.leftpanel .nav-bracket').length == 2) {
               jQuery('.nav-horizontal').insertAfter('.nav-bracket:eq(1)');
            } else {
               // only add to bottom if .nav-horizontal is not yet in the left panel
               if(jQuery('.leftpanel .nav-horizontal').length == 0)
                  jQuery('.nav-horizontal').appendTo('.leftpanelinner');
            }
            
            jQuery('.nav-horizontal').css({display: 'block'})
                                  .addClass('nav-pills nav-stacked nav-bracket');
            
            jQuery('.nav-horizontal .children').removeClass('dropdown-menu');
            jQuery('.nav-horizontal > li').each(function() { 
               
               jQuery(this).removeClass('open');
               jQuery(this).find('a').removeAttr('class');
               jQuery(this).find('a').removeAttr('data-toggle');
               
            });
            
            if(jQuery('.nav-horizontal li:last-child').has('form')) {
               jQuery('.nav-horizontal li:last-child form').addClass('searchform').appendTo('.topnav');
               jQuery('.nav-horizontal li:last-child').hide();
            }
         
         } else {
            // move nav only when .nav-horizontal is currently from leftpanel
            // that is viewed from screen size above 1024
            if(jQuery('.leftpanel .nav-horizontal').length > 0) {
               
               jQuery('.nav-horizontal').removeClass('nav-pills nav-stacked nav-bracket')
                                        .appendTo('.topnav');
               jQuery('.nav-horizontal .children').addClass('dropdown-menu').removeAttr('style');
               jQuery('.nav-horizontal li:last-child').show();
               jQuery('.searchform').removeClass('searchform').appendTo('.nav-horizontal li:last-child .dropdown-menu');
               jQuery('.nav-horizontal > li > a').each(function() {
                  
                  jQuery(this).parent().removeClass('nav-active');
                  
                  if(jQuery(this).parent().find('.dropdown-menu').length > 0) {
                     jQuery(this).attr('class','dropdown-toggle');
                     jQuery(this).attr('data-toggle','dropdown');  
                  }
                  
               });              
            }
            
         }
         
      }
   }
   
   
   // Sticky Header
   if(jQuery.cookie('sticky-header'))
      jQuery('body').addClass('stickyheader');
      
   // Sticky Left Panel
   if(jQuery.cookie('sticky-leftpanel')) {
      jQuery('body').addClass('stickyheader');
      jQuery('.leftpanel').addClass('sticky-leftpanel');
   }
   
   // Left Panel Collapsed
   if(jQuery.cookie('leftpanel-collapsed')) {
      jQuery('body').addClass('leftpanel-collapsed');
      jQuery('.menutoggle').addClass('menu-collapsed');
   }
   
   // Changing Skin
   var c = jQuery.cookie('change-skin');
   var cssSkin = 'css/style.'+c+'.css';
   if (jQuery('body').css('direction') == 'rtl') {
      cssSkin = '../css/style.'+c+'.css';
      jQuery('html').addClass('rtl');
   }
   if(c) {
      jQuery('head').append('<link id="skinswitch" rel="stylesheet" href="'+cssSkin+'" />');
   }
   
   // Changing Font
   var fnt = jQuery.cookie('change-font');
   if(fnt) {
      jQuery('head').append('<link id="fontswitch" rel="stylesheet" href="css/font.'+fnt+'.css" />');
   }
   
   // Check if leftpanel is collapsed
   if(jQuery('body').hasClass('leftpanel-collapsed'))
      jQuery('.nav-bracket .children').css({display: ''});
      
     
   // Handles form inside of dropdown 
   jQuery('.dropdown-menu').find('form').click(function (e) {
      e.stopPropagation();
   });
   
   
   // This is not actually changing color of btn-primary
   // This is like you are changing it to use btn-orange instead of btn-primary
   // This is for demo purposes only
   var c = jQuery.cookie('change-skin');
   if (c && c == 'greyjoy') {
      $('.btn-primary').removeClass('btn-primary').addClass('btn-orange');
      $('.rdio-primary').addClass('rdio-default').removeClass('rdio-primary');
      $('.text-primary').removeClass('text-primary').addClass('text-orange');
   }
      

});;
var Constant = (function () {

    var TRANG_THAI_BAO_CAO = {
        DANG_SOAN_THAO: 0,
        DA_GUI: 1,
        DA_DUYET: 2,
        TU_CHOI: 3
    };

    var COSO_PHAN_LOAI = {
        CO_SO_TIEM_CHUNG_CONG: 0,
        CO_SO_TIEM_CHUNG_TU: 1,
        BENH_VIEN: 2,
        VIETTEL_ICT: 8,
        KINH_DOANH_VTT: 9,
        KHAC: 3,
        KHONG_RO: -1
    };

    var COSO_HINH_THUC = {
        TRUNG_UONG: 0,
        KHU_VUC: 1,
        TINH: 2,
        HUYEN: 3,
        XA: 4,
        CO_SO_DICH_VU: 5,
        BENH_VIEN: 7
    };

    var RESPONSE_STATUS = {
        SUCCESS: 1,
        ERROR: 0
    };

    var DON_VI_CAN_NANG = {
        GRAM: 1,
        KILOGRAM: 2
    };

    var DON_VI_CHIEU_CAO = {
        CENTIMET: 1,
        MET: 2
    };

    var GIOI_TINH = {
        NAM: 0,
        NU: 1,
        KHONG_RO: 2
    };

    var LOAI_GIO_HEN_TIEM = {
        BUOI_SANG: 1,
        BUOI_CHIEU: 2,
        CA_NGAY: 3,
        GIO: 4
    };

    var LOAI_TIN_NHAN = {
        KICH_HOAT: 1,
        MOI_TIEM: 2,
        NHAC_LICH: 3
    };

    var LOAI_GOI_PHAN_MEM = {
        THU_NGHIEM: 0,
        THU_PHI: 1,
        TIN_NHAN: 2
    };

    var HOTLINE = '19008068';
    var TEXT_REMINDER_TEMPLATE = 'Cháu [TEN_DOI_TUONG] đến tuổi tiêm phòng mũi [TEN_VACXIN] (theo phác đồ của Bộ Y tế). Gia đình lưu ý cho trẻ đi tiêm chủng đầy đủ. Chi tiết LH [DIEN_THOAI_LIEN_HE]';
    var TEXT_INVITATION_TEMPLATE = 'Kính mời ông/bà đưa cháu [TEN_DOI_TUONG] tới [DIA_DIEM_TIEM] để tiêm vắc xin phòng bệnh [TEN_BENH] vào lúc [THOI_GIAN_HEN] ngày [NGAY_HEN]. Chi tiết LH: [DIEN_THOAI_LIEN_HE]';

    var LOC_TRUNG = {
        TIEU_CHI_LOC_TRUNG_LABEL: [
            "Họ và tên",
            "Giới tính",
            "Ngày sinh",
            "Số điện thoại",
            "CMT/CCCD",
            "Dân tộc",
            "Hộ khẩu: Tỉnh/Thành phố",
            "Hộ khẩu: TTYT khu vực",
            "Hộ khẩu: Xã/Phường",
            "Hộ khẩu: Thôn/Ấp",
            "Hộ khẩu: Địa chỉ",
            "Tạm trú: Tỉnh/Thành phố",
            "Tạm trú: TTYT khu vực",
            "Tạm trú: Xã/Phường",
            "Tạm trú: Thôn/Ấp",
            "Tạm trú: Địa chỉ",
            "Số mũi UVSS",
            "Bảo vệ UVSS",
            "Họ và tên Mẹ",
            "CMT/CCCD Mẹ",
            "Số điện thoại Mẹ",
            "Năm sinh Mẹ",
            "Họ và tên Bố",
            "CMT/CCCD Bố",
            "Điện thoại Bố",
            "Năm sinh Bố",
            "Tên người bảo hộ",
            "CMT/CCCD người bảo hộ",
            "Điện thoại người bảo hộ",
            "Năm sinh người bảo hộ"
        ],
        TIEU_CHI_LOC_TRUNG_VALUE: [
            1, // Họ và tên
            2, // Giới tính
            3, // Ngày sinh
            4, // Số điện thoại
            5, // CMT/CCCD
            6, // Dân tộc
            7, // Hộ khẩu: Tỉnh/Thành phố
            8, // Hộ khẩu: Quận/Huyện
            9, // Hộ khẩu: Xã/Phường
            10, // Hộ khẩu: Thôn/Ấp
            11, // Hộ khẩu: Địa chỉ
            12, // Tạm trú: Tỉnh/Thành phố
            13, // Tạm trú: Quận/Huyện
            14, // Tạm trú: Xã/Phường
            15, // Tạm trú: Thôn/Ấp
            16, // Tạm trú: Địa chỉ
            17, // Số mũi UVSS
            18, // Bảo vệ UVSS
            19, // Họ và tên Mẹ
            20, // CMT/CCCD Mẹ
            21, // Số điện thoại Mẹ
            22, // Năm sinh Mẹ
            23, // Họ và tên Bố
            24, // CMT/CCCD Bố
            25, // Điện thoại Bố
            26, // Năm sinh Bố
            27, // Tên người bảo hộ
            28, // CMT/CCCD người bảo hộ
            29, // Điện thoại người bảo hộ
            30 // Năm sinh người bảo hộ
        ],

        TRUONG_HIEN_THI_LABEL: [
            "Họ và tên",
            "Giới tính",
            "Ngày sinh",
            "Số điện thoại",
            "CMT/CCCD",
            "Dân tộc",
            "Hộ khẩu: Tỉnh/Thành phố",
            "Hộ khẩu: TTYT khu vực",
            "Hộ khẩu: Xã/Phường",
            "Hộ khẩu: Thôn/Ấp",
            "Hộ khẩu: Địa chỉ",
            "Tạm trú: Tỉnh/Thành phố",
            "Tạm trú: TTYT khu vực",
            "Tạm trú: Xã/Phường",
            "Tạm trú: Thôn/Ấp",
            "Tạm trú: Địa chỉ",
            "Số mũi UVSS",
            "Bảo vệ UVSS",
            "Họ và tên Mẹ",
            "CMT/CCCD Mẹ",
            "Số điện thoại Mẹ",
            "Năm sinh Mẹ",
            "Họ và tên Bố",
            "CMT/CCCD Bố",
            "Điện thoại Bố",
            "Năm sinh Bố",
            "Tên người bảo hộ",
            "CMT/CCCD người bảo hộ",
            "Điện thoại người bảo hộ",
            "Năm sinh người bảo hộ",
            "Mã đối tượng",
            "Cơ sở tạo",
            "Ngày tạo"
        ],
        TRUONG_HIEN_THI_VALUE: [
            1,
            2,
            3,
            4,
            5,
            6,
            7,
            8,
            9,
            10,
            11,
            12,
            13,
            14,
            15,
            16,
            17,
            18,
            19,
            20,
            21,
            22,
            23,
            24,
            25,
            26,
            27,
            28,
            29,
            30,
            31,
            32,
            33
        ],
    };

    return {
        TRANG_THAI_BAO_CAO: TRANG_THAI_BAO_CAO,
        COSO_PHAN_LOAI: COSO_PHAN_LOAI,
        COSO_HINH_THUC: COSO_HINH_THUC,
        RESPONSE_STATUS: RESPONSE_STATUS,
        DON_VI_CAN_NANG: DON_VI_CAN_NANG,
        DON_VI_CHIEU_CAO: DON_VI_CHIEU_CAO,
        GIOI_TINH: GIOI_TINH,
        LOAI_GIO_HEN_TIEM: LOAI_GIO_HEN_TIEM,
        HOTLINE: HOTLINE,
        TEXT_REMINDER_TEMPLATE: TEXT_REMINDER_TEMPLATE,
        TEXT_INVITATION_TEMPLATE: TEXT_INVITATION_TEMPLATE,
        LOAI_TIN_NHAN: LOAI_TIN_NHAN,
        LOAI_GOI_PHAN_MEM: LOAI_GOI_PHAN_MEM,
        LOC_TRUNG: LOC_TRUNG
    }
})();;
function Common() { };
Common.UI = function () { };
Common.UI.BlockElement = function (element) {
    $(element).block({
        message: '<div class="cssload-loader"></div>',
        overlayCSS: {
            backgroundColor: '#fff',
            opacity: 0.8,
            cursor: 'wait'
        },
        css: {
            border: 0,
            padding: 0,
            backgroundColor: 'transparent'
        }
    });
}

Common.UI.UnBlockElement = function (element) {
    $(element).unblock();
}

Common.UI.Notify = function () { };

Common.UI.Notify.Info = function (message, title) {
    jQuery.gritter.add({
        text: message,
        title: title,
        sticky: false,
        timeout: 2000,
        class_name: "growl-info"
    });
};
Common.UI.Notify.Danger = function (message, title) {
    jQuery.gritter.add({
        text: message,
        title: title != null && title != 'undefined' ? title : '',
        sticky: false,
        timeout: 2000,
        class_name: "growl-danger"
    });
};
Common.UI.Notify.Success = function (message, title) {
    jQuery.gritter.add({
        text: message,
        title: title != null && title != 'undefined' ? title : '',
        sticky: false,
        timeout: 2000,
        class_name: "growl-success"
    });
};
Common.UI.Notify.Warning = function (message, title) {
    jQuery.gritter.add({
        text: message,
        title: title != null && title != 'undefined' ? title : '',
        sticky: false,
        timeout: 2000,
        class_name: "growl-warning"
    });
};

;
function CommonValidation() {
}

CommonValidation.ValidateDate = function (dateString) {
    var regex = /^(((0[1-9]|[12]\d|3[01])\/(0[13578]|1[02])\/((1[6-9]|[2-9]\d)\d{2}))|((0[1-9]|[12]\d|30)\/(0[13456789]|1[012])\/((1[6-9]|[2-9]\d)\d{2}))|((0[1-9]|1\d|2[0-8])\/02\/((1[6-9]|[2-9]\d)\d{2}))|(29\/02\/((1[6-9]|[2-9]\d)(0[48]|[2468][048]|[13579][26])|((16|[2468][048]|[3579][26])00))))$/;
    if (!(regex.test(dateString))) {
        return false;
    }
    else {
        return true;
    }
}

CommonValidation.ParseDateStringVi = function (dateString) {

    // regular expression to match required date format
    var re = /^(\d{1,2})\/(\d{1,2})\/(\d{4})$/;
    var regs = dateString.match(re);
    if (regs) {
        // day value between 1 and 31
        if (regs[1] < 1 || regs[1] > 31) {
            return false;
        }
        // month value between 1 and 12
        if (regs[2] < 1 || regs[2] > 12) {
            return false;
        }

    } else {
        return false;
    }

    return new Date(regs[3], regs[2] - 1, regs[1]);
}

// kiem tra so dien thoai
CommonValidation.ValidatePhone = function (phoneString) {
    var pattern = /^\d+$/;
    return pattern.test(phoneString);
}

// Kiem tra so chung minh thu hoac the can cuoc cong dan
CommonValidation.ValidateIdentificationCardNumber = function (CardIdNumber) {
    if (!CardIdNumber) return false;
    CardIdNumber = CardIdNumber.trim();
    var pattern = /^(?:[0-9]{9}|[0-9]{12})$/;
    return pattern.test(CardIdNumber);
};

// kiem tra number
CommonValidation.ValidateNumber = function (Number) {
    var pattern = /^\d+$/;
    return pattern.test(Number);
};
var ImCache=function(){var n={},t={COSO_HINH_THUC:{TRUNG_UONG:0,KHU_VUC:1,TINH:2,HUYEN:3,XA:4,CO_SO_DICH_VU:5,BENH_VIEN:7},COSO_PHAN_LOAI:{CO_SO_TIEM_CHUNG_CONG:0,CO_SO_TIEM_CHUNG_TU:1,BENH_VIEN:2,VIETTEL_ICT:8,KINH_DOANH_VTT:9,KHAC:3,KHONG_RO:-1}};return n.GetDsVacxinFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("VACXIN_CACHED")==null?$.ajax({url:"/Vacxin/DsVacxin","async":!0,type:"GET",success:function(t){sessionStorage.getItem("VACXIN_CACHED")==null&&(sessionStorage.setItem("VACXIN_DATA",JSON.stringify(t)),sessionStorage.setItem("VACXIN_CACHED","true"),n(t))}}):n(JSON.parse(sessionStorage.getItem("VACXIN_DATA"))):$.ajax({url:"/Vacxin/DsVacxin","async":!0,type:"GET",success:function(t){n(t)}})},n.GetDsVacxinKhongCovidFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("VACXINKHONGCOVID_CACHED")==null?$.ajax({url:"/Vacxin/DsVacxinKhongCovid","async":!0,type:"GET",success:function(t){sessionStorage.getItem("VACXINKHONGCOVID_CACHED")==null&&(sessionStorage.setItem("VACXINKHONGCOVID_DATA",JSON.stringify(t)),sessionStorage.setItem("VACXINKHONGCOVID_CACHED","true"),n(t))}}):n(JSON.parse(sessionStorage.getItem("VACXINKHONGCOVID_DATA"))):$.ajax({url:"/Vacxin/DsVacxinKhongCovid","async":!0,type:"GET",success:function(t){n(t)}})},n.GetDsVacxinCovidFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("VACXINCOVID_CACHED")==null?$.ajax({url:"/Vacxin/DsVacxinCovid","async":!0,type:"GET",success:function(t){sessionStorage.getItem("VACXINCOVID_CACHED")==null&&(sessionStorage.setItem("VACXINCOVID_DATA",JSON.stringify(t)),sessionStorage.setItem("VACXINCOVID_CACHED","true"),n(t))}}):n(JSON.parse(sessionStorage.getItem("VACXINCOVID_DATA"))):$.ajax({url:"/Vacxin/DsVacxinCovid","async":!0,type:"GET",success:function(t){n(t)}})},n.GetVacxinByIdFromCache=function(t,i){n.GetDsVacxinFromCache(function(n){var r=n.filter(function(n){if(n.VACXIN_ID==t)return n});r!=null&&r.length>0?i(r[0]):i(null)})},n.GetDsQuocGiaFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("QUOC_GIA_CACHED")==null?$.ajax({url:"/DanToc/DsQuocGia","async":!0,type:"GET",success:function(t){sessionStorage.setItem("QUOC_GIA_DATA",JSON.stringify(t));sessionStorage.setItem("QUOC_GIA_CACHED","true");n(t)}}):n(JSON.parse(sessionStorage.getItem("QUOC_GIA_DATA"))):$.ajax({url:"/DanToc/DsQuocGia","async":!0,type:"GET",success:function(t){n(t)}})},n.GetDsDanTocFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("DAN_TOC_CACHED")==null?$.ajax({url:"/DanToc/DsDanToc","async":!0,type:"GET",success:function(t){sessionStorage.setItem("DAN_TOC_DATA",JSON.stringify(t));sessionStorage.setItem("DAN_TOC_CACHED","true");n(t)}}):n(JSON.parse(sessionStorage.getItem("DAN_TOC_DATA"))):$.ajax({url:"/DanToc/DsDanToc","async":!0,type:"GET",success:function(t){n(t)}})},n.GetDanTocByIdFromCache=function(t,i){n.GetDsDanTocFromCache(function(n){var r=n.filter(function(n){if(n.DAN_TOC_ID==t)return n});r!=null&&r.length>0?i(r[0]):i(null)})},n.GetDsTonGiaoFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("TON_GIAO_CACHED")==null?$.ajax({url:"/TonGiao/DsTonGiao","async":!0,type:"GET",success:function(t){sessionStorage.setItem("TON_GIAO_DATA",JSON.stringify(t));sessionStorage.setItem("TON_GIAO_CACHED","true");n(t)},error:function(t,i,r){console.error("Lỗi khi tải danh mục Tôn giáo:",r);n([])}}):n(JSON.parse(sessionStorage.getItem("TON_GIAO_DATA"))):$.ajax({url:"/TonGiao/DsTonGiao","async":!0,type:"GET",success:function(t){n(t)},error:function(t,i,r){console.error("Lỗi khi tải danh mục Tôn giáo:",r);n([])}})},n.GetTonGiaoByIdFromCache=function(t,i){n.GetDsTonGiaoFromCache(function(n){var r=n?.filter(function(n){return n.TON_GIAO_ID==t});r!=null&&r.length>0?i(r[0]):i(null)})},n.GetDsVungMienFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("VUNG_MIEN_CACHED")==null?$.ajax({url:"/DonViHanhChinh/DsVungMien","async":!0,type:"GET",success:function(t){sessionStorage.setItem("VUNG_MIEN_DATA",JSON.stringify(t));sessionStorage.setItem("VUNG_MIEN_CACHED","true");n(t.sort(function(n,t){return n.TEN_VUNG_MIEN>t.TEN_VUNG_MIEN}))}}):n(JSON.parse(sessionStorage.getItem("VUNG_MIEN_DATA")).sort(function(n,t){return n.TEN_VUNG_MIEN>t.TEN_VUNG_MIEN})):$.ajax({url:"/DonViHanhChinh/DsVungMien","async":!0,type:"GET",success:function(t){n(t.sort(function(n,t){return n.TEN_VUNG_MIEN>t.TEN_VUNG_MIEN}))}})},n.GetDsDonViHanhChinhFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("DMDONVIHANHCHINH_CACHED")==null?$.ajax({url:"/DonViHanhChinh/DsDonViHanhChinh","async":!0,type:"GET",success:function(t){sessionStorage.setItem("DMDONVIHANHCHINH_DATA",JSON.stringify(t));sessionStorage.setItem("DMDONVIHANHCHINH_CACHED","true");n(t.sort(function(n,t){return n.TENDAYDU>t.TENDAYDU}))}}):n(JSON.parse(sessionStorage.getItem("DMDONVIHANHCHINH_DATA")).sort(function(n,t){return n.TENDAYDU>t.TENDAYDU})):$.ajax({url:"/DonViHanhChinh/DsDonViHanhChinh","async":!0,type:"GET",success:function(t){n(t.sort(function(n,t){return n.TENDAYDU>t.TENDAYDU}))}})},n.GetDsDonViHanhChinhByVungMienFromCache=function(n,t){if(typeof Storage!="undefined")if(sessionStorage.getItem("DMDONVIHANHCHINH_CACHED")==null||sessionStorage.getItem("DMDONVIHANHCHINH_CACHED")==undefined)$.ajax({url:"/DonViHanhChinh/DsDonViHanhChinhByVungMienId","async":!0,data:{vungMienId:-1},type:"GET",success:function(i){sessionStorage.setItem("DMDONVIHANHCHINH_DATA",JSON.stringify(i));sessionStorage.setItem("DMDONVIHANHCHINH_CACHED","true");n==-1?t(i.sort(function(n,t){return n.TENDAYDU>t.TENDAYDU})):t(i.filter(function(t){if(t.NIISID==n)return t}).sort(function(n,t){return n.TENDAYDU>t.TENDAYDU}))}});else{var i=JSON.parse(sessionStorage.getItem("DMDONVIHANHCHINH_DATA"));n==-1?t(i.sort(function(n,t){return n.TENDAYDU>t.TENDAYDU})):t(i.filter(function(t){if(t.NIISID==n)return t}).sort(function(n,t){return n.TENDAYDU>t.TENDAYDU}))}else $.ajax({url:"/DonViHanhChinh/DsDonViHanhChinhByVungMienId","async":!0,type:"GET",data:{vungMienId:n},success:function(n){t(n.sort(function(n,t){return n.TENDAYDU>t.TENDAYDU}))}})},n.GetDsTinhFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("TINH_CACHED")==null?$.ajax({url:"/DonViHanhChinh/DsTinh","async":!0,type:"GET",success:function(t){sessionStorage.setItem("TINH_DATA",JSON.stringify(t));sessionStorage.setItem("TINH_CACHED","true");n(t.sort(function(n,t){return n.TENTINH>t.TENTINH}))}}):n(JSON.parse(sessionStorage.getItem("TINH_DATA")).sort(function(n,t){return n.TENTINH>t.TENTINH})):$.ajax({url:"/DonViHanhChinh/DsTinh","async":!0,type:"GET",success:function(t){n(t.sort(function(n,t){return n.TENTINH>t.TENTINH}))}})},n.GetDsTinhByVungMienFromCache=function(n,t){if(typeof Storage!="undefined")if(sessionStorage.getItem("TINH_CACHED")==null)$.ajax({url:"/DonViHanhChinh/DsTinhByVungMienId","async":!0,data:{vungMienId:-1},type:"GET",success:function(i){sessionStorage.setItem("TINH_DATA",JSON.stringify(i));sessionStorage.setItem("TINH_CACHED","true");n==-1?t(i.sort(function(n,t){return n.TENTINH>t.TENTINH})):t(i.filter(function(t){if(t.VUNG_MIEN==n)return t}).sort(function(n,t){return n.TENTINH>t.TENTINH}))}});else{var i=JSON.parse(sessionStorage.getItem("TINH_DATA"));n==-1?t(i.sort(function(n,t){return n.TENTINH>t.TENTINH})):t(i.filter(function(t){if(t.VUNG_MIEN==n)return t}).sort(function(n,t){return n.TENTINH>t.TENTINH}))}else $.ajax({url:"/DonViHanhChinh/DsTinhByVungMienId","async":!0,type:"GET",data:{vungMienId:n},success:function(n){t(n.sort(function(n,t){return n.TENTINH>t.TENTINH}))}})},n.GetTinhByIdFromCache=function(t,i){n.GetDsTinhByVungMienFromCache(-1,function(n){var r=n.filter(function(n){if(n.TINH_ID==t)return n});r!=null&&r.length>0?i(r[0]):i(null)})},n.GetDsHuyenFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("HUYEN_CACHED")==null?$.ajax({url:"/DonViHanhChinh/DsHuyen","async":!0,type:"GET",data:{tinhId:-1},success:function(t){sessionStorage.setItem("HUYEN_DATA",JSON.stringify(t));sessionStorage.setItem("HUYEN_CACHED","true");n(t.sort(function(n,t){return n.TENHUYEN>t.TENHUYEN}))}}):n(JSON.parse(sessionStorage.getItem("HUYEN_DATA")).sort(function(n,t){return n.TENHUYEN>t.TENHUYEN})):$.ajax({url:"/DonViHanhChinh/DsHuyen","async":!0,type:"GET",data:{tinhId:-1},success:function(t){n(t.sort(function(n,t){return n.TENHUYEN>t.TENHUYEN}))}})},n.GetDsHuyenByTinhFromCache=function(n,t){if(typeof Storage!="undefined")if(sessionStorage.getItem("HUYEN_CACHED")==null)$.ajax({url:"/DonViHanhChinh/DsHuyen","async":!0,data:{tinhId:-1},type:"GET",success:function(i){sessionStorage.setItem("HUYEN_DATA",JSON.stringify(i));sessionStorage.setItem("HUYEN_CACHED","true");n==-1?t(i.sort(function(n,t){return n.TENHUYEN>t.TENHUYEN})):t(i.filter(function(t){if(t.TINH_ID==n)return t}).sort(function(n,t){return n.TENHUYEN>t.TENHUYEN}))}});else{var i=JSON.parse(sessionStorage.getItem("HUYEN_DATA"));n==-1?t(i.sort(function(n,t){return n.TENHUYEN>t.TENHUYEN})):t(i.filter(function(t){if(t.TINH_ID==n)return t}).sort(function(n,t){return n.TENHUYEN>t.TENHUYEN}))}else $.ajax({url:"/DonViHanhChinh/DsHuyen","async":!0,type:"GET",data:{tinhId:n},success:function(n){t(n.sort(function(n,t){return n.TENHUYEN>t.TENHUYEN}))}})},n.GetHuyenByIdFromCache=function(t,i){n.GetDsHuyenByTinhFromCache(-1,function(n){var r=n.filter(function(n){if(n.HUYEN_ID==t)return n});r!=null&&r.length>0?i(r[0]):i(null)})},n.GetDsXaFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("XA_CACHED")==null?$.ajax({url:"/DonViHanhChinh/DsXa","async":!0,type:"GET",data:{huyenId:-1},success:function(t){sessionStorage.setItem("XA_DATA",JSON.stringify(t));sessionStorage.setItem("XA_CACHED","true");n(t.sort(function(n,t){return n.TEN_XA>t.TEN_XA}))}}):n(JSON.parse(sessionStorage.getItem("XA_DATA")).sort(function(n,t){return n.TEN_XA>t.TEN_XA})):$.ajax({url:"/DonViHanhChinh/DsXa","async":!0,type:"GET",data:{huyenId:-1},success:function(t){n(t.sort(function(n,t){return n.TEN_XA>t.TEN_XA}))}})},n.GetDsXaByHuyenFromCache=function(n,t){if(typeof Storage!="undefined")if(sessionStorage.getItem("XA_CACHED")==null)$.ajax({url:"/DonViHanhChinh/DsXa","async":!0,data:{huyenId:-1},type:"GET",success:function(i){sessionStorage.setItem("XA_DATA",JSON.stringify(i));sessionStorage.setItem("XA_CACHED","true");n==-1?t(i.sort(function(n,t){return n.TEN_XA>t.TEN_XA})):t(i.filter(function(t){if(t.HUYEN_ID==n)return t}).sort(function(n,t){return n.TEN_XA>t.TEN_XA}))}});else{var i=JSON.parse(sessionStorage.getItem("XA_DATA"));n==-1?t(i.sort(function(n,t){return n.TEN_XA>t.TEN_XA})):t(i.filter(function(t){if(t.HUYEN_ID==n)return t}).sort(function(n,t){return n.TEN_XA>t.TEN_XA}))}else $.ajax({url:"/DonViHanhChinh/DsXa","async":!0,type:"GET",data:{huyenId:n},success:function(n){t(n.sort(function(n,t){return n.TEN_XA>t.TEN_XA}))}})},n.GetDsXa715FromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("XA_CACHED")==null?$.ajax({url:"/DonViHanhChinh/DsXa715","async":!0,type:"GET",data:{tinhId:-1},success:function(t){sessionStorage.setItem("XA_DATA",JSON.stringify(t));sessionStorage.setItem("XA_CACHED","true");n(t.sort(function(n,t){return n.TEN_XA>t.TEN_XA}))}}):n(JSON.parse(sessionStorage.getItem("XA_DATA")).sort(function(n,t){return n.TEN_XA>t.TEN_XA})):$.ajax({url:"/DonViHanhChinh/DsXa715","async":!0,type:"GET",data:{tinhId:-1},success:function(t){n(t.sort(function(n,t){return n.TEN_XA>t.TEN_XA}))}})},n.GetDsXaByTinhFromCache=function(n,t){if(typeof Storage!="undefined")if(sessionStorage.getItem("XA_CACHED")==null)$.ajax({url:"/DonViHanhChinh/DsXa715","async":!0,data:{tinhId:-1},type:"GET",success:function(i){sessionStorage.setItem("XA_DATA",JSON.stringify(i));sessionStorage.setItem("XA_CACHED","true");n==-1?t(i.sort(function(n,t){return n.TEN_XA>t.TEN_XA})):t(i.filter(function(t){if(t.TINH_ID==n)return t}).sort(function(n,t){return n.TEN_XA>t.TEN_XA}))}});else{var i=JSON.parse(sessionStorage.getItem("XA_DATA"));n==-1?t(i.sort(function(n,t){return n.TEN_XA>t.TEN_XA})):t(i.filter(function(t){if(t.TINH_ID==n)return t}).sort(function(n,t){return n.TEN_XA>t.TEN_XA}))}else $.ajax({url:"/DonViHanhChinh/DsXa715","async":!0,type:"GET",data:{tinhId:n},success:function(n){t(n.sort(function(n,t){return n.TEN_XA>t.TEN_XA}))}})},n.GetXaByIdFromCache=function(t,i){if(!t){i(null);return}n.GetDsXaByHuyenFromCache(-1,function(n){if(!Array.isArray(n)||n.length===0){i(null);return}var r=n.find(function(n){return n&&n.XA_ID==t});i(r||null)})},n.GetDsCoSoTiemChungFromCache=function(n,t){console.log("tinhId: ",n);typeof Storage!="undefined"?sessionStorage.getItem("CO_SO_TIEM_CHUNG_CACHED")==null?$.ajax({url:"/DonViHanhChinh/DsCoSoTiemChung","async":!0,data:{tinhId:n},type:"GET",success:function(n){sessionStorage.setItem("CO_SO_TIEM_CHUNG_DATA",JSON.stringify(n));sessionStorage.setItem("CO_SO_TIEM_CHUNG_CACHED","true");t(n)}}):t(JSON.parse(sessionStorage.getItem("CO_SO_TIEM_CHUNG_DATA"))):(console.log("tinhId: ",n),$.ajax({url:"/DonViHanhChinh/DsCoSoTiemChung","async":!0,data:{tinhId:n},type:"GET",success:function(n){t(n)}}))},n.GetCoSoTiemChungByIdFromCache=function(t,i){n.GetDsCoSoTiemChungFromCache(function(n){var r=n.filter(function(n){if(n.COSO_ID==t)return n});r!=null&&r.length>0?i(r[0]):i(null)})},n.GetDsCoSoTiemChungByParamsFromCache=function(i,r,u,f,e,o){n.GetDsCoSoTiemChungFromCache(i,function(n){var s;f==t.COSO_HINH_THUC.TRUNG_UONG?s=n.filter(n=>n.HINHTHUC==t.COSO_HINH_THUC.TRUNG_UONG):f==t.COSO_HINH_THUC.KHU_VUC?s=n.filter(n=>n.HINHTHUC==t.COSO_HINH_THUC.KHU_VUC&&n.VUNG_MIEN_ID==e):f==t.COSO_HINH_THUC.TINH?s=n.filter(n=>{if(n.PHANLOAI==t.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_CONG&&n.HINHTHUC==t.COSO_HINH_THUC.TINH&&n.TINH_ID==i)return n}):f==t.COSO_HINH_THUC.HUYEN?s=n.filter(n=>n.PHANLOAI==t.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_CONG&&n.HINHTHUC==t.COSO_HINH_THUC.HUYEN&&n.HUYEN_ID==r):f==t.COSO_HINH_THUC.XA?s=n.filter(n=>n.PHANLOAI==t.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_CONG&&n.HINHTHUC==t.COSO_HINH_THUC.XA&&n.XA_ID==u):f==t.COSO_HINH_THUC.CO_SO_DICH_VU?s=u>0?n.filter(n=>n.PHANLOAI==t.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_TU&&n.TINH_ID==i&&n.HUYEN_ID==r&&n.XA_ID==u):r>0?n.filter(n=>n.PHANLOAI==t.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_TU&&n.TINH_ID==i&&n.HUYEN_ID==r):n.filter(n=>n.PHANLOAI==t.COSO_PHAN_LOAI.CO_SO_TIEM_CHUNG_TU&&n.TINH_ID==i):f==t.COSO_HINH_THUC.BENH_VIEN&&(s=u>0?n.filter(n=>n.PHANLOAI==t.COSO_PHAN_LOAI.BENH_VIEN&&n.TINH_ID==i&&n.HUYEN_ID==r&&n.XA_ID==u):r>0?n.filter(n=>n.PHANLOAI==t.COSO_PHAN_LOAI.BENH_VIEN&&n.TINH_ID==i&&n.HUYEN_ID==r):n.filter(n=>n.PHANLOAI==t.COSO_PHAN_LOAI.BENH_VIEN&&n.TINH_ID==i));s&&s.length>0?o(s):o(s)})},n.GetDsCoSoTiemChungByTinhHuyenXaFromCache=function(t,i,r,u){n.GetDsCoSoTiemChungFromCache(t,function(n){var f;f=r>0?n.filter(n=>n.TINH_ID==t&&n.HUYEN_ID==i&&n.XA_ID==r):i>0?n.filter(n=>n.TINH_ID==t&&n.HUYEN_ID==i):n.filter(n=>n.TINH_ID==t);f&&f.length>0?u(f):u(f)})},n.GetDsQuocGiaFromCache=function(n){typeof Storage!="undefined"?sessionStorage.getItem("QUOC_GIA_CACHED")==null?$.ajax({url:"/DonViHanhChinh/DsQuocGia","async":!0,type:"GET",success:function(t){sessionStorage.setItem("QUOC_GIA_DATA",JSON.stringify(t));sessionStorage.setItem("QUOC_GIA_CACHED","true");n(t)}}):n(JSON.parse(sessionStorage.getItem("QUOC_GIA_DATA"))):$.ajax({url:"/DonViHanhChinh/DsQuocGia","async":!0,type:"GET",success:function(t){n(t)}})},n.GetQuocGiaByIdFromCache=function(t,i){n.GetDsQuocGiaFromCache(function(n){var r=n.filter(function(n){if(n.QUOC_GIA_ID==t)return n});r!=null&&r.length>0?i(r[0]):i(null)})},n}();;
